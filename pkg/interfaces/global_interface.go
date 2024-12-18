package interfaces

type (
	// similar to [database/sql#NullBool]
	NullBool struct {
		Bool  bool
		Valid bool // Valid should be set to `true` if `Bool` is not null
	}
)
