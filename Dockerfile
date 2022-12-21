FROM golang:alpine
ARG BUF_VERSION="v1.11.0"
ADD https://github.com/bufbuild/buf/releases/download/$BUF_VERSION/buf-Linux-x86_64 /usr/local/bin/buf
RUN chmod +x /usr/local/bin/buf
WORKDIR /app
COPY . ./
RUN buf generate && go build -o grpc-gateway cmd/grpc-gateway/main.go

FROM alpine
COPY --from=0 /app/grpc-gateway ./
CMD ["/grpc-gateway"]