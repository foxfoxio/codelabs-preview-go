package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/codelabs"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdoc"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdrive"
	"github.com/foxfoxio/codelabs-preview-go/internal/gstorage"
	"github.com/foxfoxio/codelabs-preview-go/internal/stopwatch"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	"path/filepath"
	"strconv"
	"strings"
)

type Viewer interface {
	Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error)
	Draft(ctx context.Context, request *requests.ViewerDraftRequest) (*requests.ViewerDraftResponse, error)
	Publish(ctx context.Context, request *requests.ViewerPublishRequest) (*requests.ViewerPublishResponse, error)
	View(ctx context.Context, request *requests.ViewerViewRequest) (*requests.ViewerViewResponse, error)
	Meta(ctx context.Context, request *requests.ViewerMetaRequest) (*requests.ViewerMetaResponse, error)
	Media(ctx context.Context, request *requests.ViewerMediaRequest) (*requests.ViewerMediaResponse, error)
}

func NewViewer(driveClient gdrive.Client, gDocClient gdoc.Client, gStorageClient gstorage.Client, templateFileId string, driveRootId string, adminEmail string, storagePath string) Viewer {
	return &viewerUsecase{
		driveClient:    driveClient,
		gDocClient:     gDocClient,
		gStorageClient: gStorageClient,
		templateFileId: templateFileId,
		driveRootId:    driveRootId,
		adminEmail:     adminEmail,
		storagePath:    storagePath,
	}
}

type viewerUsecase struct {
	driveClient    gdrive.Client
	gDocClient     gdoc.Client
	gStorageClient gstorage.Client
	templateFileId string
	driveRootId    string
	adminEmail     string
	storagePath    string
}

func (uc *viewerUsecase) parseCodeLabs(ctx context.Context, fileId string) (*codelabs.Result, error) {
	log := cp.Log(ctx, "ViewerUsecase.parseCodeLabs").WithField("fileId", fileId)
	s, err := uc.driveClient.ExportFile(ctx, fileId, "text/html")

	if err != nil {
		log.WithError(err).Error("google drive, get file failed")
		return nil, err
	}

	return codelabs.ParseCodeLab(fileId, s.Reader)
}

func (uc *viewerUsecase) Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Parse").WithField("fileId", request.FileId)
	defer stopwatch.StartWithLogger(log).Stop()

	result, err := uc.parseCodeLabs(ctx, request.FileId)

	response := ""
	if err != nil {
		response = err.Error()
	} else {
		response = result.HtmlContentBase64()
	}

	return &requests.ViewerParseResponse{
		Response: response,
	}, nil
}

func (uc *viewerUsecase) Draft(ctx context.Context, request *requests.ViewerDraftRequest) (*requests.ViewerDraftResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Draft").WithField("title", request.Title)
	defer stopwatch.StartWithLogger(log).Stop()

	session := getSession(ctx)

	if session == nil {
		log.Errorf("get user session failed")
		return nil, errors.New("unauthorized")
	}

	log.WithField("email", session.Email).
		WithField("user_id", session.UserId).
		Info("session found")

	if !request.Valid() {
		log.Errorf("invalid request")
		return nil, errors.New("bad request")
	}

	// create new document from template
	f, err := uc.driveClient.CopyFile(ctx, uc.templateFileId, request.Title(), uc.driveRootId)
	if err != nil {
		log.WithError(err).Error("google drive, copy file failed")
		return nil, err
	}
	log.WithField("file_id", f.Id).Info("file copied")

	// override template
	if len(request.MetaData) > 0 {
		doc, err := uc.gDocClient.ReplaceTexts(ctx, f.Id, request.ReplaceTextParams())
		if err != nil {
			log.WithError(err).Error("google doc, replace template")
			return nil, err
		}
		log.WithField("doc_id", doc.Id).Info("templated created")
	} else {
		log.Info("no metadata provided, skip replacing template")
	}

	// share document
	s, err := uc.driveClient.GrantWritePermission(ctx, f.Id, session.Email)

	if err != nil {
		log.WithError(err).Error("google drive, share file failed")
		return nil, err
	}

	log.WithField("permission_id", s.Id).Info("file shared")

	// set document owner
	if uc.adminEmail != "" {
		s, err := uc.driveClient.GrantOwnerPermission(ctx, f.Id, uc.adminEmail)

		if err != nil {
			log.WithError(err).Error("google drive, set file owner failed")
		}

		log.WithField("permission_id", s.Id).Info("owner set")
	}

	// return to user
	return &requests.ViewerDraftResponse{FileId: f.Id}, nil
}

func (uc *viewerUsecase) Publish(ctx context.Context, request *requests.ViewerPublishRequest) (*requests.ViewerPublishResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Publish").WithField("fileId", request.FileId)
	defer stopwatch.StartWithLogger(log).Stop()

	// parse codelabs
	result, err := uc.parseCodeLabs(ctx, request.FileId)

	if err != nil {
		return nil, err
	}

	// get latest revisions
	latestMetaPath := fmt.Sprintf("%s/%s/latest/meta.json", uc.storagePath, request.FileId)
	latestIndexPath := fmt.Sprintf("%s/%s/latest/index.html", uc.storagePath, request.FileId)

	latestMetaBytes, err := uc.gStorageClient.Read(ctx, latestMetaPath)
	if err != nil {
		if !gstorage.IsNotExistError(err) {
			log.WithError(err).WithField("path", latestMetaPath).Error("read latest meta file failed")
			return nil, err
		}
	} else {
		latestMeta := &codelabs.Meta{}
		mm := latestMetaBytes.Bytes()
		if e := json.Unmarshal(mm, latestMeta); e != nil {
			log.WithError(err).WithField("data", string(mm)).Error("unmarshal latest meta file failed")
		} else {
			log.WithField("revision", latestMeta.Revision).Error("latest revision")
			result.Meta.Revision = latestMeta.Revision + 1
		}
	}

	revMetaPath := fmt.Sprintf("%s/%s/%d/meta.json", uc.storagePath, request.FileId, result.Meta.Revision)
	revIndexPath := fmt.Sprintf("%s/%s/%d/index.html", uc.storagePath, request.FileId, result.Meta.Revision)

	// save new revision to bucket
	size, err := uc.gStorageClient.Write(ctx, latestIndexPath, bytes.NewBufferString(result.HtmlContent))
	if err != nil {
		log.WithError(err).WithField("path", latestIndexPath).Error("write index file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", latestIndexPath).Info("latest index file created")
	size, err = uc.gStorageClient.Write(ctx, revIndexPath, bytes.NewBufferString(result.HtmlContent))
	if err != nil {
		log.WithError(err).WithField("path", revIndexPath).Error("write revision index file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", revIndexPath).Info("revision index file created")
	size, err = uc.gStorageClient.Write(ctx, latestMetaPath, bytes.NewBufferString(result.Meta.JsonString()))
	if err != nil {
		log.WithError(err).WithField("path", latestMetaPath).Error("write latest meta file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", latestMetaPath).Info("latest meta file created")
	size, err = uc.gStorageClient.Write(ctx, revMetaPath, bytes.NewBufferString(result.Meta.JsonString()))
	if err != nil {
		log.WithError(err).WithField("path", revMetaPath).Error("write revision meta file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", revMetaPath).Info("revision meta file created")

	// save images
	for _, img := range result.Images {
		imgLatestPath := filepath.Join(uc.storagePath, request.FileId, "latest", img.Path())
		imgRevPath := filepath.Join(uc.storagePath, request.FileId, strconv.Itoa(result.Meta.Revision), img.Path())

		size, err = uc.gStorageClient.Write(ctx, imgLatestPath, bytes.NewBuffer(img.Content))
		if err != nil {
			log.WithError(err).WithField("path", imgLatestPath).Error("write latest image file failed")
			return nil, err
		}
		log.WithField("size", size).WithField("path", imgLatestPath).Info("latest image file created")

		size, err = uc.gStorageClient.Write(ctx, imgRevPath, bytes.NewBuffer(img.Content))
		if err != nil {
			log.WithError(err).WithField("path", imgRevPath).Error("write revision image file failed")
			return nil, err
		}
		log.WithField("size", size).WithField("path", imgRevPath).Info("revision image file created")
	}

	return &requests.ViewerPublishResponse{
		Revision: result.Meta.Revision,
		Meta:     result.Meta,
	}, nil

}

func (uc *viewerUsecase) View(ctx context.Context, request *requests.ViewerViewRequest) (*requests.ViewerViewResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.View").WithField("fileId", request.FileId).WithField("revision", request.Revision)
	defer stopwatch.StartWithLogger(log).Stop()

	path := ""
	if request.Revision <= 0 {
		path = fmt.Sprintf("%s/%s/latest/index.html", uc.storagePath, request.FileId)
	} else {
		path = fmt.Sprintf("%s/%s/%d/index.html", uc.storagePath, request.FileId, request.Revision)
	}

	indexBytes, err := uc.gStorageClient.Read(ctx, path)

	if err != nil {
		log.WithError(err).WithField("path", path).Error("read index file failed")
		if gstorage.IsNotExistError(err) {

			return nil, errors.New("not found")
		} else {
			return nil, err
		}
	}

	return &requests.ViewerViewResponse{Response: indexBytes.String()}, nil
}

func (uc *viewerUsecase) Media(ctx context.Context, request *requests.ViewerMediaRequest) (*requests.ViewerMediaResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Media").
		WithField("fileId", request.FileId).
		WithField("revision", request.Revision).
		WithField("filename", request.Filename)
	defer stopwatch.StartWithLogger(log).Stop()

	rev := "latest"
	if request.Revision > 0 {
		rev = strconv.Itoa(request.Revision)
	}
	path := fmt.Sprintf("%s/%s/%s/img/%s", uc.storagePath, request.FileId, rev, request.Filename)

	imgBytes, err := uc.gStorageClient.Read(ctx, path)

	if err != nil {
		log.WithError(err).WithField("path", path).Error("read index file failed")
		if gstorage.IsNotExistError(err) {

			return nil, errors.New("not found")
		} else {
			return nil, err
		}
	}

	ext := strings.TrimLeft(filepath.Ext(request.Filename), ".")
	return &requests.ViewerMediaResponse{
		ContentType: fmt.Sprintf("image/%s", ext),
		Content:     imgBytes.Bytes(),
	}, nil
}

func (uc *viewerUsecase) Meta(ctx context.Context, request *requests.ViewerMetaRequest) (*requests.ViewerMetaResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Meta").WithField("fileId", request.FileId).WithField("revision", request.Revision)
	defer stopwatch.StartWithLogger(log).Stop()

	path := ""
	if request.Revision <= 0 {
		path = fmt.Sprintf("%s/%s/latest/meta.json", uc.storagePath, request.FileId)
	} else {
		path = fmt.Sprintf("%s/%s/%d/meta.json", uc.storagePath, request.FileId, request.Revision)
	}

	metaBytes, err := uc.gStorageClient.Read(ctx, path)

	if err != nil {
		log.WithError(err).WithField("path", path).Error("read meta file failed")
		if gstorage.IsNotExistError(err) {

			return nil, errors.New("not found")
		} else {
			return nil, err
		}
	}

	meta := &codelabs.Meta{}
	mm := metaBytes.Bytes()
	if e := json.Unmarshal(mm, meta); e != nil {
		log.WithError(err).WithField("data", string(mm)).Error("unmarshal meta file failed")
		return nil, err
	}

	return &requests.ViewerMetaResponse{Meta: meta}, nil
}
