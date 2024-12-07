package usecase

type Error struct {
	Code    int           `json:"-"`
	Message string        `json:"message,omitempty"`
	Err     error         `json:"-"`
	Errors  []DomainError `json:"errors,omitempty"`
}

type DomainError struct {
	Reason   string `json:"reason"`
	Message  string `json:"message"`
	Location string `json:"location"`
}

func (e Error) Error() string {
	return e.Message
}

const (
	msgUserNotRegistered = "user is not registered"
)
