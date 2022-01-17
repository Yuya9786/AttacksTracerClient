package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	"github.com/pkg/errors"
)

func prepare(c pb.MalwareSimulatorClient) error {
	fmt.Println("Start to prepare simulating environments.")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	node1 := pb.AddNodeRequest{
		Name: "cf:1",
	}
	fmt.Println("Send AddNode(cf:1)")
	resultNode, err := c.AddNode(ctx, &node1)
	if err != nil {
		return errors.Wrap(err, "failed to send AddNode")
	}

	actorMap[resultNode.GetName()] = resultNode

	network := pb.AddNetworkRequest{
		Name:       "LAN",
		Address:    "172.31.16.0",
		SubnetMask: 20,
	}
	fmt.Printf("Send AddNetwork(%+v)\n", &network)
	resultNet, err := c.AddNetwork(ctx, &network)
	if err != nil {
		return errors.Wrap(err, "failed to send AddNetwork")
	}
	actorMap[resultNet.GetNetworkName()] = resultNet

	makeConnectionRequest := pb.MakeConnectionRequest{
		NodeID:     actorMap[resultNode.GetName()].(*pb.Node).GetId(),
		NetworkID:  actorMap[resultNet.GetNetworkName()].(*pb.Network).GetId(),
		Address:    "172.31.17.125",
		SubnetMask: 20,
	}
	fmt.Printf("Send MakeConnection(%+v)\n", &makeConnectionRequest)
	resultNet, err = c.MakeConnection(ctx, &makeConnectionRequest)
	if err != nil {
		return errors.Wrap(err, "failed to send MakeConnection")
	}


	addRouteRequest := pb.AddRouteRequest {
		NodeID: actorMap[resultNode.GetName()].(*pb.Node).GetId(),
		DstAddress: "0.0.0.0",
		Mask: 0,
		Nexthop: "172.31.16.1",
		NetID: actorMap[network.Name].(*pb.Network).GetId(),
	}
	fmt.Printf("Send AddRoute(%+v)\n", &addRouteRequest)
	_, err = c.AddRoute(ctx, &addRouteRequest)
	if err != nil {
		return errors.Wrap(err, "failed to send AddRoute")
	}

	internet := pb.AddNetworkRequest{
		Name:       "Internet",
		Address:    "0.0.0.0",
		SubnetMask: 0,
	}

	fmt.Printf("Send AddNetwork(%+v)\n", &internet)
	resultNet, err = c.AddNetwork(ctx, &internet)
	if err != nil {
		return errors.Wrap(err, "failed to send AddNetwork")
	}
	actorMap[resultNet.GetNetworkName()] = resultNet

	router := pb.AddNodeRequest{
		Name: "Router",
	}

	fmt.Println("Send AddNode(Router))")
	resultNode, err = c.AddNode(ctx, &router)
	if err != nil {
		return errors.Wrap(err, "failed to send AddNode")
	}
	actorMap[resultNode.GetName()] = resultNode

	makeConnectionRequest = pb.MakeConnectionRequest{
		NodeID:     actorMap[resultNode.GetName()].(*pb.Node).GetId(),
		NetworkID:  actorMap[internet.Name].(*pb.Network).GetId(),
		Address:    "18.181.233.66",
		SubnetMask: 24,
	}

	fmt.Printf("Send MakeConnection(%+v)\n", &makeConnectionRequest)
	_, err = c.MakeConnection(ctx, &makeConnectionRequest)
	if err != nil {
		return errors.Wrap(err, "failed to send MakeConnection")
	}

	makeConnectionRequest = pb.MakeConnectionRequest{
		NodeID:     actorMap[resultNode.GetName()].(*pb.Node).GetId(),
		NetworkID:  actorMap[network.Name].(*pb.Network).GetId(),
		Address:    "172.31.16.1",
		SubnetMask: 20,
	}

	fmt.Printf("Send MakeConnection(%+v)\n", &makeConnectionRequest)
	_, err = c.MakeConnection(ctx, &makeConnectionRequest)
	if err != nil {
		return errors.Wrap(err, "failed to send MakeConnection")
	}

	return nil
}
