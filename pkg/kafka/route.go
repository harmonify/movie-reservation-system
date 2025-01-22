package kafka

import "context"

type (
	// Route handle incoming messages from a Topic.
	Route interface {
		// Identifier returns the route identifier. This identifier is useful for error debugging
		Identifier() string
		// Decode decodes incoming message value
		Decode(ctx context.Context, value []byte) (interface{}, error)
		// Match determines if the route should handle the incoming event (message that has been decoded)
		Match(ctx context.Context, event *Event) (bool, error)
		// Handle handles the incoming event. If the function returned an error, the router will NOT mark the incoming message as read.
		Handle(ctx context.Context, event *Event) error
		// AddEventListener adds listener that will be triggered on incoming event. Mainly used for testing purposes.
		AddEventListener(listener EventListener)
	}
)
