package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Data struct {
	Tag    string          `db:"tag"`
	Time   string          `db:"name"`
	Record json.RawMessage `db:"record"`
}

type Activity struct {
	ID          string              `json:"id"`
	Type        string              `json:"type"`
	Annotations ActivityAnnotations `json:"annotations"`
}

type ActivityAnnotations struct {
	VM           string `json:"vm"`
	PID          int    `json:"pid"`
	RSS          string `json:"rss"`
	VPID         int    `json:"vpid"`
	HW_VM        string `json:"hw_vm"`
	Stime        string `json:"stime"`
	Utime        string `json:"utime"`
	HW_RSS       string `json:"hw_rss"`
	Rbytes       string `json:"rbytes"`
	Secctx       string `json:"secctx"`
	Wbytes       string `json:"wbytes"`
	BootID       int    `json:"boot_id"`
	Date         string `json:"cf:date"`
	Version      int    `json:"version"`
	Epoch        int    `json:"cf:epoch"`
	Taint        string `json:"cf:taint"`
	ObjectID     string `json:"object_id"`
	Jiffies      string `json:"cf:jiffies"`
	ObjectType   string `json:"object_type"`
	CancelWbytes string `json:"cancel_wbytes"`
	MachineID    string `json:"cf:machine_id"`
}

type Relation struct {
	To          string              `json:"to"`
	From        string              `json:"from"`
	Type        string              `json:"type"`
	Annotations RelationAnnotations `json:"annotations"`
}

type RelationAnnotations struct {
	ID           string `json:"id"`
	Epoch        int    `json:"epoch"`
	Flags        string `json:"flags"`
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

	actorMap = map[string]interface{}{}

	err = prepare(c)
	if err != nil {
		return errors.Wrap(err, "failed to prepare")
	}

	err = sendPacket(db, c)
	if err != nil {
		return errors.Wrap(err, "failed to sendPacket")
	}

	return nil
}

func sendPacket(db *sql.DB, c pb.MalwareSimulatorClient) error {
	rows, err := db.Query("select * from data where record->'annotations'->>'relation_type'='send_packet';")
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "failed to query on db")
	}

	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return errors.Wrap(err, "failed to scan rows")
		}

		var relation Relation
		if err = json.Unmarshal(data.Record, &relation); err != nil {
			return errors.Wrap(err, "failed to unmarshal")
		}

		if relation.Annotations.ToType != "socket" || relation.Annotations.FromType != "packet" {
			continue
		}

		task, err := socketToTask(db, c, relation.To)
		if err != nil {
			return errors.Wrap(err, "failed to socketToTask")
		}

		packet, err := getPacket(db, relation.From)
		if err != nil {
			return errors.Wrap(err, "failed to getPacket")
		}
		sender := strings.Split(packet.Sender, ":")
		srcPort, err := strconv.Atoi(sender[1])
		if err != nil {
			return errors.Wrap(err, "failed to convert to int")
		}
		receiver := strings.Split(packet.Receiver, ":")
		dstPort, err := strconv.Atoi(receiver[1])
		if err != nil {
			return errors.Wrap(err, "failed to convert to int")
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Second,
		)
		defer cancel()

		request := &pb.SendPacketRequest{
			AppID:      int32(task),
			SrcAddress: sender[0],
			SrcPort:    int32(srcPort),
			DstAddress: receiver[0],
			DstPort:    int32(dstPort),
			Data:       []byte(string(packet.PacketID)),
		}
		fmt.Printf("SendPacket(%+v)", request)
		_, err = c.SendPacket(ctx, request)
		if err != nil {
			return errors.Wrap(err, "failed to SendPacket")
		}
		fmt.Printf("Succeeded to SendPacket(%+v)", request)
	}

	return nil
}

func socketToTask(db *sql.DB, c pb.MalwareSimulatorClient, socketID string) (int, error) {
	query := fmt.Sprintf("select * from data where record->>'from'='%v';", socketID)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return 0, errors.Wrap(err, "failed to query on db")
	}

	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return 0, errors.Wrap(err, "failed to scan rows")
		}

		var relation Relation
		if err = json.Unmarshal(data.Record, &relation); err != nil {
			return 0, errors.Wrap(err, "failed to unmarshal")
		}

		if relation.Type == "WasGeneratedBy" &&
			(relation.Annotations.RelationType == "send" || relation.Annotations.RelationType == "socket_create") {
			task, err := task(db, c, relation.To)
			if err != nil {
				return 0, errors.Wrap(err, "failed to find task")
			}
			return task, nil
		} else if relation.Type == "WasDerivedFrom" && relation.Annotations.RelationType == "version_entity" {
			return socketToTask(db, c, relation.To)
		}
	}

	return 0, errors.New(fmt.Sprintf("not found task from socketID: %s", socketID))
}

func task(db *sql.DB, c pb.MalwareSimulatorClient, id string) (int, error) {
	v, ok := actorMap[id]
	if ok {
		task := v.(*pb.Application)
		return int(task.Id), nil
	}

	query := fmt.Sprintf("select * from data where record->>'from'='%v';", id)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return 0, errors.Wrap(err, "failed to query on db")
	}

	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return 0, errors.Wrap(err, "failed to scan rows")
		}

		var relation Relation
		if err = json.Unmarshal(data.Record, &relation); err != nil {
			return 0, errors.Wrap(err, "failed to unmarshal")
		}

		if relation.Type == "WasInformedBy" && relation.Annotations.RelationType == "version_activity" {
			return task(db, c, relation.To)
		}
	}

	// Not found a parent task vertice
	query = fmt.Sprintf("select * from data where record->>'id'='%v';", id)
	rows, err = db.Query(query)
	defer rows.Close()
	if err != nil {
		return 0, errors.Wrap(err, "failed to connect to db")
	}

	if rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return 0, errors.Wrap(err, "failed to scan rows")
		}

		var vertice Activity
		if err = json.Unmarshal(data.Record, &vertice); err != nil {
			return 0, errors.Wrap(err, "failed to unmarshal")
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Second,
		)
		defer cancel()

		node := actorMap[vertice.Annotations.MachineID].(*pb.Node)
		request := pb.AddApplicationRequest{
			Name:   fmt.Sprintf("%s:%d", vertice.Annotations.MachineID, vertice.Annotations.PID),
			NodeID: node.GetId(),
		}
		fmt.Printf("Send AddApplication(%+v)\n", &request)
		application, err := c.AddApplication(ctx, &request)
		if err != nil {
			return 0, errors.Wrap(err, "failed to AddApplication")
		}
		fmt.Printf("Succeded to send AddApplication(%+v)\n", &request)

		actorMap[id] = application

		return int(application.GetId()), nil
	}

	return 0, errors.New("not found task")
}

type PacketData struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Annotations Packet `json:"annotations"`
}

type Packet struct {
	Seq        int    `json:"seq"`
	Ih_len     int    `json:"ih_len"`
	Sender     string `json:"sender"`
	Date       string `json:"cf:date"`
	Jiffies    string `json:"jiffies"`
	Receiver   string `json:"receiver"`
	PacketID   int    `json:"packet_id"`
	ObjectType string `json:"object_type"`
}

func getPacket(db *sql.DB, id string) (*Packet, error) {
	var packet *Packet

	query := fmt.Sprintf("select * from data where record->>'id'='%s';", id)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to query on db")
	}

	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan rows")
		}

		var packetData PacketData
		if err = json.Unmarshal(data.Record, &packetData); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}

		packet = &packetData.Annotations
	}

	if packet == nil {
		return nil, errors.New("not found packet")
	}

	return packet, nil
}
