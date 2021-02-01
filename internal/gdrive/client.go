package gdrive

import (
	"context"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
)

const (
	GoogleDocumentMimeType = "application/vnd.google-apps.document"
)

func createDir(ctx context.Context, service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Context(ctx).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFile(ctx context.Context, service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Context(ctx).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func copyFile(ctx context.Context, service *drive.Service, sourceFileId string, name *string, parentId string) (*drive.File, error) {
	f := &drive.File{
		Parents: []string{parentId},
	}

	if name != nil {
		f.Name = *name
	}

	file, err := service.Files.Copy(sourceFileId, f).Context(ctx).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func grantWritePermission(ctx context.Context, service *drive.Service, fileId string, userEmail string) (*drive.Permission, error) {
	perm := &drive.Permission{
		EmailAddress: userEmail,
		Role:         "writer",
		Type:         "user",
	}

	perm, err := service.Permissions.Create(fileId, perm).Context(ctx).Do()

	if err != nil {
		log.Println("Could not share file: " + err.Error())
		return nil, err
	}

	return perm, nil
}

func grantOwnerPermission(ctx context.Context, service *drive.Service, fileId string, userEmail string) (*drive.Permission, error) {
	perm := &drive.Permission{
		EmailAddress: userEmail,
		Role:         "owner",
		Type:         "user",
	}

	perm, err := service.Permissions.Create(fileId, perm).Context(ctx).TransferOwnership(true).Do()

	if err != nil {
		log.Println("Could not share file: " + err.Error())
		return nil, err
	}

	return perm, nil
}

func getClient() *drive.Service {
	ctx := context.Background()
	service, err := drive.NewService(ctx)

	if err != nil {
		panic(err)
	}

	return service
}

func getFile(ctx context.Context, service *drive.Service, fileId string) (io.ReadCloser, error) {
	response, err := service.Files.Get(fileId).Context(ctx).Download()

	if err != nil {
		log.Println("Could not get file: " + err.Error())
		return nil, err
	}

	return response.Body, nil
}

func getInfo(ctx context.Context, service *drive.Service, fileId string) (*drive.File, error) {
	response, err := service.Files.Get(fileId).Context(ctx).Do()

	if err != nil {
		log.Println("Could not get file info: " + err.Error())
		return nil, err
	}

	return response, nil
}

func exportFile(ctx context.Context, service *drive.Service, fileId string, mimeType string) (io.ReadCloser, error) {
	response, err := service.Files.Export(fileId, mimeType).Context(ctx).Download()

	if err != nil {
		log.Println("Could not export file: " + err.Error())
		return nil, err
	}

	return response.Body, nil
}
