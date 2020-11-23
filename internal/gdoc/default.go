package gdoc

import (
	"context"
	"google.golang.org/api/docs/v1"
)

type DocFile struct {
	Id string
}

type Client interface {
	ReplaceTexts(ctx context.Context, docId string, replaceParams map[string]string) (*DocFile, error)
}

func NewClient() Client {
	return &client{
		service: getClient(),
	}
}

type client struct {
	service *docs.Service
}

func (c *client) ReplaceTexts(ctx context.Context, docId string, replaceParams map[string]string) (*DocFile, error) {
	d, err := replaceTexts(ctx, c.service, docId, replaceParams)

	if err != nil {
		return nil, err
	}

	return &DocFile{Id: d.DocumentId}, nil
}
