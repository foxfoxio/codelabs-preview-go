package entities

import (
	"github.com/googlecodelabs/tools/claat/types"
	"time"
)

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
