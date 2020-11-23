package requests

type HttpDraftRequest struct {
	Data map[string]string `json:"data"`
}

type HttpDraftResponse struct {
	FileId string `json:"fileId"`
}

type HttpPublishResponse struct {
	Revision int `json:"revision"`
}

type HttpMetaResponse struct {
}
