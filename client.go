package main

import (
	"database/sql"
	"fmt"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	networks     = map[string]*pb.Network{}
	knownNodes   = map[string]*pb.Node{}
	unknownNodes = map[string]*pb.Node{}
	applications = map[string]*pb.Application{}
)

func client() error {
	var connectionString string = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", *dbaddr, *dbport, *dbuser, *dbpassword, *dbname)

	db, err := sql.Open("postgres", connectionString)
	defer db.Close()
	if err != nil {
		return errors.Wrap(err, "failed to open db")
	}

	if err = db.Ping(); err != nil {
		return errors.Wrap(err, "failed to connect to db.")
	}
	fmt.Println("Successfully created connection to database.")

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to connect")
	}
	defer conn.Close()
	fmt.Println("Successfully created connection to the AttacksTracer.")

	c := pb.NewMalwareSimulatorClient(conn)

	err = prepare(c)
	if err != nil {
		return errors.Wrap(err, "failed to prepare")
	}

	err = sendReceivePacket(db, c)
	if err != nil {
		return errors.Wrap(err, "failed to sendPacket")
	}

	return nil
}
