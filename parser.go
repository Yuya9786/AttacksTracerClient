package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	"github.com/pkg/errors"
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

func sendReceivePacket(db *sql.DB, c pb.MalwareSimulatorClient) error {
	rows, err := db.Query(
		"select * from data where record->'annotations'->>'relation_type'='send_packet' or record->'annotations'->>'relation_type'='receive_packet';")
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

		if relation.Annotations.RelationType == "send_packet" {
			sendPacketQuery(db, c, &relation)
		} else if relation.Annotations.RelationType == "receive_packet" {
			receivePacketQuey(db, c, &relation)
		}
	}

	return nil
}

func sendPacketQuery(db *sql.DB, c pb.MalwareSimulatorClient, relation *Relation) error {
	if relation.Annotations.ToType != "socket" || relation.Annotations.FromType != "packet" {
		return nil
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
	if status := addrsssCheck(receiver[0]); status == 0 {
		resultNode, err := addNode(c, receiver[0])
		if err != nil {
			return errors.Wrap(err, "failed to addNode")
		}
		unknownNodes[resultNode.GetName()] = resultNode

		senderAddr, err := inetAddress(sender[0])
		if err != nil {
			return errors.Wrap(err, "failed to inetAddress")
		}
		receiverAddr, err := inetAddress(receiver[0])
		if err != nil {
			return errors.Wrap(err, "failed to inetAddress")
		}

		receiverMask := -1
		var receiverNet *pb.Network
		senderMask := -1
		var senderNet *pb.Network
		for _, v := range networks {
			netAddr, err := inetAddress(v.Address)
			if err != nil {
				return errors.Wrap(err, "failed to inetAddress")
			}

			if netAddr.isSameNetwork(senderAddr, int(v.SubnetMask)) {
				if v.GetSubnetMask() > int32(senderMask) {
					senderNet = v
					senderMask = int(v.GetSubnetMask())
				}
			}

			if netAddr.isSameNetwork(receiverAddr, int(v.SubnetMask)) {
				if v.GetSubnetMask() > int32(receiverMask) {
					receiverNet = v
					receiverMask = int(v.GetSubnetMask())
				}
			}
		}
		if senderNet == nil {
			return errors.New("sender network is not found")
		}

		if receiverNet == nil {
			return errors.New("receiver network is not found")
		}

		_, err = makeConnection(c, int(resultNode.GetId()), int(receiverNet.GetId()), receiver[0], int(receiverNet.SubnetMask))
		if err != nil {
			return errors.Wrap(err, "failed to makeConnection")
		}
		resultNode, err = updateNodeInfo(c, int(resultNode.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to updateNodeInfo")
		}
		unknownNodes[receiver[0]] = resultNode

		router := unknownNodes["Router"]
		r, err := updateNodeInfo(c, int(router.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to updateNodeInfo")
		}
		unknownNodes["Router"] = r

		var nexthop *string
		for _, v := range r.Connections {
			if v.SubnetMask == 0 {
				nexthop = &v.Address
			}
		}
		if nexthop == nil {
			return errors.New("router address is not found")
		}
		_, err = addRoute(c, int(resultNode.GetId()), senderNet.GetAddress(), int(senderNet.GetSubnetMask()), *nexthop, int(networks["Internet"].GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to addRoute")
		}
	}

	_, err = sendPacket(c, task, sender[0], srcPort, receiver[0], dstPort, []byte(fmt.Sprint(packet.PacketID)))
	if err != nil {
		return errors.Wrap(err, "failed to SendPacket")
	}

	return nil
}

func receivePacketQuey(db *sql.DB, c pb.MalwareSimulatorClient, relation *Relation) error {
	if relation.Annotations.FromType != "socket" || relation.Annotations.ToType != "packet" {
		return nil
	}

	packet, err := getPacket(db, relation.To)
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
	if status := addrsssCheck(sender[0]); status == 0 {
		resultNode, err := addNode(c, sender[0])
		if err != nil {
			return errors.Wrap(err, "failed to addNode")
		}
		unknownNodes[resultNode.GetName()] = resultNode

		senderAddr, err := inetAddress(sender[0])
		if err != nil {
			return errors.Wrap(err, "failed to inetAddress")
		}
		receiverAddr, err := inetAddress(receiver[0])
		if err != nil {
			return errors.Wrap(err, "failed to inetAddress")
		}

		receiverMask := -1
		var receiverNet *pb.Network
		senderMask := -1
		var senderNet *pb.Network
		for _, v := range networks {
			netAddr, err := inetAddress(v.Address)
			if err != nil {
				return errors.Wrap(err, "failed to inetAddress")
			}

			if netAddr.isSameNetwork(senderAddr, int(v.SubnetMask)) {
				if v.GetSubnetMask() > int32(senderMask) {
					senderNet = v
					senderMask = int(v.GetSubnetMask())
				}
			}

			if netAddr.isSameNetwork(receiverAddr, int(v.SubnetMask)) {
				if v.GetSubnetMask() > int32(receiverMask) {
					receiverNet = v
					receiverMask = int(v.GetSubnetMask())
				}
			}
		}
		if senderNet == nil {
			return errors.New("sender network is not found")
		}

		if receiverNet == nil {
			return errors.New("receiver network is not found")
		}

		_, err = makeConnection(c, int(resultNode.GetId()), int(receiverNet.GetId()), sender[0], int(senderNet.SubnetMask))
		if err != nil {
			return errors.Wrap(err, "failed to makeConnection")
		}
		resultNode, err = updateNodeInfo(c, int(resultNode.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to updateNodeInfo")
		}
		unknownNodes[sender[0]] = resultNode

		router := unknownNodes["Router"]
		r, err := updateNodeInfo(c, int(router.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to updateNodeInfo")
		}
		unknownNodes["Router"] = r

		var nexthop *string
		for _, v := range r.Connections {
			if v.SubnetMask == 0 {
				nexthop = &v.Address
			}
		}
		if nexthop == nil {
			return errors.New("router address is not found")
		}
		_, err = addRoute(c, int(resultNode.GetId()), receiverNet.GetAddress(), int(receiverNet.GetSubnetMask()), *nexthop, int(networks["Internet"].GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to addRoute")
		}

		task, err := addApplication(c, fmt.Sprint(rand.Int()), int(resultNode.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to addApplication")
		}

		_, err = sendPacket(c, int(task.GetId()), sender[0], srcPort, receiver[0], dstPort, []byte(fmt.Sprint(packet.PacketID)))
		if err != nil {
			return errors.Wrap(err, "failed to sendPacket")
		}

		_, err = remove(c, int(task.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to remove")
		}
	} else if status == 2 {

		resultNode := unknownNodes[sender[0]]

		task, err := addApplication(c, fmt.Sprint(rand.Int()), int(resultNode.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to addApplication")
		}

		_, err = sendPacket(c, int(task.GetId()), sender[0], srcPort, receiver[0], dstPort, []byte(fmt.Sprint(packet.PacketID)))
		if err != nil {
			return errors.Wrap(err, "failed to sendPacket")
		}

		_, err = remove(c, int(task.GetId()))
		if err != nil {
			return errors.Wrap(err, "failed to remove")
		}
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
	v, ok := applications[id]
	if ok {
		return int(v.GetId()), nil
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

		node := knownNodes[vertice.Annotations.MachineID]
		application, err := addApplication(c, fmt.Sprintf("%s:%d", vertice.Annotations.MachineID, vertice.Annotations.PID), int(node.GetId()))
		if err != nil {
			return 0, errors.Wrap(err, "failed to AddApplication")
		}
		applications[id] = application

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

func addrsssCheck(address string) int {
	for _, v := range knownNodes {
		for _, av := range v.Connections {
			if av.Address == address {
				return 1
			}
		}
	}

	for _, v := range unknownNodes {
		for _, av := range v.Connections {
			if av.Address == address {
				return 2
			}
		}
	}

	return 0
}
