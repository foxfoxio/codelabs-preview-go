package gdoc

import (
	"context"
	"google.golang.org/api/docs/v1"
	"log"
)

func getClient() *docs.Service {
	ctx := context.Background()
	svc, err := docs.NewService(ctx)

	if err != nil {
		panic(err)
	}

	return svc
}

func replaceTexts(ctx context.Context, service *docs.Service, docId string, replaceParams map[string]string) (*docs.BatchUpdateDocumentResponse, error) {

	requests := make([]*docs.Request, 0)
	for k, v := range replaceParams {
		requests = append(requests, &docs.Request{
			ReplaceAllText: &docs.ReplaceAllTextRequest{
				ContainsText: &docs.SubstringMatchCriteria{
					MatchCase: true,
					Text:      k,
				},
				ReplaceText: v,
			},
		})
	}

	update := &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}

	resp, err := service.Documents.BatchUpdate(docId, update).Context(ctx).Do()

	if err != nil {
		log.Println("could not update file: " + err.Error())
		return nil, err
	}

	return resp, nil
}
