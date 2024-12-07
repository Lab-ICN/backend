package usecase

type Error struct {
	Code    int           `json:"-"`
	Message string        `json:"message"`
	Err     error         `json:"-"`
	Errors  []DomainError `json:"errors"`
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
	msgUserExist    = "user already exist"
	msgUserNotFound = "user not found"
)
