package main

import (
	"flag"
	"log"
	"os"
	"strings"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"

	"google.golang.org/grpc"

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

	}

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// r, err := c.ShowNet(ctx, &pb.ShowNetRequest{})
	// if err != nil {
	// 	log.Fatalf("could not show network: %v", err)
	// }

	// networks := r.GetNetworks()
	// log.Printf("ShowNet: %s", r.GetNetworks())
	// for n := range networks {
	// 	log.Println(n)
	// }

}
