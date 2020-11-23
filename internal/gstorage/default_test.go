package gstorage

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Skip("skip real test")

	ctx := context.Background()
	c := NewClient("codelabs-preview")
	readRes, err := c.Read(ctx, "files-dev/x/meta.json")
	fmt.Println(err)
	fmt.Println(readRes)

	size, err := c.Write(ctx, "files-dev/x/meta.json", bytes.NewBufferString(`{"a":"b"}`))
	fmt.Println(err)
	fmt.Println(size)

	readRes, err = c.Read(ctx, "files-dev/x/meta.json")
	fmt.Println(err)
	fmt.Println(readRes)
}
