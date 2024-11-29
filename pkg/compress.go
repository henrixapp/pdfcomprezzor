package pkg

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/nfnt/resize"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Compress(array []byte, Log func(a ...interface{}), imageSize int) ([]byte, error) {
	buffi := bytes.NewReader(array)
	ctx, err := pdfcpu.Read(buffi, model.NewDefaultConfiguration())
	if err != nil {
		Log("Error while loading file", err)
		return make([]byte, 0), err
	}
	//Log(ctx)
	err = ctx.EnsurePageCount()
	if err != nil {
		return nil, err
	}
	count := ctx.PageCount
	err = pdfcpu.OptimizeXRefTable(ctx)
	if err != nil {
		return make([]byte, 0), err
	}
	err = api.OptimizeContext(ctx)
	if err != nil {
		return make([]byte, 0), err
	}
	Log("Page count:", count)

	for i := 1; i <= count; i++ {
		Log("Processing page no:", i)
		imageObjNrs := pdfcpu.ImageObjNrs(ctx, i)
		Log("Images on page:", len(imageObjNrs))
		images, err := pdfcpu.ExtractPageImages(ctx, i, false)
		if err != nil {
			Log("e", err)
			return make([]byte, 0), err
		}
		for idx, i := range images {
			img, _, err := image.Decode(i)
			if err != nil {
				Log("e", err)
				return make([]byte, 0), err
			}
			if img.Bounds().Dx() > imageSize {
				Log(fmt.Sprint("Compress this image with id", idx, "....."))
				smaller := resize.Thumbnail(1240, 1740, img, resize.Lanczos2)
				Log(smaller.Bounds().Dx(), img.Bounds().Dx())
				if err != nil {
					Log("e", err)
					return make([]byte, 0), err
				}

				if err != nil {
					Log("e", err)
					return make([]byte, 0), err
				}
				buf := new(bytes.Buffer)
				err = png.Encode(buf, smaller)
				if err != nil {
					Log("e", err)
					return make([]byte, 0), err
				}
				sd2, _, _, err := model.CreateImageStreamDict(ctx.XRefTable, buf, false, false)
				if err != nil {
					Log("e", err)
					return make([]byte, 0), err
				}
				Log("idx:", idx)
				ctx.XRefTable.Table[idx].Object = *sd2
			}
		}
	}
	//err = pdfcpu.OptimizeXRefTable(ctx)
	if err != nil {
		Log("e", err)
		return make([]byte, 0), err
	}
	//api.OptimizeContext(ctx)
	ctx.EnsureVersionForWriting()
	Log("Write file...")
	wr := new(bytes.Buffer)
	err = api.WriteContext(ctx, wr)
	return wr.Bytes(), err
}
