package requests

type HttpDraftRequest struct {
	Data map[string]string `json:"data"`
}

type HttpDraftResponse struct {
	FileId string `json:"fileId"`
}
