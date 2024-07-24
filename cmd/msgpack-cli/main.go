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
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := msgpack.InitConstants(); err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	jsonData := []byte(
		`{
    "_id": "66a0d3af2f64df4a43dc28ca",
    "index": 0,
    "guid": "4b984e02-fd5c-49f9-8659-5639c242a866",
    "isActive": true,
    "balance": "$1,135.79",
    "picture": "http://placehold.it/32x32",
    "age": 34,
    "eyeColor": "blue",
    "name": "Cline Maddox",
    "gender": "male",
    "company": "VISUALIX",
    "email": "clinemaddox@visualix.com",
    "phone": "+1 (812) 560-2587",
    "address": "665 Adler Place, Mulberry, New Jersey, 443",
    "about": "Irure irure nostrud officia duis nulla laborum ipsum non qui nulla cupidatat exercitation dolore. Proident cillum consequat nulla laboris occaecat. Eiusmod commodo duis ad deserunt elit tempor labore irure aute anim nisi.\r\n",
    "registered": "2015-10-09T08:57:58 -08:00",
    "latitude": 84.989192,
    "longitude": -113.638803,
    "tags": [
      "cupidatat",
      "in",
      "magna",
      "deserunt",
      "duis",
      "elit",
      "reprehenderit"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Floyd Stone"
      },
      {
        "id": 1,
        "name": "Kirby Pearson"
      },
      {
        "id": 2,
        "name": "Jane Chapman"
      }
    ],
    "greeting": "Hello, Cline Maddox! You have 10 unread messages.",
    "favoriteFruit": "strawberry"
  }`)

	mp, err := msgpack.JSONToMessagePack(jsonData)
	fmt.Printf("mp: %v, err: %v\n", formatBytecode(mp), err)

	jsonData2, err := msgpack.MessagePackToJSON(mp)
	fmt.Printf("jsonData2: %v, err: %v\n", jsonData2, err)
}

// TODO: remove
func formatBytecode(bytecode []byte) (s string) {
	for _, b := range bytecode {
		s += fmt.Sprintf("% #03X,", b)
	}
	if len(s) > 0 { // remove leading space
		s = s[1:]
	}
	return s
}
