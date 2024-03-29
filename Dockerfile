# Multi stage build

# Build stage I : Go lang and Alpine Linux is only needed to build the program
#FROM golang:1.11-alpine AS build
FROM golang:1.15 AS build


ENV location /go/src/github.com/grpc-up-and-running/samples/ch07/grpc-docker/go

WORKDIR ${location}/server

ADD ./server ${location}/server
ADD ./proto-gen ${location}/proto-gen

#ADD main.go ${location}/server
#ADD ../proto-gen ${location}/proto-gen


# Download all the dependencies
RUN go get -d ./...
# Install the package
RUN go install ./...


RUN CGO_ENABLED=0 go build -o /bin/grpc-productinfo-server

FROM golang:1.15 AS build_test
ENV location /go/src/github.com/grpc-up-and-running/samples/ch07/grpc-docker/go
WORKDIR ${location}/server
ADD ./server ${location}/server
ADD ./proto-gen ${location}/proto-gen
RUN CGO_ENABLED=0 go test -c . -o /bin/grpc-productinfo-server-test
ARG LINKERD_AWAIT_VERSION=v0.2.3
RUN curl -sSLo /tmp/linkerd-await https://github.com/linkerd/linkerd-await/releases/download/release%2F${LINKERD_AWAIT_VERSION}/linkerd-await-${LINKERD_AWAIT_VERSION}-amd64 && \
    chmod 755 /tmp/linkerd-await

FROM scratch AS run_test
COPY --from=build_test /bin/grpc-productinfo-server-test /bin/grpc-productinfo-server-test
COPY --from=build_test /tmp/linkerd-await /linkerd-await
ENTRYPOINT ["/linkerd-await", "--shutdown", "--"]
CMD ["/bin/grpc-productinfo-server-test"]

# Build stage II : Go binaries are self-contained executables.
FROM scratch
COPY --from=build /bin/grpc-productinfo-server /bin/grpc-productinfo-server


ENTRYPOINT ["/bin/grpc-productinfo-server"]
EXPOSE 50051
