package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func isJSON(b []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(b, &js) == nil
}

func isYaml(b []byte) bool {
	var ya map[string]interface{}
	return yaml.Unmarshal(b, &ya) == nil
}

func convertJSON(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertJSON(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertJSON(v)
		}
	}
	return i
}

func main() {
	if len(os.Args) != 2 {
		panic("illegal args\ntjy foo.json or tjy bar.json")
	}

	bin, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var buffer interface{}
	if isJSON(bin) {
		err := json.Unmarshal(bin, &buffer)
		if err != nil {
			panic(err)
		}
		result, err := yaml.Marshal(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(result))
		os.Exit(0)
	} else if isYaml(bin) {
		err := yaml.Unmarshal(bin, &buffer)
		if err != nil {
			panic(err)
		}
		result, err := json.MarshalIndent(convertJSON(buffer), "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(result))
		os.Exit(0)
	}
	panic(os.Args[1] + " is neither json nor yaml.")
}
