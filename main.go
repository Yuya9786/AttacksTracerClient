package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:8080", "The address to connect to")
	file = flag.String("file", "/tmp/audit.log", "Log file.")
)

func main() {
	flag.Parse()
	if _, err := os.Stat(*file); err != nil {
		log.Fatalf("no such a file: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMalwareSimulatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ShowNet(ctx, &pb.ShowNetRequest{})
	if err != nil {
		log.Fatalf("could not show net: %v", err)
	}

	networks := r.GetNetworks()
	log.Printf("ShowNet: %s", r.GetNetworks())
	for n := range networks {
		log.Println(n)
	}

}
