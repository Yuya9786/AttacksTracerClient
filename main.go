package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	"github.com/hpcloud/tail"
	"github.com/mattn/go-scan"
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

	machineId, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		log.Fatalf("cannot open /etc/machine-id: %v", err)
	}

	_, err = c.AddNode(ctx, &pb.AddNodeRequest{Name: string(machineId), Address: "0.0.0.0"})
	if err != nil {
		log.Fatalf("could not show network: %v", err)
	}

	t, err := tail.TailFile(*file, tail.Config{ReOpen: true, Follow: true})
	if err != nil {
		log.Fatalf("cannot tail file: %v", err)
	}

	for line := range t.Lines {
		var s string
		js := strings.NewReader(line.Text)
		if err := scan.ScanJSON(js, "/type", &s); err != nil {
			log.Fatalf("failed to scan json: %v", err)
		}
		if s == "Activity" {
			var id string
			if err := scan.ScanJSON(js, "/id", &id); err != nil {
				log.Fatalf("failed to scan json: %v", err)
			}

			_, err := c.AddApplication(ctx, &pb.AddApplicationRequest{Name: id})
			if err != nil {
				log.Fatalf("failed to send a gRPC request.")
			}
		}
		log.Println(s)
	}

}
