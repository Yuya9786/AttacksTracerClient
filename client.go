package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	_ "github.com/lib/pq"
)

type Data struct {
	Tag    string          `db:"tag"`
	Time   string          `db:"name"`
	Record json.RawMessage `db:"record"`
}

type Relation struct {
	To          string      `json:"to"`
	From        string      `json:"from"`
	Type        string      `json:"type"`
	Annotations Annotations `json:"annotations"`
}

type Annotations struct {
	ID           string `json:"id"`
	Epoch        int    `json:"epoch"`
	flags        string `json:"flags"`
	Allowed      string `json:"allowed"`
	BootId       int    `json:"boot_id"`
	Date         string `json:"cf:date"`
	Jiffies      string `json:"jiffies"`
	TaskID       string `json:"task_id"`
	ToType       string `json:"to_type"`
	FromType     string `json:"from_type"`
	RelationID   string `json:"relation_id"`
	MachineID    string `json:"cf:machine_id"`
	RelationType string `json:"relation_type"`
}

var (
	actorMap map[string]interface{}
)

func client() error {
	var connectionString string = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", *dbaddr, *dbport, *dbuser, *dbpassword, *dbname)

	db, err := sql.Open("postgres", connectionString)
	defer db.Close()
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}
	fmt.Println("Successfully created connection to database.")

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	c := pb.NewMalwareSimulatorClient(conn)

	actorMap = map[string]interface{}{}

	err = prepare(c, ctx)
	if err != nil {
		return err
	}

	// err = sendPacket(db, c, ctx)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func prepare(c pb.MalwareSimulatorClient, ctx context.Context) error {
	node1 := pb.AddNodeRequest{
		Name: "cf:1",
	}

	resultNode, err := c.AddNode(ctx, &node1)
	if err != nil {
		return err
	}

	println(resultNode.String())

	actorMap[resultNode.GetName()] = resultNode

	network := pb.AddNetworkRequest{
		Name:       "network",
		Address:    "172.31.17.125",
		SubnetMask: 20,
	}

	resultNet, err := c.AddNetwork(ctx, &network)
	if err != nil {
		return err
	}

	println(resultNet.String())

	actorMap[resultNet.GetNetworkName()] = resultNet

	return nil
}

// func sendPacket(db *sql.DB, c pb.MalwareSimulatorClient, ctx context.Context) error {
// 	rows, err := db.Query("select * from data where record->'annotations'->>'relation_type'='send_packet';")
// 	if err != nil {
// 		return err
// 	}

// 	for rows.Next() {
// 		var data Data
// 		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
// 		if err != nil {
// 			return err
// 		}

// 		var relation Relation
// 		if err = json.Unmarshal(data.Record, &relation); err != nil {
// 			return err
// 		}

// 		if relation.Annotations.ToType == "socket" {
// 			err := socket(db, c, ctx, relation.To)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func socket(db *sql.DB, c pb.MalwareSimulatorClient, ctx context.Context, socketID string) error {
// 	query := fmt.Sprintf("select * from data where record->>'from'='%v';", socketID)
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return err
// 	}

// 	for rows.Next() {
// 		var data Data
// 		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
// 		if err != nil {
// 			return err
// 		}

// 		var relation Relation
// 		if err = json.Unmarshal(data.Record, &relation); err != nil {
// 			return err
// 		}

// 		if relation.Type == "WasGenaratedBy" && relation.Annotations.RelationType == "send" {
// 			task, err := task(db, c, ctx, relation.To)
// 			if err != nil {
// 				return err
// 			}

// 		}
// 	}
// }

// func task(db *sql.DB, c pb.MalwareSimulatorClient, ctx context.Context, id string) (int, error) {
// 	v, ok := actorMap[id]
// 	if ok {
// 		return v, nil
// 	}

// 	query := fmt.Sprintf("select * from data where record->>'from'='%v';", id)
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return 0, err
// 	}

// 	for rows.Next() {
// 		var data Data
// 		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
// 		if err != nil {
// 			return 0, err
// 		}

// 		var relation Relation
// 		if err = json.Unmarshal(data.Record, &relation); err != nil {
// 			return 0, err
// 		}

// 		if relation.Type == "WasInformedBy" && relation.Annotations.RelationType == "version_activity" {
// 			return task(db, relation.To)
// 		}
// 	}

// 	request := pb.AddApplicationRequest{
// 		Name: +"",
// 	}
// 	c.AddApplication(ctx, &pb.AddApplicationRequest{})

// 	return v, nil
// }
