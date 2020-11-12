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

func sendResponse(w http.ResponseWriter, response *apiResponse) {
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// explicitly specify cache-control here to prevent gcp-frontend server caching
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
