# Adapter for connecting Ory Oathkeeper with SpiceDB

It is not possible to directly connect [Ory Oathkeeper](https://github.com/ory/oathkeeper) to [SpiceDB](https://github.com/authzed/spicedb). This adapter allows you to do so.

Oathkeeper [relies](https://www.ory.sh/docs/oathkeeper/pipeline/authz#remote_json) on the http response code from the authorization server (200 OK, 403 Forbidden). The SpiceDB HTTP server does not return the permissions state with the HTTP Response code, but does so in the body of the response.

## Getting Started

TODO

### Environment Variables

| Name            | Default | Description                                                           |
| --------------- | ------- | --------------------------------------------------------------------- |
| PORT            | 50150   | Port to listen on                                                     |
| SPICE_DB_SECRET |         | **(required)** GRPC preshared key for the SpiceDB instance            |
| SPICE_DB_URL    |         | **(required)** URL to the SpiceDB instance. *Example: localhost:50051 |