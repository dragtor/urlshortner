package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseNewData struct {
	Message string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Success bool        `json:"success"`
}

func HTTPResponseData(w http.ResponseWriter, success bool, result interface{}, errmsg string, statusCode int) {
	var respData ResponseNewData
	respData.Success = success
	if success {
		respData.Result = result
	}
	if !success {
		respData.Message = errmsg
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResponse, err := json.Marshal(respData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Failed to write JSON response:", err)
	}
}
