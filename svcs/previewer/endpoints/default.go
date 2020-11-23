/*
 * Copyright (c) adixity 2020. https://github.com/adixity
 */

package endpoints

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func newResponse(code int32, message string, data interface{}) *apiResponse {
	return &apiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func successResponse(data interface{}) *apiResponse {
	return &apiResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func sendNotFound(w http.ResponseWriter) {
	http.Error(w, "not found", http.StatusNotFound)
	// explicitly specify cache-control here to prevent gcp-frontend server caching
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte("not found"))
}

func sendResponse(w http.ResponseWriter, response *apiResponse) {

	// explicitly specify cache-control here to prevent gcp-frontend server caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if response != nil {
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(js)
	}
}
