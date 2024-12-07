package repository

type _error uint

const (
	ErrDuplicateRow _error = iota
	ErrNoRowAffected
	ErrNoRow
)

var _errors = []string{
	ErrDuplicateRow:  "DUPLICATE_ROW",
	ErrNoRowAffected: "NO_ROW_AFFECTED",
	ErrNoRow:         "NO_ROW",
}

func (e _error) Error() string {
	return _errors[e]
}
