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
	"time"

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
	} else {
		/*div := doc.Call("createElement", "div")
		div.Set("textContent", fmt.Sprint(a))
		body.Call("appendChild", div)*/
	}
	time.Sleep(200 * time.Millisecond)
}
func main() {
	onCompress := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		LogCallback = args[len(args)-1:][0]
		Log("called")
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
				Log("Image file name:", imageF.Name)
				image1, format, err := image.Decode(imageF)
				Log("Image format:", format)
				if err != nil {
					Log("e", err)
				}
				Log("Image size:", image1.Bounds().Dx(), image1.Bounds().Dy())
				if image1.Bounds().Dx() > 1000 {
					Log("Compress this image.....")
					smaller := resize.Thumbnail(1240, 1754, image1, resize.Lanczos2)
					//smaller := resize.Resize(uint(image1.Bounds().Dx()/2.0), 0, image1, resize.Lanczos2)
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
					//imageO.ImageDict.StreamLength
				}
				//imageO.ImageDict.Raw =
			}

		}
		/*fmt.Println(ctx.XRefTable.FindObject(1))
		root, _ := ctx.XRefTable.FindObject(int(ctx.Root.ObjectNumber))
		fmt.Println((root.(pdfcpu.Dict)["Pages"]).(pdfcpu.IndirectRef).ObjectNumber)
		rootPages, _ := ctx.XRefTable.FindObject(int((root.(pdfcpu.Dict)["Pages"]).(pdfcpu.IndirectRef).ObjectNumber))
		fmt.Println(rootPages.(pdfcpu.Dict).ArrayEntry("Kids"))
		for _, pageR := range rootPages.(pdfcpu.Dict).ArrayEntry("Kids") {
			page, _ := ctx.FindObject(int(pageR.(pdfcpu.IndirectRef).ObjectNumber))
			fmt.Println(page)
			resources, _ := ctx.FindObject(int(page.(pdfcpu.Dict)["Resources"].(pdfcpu.IndirectRef).ObjectNumber))
			fmt.Println(resources)
			n := resources.(pdfcpu.Dict)["XObject"]
			fmt.Println(reflect.TypeOf(n.(pdfcpu.Dict)["Img1"]))
		}*/
		ctx.EnsureVersionForWriting()
		Log("Write file...")
		wr := new(bytes.Buffer)
		api.WriteContext(ctx, wr)
		js.CopyBytesToJS(args[0], wr.Bytes())
		args[1].Set("l", len(wr.Bytes()))
		return len(wr.Bytes())
	})

	js.Global().Set("compress", onCompress)
	<-done
	if false {

	}
}
