package main

import (
	"fmt"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	"github.com/pkg/errors"
)

func prepare(c pb.MalwareSimulatorClient) error {
	fmt.Println("Start to prepare simulating environments.")

	resultNode, err := addNode(c, "cf:1")
	if err != nil {
		return errors.Wrap(err, "failed to send AddNode")
	}

	knownNodes[resultNode.GetName()] = resultNode

	resultNet, err := addNetwork(c, "LAN", "172.31.16.0", 20)
	if err != nil {
		return errors.Wrap(err, "failed to send AddNetwork")
	}
	networks[resultNet.GetNetworkName()] = resultNet

	nodeID := knownNodes["cf:1"].GetId()
	netID := networks["LAN"].GetId()
	_, err = makeConnection(c, int(nodeID), int(netID), "172.31.17.125", 20)
	if err != nil {
		return errors.Wrap(err, "failed to send MakeConnection")
	}

	_, err = addRoute(c, int(nodeID), "0.0.0.0", 0, "172.31.16.1", int(netID))
	if err != nil {
		return errors.Wrap(err, "failed to send AddRoute")
	}

	resultNet, err = addNetwork(c, "Internet", "0.0.0.0", 0)
	if err != nil {
		return errors.Wrap(err, "failed to send AddNetwork")
	}
	networks[resultNet.GetNetworkName()] = resultNet

	resultNode, err = addNode(c, "Router")
	if err != nil {
		return errors.Wrap(err, "failed to send AddNode")
	}
	unknownNodes[resultNode.GetName()] = resultNode

	nodeID = unknownNodes["Router"].GetId()
	netID = networks["Internet"].GetId()
	_, err = makeConnection(c, int(nodeID), int(netID), "18.181.233.66", 0)
	if err != nil {
		return errors.Wrap(err, "failed to send MakeConnection")
	}

	netID = networks["LAN"].GetId()
	_, err = makeConnection(c, int(nodeID), int(netID), "172.31.16.1", 20)
	if err != nil {
		return errors.Wrap(err, "failed to send MakeConnection")
	}

	return nil
}
