package gdrive

import (
	"context"
	"fmt"
	"testing"
)

func TestGDrive(t *testing.T) {
	svc := getClient()
	ctx := context.Background()

	templateId := "1oZh5YrbA54pX9WfolES9MD5NvPdR_haEVeI3D56rHzM"

	folderId := "1uH1lq__vo-PTusArFsOduKfHk6ZhW1gX"
	fileName := "Introduction to everything else"
	//mimeType := "application/vnd.google-apps.document"

	//f := strings.NewReader("CONFIRMED.")
	file, err := copyFile(ctx, svc, templateId, fileName, folderId)

	//file, err := createFile(ctx, svc, fileName, mimeType, f, folderId)

	if err != nil {
		panic(err)
	}

	userEmailAddress := "easy.adix@gmail.com"

	perm, err := grantWritePermission(ctx, svc, file.Id, userEmailAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println(perm)
}
