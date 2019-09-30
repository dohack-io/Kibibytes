package secrets

import (
	"encoding/json"
	"io/ioutil"
)

var secrets map[string]string

func init() {
	text, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(text, &secrets)
	if err != nil {
		panic(err)
	}
}

func Get(key string) string {
	return secrets[key]
}
