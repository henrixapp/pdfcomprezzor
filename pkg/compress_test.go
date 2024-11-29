package pkg_test

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/henrixapp/pdfcomprezzor/pkg"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/stretchr/testify/assert"
)

//go:embed test_pdfs/empty.pdf
var emptyPDF []byte

//go:embed test_pdfs/big_image.pdf
var bigImagePDF []byte

func LogMock(s *[]string) func(a ...interface{}) {
	return func(a ...interface{}) {
		*s = append(*s, fmt.Sprint(a))
	}
}
func TestCompressEmpty(t *testing.T) {
	l := make([]string, 0)
	comp, err := pkg.Compress(emptyPDF, LogMock(&l), 10)
	assert.Equal(t, err == nil, true, fmt.Sprintln(l))
	assert.Equal(t, l[0], "[Page count: 1]")
	assert.Equal(t, l[1], "[Processing page no: 1]")
	assert.Equal(t, l[2], "[Images on page: 0]")
	err = api.Validate(bytes.NewReader(comp), nil)
	assert.Equal(t, err == nil, true)
}

func TestCompressImage(t *testing.T) {
	l := make([]string, 0)
	comp, err := pkg.Compress(bigImagePDF, LogMock(&l), 10)
	assert.Equal(t, err == nil, true, fmt.Sprintln(l))
	assert.Equal(t, l[0], "[Page count: 1]")
	assert.Equal(t, l[1], "[Processing page no: 1]")
	assert.Equal(t, l[2], "[Images on page: 1]")
	assert.Equal(t, l[3], "[Compress this image with id14.....]")
	err = api.Validate(bytes.NewReader(comp), nil)
	assert.Equal(t, err == nil, true)
}
