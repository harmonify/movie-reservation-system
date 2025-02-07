package kafka

import "context"

type (
	// Route handle incoming messages from a Topic.
	Route interface {
		// Identifier returns the route identifier. This identifier is useful for error debugging
		Identifier() string
		// Match determines if the route should handle the incoming event.
		// Match should return true if the route should handle the incoming event.
		Match(ctx context.Context, event *Event) (bool, error)
		// Handle handles the incoming event. If Handle returned an error, the router will mark the incoming message as read and send it to DLQ.
		Handle(ctx context.Context, event *Event) error
		// AddEventListener adds listener that will be triggered on incoming event. Mainly used for testing purposes.
		AddEventListener(listener EventListener)
	}
)
