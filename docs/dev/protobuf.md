# Protobuf

An example to import another Protobuf in a `.proto` file:

```proto
// proto/harmonify/movie_reservation_system/notification/service.proto
syntax = "proto3";

package harmonify.movie_reservation_system.notification;

import "notification/email.proto";
import "notification/sms.proto";
import "shared/response.proto";

service NotificationService {
    rpc SendEmail(SendEmailRequest) returns (shared.Response);
    rpc SendSms(SendSmsRequest) returns (shared.Response);
    rpc BulkSendSms(BulkSendSmsRequest) returns (shared.Response);
}

```

## Resources

- [How to Go + Kafka + Protobuf](https://hashnode.com/post/how-to-use-protobuf-with-golang-and-kafka-cl0e1elcl04mgo5nv93et9uew)
