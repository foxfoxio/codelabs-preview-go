package codelabs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/googlecodelabs/tools/claat/fetch"
	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/types"
	"golang.org/x/xerrors"
	"hash/crc64"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

func ParseCodeLabWithExtractImage(fileId string, reader io.ReadCloser) (*Result, error) {
	fetcher := fetch.NewGoogleDocMemoryFetcher(map[string]bool{}, parser.Blackfriday)
	codelabs, err := fetcher.SlurpCodelab(reader)

	if err != nil {
		return nil, xerrors.New("parse codelabs failed: " + err.Error())
	}

	images, err := slurpImages(ImageDir, codelabs.Steps)

	if err != nil {
		return nil, xerrors.New("download images failed: " + err.Error())
	}

	var buffer bytes.Buffer
	err = renderOutput(&buffer, codelabs.Codelab)

	if err != nil {
		return nil, xerrors.New("render output failed: " + err.Error())
	}

	meta := &Meta{
		FileId:       fileId,
		Revision:     1, // default revision
		ExportedDate: time.Now(),
		Meta: &MetaEx{
			Meta:          &codelabs.Meta,
			TotalChapters: len(codelabs.Steps),
		},
	}

	return &Result{
		HtmlContent: buffer.String(),
		Images:      images,
		Meta:        meta,
	}, nil
}

func ParseCodeLab(fileId string, reader io.ReadCloser) (*Result, error) {
	fetcher := fetch.NewGoogleDocMemoryFetcher(map[string]bool{}, parser.Blackfriday)
	codelabs, err := fetcher.SlurpCodelab(reader)

	if err != nil {
		return nil, xerrors.New("parse codelabs failed: " + err.Error())
	}

	var buffer bytes.Buffer
	err = renderOutput(&buffer, codelabs.Codelab)

	if err != nil {
		return nil, xerrors.New("render output failed: " + err.Error())
	}

	meta := &Meta{
		FileId:       fileId,
		Revision:     1, // default revision
		ExportedDate: time.Now(),
		Meta: &MetaEx{
			Meta:          &codelabs.Meta,
			TotalChapters: len(codelabs.Steps),
		},
	}

	return &Result{
		HtmlContent: buffer.String(),
		Images:      nil,
		Meta:        meta,
	}, nil
}

func slurpImages(imgDir string, steps []*types.Step) (ImageBuffers, error) {
	type res struct {
		url   string
		Image *ImageBuffer
		err   error
	}

	ch := make(chan *res, 100)
	defer close(ch)

	var count int
	for _, st := range steps {
		nodes := types.ImageNodes(st.Content.Nodes)
		count += len(nodes)
		for _, n := range nodes {
			go func(n *types.ImageNode) {
				imgUrl := n.Src
				file, err := slurpBytes(imgUrl)
				if err == nil {
					n.Src = filepath.Join(imgDir, file.Filename)
				}
				ch <- &res{imgUrl, file, err}
			}(n)
		}
	}

	imgBuffers := make(ImageBuffers, 0)
	var errStr string
	for i := 0; i < count; i++ {
		r := <-ch
		if r.err != nil {
			errStr += fmt.Sprintf("%s => %v\n", r.url, r.err)
		}
		if r.Image != nil {
			imgBuffers = append(imgBuffers, r.Image)
		}
	}
	if len(errStr) > 0 {
		return nil, errors.New(errStr)
	}

	return imgBuffers, nil
}

// slurpBytes this method assume all images are from network
func slurpBytes(imgURL string) (*ImageBuffer, error) {
	// images can be local in Markdown cases or remote.
	// Only proceed a simple copy on local reference.
	var b []byte
	var ext string
	u, err := url.Parse(imgURL)
	if err != nil {
		return nil, err
	}

	b, err = slurpRemoteBytes(u.String(), 2)
	if string(b[6:10]) == "JFIF" {
		ext = ".jpeg"
	} else if string(b[0:3]) == "GIF" {
		ext = ".gif"
	} else {
		ext = ".png"
	}

	if err != nil {
		return nil, err
	}

	crcTable := crc64.MakeTable(crc64.ECMA)
	crc := crc64.Checksum(b, crcTable)
	file := fmt.Sprintf("%x%s", crc, ext)
	return &ImageBuffer{
		Url:       imgURL,
		Filename:  file,
		Extension: ext,
		Content:   b,
	}, nil
}

func slurpRemoteBytes(url string, n int) ([]byte, error) {
	res, err := retryGet(http.DefaultClient, url, n)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	return ioutil.ReadAll(res.Body)
}

// retryGet tries to GET specified url up to n times.
// Attempts are spaced out with exponential backoff.
// Default client will be used if not provided.
func retryGet(client *http.Client, url string, n int) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	for i := 0; i <= n; i++ {
		if i > 0 {
			t := time.Duration((math.Pow(2, float64(i)) + rand.Float64()) * float64(time.Second))
			time.Sleep(t)
		}
		res, err := client.Get(url)
		// return early with a good response
		// the rest is error handling
		if err == nil && res.StatusCode == http.StatusOK {
			return res, nil
		}

		// sometimes Drive API wouldn't even start a response,
		// we get net/http: TLS handshake timeout instead:
		// consider this a temporary failure and retry again
		if err != nil {
			continue
		}
		// otherwise, decode error response and check for "rate limit"
		defer func() { _ = res.Body.Close() }()
		var erres struct {
			Error struct {
				Errors []struct{ Reason string }
			}
		}
		b, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(b, &erres)
		var rateLimit bool
		for _, e := range erres.Error.Errors {
			if e.Reason == "rateLimitExceeded" || e.Reason == "userRateLimitExceeded" {
				rateLimit = true
				break
			}
		}
		// this is neither a rate limit error, nor a server error:
		// retrying is useless
		if !rateLimit && res.StatusCode < http.StatusInternalServerError {
			return nil, fmt.Errorf("fetch %s: %s; %s", url, res.Status, b)
		}
	}
	return nil, fmt.Errorf("%s: failed after %d retries", url, n)
}
