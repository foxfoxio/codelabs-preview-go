package entities

import (
	"github.com/googlecodelabs/tools/claat/types"
	"time"
)

type Meta struct {
	FileId       string      `json:"fileId"`
	Revision     int         `json:"revision"`
	ExportedDate time.Time   `json:"exportedDate"`
	Meta         *types.Meta `json:"meta"`
}
