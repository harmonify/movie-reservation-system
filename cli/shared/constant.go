package shared

type Topic int

var (
	topics = [...]string{
		"new-order",
	}
)

const (
	NewOrderTopic Topic = iota
)

func (t Topic) String() string {
	return topics[t]
}
