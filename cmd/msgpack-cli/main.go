package main

import (
	"fmt"

	"github.com/mu8086/msgpack"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := msgpack.InitConstants(); err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	jsonData := []byte(
		`{
			"friends": [
				{
					"id": 0,
					"name": "Swanson Ayers"
				},
				true,
				"string",
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
				12,
				13,
				14,
				15,
				16
			]
		}`)

	mp, err := msgpack.JSONToMessagePack(jsonData)
	fmt.Printf("mp: %v, err: %v\n", formatBytecode(mp), err)
}

// TODO: remove
func formatBytecode(bytecode []byte) (s string) {
	for _, b := range bytecode {
		s += fmt.Sprintf("% 03X", b)
	}
	if len(s) > 0 { // remove leading space
		s = s[1:]
	}
	return s
}
