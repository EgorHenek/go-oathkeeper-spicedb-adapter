# Adapter for connecting Ory Oathkeeper with SpiceDB

It is not possible to directly connect [Ory Oathkeeper](https://github.com/ory/oathkeeper) to [SpiceDB](https://github.com/authzed/spicedb). This adapter allows you to do so.

Oathkeeper [relies](https://www.ory.sh/docs/oathkeeper/pipeline/authz#remote_json) on the http response code from the authorization server (200 OK, 403 Forbidden). The SpiceDB HTTP server does not return the permissions state with the HTTP Response code, but does so in the body of the response.

## Getting Started

### Installation

#### Binaries

Go to the [release page](https://github.com/EgorHenek/go-oathkeeper-spicedb-adapter/releases) and download the binary for your platform.

For example, on Linux:

```bash
curl -L https://github.com/EgorHenek/go-oathkeeper-spicedb-adapter/releases/download/v1.0.1/adapter_1.0.1_linux_amd64 -o go-oathkeeper-spicedb-adapter
chmod +x go-oathkeeper-spicedb-adapter
```

and run:

```bash
./go-oathkeeper-spicedb-adapter
```

#### Docker

```bash
docker pull henek/go-oathkeeper-spicedb-adapter:latest
docker run -p 50150:50150 --name osadapter -e SPICE_DB_URL=http://spicedb:50051 -e SPICE_DB_SECRET=topsecret henek/go-oathkeeper-spicedb-adapter
```

**Note:** Don't use the `latest` tag in a production environment.

### Docker Compose

You can see an example environment deployment in the /deployments directory.

## Usage examples

```bash
curl -i -X POST --json '{"resource": {"object_type": "beer", "object_id": "1"}, "permission": "drink", "subject": {"object": {"object_type": "user", "object_id": "1"}}}' http://localhost:50150/permissions/check
```

```bash
HTTP/1.1 403 Forbidden
Content-Type: application/json
Date: Mon, 21 Aug 2023 06:24:31 GMT
Content-Length: 18
```

## Configuration

An example of integrating the adapter into oathkeeper for the whoami service can be seen in the [deployments](go-oathkeeper-spicedb-adapter/tree/master/deployments).

### Environment Variables

| Name            | Default | Description                                                                                                   |
| --------------- | ------- | ------------------------------------------------------------------------------------------------------------- |
| PORT            | 50150   | Port to listen on                                                                                             |
| SPICE_DB_SECRET |         | **(required)** GRPC preshared key for the SpiceDB instance                                                    |
| SPICE_DB_URL    |         | **(required)** URL to the SpiceDB instance. *Example*: localhost:50051                                        |
| TLS_CERT_PATH   | empty   | **(optional)** Path to a TLS certificate files. If this value is empty, an unencrypted connection will be used. *Example*: /mnt/tls/ca.crt,/mnt/tls/tls.crt |