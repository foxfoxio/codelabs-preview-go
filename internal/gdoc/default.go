package gdoc

import (
	"context"
	"google.golang.org/api/docs/v1"
)

func getService() *docs.Service {
	ctx := context.Background()
	svc, err := docs.NewService(ctx)

	if err != nil {
		panic(err)
	}

	return svc
}
