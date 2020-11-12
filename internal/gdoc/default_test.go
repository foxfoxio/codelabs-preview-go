package gdoc

import (
	"context"
	"fmt"
	"testing"
)

func TestGDoc(t *testing.T) {
	t.Skip("ignore real test")

	svc := getService()

	ctx := context.Background()

	req := svc.Documents.Get("1KdO9GjdiN8aFLdNJpHgJDGp2E4fxtErR5LiLD6GtMTg")
	outDoc, err := req.Context(ctx).Do()

	if err != nil {
		panic(err)
	}

	fmt.Println(outDoc)
}
