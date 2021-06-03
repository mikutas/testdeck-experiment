# Multi stage build

# Build stage I : Go lang and Alpine Linux is only needed to build the program
#FROM golang:1.11-alpine AS build
FROM golang AS build


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

FROM golang AS build_test
ENV location /go/src/github.com/grpc-up-and-running/samples/ch07/grpc-docker/go
WORKDIR ${location}/server
ADD ./server ${location}/server
ADD ./proto-gen ${location}/proto-gen
RUN CGO_ENABLED=0 go test -c . -o /bin/grpc-productinfo-server-test

FROM scratch AS run_test
COPY --from=build_test /bin/grpc-productinfo-server-test /bin/grpc-productinfo-server-test
ENTRYPOINT ["/bin/grpc-productinfo-server-test"]

# Build stage II : Go binaries are self-contained executables.
FROM scratch
COPY --from=build /bin/grpc-productinfo-server /bin/grpc-productinfo-server


ENTRYPOINT ["/bin/grpc-productinfo-server"]
EXPOSE 50051

