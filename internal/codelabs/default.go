package codelabs

import (
	"encoding/base64"
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
	"io"
	"path/filepath"
	"strings"
	"time"
)

const ImageDir = "img"

type MetaEx struct {
	*types.Meta
	TotalChapters int `json:"totalChapters"`
}

type Meta struct {
	FileId       string    `json:"fileId"`
	Revision     int       `json:"revision"`
	ExportedDate time.Time `json:"exportedDate"`
	Meta         *MetaEx   `json:"meta"`
}

func (m *Meta) JsonString() string {
	return utils.StringifyIndent(m)
}

type ImageBuffer struct {
	Url       string
	Filename  string
	Extension string
	Content   []byte
}

func (i *ImageBuffer) Path() string {
	if i == nil {
		return ""
	}

	return filepath.Join(ImageDir, i.Filename)
}

type ImageBuffers []*ImageBuffer

type Result struct {
	HtmlContent string
	Images      ImageBuffers
	Meta        *Meta
}

func (r *Result) HtmlContentBase64() string {
	if len(r.Images) <= 0 {
		return r.HtmlContent
	}
	content := r.HtmlContent
	for _, i := range r.Images {
		path := i.Path()
		img64 := base64.StdEncoding.EncodeToString(i.Content)
		prefix := fmt.Sprintf("data:image/%s;base64, ", strings.TrimLeft(".", i.Extension))
		base64Img := fmt.Sprintf(`"%s%s"`, prefix, img64)

		content = strings.ReplaceAll(content, fmt.Sprintf(`"%s"`, path), base64Img)
	}

	return content
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
