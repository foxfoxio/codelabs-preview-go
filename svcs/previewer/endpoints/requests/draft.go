package requests

type HttpDraftRequest struct {
	Data struct {
		Title string `json:"title"`
	} `json:"data"`
}

type HttpDraftResponse struct {
	FileId string `json:"fileId"`
}
