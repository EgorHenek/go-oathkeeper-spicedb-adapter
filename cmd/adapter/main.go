package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EgorHenek/go-oathkeeper-spicedb-adapter/configs"
	"github.com/EgorHenek/go-oathkeeper-spicedb-adapter/internal/domain"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ReadHeaderTimeout = 5 * time.Second
	ShutdownTimeOut   = 30 * time.Second
)

func main() {
	config := configs.NewConfig()
	var certs grpc.DialOption
	var token grpc.DialOption
	var err error

	if len(config.TLSCertPath) == 0 {
		certs = grpc.WithTransportCredentials(insecure.NewCredentials())
		token = grpcutil.WithInsecureBearerToken(config.SpiceDBSecret)
	} else {
		certs, err = grpcutil.WithCustomCerts(grpcutil.VerifyCA, config.TLSCertPath...)
		if err != nil {
			log.Fatal(err)
		}
		token = grpcutil.WithBearerToken(config.SpiceDBSecret)
	}

	sdbClient, err := authzed.NewClient(
		config.SpiceDBURL,
		certs,
		token,
	)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           run(sdbClient),
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, ShutdownTimeOut)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out... forcing exit.")
			}
		}()

		err = server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func run(client *authzed.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)

	r.Post("/permissions/check", PermissionsCheckHandler(client))

	return r
}

func PermissionsCheckHandler(client *authzed.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &domain.CheckPermissionRequest{}
		err := render.Bind(r, data)
		if err != nil {
			_ = render.Render(w, r, domain.ErrInvalidRequest(err))
			return
		}

		resp, err := client.CheckPermission(context.Background(), &v1.CheckPermissionRequest{
			Consistency: data.Consistency,
			Subject:     data.Subject,
			Resource:    data.Resource,
			Permission:  data.Permission,
		})
		if err != nil {
			_ = render.Render(w, r, domain.ErrInvalidRequest(err))
			return
		}

		allowed := resp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION
		_ = render.Render(w, r, domain.NewCheckPermissionResponse(allowed))
	}
}
