// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: github.com/henry118/containerd/api/services/networks/v1/networks.proto

package networks

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NetworkClient is the client API for Network service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkClient interface {
	CreateNetwork(ctx context.Context, in *CreateNetworkRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteNetwork(ctx context.Context, in *DeleteNetworkRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetNetwork(ctx context.Context, in *GetNetworkRequest, opts ...grpc.CallOption) (*GetNetworkResponse, error)
	ListNetworks(ctx context.Context, in *ListNetworksRequest, opts ...grpc.CallOption) (*ListNetworksResponse, error)
	AttachNetwork(ctx context.Context, in *AttachNetworkRequest, opts ...grpc.CallOption) (*AttachNetworkResponse, error)
	DetachNetwork(ctx context.Context, in *DetachNetworkRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAttachment(ctx context.Context, in *GetAttachmentRequest, opts ...grpc.CallOption) (*GetAttachmentResponse, error)
	CheckAttachment(ctx context.Context, in *CheckAttachmentRequest, opts ...grpc.CallOption) (*CheckAttachmentResponse, error)
	ListAttachments(ctx context.Context, in *ListAttachmentsRequest, opts ...grpc.CallOption) (*ListAttachmentsResponse, error)
}

type networkClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkClient(cc grpc.ClientConnInterface) NetworkClient {
	return &networkClient{cc}
}

func (c *networkClient) CreateNetwork(ctx context.Context, in *CreateNetworkRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/CreateNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) DeleteNetwork(ctx context.Context, in *DeleteNetworkRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/DeleteNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) GetNetwork(ctx context.Context, in *GetNetworkRequest, opts ...grpc.CallOption) (*GetNetworkResponse, error) {
	out := new(GetNetworkResponse)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/GetNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) ListNetworks(ctx context.Context, in *ListNetworksRequest, opts ...grpc.CallOption) (*ListNetworksResponse, error) {
	out := new(ListNetworksResponse)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/ListNetworks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) AttachNetwork(ctx context.Context, in *AttachNetworkRequest, opts ...grpc.CallOption) (*AttachNetworkResponse, error) {
	out := new(AttachNetworkResponse)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/AttachNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) DetachNetwork(ctx context.Context, in *DetachNetworkRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/DetachNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) GetAttachment(ctx context.Context, in *GetAttachmentRequest, opts ...grpc.CallOption) (*GetAttachmentResponse, error) {
	out := new(GetAttachmentResponse)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/GetAttachment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) CheckAttachment(ctx context.Context, in *CheckAttachmentRequest, opts ...grpc.CallOption) (*CheckAttachmentResponse, error) {
	out := new(CheckAttachmentResponse)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/CheckAttachment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkClient) ListAttachments(ctx context.Context, in *ListAttachmentsRequest, opts ...grpc.CallOption) (*ListAttachmentsResponse, error) {
	out := new(ListAttachmentsResponse)
	err := c.cc.Invoke(ctx, "/containerd.services.networks.v1.Network/ListAttachments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServer is the server API for Network service.
// All implementations must embed UnimplementedNetworkServer
// for forward compatibility
type NetworkServer interface {
	CreateNetwork(context.Context, *CreateNetworkRequest) (*emptypb.Empty, error)
	DeleteNetwork(context.Context, *DeleteNetworkRequest) (*emptypb.Empty, error)
	GetNetwork(context.Context, *GetNetworkRequest) (*GetNetworkResponse, error)
	ListNetworks(context.Context, *ListNetworksRequest) (*ListNetworksResponse, error)
	AttachNetwork(context.Context, *AttachNetworkRequest) (*AttachNetworkResponse, error)
	DetachNetwork(context.Context, *DetachNetworkRequest) (*emptypb.Empty, error)
	GetAttachment(context.Context, *GetAttachmentRequest) (*GetAttachmentResponse, error)
	CheckAttachment(context.Context, *CheckAttachmentRequest) (*CheckAttachmentResponse, error)
	ListAttachments(context.Context, *ListAttachmentsRequest) (*ListAttachmentsResponse, error)
	mustEmbedUnimplementedNetworkServer()
}

// UnimplementedNetworkServer must be embedded to have forward compatible implementations.
type UnimplementedNetworkServer struct {
}

func (UnimplementedNetworkServer) CreateNetwork(context.Context, *CreateNetworkRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNetwork not implemented")
}
func (UnimplementedNetworkServer) DeleteNetwork(context.Context, *DeleteNetworkRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNetwork not implemented")
}
func (UnimplementedNetworkServer) GetNetwork(context.Context, *GetNetworkRequest) (*GetNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNetwork not implemented")
}
func (UnimplementedNetworkServer) ListNetworks(context.Context, *ListNetworksRequest) (*ListNetworksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNetworks not implemented")
}
func (UnimplementedNetworkServer) AttachNetwork(context.Context, *AttachNetworkRequest) (*AttachNetworkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttachNetwork not implemented")
}
func (UnimplementedNetworkServer) DetachNetwork(context.Context, *DetachNetworkRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetachNetwork not implemented")
}
func (UnimplementedNetworkServer) GetAttachment(context.Context, *GetAttachmentRequest) (*GetAttachmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAttachment not implemented")
}
func (UnimplementedNetworkServer) CheckAttachment(context.Context, *CheckAttachmentRequest) (*CheckAttachmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAttachment not implemented")
}
func (UnimplementedNetworkServer) ListAttachments(context.Context, *ListAttachmentsRequest) (*ListAttachmentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAttachments not implemented")
}
func (UnimplementedNetworkServer) mustEmbedUnimplementedNetworkServer() {}

// UnsafeNetworkServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkServer will
// result in compilation errors.
type UnsafeNetworkServer interface {
	mustEmbedUnimplementedNetworkServer()
}

func RegisterNetworkServer(s grpc.ServiceRegistrar, srv NetworkServer) {
	s.RegisterService(&Network_ServiceDesc, srv)
}

func _Network_CreateNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).CreateNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/CreateNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).CreateNetwork(ctx, req.(*CreateNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_DeleteNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).DeleteNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/DeleteNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).DeleteNetwork(ctx, req.(*DeleteNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_GetNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).GetNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/GetNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).GetNetwork(ctx, req.(*GetNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_ListNetworks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNetworksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).ListNetworks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/ListNetworks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).ListNetworks(ctx, req.(*ListNetworksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_AttachNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttachNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).AttachNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/AttachNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).AttachNetwork(ctx, req.(*AttachNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_DetachNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetachNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).DetachNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/DetachNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).DetachNetwork(ctx, req.(*DetachNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_GetAttachment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).GetAttachment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/GetAttachment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).GetAttachment(ctx, req.(*GetAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_CheckAttachment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).CheckAttachment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/CheckAttachment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).CheckAttachment(ctx, req.(*CheckAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Network_ListAttachments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAttachmentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServer).ListAttachments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/containerd.services.networks.v1.Network/ListAttachments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServer).ListAttachments(ctx, req.(*ListAttachmentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Network_ServiceDesc is the grpc.ServiceDesc for Network service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Network_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "containerd.services.networks.v1.Network",
	HandlerType: (*NetworkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNetwork",
			Handler:    _Network_CreateNetwork_Handler,
		},
		{
			MethodName: "DeleteNetwork",
			Handler:    _Network_DeleteNetwork_Handler,
		},
		{
			MethodName: "GetNetwork",
			Handler:    _Network_GetNetwork_Handler,
		},
		{
			MethodName: "ListNetworks",
			Handler:    _Network_ListNetworks_Handler,
		},
		{
			MethodName: "AttachNetwork",
			Handler:    _Network_AttachNetwork_Handler,
		},
		{
			MethodName: "DetachNetwork",
			Handler:    _Network_DetachNetwork_Handler,
		},
		{
			MethodName: "GetAttachment",
			Handler:    _Network_GetAttachment_Handler,
		},
		{
			MethodName: "CheckAttachment",
			Handler:    _Network_CheckAttachment_Handler,
		},
		{
			MethodName: "ListAttachments",
			Handler:    _Network_ListAttachments_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/henry118/containerd/api/services/networks/v1/networks.proto",
}
