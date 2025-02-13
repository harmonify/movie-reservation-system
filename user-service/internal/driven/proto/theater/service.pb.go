// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.29.3
// source: theater/service.proto

package theater_proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetActiveMoviesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TheaterId       string `protobuf:"bytes,1,opt,name=theater_id,json=theaterId,proto3" json:"theater_id,omitempty"`
	IncludeUpcoming bool   `protobuf:"varint,2,opt,name=include_upcoming,json=includeUpcoming,proto3" json:"include_upcoming,omitempty"`
}

func (x *GetActiveMoviesRequest) Reset() {
	*x = GetActiveMoviesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActiveMoviesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveMoviesRequest) ProtoMessage() {}

func (x *GetActiveMoviesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveMoviesRequest.ProtoReflect.Descriptor instead.
func (*GetActiveMoviesRequest) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetActiveMoviesRequest) GetTheaterId() string {
	if x != nil {
		return x.TheaterId
	}
	return ""
}

func (x *GetActiveMoviesRequest) GetIncludeUpcoming() bool {
	if x != nil {
		return x.IncludeUpcoming
	}
	return false
}

type GetActiveMoviesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of unique movies with active showtimes
	Movies []*GetActiveMoviesResponse_Movie `protobuf:"bytes,1,rep,name=movies,proto3" json:"movies,omitempty"`
}

func (x *GetActiveMoviesResponse) Reset() {
	*x = GetActiveMoviesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActiveMoviesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveMoviesResponse) ProtoMessage() {}

func (x *GetActiveMoviesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveMoviesResponse.ProtoReflect.Descriptor instead.
func (*GetActiveMoviesResponse) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetActiveMoviesResponse) GetMovies() []*GetActiveMoviesResponse_Movie {
	if x != nil {
		return x.Movies
	}
	return nil
}

type GetActiveShowtimesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TheaterId string `protobuf:"bytes,1,opt,name=theater_id,json=theaterId,proto3" json:"theater_id,omitempty"`
	MovieId   string `protobuf:"bytes,2,opt,name=movie_id,json=movieId,proto3" json:"movie_id,omitempty"`
}

func (x *GetActiveShowtimesRequest) Reset() {
	*x = GetActiveShowtimesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActiveShowtimesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveShowtimesRequest) ProtoMessage() {}

func (x *GetActiveShowtimesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveShowtimesRequest.ProtoReflect.Descriptor instead.
func (*GetActiveShowtimesRequest) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetActiveShowtimesRequest) GetTheaterId() string {
	if x != nil {
		return x.TheaterId
	}
	return ""
}

func (x *GetActiveShowtimesRequest) GetMovieId() string {
	if x != nil {
		return x.MovieId
	}
	return ""
}

type GetActiveShowtimesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Showtimes []*GetActiveShowtimesResponse_Showtime `protobuf:"bytes,1,rep,name=showtimes,proto3" json:"showtimes,omitempty"`
}

func (x *GetActiveShowtimesResponse) Reset() {
	*x = GetActiveShowtimesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActiveShowtimesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveShowtimesResponse) ProtoMessage() {}

func (x *GetActiveShowtimesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveShowtimesResponse.ProtoReflect.Descriptor instead.
func (*GetActiveShowtimesResponse) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetActiveShowtimesResponse) GetShowtimes() []*GetActiveShowtimesResponse_Showtime {
	if x != nil {
		return x.Showtimes
	}
	return nil
}

type GetAvailableSeatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShowtimeId string `protobuf:"bytes,1,opt,name=showtime_id,json=showtimeId,proto3" json:"showtime_id,omitempty"`
}

func (x *GetAvailableSeatsRequest) Reset() {
	*x = GetAvailableSeatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAvailableSeatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableSeatsRequest) ProtoMessage() {}

func (x *GetAvailableSeatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableSeatsRequest.ProtoReflect.Descriptor instead.
func (*GetAvailableSeatsRequest) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetAvailableSeatsRequest) GetShowtimeId() string {
	if x != nil {
		return x.ShowtimeId
	}
	return ""
}

type GetAvailableSeatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seats []*GetAvailableSeatsResponse_Seat `protobuf:"bytes,1,rep,name=seats,proto3" json:"seats,omitempty"`
}

func (x *GetAvailableSeatsResponse) Reset() {
	*x = GetAvailableSeatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAvailableSeatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableSeatsResponse) ProtoMessage() {}

func (x *GetAvailableSeatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableSeatsResponse.ProtoReflect.Descriptor instead.
func (*GetAvailableSeatsResponse) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{5}
}

func (x *GetAvailableSeatsResponse) GetSeats() []*GetAvailableSeatsResponse_Seat {
	if x != nil {
		return x.Seats
	}
	return nil
}

type GetActiveMoviesResponse_Movie struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MovieId string `protobuf:"bytes,1,opt,name=movie_id,json=movieId,proto3" json:"movie_id,omitempty"`
}

func (x *GetActiveMoviesResponse_Movie) Reset() {
	*x = GetActiveMoviesResponse_Movie{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActiveMoviesResponse_Movie) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveMoviesResponse_Movie) ProtoMessage() {}

func (x *GetActiveMoviesResponse_Movie) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveMoviesResponse_Movie.ProtoReflect.Descriptor instead.
func (*GetActiveMoviesResponse_Movie) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{1, 0}
}

func (x *GetActiveMoviesResponse_Movie) GetMovieId() string {
	if x != nil {
		return x.MovieId
	}
	return ""
}

type GetActiveShowtimesResponse_Showtime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShowtimeId     string `protobuf:"bytes,1,opt,name=showtime_id,json=showtimeId,proto3" json:"showtime_id,omitempty"`
	StartTime      uint32 `protobuf:"varint,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	AvailableSeats uint32 `protobuf:"varint,3,opt,name=available_seats,json=availableSeats,proto3" json:"available_seats,omitempty"`
}

func (x *GetActiveShowtimesResponse_Showtime) Reset() {
	*x = GetActiveShowtimesResponse_Showtime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActiveShowtimesResponse_Showtime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveShowtimesResponse_Showtime) ProtoMessage() {}

func (x *GetActiveShowtimesResponse_Showtime) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveShowtimesResponse_Showtime.ProtoReflect.Descriptor instead.
func (*GetActiveShowtimesResponse_Showtime) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{3, 0}
}

func (x *GetActiveShowtimesResponse_Showtime) GetShowtimeId() string {
	if x != nil {
		return x.ShowtimeId
	}
	return ""
}

func (x *GetActiveShowtimesResponse_Showtime) GetStartTime() uint32 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GetActiveShowtimesResponse_Showtime) GetAvailableSeats() uint32 {
	if x != nil {
		return x.AvailableSeats
	}
	return 0
}

type GetAvailableSeatsResponse_Seat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SeatId     string `protobuf:"bytes,1,opt,name=seat_id,json=seatId,proto3" json:"seat_id,omitempty"`
	SeatRow    string `protobuf:"bytes,2,opt,name=seat_row,json=seatRow,proto3" json:"seat_row,omitempty"`
	SeatColumn string `protobuf:"bytes,3,opt,name=seat_column,json=seatColumn,proto3" json:"seat_column,omitempty"`
}

func (x *GetAvailableSeatsResponse_Seat) Reset() {
	*x = GetAvailableSeatsResponse_Seat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_theater_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAvailableSeatsResponse_Seat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAvailableSeatsResponse_Seat) ProtoMessage() {}

func (x *GetAvailableSeatsResponse_Seat) ProtoReflect() protoreflect.Message {
	mi := &file_theater_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAvailableSeatsResponse_Seat.ProtoReflect.Descriptor instead.
func (*GetAvailableSeatsResponse_Seat) Descriptor() ([]byte, []int) {
	return file_theater_service_proto_rawDescGZIP(), []int{5, 0}
}

func (x *GetAvailableSeatsResponse_Seat) GetSeatId() string {
	if x != nil {
		return x.SeatId
	}
	return ""
}

func (x *GetAvailableSeatsResponse_Seat) GetSeatRow() string {
	if x != nil {
		return x.SeatRow
	}
	return ""
}

func (x *GetAvailableSeatsResponse_Seat) GetSeatColumn() string {
	if x != nil {
		return x.SeatColumn
	}
	return ""
}

var File_theater_service_proto protoreflect.FileDescriptor

var file_theater_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2a, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69,
	0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61,
	0x74, 0x65, 0x72, 0x22, 0x62, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x10,
	0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x75, 0x70, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x55,
	0x70, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x22, 0xa0, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x06, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x49, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x06,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x1a, 0x22, 0x0a, 0x05, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49, 0x64, 0x22, 0x55, 0x0a, 0x19, 0x47, 0x65,
	0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x68, 0x65, 0x61, 0x74,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x68, 0x65,
	0x61, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x49,
	0x64, 0x22, 0x80, 0x02, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53,
	0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x6d, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x4f, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x77, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x68, 0x6f, 0x77,
	0x74, 0x69, 0x6d, 0x65, 0x52, 0x09, 0x73, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x1a,
	0x73, 0x0a, 0x08, 0x53, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x61,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x73, 0x65, 0x61, 0x74, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x53,
	0x65, 0x61, 0x74, 0x73, 0x22, 0x3b, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x53, 0x65, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x49,
	0x64, 0x22, 0xda, 0x01, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x6c, 0x65, 0x53, 0x65, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x60, 0x0a, 0x05, 0x73, 0x65, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x4a,
	0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x65, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x74, 0x52, 0x05, 0x73, 0x65, 0x61, 0x74,
	0x73, 0x1a, 0x5b, 0x0a, 0x04, 0x53, 0x65, 0x61, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x65, 0x61,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x61, 0x74,
	0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x65, 0x61, 0x74, 0x5f, 0x72, 0x6f, 0x77, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x61, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x1f, 0x0a,
	0x0b, 0x73, 0x65, 0x61, 0x74, 0x5f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x61, 0x74, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x32, 0xfc,
	0x03, 0x0a, 0x0e, 0x54, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x9c, 0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x73, 0x12, 0x42, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66,
	0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74,
	0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x43, 0x2e, 0x68, 0x61, 0x72, 0x6d,
	0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65,
	0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74,
	0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0xa5, 0x01, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x68,
	0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x12, 0x45, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e,
	0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65,
	0x61, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x68,
	0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x46,
	0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x77, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0xa2, 0x01, 0x0a, 0x11, 0x47, 0x65, 0x74,
	0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x65, 0x61, 0x74, 0x73, 0x12, 0x44,
	0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x65, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x45, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79,
	0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x74, 0x68, 0x65, 0x61, 0x74, 0x65,
	0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x65,
	0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_theater_service_proto_rawDescOnce sync.Once
	file_theater_service_proto_rawDescData = file_theater_service_proto_rawDesc
)

func file_theater_service_proto_rawDescGZIP() []byte {
	file_theater_service_proto_rawDescOnce.Do(func() {
		file_theater_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_theater_service_proto_rawDescData)
	})
	return file_theater_service_proto_rawDescData
}

var file_theater_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_theater_service_proto_goTypes = []interface{}{
	(*GetActiveMoviesRequest)(nil),              // 0: harmonify.movie_reservation_system.theater.GetActiveMoviesRequest
	(*GetActiveMoviesResponse)(nil),             // 1: harmonify.movie_reservation_system.theater.GetActiveMoviesResponse
	(*GetActiveShowtimesRequest)(nil),           // 2: harmonify.movie_reservation_system.theater.GetActiveShowtimesRequest
	(*GetActiveShowtimesResponse)(nil),          // 3: harmonify.movie_reservation_system.theater.GetActiveShowtimesResponse
	(*GetAvailableSeatsRequest)(nil),            // 4: harmonify.movie_reservation_system.theater.GetAvailableSeatsRequest
	(*GetAvailableSeatsResponse)(nil),           // 5: harmonify.movie_reservation_system.theater.GetAvailableSeatsResponse
	(*GetActiveMoviesResponse_Movie)(nil),       // 6: harmonify.movie_reservation_system.theater.GetActiveMoviesResponse.Movie
	(*GetActiveShowtimesResponse_Showtime)(nil), // 7: harmonify.movie_reservation_system.theater.GetActiveShowtimesResponse.Showtime
	(*GetAvailableSeatsResponse_Seat)(nil),      // 8: harmonify.movie_reservation_system.theater.GetAvailableSeatsResponse.Seat
}
var file_theater_service_proto_depIdxs = []int32{
	6, // 0: harmonify.movie_reservation_system.theater.GetActiveMoviesResponse.movies:type_name -> harmonify.movie_reservation_system.theater.GetActiveMoviesResponse.Movie
	7, // 1: harmonify.movie_reservation_system.theater.GetActiveShowtimesResponse.showtimes:type_name -> harmonify.movie_reservation_system.theater.GetActiveShowtimesResponse.Showtime
	8, // 2: harmonify.movie_reservation_system.theater.GetAvailableSeatsResponse.seats:type_name -> harmonify.movie_reservation_system.theater.GetAvailableSeatsResponse.Seat
	0, // 3: harmonify.movie_reservation_system.theater.TheaterService.GetActiveMovies:input_type -> harmonify.movie_reservation_system.theater.GetActiveMoviesRequest
	2, // 4: harmonify.movie_reservation_system.theater.TheaterService.GetActiveShowtimes:input_type -> harmonify.movie_reservation_system.theater.GetActiveShowtimesRequest
	4, // 5: harmonify.movie_reservation_system.theater.TheaterService.GetAvailableSeats:input_type -> harmonify.movie_reservation_system.theater.GetAvailableSeatsRequest
	1, // 6: harmonify.movie_reservation_system.theater.TheaterService.GetActiveMovies:output_type -> harmonify.movie_reservation_system.theater.GetActiveMoviesResponse
	3, // 7: harmonify.movie_reservation_system.theater.TheaterService.GetActiveShowtimes:output_type -> harmonify.movie_reservation_system.theater.GetActiveShowtimesResponse
	5, // 8: harmonify.movie_reservation_system.theater.TheaterService.GetAvailableSeats:output_type -> harmonify.movie_reservation_system.theater.GetAvailableSeatsResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_theater_service_proto_init() }
func file_theater_service_proto_init() {
	if File_theater_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_theater_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActiveMoviesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActiveMoviesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActiveShowtimesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActiveShowtimesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAvailableSeatsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAvailableSeatsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActiveMoviesResponse_Movie); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActiveShowtimesResponse_Showtime); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_theater_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAvailableSeatsResponse_Seat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_theater_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_theater_service_proto_goTypes,
		DependencyIndexes: file_theater_service_proto_depIdxs,
		MessageInfos:      file_theater_service_proto_msgTypes,
	}.Build()
	File_theater_service_proto = out.File
	file_theater_service_proto_rawDesc = nil
	file_theater_service_proto_goTypes = nil
	file_theater_service_proto_depIdxs = nil
}
