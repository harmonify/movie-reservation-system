package shared

const (
	PublicUserRegisteredV1Topic Topic = "public.user.registered.v1"
)

type Topic string

func (t Topic) String() string {
	return string(t)
}
