package gstorage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"io/ioutil"
)

type Client interface {
	Read(ctx context.Context, object string) (*bytes.Buffer, error)
	Write(ctx context.Context, object string, content io.Reader) (int64, error)
}

func NewClient(bucketName string) Client {
	return &client{bucketName: bucketName}
}

type client struct {
	bucketName string
}

func (c *client) new(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx)
}

func (c *client) Read(ctx context.Context, object string) (*bytes.Buffer, error) {
	client, err := c.new(ctx)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	reader, err := client.Bucket(c.bucketName).Object(object).NewReader(ctx)

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func (c *client) Write(ctx context.Context, object string, content io.Reader) (int64, error) {
	client, err := c.new(ctx)
	if err != nil {
		return 0, err
	}

	defer client.Close()

	writer := client.Bucket(c.bucketName).Object(object).NewWriter(ctx)
	defer writer.Close()

	return io.Copy(writer, content)
}

func IsNotExistError(err error) bool {
	return err != nil && storage.ErrObjectNotExist.Error() == err.Error()
}
