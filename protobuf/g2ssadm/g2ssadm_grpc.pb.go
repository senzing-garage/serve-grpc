// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: g2ssadm.proto

package g2ssadm

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

// G2SsadmClient is the client API for G2Ssadm service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type G2SsadmClient interface {
	ClearLastException(ctx context.Context, in *ClearLastExceptionRequest, opts ...grpc.CallOption) (*ClearLastExceptionResponse, error)
	CreateSaltInStore(ctx context.Context, in *CreateSaltInStoreRequest, opts ...grpc.CallOption) (*CreateSaltInStoreResponse, error)
	Destroy(ctx context.Context, in *DestroyRequest, opts ...grpc.CallOption) (*DestroyResponse, error)
	GetLastException(ctx context.Context, in *GetLastExceptionRequest, opts ...grpc.CallOption) (*GetLastExceptionResponse, error)
	GetLastExceptionCode(ctx context.Context, in *GetLastExceptionCodeRequest, opts ...grpc.CallOption) (*GetLastExceptionCodeResponse, error)
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error)
	InitializeNewToken(ctx context.Context, in *InitializeNewTokenRequest, opts ...grpc.CallOption) (*InitializeNewTokenResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	ReinitializeToken(ctx context.Context, in *ReinitializeTokenRequest, opts ...grpc.CallOption) (*ReinitializeTokenResponse, error)
	SetupStore(ctx context.Context, in *SetupStoreRequest, opts ...grpc.CallOption) (*SetupStoreResponse, error)
}

type g2SsadmClient struct {
	cc grpc.ClientConnInterface
}

func NewG2SsadmClient(cc grpc.ClientConnInterface) G2SsadmClient {
	return &g2SsadmClient{cc}
}

func (c *g2SsadmClient) ClearLastException(ctx context.Context, in *ClearLastExceptionRequest, opts ...grpc.CallOption) (*ClearLastExceptionResponse, error) {
	out := new(ClearLastExceptionResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/ClearLastException", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) CreateSaltInStore(ctx context.Context, in *CreateSaltInStoreRequest, opts ...grpc.CallOption) (*CreateSaltInStoreResponse, error) {
	out := new(CreateSaltInStoreResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/CreateSaltInStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) Destroy(ctx context.Context, in *DestroyRequest, opts ...grpc.CallOption) (*DestroyResponse, error) {
	out := new(DestroyResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/Destroy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) GetLastException(ctx context.Context, in *GetLastExceptionRequest, opts ...grpc.CallOption) (*GetLastExceptionResponse, error) {
	out := new(GetLastExceptionResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/GetLastException", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) GetLastExceptionCode(ctx context.Context, in *GetLastExceptionCodeRequest, opts ...grpc.CallOption) (*GetLastExceptionCodeResponse, error) {
	out := new(GetLastExceptionCodeResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/GetLastExceptionCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	out := new(InitResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) InitializeNewToken(ctx context.Context, in *InitializeNewTokenRequest, opts ...grpc.CallOption) (*InitializeNewTokenResponse, error) {
	out := new(InitializeNewTokenResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/InitializeNewToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) ReinitializeToken(ctx context.Context, in *ReinitializeTokenRequest, opts ...grpc.CallOption) (*ReinitializeTokenResponse, error) {
	out := new(ReinitializeTokenResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/ReinitializeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *g2SsadmClient) SetupStore(ctx context.Context, in *SetupStoreRequest, opts ...grpc.CallOption) (*SetupStoreResponse, error) {
	out := new(SetupStoreResponse)
	err := c.cc.Invoke(ctx, "/g2ssadm.G2Ssadm/SetupStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// G2SsadmServer is the server API for G2Ssadm service.
// All implementations must embed UnimplementedG2SsadmServer
// for forward compatibility
type G2SsadmServer interface {
	ClearLastException(context.Context, *ClearLastExceptionRequest) (*ClearLastExceptionResponse, error)
	CreateSaltInStore(context.Context, *CreateSaltInStoreRequest) (*CreateSaltInStoreResponse, error)
	Destroy(context.Context, *DestroyRequest) (*DestroyResponse, error)
	GetLastException(context.Context, *GetLastExceptionRequest) (*GetLastExceptionResponse, error)
	GetLastExceptionCode(context.Context, *GetLastExceptionCodeRequest) (*GetLastExceptionCodeResponse, error)
	Init(context.Context, *InitRequest) (*InitResponse, error)
	InitializeNewToken(context.Context, *InitializeNewTokenRequest) (*InitializeNewTokenResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Put(context.Context, *PutRequest) (*PutResponse, error)
	ReinitializeToken(context.Context, *ReinitializeTokenRequest) (*ReinitializeTokenResponse, error)
	SetupStore(context.Context, *SetupStoreRequest) (*SetupStoreResponse, error)
	mustEmbedUnimplementedG2SsadmServer()
}

// UnimplementedG2SsadmServer must be embedded to have forward compatible implementations.
type UnimplementedG2SsadmServer struct {
}

func (UnimplementedG2SsadmServer) ClearLastException(context.Context, *ClearLastExceptionRequest) (*ClearLastExceptionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearLastException not implemented")
}
func (UnimplementedG2SsadmServer) CreateSaltInStore(context.Context, *CreateSaltInStoreRequest) (*CreateSaltInStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSaltInStore not implemented")
}
func (UnimplementedG2SsadmServer) Destroy(context.Context, *DestroyRequest) (*DestroyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Destroy not implemented")
}
func (UnimplementedG2SsadmServer) GetLastException(context.Context, *GetLastExceptionRequest) (*GetLastExceptionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastException not implemented")
}
func (UnimplementedG2SsadmServer) GetLastExceptionCode(context.Context, *GetLastExceptionCodeRequest) (*GetLastExceptionCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastExceptionCode not implemented")
}
func (UnimplementedG2SsadmServer) Init(context.Context, *InitRequest) (*InitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedG2SsadmServer) InitializeNewToken(context.Context, *InitializeNewTokenRequest) (*InitializeNewTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitializeNewToken not implemented")
}
func (UnimplementedG2SsadmServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedG2SsadmServer) Put(context.Context, *PutRequest) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedG2SsadmServer) ReinitializeToken(context.Context, *ReinitializeTokenRequest) (*ReinitializeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReinitializeToken not implemented")
}
func (UnimplementedG2SsadmServer) SetupStore(context.Context, *SetupStoreRequest) (*SetupStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetupStore not implemented")
}
func (UnimplementedG2SsadmServer) mustEmbedUnimplementedG2SsadmServer() {}

// UnsafeG2SsadmServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to G2SsadmServer will
// result in compilation errors.
type UnsafeG2SsadmServer interface {
	mustEmbedUnimplementedG2SsadmServer()
}

func RegisterG2SsadmServer(s grpc.ServiceRegistrar, srv G2SsadmServer) {
	s.RegisterService(&G2Ssadm_ServiceDesc, srv)
}

func _G2Ssadm_ClearLastException_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearLastExceptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).ClearLastException(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/ClearLastException",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).ClearLastException(ctx, req.(*ClearLastExceptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_CreateSaltInStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSaltInStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).CreateSaltInStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/CreateSaltInStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).CreateSaltInStore(ctx, req.(*CreateSaltInStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_Destroy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DestroyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).Destroy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/Destroy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).Destroy(ctx, req.(*DestroyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_GetLastException_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLastExceptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).GetLastException(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/GetLastException",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).GetLastException(ctx, req.(*GetLastExceptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_GetLastExceptionCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLastExceptionCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).GetLastExceptionCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/GetLastExceptionCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).GetLastExceptionCode(ctx, req.(*GetLastExceptionCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_InitializeNewToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitializeNewTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).InitializeNewToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/InitializeNewToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).InitializeNewToken(ctx, req.(*InitializeNewTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_ReinitializeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReinitializeTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).ReinitializeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/ReinitializeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).ReinitializeToken(ctx, req.(*ReinitializeTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _G2Ssadm_SetupStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetupStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(G2SsadmServer).SetupStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/g2ssadm.G2Ssadm/SetupStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(G2SsadmServer).SetupStore(ctx, req.(*SetupStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// G2Ssadm_ServiceDesc is the grpc.ServiceDesc for G2Ssadm service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var G2Ssadm_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "g2ssadm.G2Ssadm",
	HandlerType: (*G2SsadmServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClearLastException",
			Handler:    _G2Ssadm_ClearLastException_Handler,
		},
		{
			MethodName: "CreateSaltInStore",
			Handler:    _G2Ssadm_CreateSaltInStore_Handler,
		},
		{
			MethodName: "Destroy",
			Handler:    _G2Ssadm_Destroy_Handler,
		},
		{
			MethodName: "GetLastException",
			Handler:    _G2Ssadm_GetLastException_Handler,
		},
		{
			MethodName: "GetLastExceptionCode",
			Handler:    _G2Ssadm_GetLastExceptionCode_Handler,
		},
		{
			MethodName: "Init",
			Handler:    _G2Ssadm_Init_Handler,
		},
		{
			MethodName: "InitializeNewToken",
			Handler:    _G2Ssadm_InitializeNewToken_Handler,
		},
		{
			MethodName: "List",
			Handler:    _G2Ssadm_List_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _G2Ssadm_Put_Handler,
		},
		{
			MethodName: "ReinitializeToken",
			Handler:    _G2Ssadm_ReinitializeToken_Handler,
		},
		{
			MethodName: "SetupStore",
			Handler:    _G2Ssadm_SetupStore_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "g2ssadm.proto",
}