package msgpack

type ErrorType struct {
	ErrCode uint16
	ErrStr  string
}

const (
	ErrCodeUnsupportedType uint16 = iota + 1
	ErrCodeStringTooLong
	ErrCodeValueOutOfRange
	ErrCodeBinaryTooLong
	ErrCodeBinaryDataInvalid
	ErrCodeInitConstants
	ErrCodeArrayTooLong
)

const (
	ErrStrUnsupportedType   = "UnsupportedType"
	ErrStrStringTooLong     = "StringTooLong"
	ErrStrValueOutOfRange   = "ValueOutOfRange"
	ErrStrBinaryTooLong     = "BinaryTooLong"
	ErrStrBinaryDataInvalid = "BinaryDataInvalid"
	ErrStrInitConstants     = "InitConstants"
	ErrStrArrayTooLong      = "ArrayTooLong"
)

var (
	ErrUnsupportedType   = ErrorType{ErrCode: ErrCodeUnsupportedType, ErrStr: ErrStrUnsupportedType}
	ErrStringTooLong     = ErrorType{ErrCode: ErrCodeStringTooLong, ErrStr: ErrStrStringTooLong}
	ErrValueOutOfRange   = ErrorType{ErrCode: ErrCodeValueOutOfRange, ErrStr: ErrStrValueOutOfRange}
	ErrBinaryTooLong     = ErrorType{ErrCode: ErrCodeBinaryTooLong, ErrStr: ErrStrBinaryTooLong}
	ErrBinaryDataInvalid = ErrorType{ErrCode: ErrCodeBinaryDataInvalid, ErrStr: ErrStrBinaryDataInvalid}
	ErrInitConstants     = ErrorType{ErrCode: ErrCodeInitConstants, ErrStr: ErrStrInitConstants}
	ErrArrayTooLong      = ErrorType{ErrCode: ErrCodeArrayTooLong, ErrStr: ErrStrArrayTooLong}
)

func (e ErrorType) Error() string {
	return e.ErrStr
}
