package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"syscall/js"

	"github.com/nfnt/resize"
	api "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

var done = make(chan struct{})

//var body js.Value
//var doc js.Value
var LogCallback js.Value

func Log(a ...interface{}) {

	if !LogCallback.IsUndefined() {
		LogCallback.Invoke(js.Null(), fmt.Sprint(a))
	}
}

var onCompress = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	LogCallback = args[len(args)-1:][0]
	Log("pdfcomprezzor/compress")
	Log("File size in Bytes:", args[0].Get("length").Int())
	array := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(array, args[0])
	buffi := bytes.NewReader(array)
	ctx, err := pdfcpu.Read(buffi, model.NewDefaultConfiguration())
	if err != nil {
		Log("Error while loading file", err)
	}
	//Log(ctx)
	ctx.EnsurePageCount()
	count := ctx.PageCount
	pdfcpu.OptimizeXRefTable(ctx)
	api.OptimizeContext(ctx)
	Log("Page count:", count)

	for i := 1; i <= count; i++ {
		Log("Processing page no:", i)
		imageObjNrs := pdfcpu.ImageObjNrs(ctx, i)
		Log("Images on page:", len(imageObjNrs))
		images, err := pdfcpu.ExtractPageImages(ctx, i, false)
		if err != nil {
			Log("e", err)
		}
		objs := pdfcpu.ImageObjNrs(ctx, i)
		for idx, i := range images {
			img, _, err := image.Decode(i)
			if err != nil {
				Log("e", err)
			}
			if img.Bounds().Dx() > 1000 {
				Log(fmt.Sprint("Compress this image", objs[idx], "....."))
				smaller := resize.Thumbnail(1240, 1740, img, resize.Lanczos2)
				Log(smaller.Bounds().Dx(), img.Bounds().Dx())
				//obj, _ := ctx.Find(objs[idx])
				//obj, _ := ctx.Optimize.ImageObjects[objs[idx]]
				if err != nil {
					Log("e", err)
				}

				if err != nil {
					Log("e", err)
				}
				buf := new(bytes.Buffer)
				png.Encode(buf, smaller)
				sd2, _, _, _ := model.CreateImageStreamDict(ctx.XRefTable, buf, false, false)
				ctx.XRefTable.Table[objs[idx]].Object = *sd2
			}
		}
	}
	pdfcpu.OptimizeXRefTable(ctx)
	//api.OptimizeContext(ctx)
	ctx.EnsureVersionForWriting()
	Log("Write file...")
	wr := new(bytes.Buffer)
	err = api.WriteContext(ctx, wr)
	Bytes = wr.Bytes()
	Log(len(Bytes), " Bytes")
	return len(Bytes)
})
var Bytes []byte
var onMerge = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	LogCallback = args[len(args)-1:][0]
	Log("pdfcomprezzor/Merge")
	Log("Files in Array: ", args[0].Length())
	files := make([]*bytes.Reader, args[0].Length())
	for i := 0; i < len(files); i++ {
		array := make([]byte, args[0].Index(i).Length())
		js.CopyBytesToGo(array, args[0].Index(i))
		files[i] = bytes.NewReader(array)
	}
	seekers := make([]io.ReadSeeker, len(files))
	for i, f := range files {
		seekers[i] = f
	}
	wr := new(bytes.Buffer)
	api.MergeRaw(seekers, wr, nil)
	Log("Write file...")
	Bytes = wr.Bytes()
	return len(Bytes)
})

//onReadBack, reads back value to JS.
var onReadBack = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	js.CopyBytesToJS(args[0], Bytes)
	return 0
})

func main() {
	model.ConfigPath = "disable"
	js.Global().Set("compress", onCompress)
	js.Global().Set("merge", onMerge)
	js.Global().Set("readBack", onReadBack)
	<-done
}
