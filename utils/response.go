package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		data = "{\"status\" : \"success\"}"
	}
	json.NewEncoder(w).Encode(data)

}

func ResponseError(w http.ResponseWriter, err error) {
	fmt.Printf("ERROR: %v", err)
	ResponseJSON(w, 500, "{Error: \"Internal Server Error\"}")
}

func ResponseErrors(w http.ResponseWriter, errs []error) {
	res := map[string][]string{}
	for _, err := range errs {
		fmt.Printf("ERROR: %v\n", err)
		res["errors"] = append(res["errors"], err.Error())
	}
	ResponseJSON(w, 500, res)
}
