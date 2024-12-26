package grpc_constant

// https://grpc.io/docs/guides/status-codes/
const (
	// Not an error; returned on success.
	GrpcOK = 0

	// The operation was cancelled, typically by the caller.
	GrpcCanceled = 1

	// Unknown error.
	GrpcUnknown = 2

	// The client specified an invalid argument.
	GrpcInvalidArgument = 3

	// The deadline expired before the operation could complete.
	GrpcDeadlineExceeded = 4

	// Some requested entity (e.g., file or directory) was not found.
	GrpcNotFound = 5

	// The entity that a client attempted to create (e.g., file or directory) already exists.
	GrpcAlreadyExists = 6

	// The caller does not have permission to execute the specified operation.
	GrpcPermissionDenied = 7

	// Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
	GrpcResourceExhausted = 8

	// 	The operation was rejected because the system is not in a state required for the operation's execution.
	GrpcFailedPrecondition = 9

	// The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort.
	GrpcAborted = 10

	// 	The operation was attempted past the valid range. E.g., seeking or reading past end-of-file.
	GrpcOutOfRange = 11

	// The operation is not implemented or is not supported/enabled in this service.
	GrpcUnimplemented = 12

	// Internal errors. This means that some invariants expected by the underlying system have been broken.

	// This error code is reserved for serious errors.
	GrpcInternal = 13

	// The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff.

	// Note that it is not always safe to retry non-idempotent operations.
	GrpcUnavailable = 14

	// Unrecoverable data loss or corruption.
	GrpcDataLoss = 15

	// The request does not have valid authentication credentials for the operation.
	GrpcUnauthenticated = 16
)
