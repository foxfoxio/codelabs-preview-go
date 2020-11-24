package requests

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
)

type ViewerParseRequest struct {
	FileId string
}

type ViewerParseResponse struct {
	Response string
}

type ViewerPublishRequest struct {
	FileId string
}

type ViewerPublishResponse struct {
	Revision int
	Meta     *entities.Meta
}

type ViewerMetaRequest struct {
	FileId   string
	Revision int
}

type ViewerMetaResponse struct {
	Meta *entities.Meta
}

type ViewerViewRequest struct {
	FileId   string
	Revision int
}

type ViewerViewResponse struct {
	Response string
}

type ViewerDraftRequest struct {
	MetaData map[string]string
}

func (e *ViewerDraftRequest) Title() string {
	if t, ok := e.MetaData[ViewerDraftKeyTitle]; ok {
		return t
	}

	return ""
}

func (e *ViewerDraftRequest) Valid() bool {
	return e.Title() != ""
}

func (e *ViewerDraftRequest) ReplaceTextParams() map[string]string {
	params := make(map[string]string)

	for k, v := range e.MetaData {
		newK := fmt.Sprintf(`${%s}`, k)
		params[newK] = v
	}

	return params
}

const (
	ViewerDraftKeyTitle            string = "title"
	ViewerDraftKeySummary          string = "summary"
	ViewerDraftKeySlug             string = "slug"
	ViewerDraftKeyType             string = "type"
	ViewerDraftKeyTags             string = "tags"
	ViewerDraftKeyStatus           string = "status"
	ViewerDraftKeyFeedbackLink     string = "feedbackLink"
	ViewerDraftKeyAuthor           string = "author"
	ViewerDraftKeyAuthorLDAP       string = "authorLDAP"
	ViewerDraftKeyAnalyticsAccount string = "analyticsAccount"
)

type ViewerDraftResponse struct {
	FileId string
}
