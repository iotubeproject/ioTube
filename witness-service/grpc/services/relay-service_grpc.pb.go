// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package services

import (
	context "context"
	types "github.com/iotexproject/ioTube/witness-service/grpc/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	// rpcs below are BTC validation related
	ListUnsignedBTCTXWithoutNonces(ctx context.Context, in *ExcludedTransactions, opts ...grpc.CallOption) (*ListUnsignedBTCTXWithoutNoncesResponse, error)
	SubmitMusigNonces(ctx context.Context, in *MusigNonceMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListUnsignedBTCTXWithNonces(ctx context.Context, in *ExcludedTransactions, opts ...grpc.CallOption) (*ListUnsignedBTCTXWithNoncesResponse, error)
	SubmitMusigSignatures(ctx context.Context, in *MusigSignatureMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
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

func (c *relayServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) ListUnsignedBTCTXWithoutNonces(ctx context.Context, in *ExcludedTransactions, opts ...grpc.CallOption) (*ListUnsignedBTCTXWithoutNoncesResponse, error) {
	out := new(ListUnsignedBTCTXWithoutNoncesResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/ListUnsignedBTCTXWithoutNonces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) SubmitMusigNonces(ctx context.Context, in *MusigNonceMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/services.RelayService/SubmitMusigNonces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) ListUnsignedBTCTXWithNonces(ctx context.Context, in *ExcludedTransactions, opts ...grpc.CallOption) (*ListUnsignedBTCTXWithNoncesResponse, error) {
	out := new(ListUnsignedBTCTXWithNoncesResponse)
	err := c.cc.Invoke(ctx, "/services.RelayService/ListUnsignedBTCTXWithNonces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) SubmitMusigSignatures(ctx context.Context, in *MusigSignatureMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/services.RelayService/SubmitMusigSignatures", in, out, opts...)
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
	List(context.Context, *ListRequest) (*ListResponse, error)
	// rpcs below are BTC validation related
	ListUnsignedBTCTXWithoutNonces(context.Context, *ExcludedTransactions) (*ListUnsignedBTCTXWithoutNoncesResponse, error)
	SubmitMusigNonces(context.Context, *MusigNonceMessage) (*emptypb.Empty, error)
	ListUnsignedBTCTXWithNonces(context.Context, *ExcludedTransactions) (*ListUnsignedBTCTXWithNoncesResponse, error)
	SubmitMusigSignatures(context.Context, *MusigSignatureMessage) (*emptypb.Empty, error)
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
func (UnimplementedRelayServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedRelayServiceServer) ListUnsignedBTCTXWithoutNonces(context.Context, *ExcludedTransactions) (*ListUnsignedBTCTXWithoutNoncesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUnsignedBTCTXWithoutNonces not implemented")
}
func (UnimplementedRelayServiceServer) SubmitMusigNonces(context.Context, *MusigNonceMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitMusigNonces not implemented")
}
func (UnimplementedRelayServiceServer) ListUnsignedBTCTXWithNonces(context.Context, *ExcludedTransactions) (*ListUnsignedBTCTXWithNoncesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUnsignedBTCTXWithNonces not implemented")
}
func (UnimplementedRelayServiceServer) SubmitMusigSignatures(context.Context, *MusigSignatureMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitMusigSignatures not implemented")
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

func _RelayService_ListUnsignedBTCTXWithoutNonces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExcludedTransactions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).ListUnsignedBTCTXWithoutNonces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/ListUnsignedBTCTXWithoutNonces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).ListUnsignedBTCTXWithoutNonces(ctx, req.(*ExcludedTransactions))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_SubmitMusigNonces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MusigNonceMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).SubmitMusigNonces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/SubmitMusigNonces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).SubmitMusigNonces(ctx, req.(*MusigNonceMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_ListUnsignedBTCTXWithNonces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExcludedTransactions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).ListUnsignedBTCTXWithNonces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/ListUnsignedBTCTXWithNonces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).ListUnsignedBTCTXWithNonces(ctx, req.(*ExcludedTransactions))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_SubmitMusigSignatures_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MusigSignatureMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).SubmitMusigSignatures(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.RelayService/SubmitMusigSignatures",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).SubmitMusigSignatures(ctx, req.(*MusigSignatureMessage))
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
			MethodName: "List",
			Handler:    _RelayService_List_Handler,
		},
		{
			MethodName: "ListUnsignedBTCTXWithoutNonces",
			Handler:    _RelayService_ListUnsignedBTCTXWithoutNonces_Handler,
		},
		{
			MethodName: "SubmitMusigNonces",
			Handler:    _RelayService_SubmitMusigNonces_Handler,
		},
		{
			MethodName: "ListUnsignedBTCTXWithNonces",
			Handler:    _RelayService_ListUnsignedBTCTXWithNonces_Handler,
		},
		{
			MethodName: "SubmitMusigSignatures",
			Handler:    _RelayService_SubmitMusigSignatures_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/relay-service.proto",
}
