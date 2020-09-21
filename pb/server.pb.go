// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("server.proto", fileDescriptor_ad098daeda4239f7) }

var fileDescriptor_ad098daeda4239f7 = []byte{
	// 146 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x92, 0x92, 0x4e, 0xcf,
	0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x07, 0x8b, 0x24, 0x95, 0xa6, 0xe9, 0xa7, 0xe6, 0x16, 0x94, 0x54,
	0x42, 0x14, 0x18, 0x79, 0x73, 0xf1, 0x85, 0xa7, 0x16, 0xa5, 0x96, 0xe7, 0xe7, 0xa4, 0x05, 0x83,
	0x35, 0x0a, 0x59, 0x72, 0xb1, 0x06, 0x97, 0x24, 0x16, 0x95, 0x08, 0x89, 0xe9, 0x41, 0x34, 0xea,
	0xc1, 0x34, 0xea, 0xb9, 0x82, 0x34, 0x4a, 0xe1, 0x10, 0x57, 0x62, 0x70, 0xd2, 0x88, 0x52, 0x4b,
	0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xf7, 0x48, 0xcd, 0x4d, 0x2d, 0xaa,
	0xf4, 0x2a, 0xd5, 0x0f, 0x2c, 0x4d, 0xcc, 0x2b, 0x29, 0xcd, 0x85, 0xd9, 0xa3, 0x5f, 0x90, 0x94,
	0xc4, 0x06, 0xd6, 0x6b, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xd1, 0xc8, 0x65, 0x58, 0xae, 0x00,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WerewolfServerClient is the client API for WerewolfServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WerewolfServerClient interface {
	Start(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
}

type werewolfServerClient struct {
	cc *grpc.ClientConn
}

func NewWerewolfServerClient(cc *grpc.ClientConn) WerewolfServerClient {
	return &werewolfServerClient{cc}
}

func (c *werewolfServerClient) Start(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/pb.WerewolfServer/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WerewolfServerServer is the server API for WerewolfServer service.
type WerewolfServerServer interface {
	Start(context.Context, *empty.Empty) (*empty.Empty, error)
}

func RegisterWerewolfServerServer(s *grpc.Server, srv WerewolfServerServer) {
	s.RegisterService(&_WerewolfServer_serviceDesc, srv)
}

func _WerewolfServer_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WerewolfServerServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.WerewolfServer/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WerewolfServerServer).Start(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _WerewolfServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.WerewolfServer",
	HandlerType: (*WerewolfServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _WerewolfServer_Start_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
