// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MalwareSimulatorClient is the client API for MalwareSimulator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MalwareSimulatorClient interface {
	AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*Node, error)
	UpdateNodeInfo(ctx context.Context, in *UpdateNodeInfoRequest, opts ...grpc.CallOption) (*Node, error)
	AddNetwork(ctx context.Context, in *AddNetworkRequest, opts ...grpc.CallOption) (*Network, error)
	AddApplication(ctx context.Context, in *AddApplicationRequest, opts ...grpc.CallOption) (*Application, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveReply, error)
	MakeConnection(ctx context.Context, in *MakeConnectionRequest, opts ...grpc.CallOption) (*Network, error)
	Infect(ctx context.Context, in *InfectRequest, opts ...grpc.CallOption) (*Node, error)
	AddRoute(ctx context.Context, in *AddRouteRequest, opts ...grpc.CallOption) (*AddRouteReply, error)
	SendPacket(ctx context.Context, in *SendPacketRequest, opts ...grpc.CallOption) (*SendPacketReply, error)
}

type malwareSimulatorClient struct {
	cc grpc.ClientConnInterface
}

func NewMalwareSimulatorClient(cc grpc.ClientConnInterface) MalwareSimulatorClient {
	return &malwareSimulatorClient{cc}
}

func (c *malwareSimulatorClient) AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/AddNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) UpdateNodeInfo(ctx context.Context, in *UpdateNodeInfoRequest, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/UpdateNodeInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) AddNetwork(ctx context.Context, in *AddNetworkRequest, opts ...grpc.CallOption) (*Network, error) {
	out := new(Network)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/AddNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) AddApplication(ctx context.Context, in *AddApplicationRequest, opts ...grpc.CallOption) (*Application, error) {
	out := new(Application)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/AddApplication", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveReply, error) {
	out := new(RemoveReply)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) MakeConnection(ctx context.Context, in *MakeConnectionRequest, opts ...grpc.CallOption) (*Network, error) {
	out := new(Network)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/MakeConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) Infect(ctx context.Context, in *InfectRequest, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/Infect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) AddRoute(ctx context.Context, in *AddRouteRequest, opts ...grpc.CallOption) (*AddRouteReply, error) {
	out := new(AddRouteReply)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/AddRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *malwareSimulatorClient) SendPacket(ctx context.Context, in *SendPacketRequest, opts ...grpc.CallOption) (*SendPacketReply, error) {
	out := new(SendPacketReply)
	err := c.cc.Invoke(ctx, "/malwaresimulator.MalwareSimulator/SendPacket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MalwareSimulatorServer is the server API for MalwareSimulator service.
// All implementations must embed UnimplementedMalwareSimulatorServer
// for forward compatibility
type MalwareSimulatorServer interface {
	AddNode(context.Context, *AddNodeRequest) (*Node, error)
	UpdateNodeInfo(context.Context, *UpdateNodeInfoRequest) (*Node, error)
	AddNetwork(context.Context, *AddNetworkRequest) (*Network, error)
	AddApplication(context.Context, *AddApplicationRequest) (*Application, error)
	Remove(context.Context, *RemoveRequest) (*RemoveReply, error)
	MakeConnection(context.Context, *MakeConnectionRequest) (*Network, error)
	Infect(context.Context, *InfectRequest) (*Node, error)
	AddRoute(context.Context, *AddRouteRequest) (*AddRouteReply, error)
	SendPacket(context.Context, *SendPacketRequest) (*SendPacketReply, error)
	mustEmbedUnimplementedMalwareSimulatorServer()
}

// UnimplementedMalwareSimulatorServer must be embedded to have forward compatible implementations.
type UnimplementedMalwareSimulatorServer struct {
}

func (UnimplementedMalwareSimulatorServer) AddNode(context.Context, *AddNodeRequest) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNode not implemented")
}
func (UnimplementedMalwareSimulatorServer) UpdateNodeInfo(context.Context, *UpdateNodeInfoRequest) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNodeInfo not implemented")
}
func (UnimplementedMalwareSimulatorServer) AddNetwork(context.Context, *AddNetworkRequest) (*Network, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNetwork not implemented")
}
func (UnimplementedMalwareSimulatorServer) AddApplication(context.Context, *AddApplicationRequest) (*Application, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddApplication not implemented")
}
func (UnimplementedMalwareSimulatorServer) Remove(context.Context, *RemoveRequest) (*RemoveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedMalwareSimulatorServer) MakeConnection(context.Context, *MakeConnectionRequest) (*Network, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeConnection not implemented")
}
func (UnimplementedMalwareSimulatorServer) Infect(context.Context, *InfectRequest) (*Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Infect not implemented")
}
func (UnimplementedMalwareSimulatorServer) AddRoute(context.Context, *AddRouteRequest) (*AddRouteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRoute not implemented")
}
func (UnimplementedMalwareSimulatorServer) SendPacket(context.Context, *SendPacketRequest) (*SendPacketReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPacket not implemented")
}
func (UnimplementedMalwareSimulatorServer) mustEmbedUnimplementedMalwareSimulatorServer() {}

// UnsafeMalwareSimulatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MalwareSimulatorServer will
// result in compilation errors.
type UnsafeMalwareSimulatorServer interface {
	mustEmbedUnimplementedMalwareSimulatorServer()
}

func RegisterMalwareSimulatorServer(s grpc.ServiceRegistrar, srv MalwareSimulatorServer) {
	s.RegisterService(&MalwareSimulator_ServiceDesc, srv)
}

func _MalwareSimulator_AddNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).AddNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/AddNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).AddNode(ctx, req.(*AddNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_UpdateNodeInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNodeInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).UpdateNodeInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/UpdateNodeInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).UpdateNodeInfo(ctx, req.(*UpdateNodeInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_AddNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).AddNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/AddNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).AddNetwork(ctx, req.(*AddNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_AddApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).AddApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/AddApplication",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).AddApplication(ctx, req.(*AddApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_MakeConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MakeConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).MakeConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/MakeConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).MakeConnection(ctx, req.(*MakeConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_Infect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).Infect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/Infect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).Infect(ctx, req.(*InfectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_AddRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).AddRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/AddRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).AddRoute(ctx, req.(*AddRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MalwareSimulator_SendPacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPacketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MalwareSimulatorServer).SendPacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/malwaresimulator.MalwareSimulator/SendPacket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MalwareSimulatorServer).SendPacket(ctx, req.(*SendPacketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MalwareSimulator_ServiceDesc is the grpc.ServiceDesc for MalwareSimulator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MalwareSimulator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "malwaresimulator.MalwareSimulator",
	HandlerType: (*MalwareSimulatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddNode",
			Handler:    _MalwareSimulator_AddNode_Handler,
		},
		{
			MethodName: "UpdateNodeInfo",
			Handler:    _MalwareSimulator_UpdateNodeInfo_Handler,
		},
		{
			MethodName: "AddNetwork",
			Handler:    _MalwareSimulator_AddNetwork_Handler,
		},
		{
			MethodName: "AddApplication",
			Handler:    _MalwareSimulator_AddApplication_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _MalwareSimulator_Remove_Handler,
		},
		{
			MethodName: "MakeConnection",
			Handler:    _MalwareSimulator_MakeConnection_Handler,
		},
		{
			MethodName: "Infect",
			Handler:    _MalwareSimulator_Infect_Handler,
		},
		{
			MethodName: "AddRoute",
			Handler:    _MalwareSimulator_AddRoute_Handler,
		},
		{
			MethodName: "SendPacket",
			Handler:    _MalwareSimulator_SendPacket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/AttacksTracer.proto",
}
