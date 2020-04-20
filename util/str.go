package util

import "encoding/json"

func ToJsonString(i interface{}) string {
	ibytes, _ := json.Marshal(i)
	return string(ibytes)
}
