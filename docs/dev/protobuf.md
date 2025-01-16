# Protobuf

An example to import another Protobuf in a `.proto` file:

```proto
// shared/index.proto
syntax = "proto3";

package harmonify.movie_reservation_system.shared;

import public "shared/test.proto";
```

## Resources

- [How to Go + Kafka + Protobuf](https://hashnode.com/post/how-to-use-protobuf-with-golang-and-kafka-cl0e1elcl04mgo5nv93et9uew)
