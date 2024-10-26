// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.0
// source: services/relay-service.proto

package services

import (
	context "context"
	types "github.com/iotexproject/ioTube/witness-service/grpc/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RelayServiceClient is the client API for RelayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelayServiceClient interface {
	Submit(ctx context.Context, in *types.Witness, opts ...grpc.CallOption) (*WitnessSubmissionResponse, error)
	Reset(ctx context.Context, in *ResetTransferRequest, opts ...grpc.CallOption) (*ResetTransferResponse, error)
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	StaleHeights(ctx context.Context, in *StaleHeightsRequest, opts ...grpc.CallOption) (*StaleHeightsResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	SubmitNewTX(ctx context.Context, in *SubmitNewTXRequest, opts ...grpc.CallOption) (*SubmitNewTXResponse, error)
	ListNewTX(ctx context.Context, in *ListNewTXRequest, opts ...grpc.CallOption) (*ListNewTXResponse, error)
}

type relayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayServiceClient(cc grpc.ClientConnInterface) RelayServiceClient {
	return &relayServiceClient{cc}
}

func (c *relayServiceClient) Submit(ctx context.Context, in *types.Witness, opts ...grpc.CallOption) (*WitnessSubmissionResponse, error) {
	out := new(WitnessSubmissionResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/Submit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) Reset(ctx context.Context, in *ResetTransferRequest, opts ...grpc.CallOption) (*ResetTransferResponse, error) {
	out := new(ResetTransferResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/Reset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error) {
	out := new(LookupResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/Lookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) StaleHeights(ctx context.Context, in *StaleHeightsRequest, opts ...grpc.CallOption) (*StaleHeightsResponse, error) {
	out := new(StaleHeightsResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/StaleHeights", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) SubmitNewTX(ctx context.Context, in *SubmitNewTXRequest, opts ...grpc.CallOption) (*SubmitNewTXResponse, error) {
	out := new(SubmitNewTXResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/SubmitNewTX", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) ListNewTX(ctx context.Context, in *ListNewTXRequest, opts ...grpc.CallOption) (*ListNewTXResponse, error) {
	out := new(ListNewTXResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/ListNewTX", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelayServiceServer is the server API for RelayService service.
// All implementations must embed UnimplementedRelayServiceServer
// for forward compatibility
type RelayServiceServer interface {
	Submit(context.Context, *types.Witness) (*WitnessSubmissionResponse, error)
	Reset(context.Context, *ResetTransferRequest) (*ResetTransferResponse, error)
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	Lookup(context.Context, *LookupRequest) (*LookupResponse, error)
	StaleHeights(context.Context, *StaleHeightsRequest) (*StaleHeightsResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	SubmitNewTX(context.Context, *SubmitNewTXRequest) (*SubmitNewTXResponse, error)
	ListNewTX(context.Context, *ListNewTXRequest) (*ListNewTXResponse, error)
	mustEmbedUnimplementedRelayServiceServer()
}

// UnimplementedRelayServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRelayServiceServer struct {
}

func (UnimplementedRelayServiceServer) Submit(context.Context, *types.Witness) (*WitnessSubmissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Submit not implemented")
}
func (UnimplementedRelayServiceServer) Reset(context.Context, *ResetTransferRequest) (*ResetTransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}
func (UnimplementedRelayServiceServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedRelayServiceServer) Lookup(context.Context, *LookupRequest) (*LookupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}
func (UnimplementedRelayServiceServer) StaleHeights(context.Context, *StaleHeightsRequest) (*StaleHeightsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StaleHeights not implemented")
}
func (UnimplementedRelayServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedRelayServiceServer) SubmitNewTX(context.Context, *SubmitNewTXRequest) (*SubmitNewTXResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitNewTX not implemented")
}
func (UnimplementedRelayServiceServer) ListNewTX(context.Context, *ListNewTXRequest) (*ListNewTXResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNewTX not implemented")
}
func (UnimplementedRelayServiceServer) mustEmbedUnimplementedRelayServiceServer() {}

// UnsafeRelayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelayServiceServer will
// result in compilation errors.
type UnsafeRelayServiceServer interface {
	mustEmbedUnimplementedRelayServiceServer()
}

func RegisterRelayServiceServer(s grpc.ServiceRegistrar, srv RelayServiceServer) {
	s.RegisterService(&RelayService_ServiceDesc, srv)
}

func _RelayService_Submit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.Witness)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).Submit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/Submit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).Submit(ctx, req.(*types.Witness))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_Reset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetTransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).Reset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/Reset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).Reset(ctx, req.(*ResetTransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).Lookup(ctx, req.(*LookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_StaleHeights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StaleHeightsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).StaleHeights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/StaleHeights",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).StaleHeights(ctx, req.(*StaleHeightsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_SubmitNewTX_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitNewTXRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).SubmitNewTX(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/SubmitNewTX",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).SubmitNewTX(ctx, req.(*SubmitNewTXRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_ListNewTX_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNewTXRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).ListNewTX(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/ListNewTX",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).ListNewTX(ctx, req.(*ListNewTXRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RelayService_ServiceDesc is the grpc.ServiceDesc for RelayService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.RelayService",
	HandlerType: (*RelayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Submit",
			Handler:    _RelayService_Submit_Handler,
		},
		{
			MethodName: "Reset",
			Handler:    _RelayService_Reset_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _RelayService_Check_Handler,
		},
		{
			MethodName: "Lookup",
			Handler:    _RelayService_Lookup_Handler,
		},
		{
			MethodName: "StaleHeights",
			Handler:    _RelayService_StaleHeights_Handler,
		},
		{
			MethodName: "List",
			Handler:    _RelayService_List_Handler,
		},
		{
			MethodName: "SubmitNewTX",
			Handler:    _RelayService_SubmitNewTX_Handler,
		},
		{
			MethodName: "ListNewTX",
			Handler:    _RelayService_ListNewTX_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/relay-service.proto",
}
