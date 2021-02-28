package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"syscall/js"

	"github.com/nfnt/resize"
	api "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type CompressedImage struct {
	ObjNr int
	Image image.Image
	Ref   *pdfcpu.IndirectRef
}

//bit hacky, because was not marked for export
func createImageResource(xRefTable *pdfcpu.XRefTable, r io.Reader) (*pdfcpu.IndirectRef, int, int, error) {

	bb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, 0, err
	}

	var sd *pdfcpu.StreamDict
	r = bytes.NewReader(bb)

	// We identify JPG via its magic bytes.
	if bytes.HasPrefix(bb, []byte("\xff\xd8")) {
		// Process JPG by wrapping byte stream into DCTEncoded object stream.
		c, _, err := image.DecodeConfig(r)
		if err != nil {
			return nil, 0, 0, err
		}

		sd, err = pdfcpu.ReadJPEG(xRefTable, bb, c)
		if err != nil {
			return nil, 0, 0, err
		}

	} else {
		// Process other formats by decoding into an image
		// and subsequent object stream encoding,
		/*img, _, err := image.Decode(r)
		if err != nil {
			return nil, 0, 0, err
		}

		sd, err = imgToImageDict(xRefTable, img)
		if err != nil {
			return nil, 0, 0, err
		}*/
		log.Panicln("not supported")
	}

	w := *sd.IntEntry("Width")
	h := *sd.IntEntry("Height")

	indRef, err := xRefTable.IndRefForNewObject(*sd)
	if err != nil {
		return nil, 0, 0, err
	}

	return indRef, w, h, nil
}

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
	ctx, err := pdfcpu.Read(buffi, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		Log("Error while loading file", err)
	}
	//Log(ctx)
	ctx.EnsurePageCount()
	count := ctx.PageCount
	pdfcpu.OptimizeXRefTable(ctx)
	Log("Page count:", count)
	images := make([]CompressedImage, 0)

	for i := 1; i <= count; i++ {
		Log("Processing page no:", i)
		imageObjNrs := ctx.ImageObjNrs(i)
		Log("Images on page:", len(imageObjNrs))
		for _, objNr := range imageObjNrs {
			imageF, _ := ctx.ExtractImage(objNr)
			image1, format, err := image.Decode(imageF)
			Log("Image format:", format)
			if err != nil {
				Log("e", err)
			}
			Log("Image size:", image1.Bounds().Dx(), image1.Bounds().Dy())
			if image1.Bounds().Dx() > 1000 {
				Log("Compress this image.....")
				smaller := resize.Thumbnail(1240, 1754, image1, resize.Lanczos2)
				var b bytes.Buffer
				w := bufio.NewWriter(&b)
				err := jpeg.Encode(w, smaller, nil)
				images = append(images, CompressedImage{ObjNr: objNr})
				if err != nil {
					Log("Error enconding", err)
				}
				ctx.DeleteObject(objNr)
				buf := new(bytes.Buffer)
				err = jpeg.Encode(buf, smaller, nil)
				images[len(images)-1].Ref, _, _, _ = createImageResource(ctx.XRefTable, buf)
			}
		}
	}
	ctx.EnsureVersionForWriting()
	Log("Write file...")
	wr := new(bytes.Buffer)
	api.WriteContext(ctx, wr)
	Bytes = wr.Bytes()
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
	api.Merge(seekers, wr, nil)
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
	pdfcpu.ConfigPath = "disable"
	js.Global().Set("compress", onCompress)
	js.Global().Set("merge", onMerge)
	js.Global().Set("readBack", onReadBack)
	<-done
}
