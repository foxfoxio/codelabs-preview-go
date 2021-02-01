package gdrive

import (
	"context"
	"fmt"
	"testing"
)

func TestGDrive(t *testing.T) {
	t.Skip("ignore real test")
	svc := getClient()
	ctx := context.Background()

	templateId := "1oZh5YrbA54pX9WfolES9MD5NvPdR_haEVeI3D56rHzM"

	folderId := "1uH1lq__vo-PTusArFsOduKfHk6ZhW1gX"
	//fileName := "Introduction to everything else"
	//mimeType := "application/vnd.google-apps.document"

	////f := strings.NewReader("CONFIRMED.")
	file, err := copyFile(ctx, svc, templateId, nil, folderId)
	//
	////file, err := createFile(ctx, svc, fileName, mimeType, f, folderId)
	//
	if err != nil {
		panic(err)
	}
	fmt.Println(file)
	//
	//userEmailAddress := "easy.adix@gmail.com"
	//
	//perm, err := grantWritePermission(ctx, svc, file.Id, userEmailAddress)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(perm)
	//
	//r, err := exportFile(ctx, svc, templateId, "text/html")
	////r, err := getFile(ctx, svc, templateId)
	//if err != nil {
	//	panic(err)
	//}
	//result, _ := ioutil.ReadAll(r)
	//fmt.Println(string(result))
}
