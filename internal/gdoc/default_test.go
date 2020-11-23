package gdoc

import (
	"context"
	"fmt"
	"testing"
)

func TestGDoc(t *testing.T) {
	t.Skip("ignore real test")

	svc := getClient()

	ctx := context.Background()

	params := map[string]string{
		"status": "draft",
	}

	res, err := replaceTexts(ctx, svc, "1I05wp5zv0fao-rpj7WoFmCl2Pd1xEINbKdnX0s9-UZE", params)

	if err != nil {
		panic(err)
	}

	fmt.Println(res.DocumentId)
}
