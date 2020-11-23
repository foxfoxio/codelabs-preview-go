package usecases

import (
	"bytes"
	"context"
	"errors"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdoc"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdrive"
	"github.com/foxfoxio/codelabs-preview-go/internal/stopwatch"
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
}

func NewViewer(driveClient gdrive.Client, gDocClient gdoc.Client, templateFileId string, driveRootId string, adminEmail string) Viewer {
	return &viewerUsecase{
		driveClient:    driveClient,
		gDocClient:     gDocClient,
		templateFileId: templateFileId,
		driveRootId:    driveRootId,
		adminEmail:     adminEmail,
	}
}

type viewerUsecase struct {
	driveClient    gdrive.Client
	gDocClient     gdoc.Client
	templateFileId string
	driveRootId    string
	adminEmail     string
}

func (uc *viewerUsecase) Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error) {
	log := cp.Log(ctx, "ViewerUsecase.Parse").WithField("file_id", request.FileId)
	defer stopwatch.StartWithLogger(log).Stop()

	s, err := uc.driveClient.ExportFile(ctx, request.FileId, "text/html")

	if err != nil {
		log.WithError(err).Error("google drive, get file failed")
		return nil, err
	}

	fetcher := fetch.NewGoogleDocMemoryFetcher(map[string]bool{}, parser.Blackfriday)
	codelabs, err := fetcher.SlurpCodelab(s.Reader)

	if err != nil {
		return nil, errors.New("bad bad: " + err.Error())
	}

	var buffer bytes.Buffer
	response := ""

	err = renderOutput(&buffer, codelabs.Codelab)

	if err != nil {
		response = err.Error()
	} else {
		response = string(buffer.Bytes())
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
		GlobalGA: "ga-001002003",
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
