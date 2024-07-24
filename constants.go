package msgpack

import (
	"github.com/spf13/viper"
)

var binaryKeyword string

func InitConstants() error {
	binaryKeyword = viper.GetString("binary_keyword")
	if binaryKeyword == "" {
		return ErrInitConstants
	}

	return nil
}
