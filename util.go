package gopcp_service

import (
	"encoding/json"
	"net/http"
)

func GetJsonBody(r *http.Request) (interface{}, error) {
	decorder := json.NewDecoder(r.Body)
	var arr interface{}
	if derr := decorder.Decode(&arr); derr != nil {
		return nil, derr
	} else {
		return arr, nil
	}
}
