package usecases

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
)

type Viewer interface {
	Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error)
}

func NewViewer() Viewer {
	return &viewerUsecase{}
}

type viewerUsecase struct {
}

func (uc *viewerUsecase) Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error) {
	return &requests.ViewerParseResponse{}, nil
}
