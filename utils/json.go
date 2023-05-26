package utils

import "encoding/json"

func MustMarshalString(v interface{}) string {
	result, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(result)
}
