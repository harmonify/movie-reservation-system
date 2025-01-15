# Protobuf

An example to import another Protobuf in a `.proto` file:

```proto
// shared/index.proto
syntax = "proto3";

package harmonify.movie_reservation_system.shared;

import public "shared/test.proto";
```
