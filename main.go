package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	res, err := SearchRepositories("miroir")
	if err != nil {
		panic(err)
	}

	jsonStr, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", string(jsonStr))
}
