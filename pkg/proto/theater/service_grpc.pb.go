// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: theater/service.proto

package theater_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TheaterService_GetActiveMovies_FullMethodName    = "/harmonify.movie_reservation_system.theater.TheaterService/GetActiveMovies"
	TheaterService_GetActiveShowtimes_FullMethodName = "/harmonify.movie_reservation_system.theater.TheaterService/GetActiveShowtimes"
	TheaterService_GetAvailableSeats_FullMethodName  = "/harmonify.movie_reservation_system.theater.TheaterService/GetAvailableSeats"
)

// TheaterServiceClient is the client API for TheaterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TheaterServiceClient interface {
	// Get movies with active showtimes
	GetActiveMovies(ctx context.Context, in *GetActiveMoviesRequest, opts ...grpc.CallOption) (*GetActiveMoviesResponse, error)
	// Get active showtimes for a movie
	GetActiveShowtimes(ctx context.Context, in *GetActiveShowtimesRequest, opts ...grpc.CallOption) (*GetActiveShowtimesResponse, error)
	// Get available seats for a showtime
	GetAvailableSeats(ctx context.Context, in *GetAvailableSeatsRequest, opts ...grpc.CallOption) (*GetAvailableSeatsResponse, error)
}

type theaterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTheaterServiceClient(cc grpc.ClientConnInterface) TheaterServiceClient {
	return &theaterServiceClient{cc}
}

func (c *theaterServiceClient) GetActiveMovies(ctx context.Context, in *GetActiveMoviesRequest, opts ...grpc.CallOption) (*GetActiveMoviesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetActiveMoviesResponse)
	err := c.cc.Invoke(ctx, TheaterService_GetActiveMovies_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *theaterServiceClient) GetActiveShowtimes(ctx context.Context, in *GetActiveShowtimesRequest, opts ...grpc.CallOption) (*GetActiveShowtimesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetActiveShowtimesResponse)
	err := c.cc.Invoke(ctx, TheaterService_GetActiveShowtimes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *theaterServiceClient) GetAvailableSeats(ctx context.Context, in *GetAvailableSeatsRequest, opts ...grpc.CallOption) (*GetAvailableSeatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAvailableSeatsResponse)
	err := c.cc.Invoke(ctx, TheaterService_GetAvailableSeats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TheaterServiceServer is the server API for TheaterService service.
// All implementations must embed UnimplementedTheaterServiceServer
// for forward compatibility.
type TheaterServiceServer interface {
	// Get movies with active showtimes
	GetActiveMovies(context.Context, *GetActiveMoviesRequest) (*GetActiveMoviesResponse, error)
	// Get active showtimes for a movie
	GetActiveShowtimes(context.Context, *GetActiveShowtimesRequest) (*GetActiveShowtimesResponse, error)
	// Get available seats for a showtime
	GetAvailableSeats(context.Context, *GetAvailableSeatsRequest) (*GetAvailableSeatsResponse, error)
	mustEmbedUnimplementedTheaterServiceServer()
}

// UnimplementedTheaterServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTheaterServiceServer struct{}

func (UnimplementedTheaterServiceServer) GetActiveMovies(context.Context, *GetActiveMoviesRequest) (*GetActiveMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveMovies not implemented")
}
func (UnimplementedTheaterServiceServer) GetActiveShowtimes(context.Context, *GetActiveShowtimesRequest) (*GetActiveShowtimesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveShowtimes not implemented")
}
func (UnimplementedTheaterServiceServer) GetAvailableSeats(context.Context, *GetAvailableSeatsRequest) (*GetAvailableSeatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableSeats not implemented")
}
func (UnimplementedTheaterServiceServer) mustEmbedUnimplementedTheaterServiceServer() {}
func (UnimplementedTheaterServiceServer) testEmbeddedByValue()                        {}

// UnsafeTheaterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TheaterServiceServer will
// result in compilation errors.
type UnsafeTheaterServiceServer interface {
	mustEmbedUnimplementedTheaterServiceServer()
}

func RegisterTheaterServiceServer(s grpc.ServiceRegistrar, srv TheaterServiceServer) {
	// If the following call pancis, it indicates UnimplementedTheaterServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TheaterService_ServiceDesc, srv)
}

func _TheaterService_GetActiveMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActiveMoviesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TheaterServiceServer).GetActiveMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TheaterService_GetActiveMovies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TheaterServiceServer).GetActiveMovies(ctx, req.(*GetActiveMoviesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TheaterService_GetActiveShowtimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActiveShowtimesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TheaterServiceServer).GetActiveShowtimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TheaterService_GetActiveShowtimes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TheaterServiceServer).GetActiveShowtimes(ctx, req.(*GetActiveShowtimesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TheaterService_GetAvailableSeats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableSeatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TheaterServiceServer).GetAvailableSeats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TheaterService_GetAvailableSeats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TheaterServiceServer).GetAvailableSeats(ctx, req.(*GetAvailableSeatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TheaterService_ServiceDesc is the grpc.ServiceDesc for TheaterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TheaterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "harmonify.movie_reservation_system.theater.TheaterService",
	HandlerType: (*TheaterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetActiveMovies",
			Handler:    _TheaterService_GetActiveMovies_Handler,
		},
		{
			MethodName: "GetActiveShowtimes",
			Handler:    _TheaterService_GetActiveShowtimes_Handler,
		},
		{
			MethodName: "GetAvailableSeats",
			Handler:    _TheaterService_GetAvailableSeats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "theater/service.proto",
}
