package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/grpc-up-and-running/samples/ch07/grpc-docker/go/proto-gen"
	"github.com/mercari/testdeck"
	"github.com/mercari/testdeck/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestMain(m *testing.M) {
	service.Start(m)
}

const (
	// ローカルで動かすときこちら有効化
	//address = "localhost:50051"
	// クラスターで動かすときこちら有効化
	address = "productinfo.default:50051"
	bufSize = 1024 * 1024
)

var listener *bufconn.Listener

func initGRPCServerHTTP2() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// Initialization of BufConn.
// Package bufconn provides a net.Conn implemented by a buffer and related dialing and listening functionality.
// func initGRPCServerBuffConn() {
// 	listener = bufconn.Listen(bufSize)
// 	s := grpc.NewServer()
// 	pb.RegisterProductInfoServer(s, &server{})
// 	// Register reflection service on gRPC server.
// 	reflection.Register(s)
// 	go func() {
// 		if err := s.Serve(listener); err != nil {
// 			log.Fatalf("failed to serve: %v", err)
// 		}
// 	}()

// }

// Conventional test that starts a gRPC server and client test the service with RPC
// func TestServer_AddProduct(t *testing.T) {
// 	initGRPCServerHTTP2() // Starting a conventional gRPC server runs on HTTP2
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewProductInfoClient(conn)

// 	// Contact the server and print out its response.
// 	name := "Sumsung S10"
// 	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
// 	price := float32(700.0)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()
// 	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
// 	if err != nil {
// 		log.Fatalf("Could not add product: %v", err)
// 	}
// 	log.Printf("Res %s", r.Value)
// }

func TestServer_AddProduct(t *testing.T) {
	var conn *grpc.ClientConn
	var c pb.ProductInfoClient
	var name, description string
	var price float32
	var res *wrapperspb.StringValue
	var addProductError error

	test := testdeck.TestCase{}
	test.Arrange = func(t *testdeck.TD) {
		//initGRPCServerHTTP2() // Starting a conventional gRPC server runs on HTTP2
		var err error
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c = pb.NewProductInfoClient(conn)

		name = "Sumsung S10"
		description = "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
		price = float32(700.0)
	}
	test.Act = func(t *testdeck.TD) {
		defer conn.Close()
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		res, addProductError = c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
	}
	test.Assert = func(t *testdeck.TD) {
		if addProductError != nil {
			t.Errorf("want: %v, got: %v", nil, addProductError)
		}
		log.Printf("Res %s", res.Value)
	}

	// Finally, call Test() to start the test
	testdeck.Test(t, &test)
}

// Test written using Buffconn
// func TestServer_AddProductBufConn(t *testing.T) {
// 	ctx := context.Background()
// 	initGRPCServerBuffConn()
// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewProductInfoClient(conn)

// 	// Contact the server and print out its response.
// 	name := "Sumsung S10"
// 	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
// 	price := float32(700.0)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()
// 	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
// 	if err != nil {
// 		log.Fatalf("Could not add product: %v", err)
// 	}
// 	log.Printf(r.Value)
// }
