/*
 * Copyright (c) adixity 2020. https://github.com/adixity
 */

package endpoints

import (
	"compress/gzip"
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
	if response != nil {
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseGZip(w, js, "application/json")
	} else {
		w.WriteHeader(200)
	}
}

func responseGZip(w http.ResponseWriter, content []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		// fallback to original way
		_, _ = w.Write(content)
		return
	}

	defer writer.Close()
	w.Header().Set("Content-Encoding", "gzip")
	_, _ = writer.Write(content)
}
