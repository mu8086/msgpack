package msgpack

type ErrorType struct {
	ErrCode uint16
	ErrStr  string
}

const (
	ErrCodeUnsupportedType uint16 = iota + 1
	ErrCodeStringTooLong
)

const (
	ErrStrUnsupportedType = "UnsupportedType"
	ErrStrStringTooLong   = "StringTooLong"
)

var (
	ErrUnsupportedType = ErrorType{ErrCode: ErrCodeUnsupportedType, ErrStr: ErrStrUnsupportedType}
	ErrStringTooLong   = ErrorType{ErrCode: ErrCodeStringTooLong, ErrStr: ErrStrStringTooLong}
)

func (e ErrorType) Error() string {
	return e.ErrStr
}
