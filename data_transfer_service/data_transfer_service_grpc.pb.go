// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package data_transfer_service

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

// DataTransferClient is the client API for DataTransfer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataTransferClient interface {
	GetDataStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (DataTransfer_GetDataStreamClient, error)
}

type dataTransferClient struct {
	cc grpc.ClientConnInterface
}

func NewDataTransferClient(cc grpc.ClientConnInterface) DataTransferClient {
	return &dataTransferClient{cc}
}

func (c *dataTransferClient) GetDataStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (DataTransfer_GetDataStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &DataTransfer_ServiceDesc.Streams[0], "/DataTransfer/GetDataStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataTransferGetDataStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DataTransfer_GetDataStreamClient interface {
	Recv() (*Data, error)
	grpc.ClientStream
}

type dataTransferGetDataStreamClient struct {
	grpc.ClientStream
}

func (x *dataTransferGetDataStreamClient) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DataTransferServer is the server API for DataTransfer service.
// All implementations must embed UnimplementedDataTransferServer
// for forward compatibility
type DataTransferServer interface {
	GetDataStream(*Request, DataTransfer_GetDataStreamServer) error
	mustEmbedUnimplementedDataTransferServer()
}

// UnimplementedDataTransferServer must be embedded to have forward compatible implementations.
type UnimplementedDataTransferServer struct {
}

func (UnimplementedDataTransferServer) GetDataStream(*Request, DataTransfer_GetDataStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetDataStream not implemented")
}
func (UnimplementedDataTransferServer) mustEmbedUnimplementedDataTransferServer() {}

// UnsafeDataTransferServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataTransferServer will
// result in compilation errors.
type UnsafeDataTransferServer interface {
	mustEmbedUnimplementedDataTransferServer()
}

func RegisterDataTransferServer(s grpc.ServiceRegistrar, srv DataTransferServer) {
	s.RegisterService(&DataTransfer_ServiceDesc, srv)
}

func _DataTransfer_GetDataStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataTransferServer).GetDataStream(m, &dataTransferGetDataStreamServer{stream})
}

type DataTransfer_GetDataStreamServer interface {
	Send(*Data) error
	grpc.ServerStream
}

type dataTransferGetDataStreamServer struct {
	grpc.ServerStream
}

func (x *dataTransferGetDataStreamServer) Send(m *Data) error {
	return x.ServerStream.SendMsg(m)
}

// DataTransfer_ServiceDesc is the grpc.ServiceDesc for DataTransfer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataTransfer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "DataTransfer",
	HandlerType: (*DataTransferServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDataStream",
			Handler:       _DataTransfer_GetDataStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protobuf/data_transfer_service.proto",
}
