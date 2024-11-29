package pkg

import (
	"bytes"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Merge(seekers []io.ReadSeeker, Log func(a ...interface{})) []byte {
	wr := new(bytes.Buffer)
	api.MergeRaw(seekers, wr, false, nil)
	Log("Write file...")
	return wr.Bytes()
}
