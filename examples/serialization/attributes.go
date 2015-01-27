// Demonstrates decision tree classification

package main

import (
	"encoding/json"
	"fmt"
	"github.com/sjwhitworth/golearn/base"
)

func main() {

	// Load in the iris dataset
	iris, err := base.ParseCSVToInstances("../datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	for _, a := range iris.AllAttributes() {
		// var ac base.CategoricalAttribute
		var af base.FloatAttribute
		s, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(s))
		err = json.Unmarshal(s, &af)
		fmt.Println(af.String())
		// err = json.Unmarshal(s, &ac)
		// fmt.Println(ac.String())
	}
}
