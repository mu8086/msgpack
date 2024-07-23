package msgpack

type ErrorType struct {
	ErrCode uint16
	ErrStr  string
}

const (
	ErrCodeUnsupportedType uint16 = iota + 1
	ErrCodeStringTooLong
	ErrCodeValueOutOfRange
)

const (
	ErrStrUnsupportedType = "UnsupportedType"
	ErrStrStringTooLong   = "StringTooLong"
	ErrStrValueOutOfRange = "ValueOutOfRange"
)

var (
	ErrUnsupportedType = ErrorType{ErrCode: ErrCodeUnsupportedType, ErrStr: ErrStrUnsupportedType}
	ErrStringTooLong   = ErrorType{ErrCode: ErrCodeStringTooLong, ErrStr: ErrStrStringTooLong}
	ErrValueOutOfRange = ErrorType{ErrCode: ErrCodeValueOutOfRange, ErrStr: ErrStrValueOutOfRange}
)

func (e ErrorType) Error() string {
	return e.ErrStr
}
