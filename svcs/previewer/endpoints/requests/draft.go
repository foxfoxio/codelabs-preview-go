package requests

type HttpDraftRequest struct {
	Data map[string]string `json:"data"`
}

type HttpDraftResponse struct {
	FileId string `json:"fileId"`
}

type HttpPublishResponse struct {
	Revision int                    `json:"revision"`
	Meta     map[string]interface{} `json:"meta"`
}

type HttpMetaResponse struct {
	Meta map[string]interface{} `json:"meta"`
}

type HttpCopyResponse struct {
	GoogleDocPath string `json:"googleDocPath"`
}
