package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdoc"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdrive"
	"github.com/foxfoxio/codelabs-preview-go/internal/gstorage"
	"github.com/foxfoxio/codelabs-preview-go/internal/stopwatch"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/googlecodelabs/tools/claat/fetch"
	"github.com/googlecodelabs/tools/claat/parser"
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
	"io"
	"time"
)

type Viewer interface {
	Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error)
	Draft(ctx context.Context, request *requests.ViewerDraftRequest) (*requests.ViewerDraftResponse, error)
	Publish(ctx context.Context, request *requests.ViewerPublishRequest) (*requests.ViewerPublishResponse, error)
	View(ctx context.Context, request *requests.ViewerViewRequest) (*requests.ViewerViewResponse, error)
	Meta(ctx context.Context, request *requests.ViewerMetaRequest) (*requests.ViewerMetaResponse, error)
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

func (uc *viewerUsecase) parseCodeLabs(ctx context.Context, fileId string) ([]byte, *entities.Meta, error) {
	log := cp.Log(ctx, "ViewerUsecase.parseCodeLabs").WithField("fileId", fileId)
	s, err := uc.driveClient.ExportFile(ctx, fileId, "text/html")

	if err != nil {
		log.WithError(err).Error("google drive, get file failed")
		return nil, nil, err
	}

	fetcher := fetch.NewGoogleDocMemoryFetcher(map[string]bool{}, parser.Blackfriday)
	codelabs, err := fetcher.SlurpCodelab(s.Reader)

	if err != nil {
		return nil, nil, errors.New("bad bad: " + err.Error())
	}

	var buffer bytes.Buffer
	err = renderOutput(&buffer, codelabs.Codelab)

	meta := &entities.Meta{
		FileId:       fileId,
		Revision:     1, // default revision
		ExportedDate: time.Now(),
		Meta:         &codelabs.Meta,
	}

	return buffer.Bytes(), meta, err
}

func (uc *viewerUsecase) Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Parse").WithField("fileId", request.FileId)
	defer stopwatch.StartWithLogger(log).Stop()

	res, _, err := uc.parseCodeLabs(ctx, request.FileId)

	response := ""
	if err != nil {
		response = err.Error()
	} else {
		response = string(res)
	}

	return &requests.ViewerParseResponse{
		Response: response,
	}, nil
}

func renderOutput(w io.Writer, codelabs *types.Codelab) error {
	data := &struct {
		render.Context
		Current *types.Step
		StepNum int
		Prev    bool
		Next    bool
	}{Context: render.Context{
		Env:      "web",
		Prefix:   "https://storage.googleapis.com",
		Format:   "html",
		GlobalGA: codelabs.GA,
		Updated:  time.Now().Format(time.RFC3339),
		Meta:     &codelabs.Meta,
		Steps:    codelabs.Steps,
		Extra:    map[string]string{},
	}}

	return render.Execute(w, "html", data)
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
	resBytes, meta, err := uc.parseCodeLabs(ctx, request.FileId)

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
		latestMeta := &entities.Meta{}
		mm := latestMetaBytes.Bytes()
		if e := json.Unmarshal(latestMetaBytes.Bytes(), latestMeta); e != nil {
			log.WithError(err).WithField("data", string(mm)).Error("unmarshal latest meta file failed")
		} else {
			log.WithField("revision", latestMeta.Revision).Error("latest revision")
			meta.Revision = latestMeta.Revision + 1
		}
	}

	revMetaPath := fmt.Sprintf("%s/%s/%d/meta.json", uc.storagePath, request.FileId, meta.Revision)
	revIndexPath := fmt.Sprintf("%s/%s/%d/index.html", uc.storagePath, request.FileId, meta.Revision)

	// save new revision to bucket
	size, err := uc.gStorageClient.Write(ctx, latestIndexPath, bytes.NewBuffer(resBytes))
	if err != nil {
		log.WithError(err).WithField("path", latestIndexPath).Error("write index file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", latestIndexPath).Info("latest index file created")
	size, err = uc.gStorageClient.Write(ctx, revIndexPath, bytes.NewBuffer(resBytes))
	if err != nil {
		log.WithError(err).WithField("path", revIndexPath).Error("write revision index file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", revIndexPath).Info("revision index file created")
	size, err = uc.gStorageClient.Write(ctx, latestMetaPath, bytes.NewBufferString(utils.StringifyIndent(meta)))
	if err != nil {
		log.WithError(err).WithField("path", latestMetaPath).Error("write latest meta file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", latestMetaPath).Info("latest meta file created")
	size, err = uc.gStorageClient.Write(ctx, revMetaPath, bytes.NewBufferString(utils.StringifyIndent(meta)))
	if err != nil {
		log.WithError(err).WithField("path", revMetaPath).Error("write revision meta file failed")
		return nil, err
	}
	log.WithField("size", size).WithField("path", revMetaPath).Info("revision meta file created")

	return &requests.ViewerPublishResponse{
		Revision: meta.Revision,
	}, nil

}

func (uc *viewerUsecase) View(ctx context.Context, request *requests.ViewerViewRequest) (*requests.ViewerViewResponse, error) {
	panic("implement me")
}

func (uc *viewerUsecase) Meta(ctx context.Context, request *requests.ViewerMetaRequest) (*requests.ViewerMetaResponse, error) {
	panic("implement me")
}
