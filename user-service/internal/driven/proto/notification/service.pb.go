// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.29.3
// source: notification/service.proto

package notification_proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_notification_service_proto protoreflect.FileDescriptor

var file_notification_service_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2f, 0x68, 0x61,
	0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x18, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32,
	0xd4, 0x03, 0x0a, 0x13, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x92, 0x01, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x41, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66,
	0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x42, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f,
	0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x8c, 0x01, 0x0a,
	0x07, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x6d, 0x73, 0x12, 0x3f, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f,
	0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x53,
	0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x40, 0x2e, 0x68, 0x61, 0x72, 0x6d,
	0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65,
	0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x53, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x98, 0x01, 0x0a, 0x0b,
	0x42, 0x75, 0x6c, 0x6b, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x6d, 0x73, 0x12, 0x43, 0x2e, 0x68, 0x61,
	0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x5f, 0x72, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x42, 0x75,
	0x6c, 0x6b, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x44, 0x2e, 0x68, 0x61, 0x72, 0x6d, 0x6f, 0x6e, 0x69, 0x66, 0x79, 0x2e, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x42, 0x75, 0x6c, 0x6b, 0x53, 0x65, 0x6e, 0x64, 0x53, 0x6d, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_notification_service_proto_goTypes = []interface{}{
	(*SendEmailRequest)(nil),    // 0: harmonify.movie_reservation_system.notification.SendEmailRequest
	(*SendSmsRequest)(nil),      // 1: harmonify.movie_reservation_system.notification.SendSmsRequest
	(*BulkSendSmsRequest)(nil),  // 2: harmonify.movie_reservation_system.notification.BulkSendSmsRequest
	(*SendEmailResponse)(nil),   // 3: harmonify.movie_reservation_system.notification.SendEmailResponse
	(*SendSmsResponse)(nil),     // 4: harmonify.movie_reservation_system.notification.SendSmsResponse
	(*BulkSendSmsResponse)(nil), // 5: harmonify.movie_reservation_system.notification.BulkSendSmsResponse
}
var file_notification_service_proto_depIdxs = []int32{
	0, // 0: harmonify.movie_reservation_system.notification.NotificationService.SendEmail:input_type -> harmonify.movie_reservation_system.notification.SendEmailRequest
	1, // 1: harmonify.movie_reservation_system.notification.NotificationService.SendSms:input_type -> harmonify.movie_reservation_system.notification.SendSmsRequest
	2, // 2: harmonify.movie_reservation_system.notification.NotificationService.BulkSendSms:input_type -> harmonify.movie_reservation_system.notification.BulkSendSmsRequest
	3, // 3: harmonify.movie_reservation_system.notification.NotificationService.SendEmail:output_type -> harmonify.movie_reservation_system.notification.SendEmailResponse
	4, // 4: harmonify.movie_reservation_system.notification.NotificationService.SendSms:output_type -> harmonify.movie_reservation_system.notification.SendSmsResponse
	5, // 5: harmonify.movie_reservation_system.notification.NotificationService.BulkSendSms:output_type -> harmonify.movie_reservation_system.notification.BulkSendSmsResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_notification_service_proto_init() }
func file_notification_service_proto_init() {
	if File_notification_service_proto != nil {
		return
	}
	file_notification_email_proto_init()
	file_notification_sms_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_notification_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_notification_service_proto_goTypes,
		DependencyIndexes: file_notification_service_proto_depIdxs,
	}.Build()
	File_notification_service_proto = out.File
	file_notification_service_proto_rawDesc = nil
	file_notification_service_proto_goTypes = nil
	file_notification_service_proto_depIdxs = nil
}
