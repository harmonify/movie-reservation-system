package validation

type ValidationKey string

func (k ValidationKey) String() string {
	return string(k)
}

const (
	PhoneNumberKey ValidationKey = "phone_number"
	AlphaSpaceKey ValidationKey = "alpha_space"
)
