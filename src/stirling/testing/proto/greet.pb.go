// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: src/stirling/testing/proto/greet.proto

package greetpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type HelloRequest struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Count int32  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *HelloRequest) Reset()      { *m = HelloRequest{} }
func (*HelloRequest) ProtoMessage() {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd3055eacbf9153c, []int{0}
}
func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(m, src)
}
func (m *HelloRequest) XXX_Size() int {
	return m.Size()
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *HelloRequest) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type HelloReply struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *HelloReply) Reset()      { *m = HelloReply{} }
func (*HelloReply) ProtoMessage() {}
func (*HelloReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd3055eacbf9153c, []int{1}
}
func (m *HelloReply) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HelloReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HelloReply.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HelloReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReply.Merge(m, src)
}
func (m *HelloReply) XXX_Size() int {
	return m.Size()
}
func (m *HelloReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReply.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReply proto.InternalMessageInfo

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloRequest)(nil), "pl.stirling.testing.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "pl.stirling.testing.HelloReply")
}

func init() {
	proto.RegisterFile("src/stirling/testing/proto/greet.proto", fileDescriptor_cd3055eacbf9153c)
}

var fileDescriptor_cd3055eacbf9153c = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xb1, 0x4a, 0x33, 0x41,
	0x14, 0x85, 0xe7, 0xfe, 0xfc, 0x31, 0xf1, 0xa2, 0x20, 0xa3, 0x45, 0xb0, 0xb8, 0xc6, 0x14, 0x21,
	0xd5, 0xae, 0xc4, 0x46, 0xb1, 0xd2, 0x46, 0x1b, 0x41, 0x36, 0xd8, 0xd8, 0x4d, 0xc2, 0xb8, 0x0c,
	0x4c, 0x76, 0xd7, 0x9d, 0x09, 0x98, 0xce, 0x47, 0xf0, 0x31, 0x6c, 0x6c, 0x7c, 0x0a, 0xcb, 0x94,
	0x29, 0xcd, 0xa4, 0xb1, 0xcc, 0x23, 0x48, 0xc6, 0x5d, 0xb1, 0x50, 0xab, 0xad, 0xe6, 0x1e, 0x38,
	0xe7, 0xcc, 0xc7, 0xe5, 0x62, 0xc7, 0xe4, 0xc3, 0xd0, 0x58, 0x95, 0x6b, 0x95, 0xc4, 0xa1, 0x95,
	0xc6, 0xae, 0xde, 0x2c, 0x4f, 0x6d, 0x1a, 0xc6, 0xb9, 0x94, 0x36, 0xf0, 0x33, 0xdf, 0xce, 0x74,
	0x50, 0xda, 0x82, 0xc2, 0xd6, 0x3e, 0xc2, 0x8d, 0x0b, 0xa9, 0x75, 0x1a, 0xc9, 0xbb, 0xb1, 0x34,
	0x96, 0x73, 0xfc, 0x9f, 0x88, 0x91, 0x6c, 0x42, 0x0b, 0xba, 0xeb, 0x91, 0x9f, 0xf9, 0x0e, 0xd6,
	0x86, 0xe9, 0x38, 0xb1, 0xcd, 0x7f, 0x2d, 0xe8, 0xd6, 0xa2, 0x4f, 0xd1, 0xee, 0x20, 0x16, 0xc9,
	0x4c, 0x4f, 0x78, 0x13, 0xeb, 0x23, 0x69, 0x8c, 0x88, 0xcb, 0x68, 0x29, 0x7b, 0x2f, 0x80, 0xf5,
	0xf3, 0x15, 0x86, 0xcc, 0xf9, 0x15, 0x36, 0xfa, 0x62, 0xe2, 0x63, 0x7c, 0x3f, 0xf8, 0x81, 0x27,
	0xf8, 0x0e, 0xb3, 0xbb, 0xf7, 0x97, 0x25, 0xd3, 0x93, 0x36, 0xe3, 0xd7, 0xb8, 0x59, 0x36, 0x9e,
	0xc6, 0x42, 0x25, 0xd5, 0xd4, 0xf6, 0x9e, 0x01, 0x1b, 0x05, 0x74, 0x8f, 0x5f, 0x62, 0x6d, 0xf5,
	0x87, 0xaa, 0x08, 0x39, 0x42, 0xf4, 0x75, 0x55, 0xf2, 0xde, 0xe2, 0x56, 0xdf, 0xe6, 0x52, 0x8c,
	0x54, 0x12, 0x97, 0xcb, 0x8e, 0xaa, 0x5e, 0xf6, 0x01, 0x9c, 0xa5, 0xd3, 0x39, 0xb1, 0xd9, 0x9c,
	0xd8, 0x72, 0x4e, 0xf0, 0xe0, 0x08, 0x9e, 0x1c, 0xc1, 0xab, 0x23, 0x98, 0x3a, 0x82, 0x37, 0x47,
	0xf0, 0xee, 0x88, 0x2d, 0x1d, 0xc1, 0xe3, 0x82, 0xd8, 0x74, 0x41, 0x6c, 0xb6, 0x20, 0x76, 0x73,
	0x9c, 0xa9, 0x7b, 0x25, 0xb5, 0x18, 0x98, 0x40, 0xa8, 0xf0, 0x4b, 0x84, 0xbf, 0x1f, 0xee, 0x89,
	0x3f, 0xdc, 0x6c, 0x30, 0x58, 0xf3, 0xf2, 0xf0, 0x23, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x43, 0xee,
	0x86, 0xe5, 0x02, 0x00, 0x00,
}

func (this *HelloRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HelloRequest)
	if !ok {
		that2, ok := that.(HelloRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Count != that1.Count {
		return false
	}
	return true
}
func (this *HelloReply) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HelloReply)
	if !ok {
		that2, ok := that.(HelloReply)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	return true
}
func (this *HelloRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&greetpb.HelloRequest{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	s = append(s, "Count: "+fmt.Sprintf("%#v", this.Count)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *HelloReply) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&greetpb.HelloReply{")
	s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringGreet(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	SayHelloAgain(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/pl.stirling.testing.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) SayHelloAgain(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/pl.stirling.testing.Greeter/SayHelloAgain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterServer is the server API for Greeter service.
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	SayHelloAgain(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pl.stirling.testing.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_SayHelloAgain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHelloAgain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pl.stirling.testing.Greeter/SayHelloAgain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHelloAgain(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pl.stirling.testing.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
		{
			MethodName: "SayHelloAgain",
			Handler:    _Greeter_SayHelloAgain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/stirling/testing/proto/greet.proto",
}

// Greeter2Client is the client API for Greeter2 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type Greeter2Client interface {
	SayHi(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	SayHiAgain(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeter2Client struct {
	cc *grpc.ClientConn
}

func NewGreeter2Client(cc *grpc.ClientConn) Greeter2Client {
	return &greeter2Client{cc}
}

func (c *greeter2Client) SayHi(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/pl.stirling.testing.Greeter2/SayHi", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeter2Client) SayHiAgain(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/pl.stirling.testing.Greeter2/SayHiAgain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Greeter2Server is the server API for Greeter2 service.
type Greeter2Server interface {
	SayHi(context.Context, *HelloRequest) (*HelloReply, error)
	SayHiAgain(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeter2Server(s *grpc.Server, srv Greeter2Server) {
	s.RegisterService(&_Greeter2_serviceDesc, srv)
}

func _Greeter2_SayHi_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Greeter2Server).SayHi(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pl.stirling.testing.Greeter2/SayHi",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Greeter2Server).SayHi(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter2_SayHiAgain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Greeter2Server).SayHiAgain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pl.stirling.testing.Greeter2/SayHiAgain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Greeter2Server).SayHiAgain(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter2_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pl.stirling.testing.Greeter2",
	HandlerType: (*Greeter2Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHi",
			Handler:    _Greeter2_SayHi_Handler,
		},
		{
			MethodName: "SayHiAgain",
			Handler:    _Greeter2_SayHiAgain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/stirling/testing/proto/greet.proto",
}

// StreamingGreeterClient is the client API for StreamingGreeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StreamingGreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (StreamingGreeter_SayHelloClient, error)
}

type streamingGreeterClient struct {
	cc *grpc.ClientConn
}

func NewStreamingGreeterClient(cc *grpc.ClientConn) StreamingGreeterClient {
	return &streamingGreeterClient{cc}
}

func (c *streamingGreeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (StreamingGreeter_SayHelloClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StreamingGreeter_serviceDesc.Streams[0], "/pl.stirling.testing.StreamingGreeter/SayHello", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamingGreeterSayHelloClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamingGreeter_SayHelloClient interface {
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type streamingGreeterSayHelloClient struct {
	grpc.ClientStream
}

func (x *streamingGreeterSayHelloClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamingGreeterServer is the server API for StreamingGreeter service.
type StreamingGreeterServer interface {
	SayHello(*HelloRequest, StreamingGreeter_SayHelloServer) error
}

func RegisterStreamingGreeterServer(s *grpc.Server, srv StreamingGreeterServer) {
	s.RegisterService(&_StreamingGreeter_serviceDesc, srv)
}

func _StreamingGreeter_SayHello_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamingGreeterServer).SayHello(m, &streamingGreeterSayHelloServer{stream})
}

type StreamingGreeter_SayHelloServer interface {
	Send(*HelloReply) error
	grpc.ServerStream
}

type streamingGreeterSayHelloServer struct {
	grpc.ServerStream
}

func (x *streamingGreeterSayHelloServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

var _StreamingGreeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pl.stirling.testing.StreamingGreeter",
	HandlerType: (*StreamingGreeterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SayHello",
			Handler:       _StreamingGreeter_SayHello_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "src/stirling/testing/proto/greet.proto",
}

func (m *HelloRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HelloRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintGreet(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.Count != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintGreet(dAtA, i, uint64(m.Count))
	}
	return i, nil
}

func (m *HelloReply) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HelloReply) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintGreet(dAtA, i, uint64(len(m.Message)))
		i += copy(dAtA[i:], m.Message)
	}
	return i, nil
}

func encodeVarintGreet(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *HelloRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovGreet(uint64(l))
	}
	if m.Count != 0 {
		n += 1 + sovGreet(uint64(m.Count))
	}
	return n
}

func (m *HelloReply) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovGreet(uint64(l))
	}
	return n
}

func sovGreet(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGreet(x uint64) (n int) {
	return sovGreet(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *HelloRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HelloRequest{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Count:` + fmt.Sprintf("%v", this.Count) + `,`,
		`}`,
	}, "")
	return s
}
func (this *HelloReply) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&HelloReply{`,
		`Message:` + fmt.Sprintf("%v", this.Message) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGreet(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *HelloRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGreet
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HelloRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HelloRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGreet
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGreet
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGreet
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGreet
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGreet(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGreet
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthGreet
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HelloReply) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGreet
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HelloReply: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HelloReply: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGreet
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGreet
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGreet
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGreet(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGreet
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthGreet
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGreet(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGreet
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGreet
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGreet
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGreet
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthGreet
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGreet
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGreet(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthGreet
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGreet = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGreet   = fmt.Errorf("proto: integer overflow")
)
