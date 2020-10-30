package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/googlecodelabs/tools/claat/fetch"
	"github.com/googlecodelabs/tools/claat/parser"
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Viewer interface {
	Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error)
}

func NewViewer(config *oauth2.Config) Viewer {
	return &viewerUsecase{config}
}

type viewerUsecase struct {
	config *oauth2.Config
}

func (uc *viewerUsecase) Parse(ctx context.Context, request *requests.ViewerParseRequest) (*requests.ViewerParseResponse, error) {
	session := getSession(ctx)

	if session == nil {
		return nil, errors.New("unauthorized")
	}

	token := session.Oauth2Token()

	if token == nil {
		return nil, errors.New("unauthorized")
	}

	//client := uc.config.Client(ctx, session.Oauth2Token())
	//userInfo, err := FetchUserInfo(client)
	//
	//response := ""
	//if err != nil {
	//	response = err.Error()
	//} else {
	//	response = utils.Stringify(userInfo)
	//}

	fetcher, err := fetch.NewFetcher(token.AccessToken, map[string]bool{}, nil, parser.Blackfriday)

	if err != nil {
		return nil, errors.New("bad bad: " + err.Error())
	}

	codelabs, err := fetcher.SlurpCodelab(request.FileId)

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

	fmt.Println("codelabs", codelabs)

	return &requests.ViewerParseResponse{
		Response: response,
	}, nil
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func FetchUserInfo(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
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
