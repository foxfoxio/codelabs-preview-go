package gdrive

import (
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"io"
)

type DriveFile struct {
	Id string
}

type DriveFileReader struct {
	Reader io.ReadCloser
}

type DrivePermission struct {
	Id string
}

type Client interface {
	CreateDir(ctx context.Context, name string, parentId string) (*DriveFile, error)
	CreateFile(ctx context.Context, name string, mimeType string, content io.Reader, parentId string) (*DriveFile, error)
	CopyFile(ctx context.Context, sourceFileId string, destinationName *string, parentId string) (*DriveFile, error)
	GrantWritePermission(ctx context.Context, fileId string, userEmail string) (*DrivePermission, error)
	GrantOwnerPermission(ctx context.Context, fileId string, userEmail string) (*DrivePermission, error)
	GetFile(ctx context.Context, fileId string) (*DriveFileReader, error)
	ExportFile(ctx context.Context, fileId string, mimeType string) (*DriveFileReader, error)
}

func NewClient() Client {
	return &client{
		service: getClient(),
	}
}

type client struct {
	service *drive.Service
}

func (c *client) CreateDir(ctx context.Context, name string, parentId string) (*DriveFile, error) {
	f, err := createDir(ctx, c.service, name, parentId)

	if err != nil {
		return nil, err
	}

	return &DriveFile{Id: f.Id}, nil
}

func (c *client) CreateFile(ctx context.Context, name string, mimeType string, content io.Reader, parentId string) (*DriveFile, error) {
	f, err := createFile(ctx, c.service, name, mimeType, content, parentId)

	if err != nil {
		return nil, err
	}

	return &DriveFile{Id: f.Id}, nil
}

func (c *client) CopyFile(ctx context.Context, sourceFileId string, destinationName *string, parentId string) (*DriveFile, error) {
	f, err := copyFile(ctx, c.service, sourceFileId, destinationName, parentId)

	if err != nil {
		return nil, err
	}

	return &DriveFile{Id: f.Id}, nil
}

func (c *client) GrantWritePermission(ctx context.Context, fileId string, userEmail string) (*DrivePermission, error) {
	f, err := grantWritePermission(ctx, c.service, fileId, userEmail)

	if err != nil {
		return nil, err
	}

	return &DrivePermission{Id: f.Id}, nil
}

func (c *client) GrantOwnerPermission(ctx context.Context, fileId string, userEmail string) (*DrivePermission, error) {
	f, err := grantOwnerPermission(ctx, c.service, fileId, userEmail)

	if err != nil {
		return nil, err
	}

	return &DrivePermission{Id: f.Id}, nil
}

func (c *client) GetFile(ctx context.Context, fileId string) (*DriveFileReader, error) {
	reader, err := getFile(ctx, c.service, fileId)

	if err != nil {
		return nil, err
	}

	return &DriveFileReader{
		Reader: reader,
	}, nil
}

func (c *client) ExportFile(ctx context.Context, fileId string, mimeType string) (*DriveFileReader, error) {
	reader, err := exportFile(ctx, c.service, fileId, mimeType)

	if err != nil {
		return nil, err
	}

	return &DriveFileReader{
		Reader: reader,
	}, nil
}
