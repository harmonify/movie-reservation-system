package opa

// OpaRequestBody is the HTTP request body sent to OPA server
// when requesting for access control decision.
// i.e. POST /v1/data/<policy_paths:.+>/allow
type OpaRequestBody struct {
	Input interface{} `json:"input"`
}

// OpaResponseBody is the HTTP response returned by OPA server
// when requesting for access control decision.
// i.e. POST /v1/data/<policy_paths:.+>/allow
type OpaResponseBody struct {
	Result bool `json:"result"`
}
