# docker-grpc-gateway

This is a gRPC to JSON proxy docker image, automatically providing a REST API for your gRPC services.

## Requirements

`docker` client.

## Usage

Copy proto files in `/proto` directory.

The proto files must be annotated with [google.api.http](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto) for URL to gRPC service mapping.

This project uses [buf](https://buf.build/) for dependencies management. Dependencies are declared in `proto/buf.yaml`.

## Building

```
docker build -t grpc-gateway .
```

## Running

Environment variables :

| Variable           | Description |
| ------------------ | ------------|
| `BACKEND_URI`      | Uri of the gRPC server, following [gRPC name syntax](https://github.com/grpc/grpc/blob/master/doc/naming.md) |
| `BACKEND_TLS`      | "true" => use TLS, else unsecured gRPC. |
| `GATEWAY_ADDRESS`  | Address of this gateway. |
| `GATEWAY_TLS_CERT` | Path to the x509 server certificate used for TLS by the gateway. |
| `GATEWAY_TLS_KEY`  | Path to the key file used for TLS by the gateway. |

```
docker run --rm -v "$(pwd)/test/cert:/cert" \
    -e GATEWAY_ADDRESS="0.0.0.0:8443" \
    -e GATEWAY_TLS_CERT=/cert/server.crt \
    -e GATEWAY_TLS_KEY=/cert/server.key \
    -e BACKEND_URI="dns:///my.grpc.service:1337" \
    -e BACKEND_TLS=true \
    -p 8443:8443 grpc-gateway

curl -k https://localhost:8443/api/v1/collections/User
```

## Forwarding headers

Headers prefixed with `Grpc-Metadata-` will be forwarded to the gRPC server.

```
curl -k -H "Grpc-Metadata-x-user-id: me" https://localhost:8443/api/v1/collections/User
```
