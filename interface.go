package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	"github.com/pkg/errors"
)

func addNode(c pb.MalwareSimulatorClient, name string) (*pb.Node, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.AddNodeRequest{
		Name: name,
	}

	fmt.Printf("AddNode(%+v)\n", request)
	result, err := c.AddNode(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send AddNode")
	}

	return result, nil
}

func updateNodeInfo(c pb.MalwareSimulatorClient, id int) (*pb.Node, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.UpdateNodeInfoRequest{
		Id: int32(id),
	}

	fmt.Printf("UpdateNodeInfo(%+v)\n", request)
	result, err := c.UpdateNodeInfo(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send AddNode")
	}

	return result, nil
}

func addNetwork(c pb.MalwareSimulatorClient, name string, address string, mask int) (*pb.Network, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.AddNetworkRequest{
		Name:       name,
		Address:    address,
		SubnetMask: int32(mask),
	}

	fmt.Printf("AddNetwork(%+v)\n", request)
	result, err := c.AddNetwork(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send AddNetwork")
	}

	return result, nil
}

func addApplication(c pb.MalwareSimulatorClient, name string, nodeID int) (*pb.Application, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.AddApplicationRequest{
		Name:   name,
		NodeID: int32(nodeID),
	}

	fmt.Printf("AddApplication(%+v)\n", request)
	result, err := c.AddApplication(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send AddApplication")
	}

	return result, nil
}

func remove(c pb.MalwareSimulatorClient, id int) (*pb.RemoveReply, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.RemoveRequest{
		Id: int32(id),
	}

	fmt.Printf("Remove(%+v)\n", request)
	result, err := c.Remove(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send Remove")
	}

	return result, nil
}

func makeConnection(c pb.MalwareSimulatorClient, nodeID int, networkID int, address string, mask int) (*pb.Network, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.MakeConnectionRequest{
		NodeID:     int32(nodeID),
		NetworkID:  int32(networkID),
		Address:    address,
		SubnetMask: int32(mask),
	}

	fmt.Printf("MakeConnection(%+v)\n", request)
	result, err := c.MakeConnection(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send MakeConnection")
	}

	return result, nil
}

func addRoute(c pb.MalwareSimulatorClient, nodeID int, dstAddress string, mask int, nexthop string, netID int) (*pb.AddRouteReply, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.AddRouteRequest{
		NodeID:     int32(nodeID),
		DstAddress: dstAddress,
		Mask:       int32(mask),
		Nexthop:    nexthop,
		NetID:      int32(netID),
	}

	fmt.Printf("AddRoute(%+v)\n", request)
	result, err := c.AddRoute(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send AddAddRoute")
	}

	return result, nil
}

func sendPacket(c pb.MalwareSimulatorClient, appID int, srcAddress string, srcPort int,
	dstAddress string, dstPort int, data []byte) (*pb.SendPacketReply, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	request := &pb.SendPacketRequest{
		AppID:      int32(appID),
		SrcAddress: srcAddress,
		SrcPort:    int32(srcPort),
		DstAddress: dstAddress,
		DstPort:    int32(dstPort),
		Data:       data,
	}

	fmt.Printf("SendPacket(%+v)\n", request)
	result, err := c.SendPacket(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send SendPacket")
	}

	return result, nil
}
