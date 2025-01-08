package shared

var (
	topics = [...]string{
		"new-order",
	}
)

const (
	NewOrderTopic Topic = iota
)

type Topic int

func (t Topic) String() string {
	return topics[t]
}
