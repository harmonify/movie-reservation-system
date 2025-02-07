package grpc_driver

type (
	// Create type for validate incoming protobuf request
	SendEmailRequest struct {
		Email string `json:"email"`
	}
		
)