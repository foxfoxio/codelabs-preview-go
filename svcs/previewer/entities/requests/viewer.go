package requests

type ViewerParseRequest struct {
	FileId string
}

type ViewerParseResponse struct {
	Response string
}

type ViewerDraftRequest struct {
	Title string
}

type ViewerDraftResponse struct {
	FileId string
}
