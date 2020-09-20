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
	"os"

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
func main() {
	ctx, err := pdfcpu.ReadFile("test2.pdf", pdfcpu.NewDefaultConfiguration())
	if err != nil {
		fmt.Println("Error", err)
	}
	//fmt.Println(ctx)
	ctx.EnsurePageCount()
	count := ctx.PageCount
	pdfcpu.OptimizeXRefTable(ctx)
	fmt.Println(count)
	images := make([]CompressedImage, 0)
	for i := 1; i <= count; i++ {
		fmt.Println(i)
		for _, objNr := range ctx.ImageObjNrs(i) {
			fmt.Println(objNr)
			imageO := ctx.Optimize.ImageObjects[objNr]
			imageF, _ := ctx.ExtractImage(objNr)
			fmt.Println(imageF.Name)
			image1, format, err := image.Decode(imageF)
			fmt.Println(format)
			if err != nil {
				fmt.Println("e", err)
			}
			fmt.Println(image1.Bounds().Dx(), image1.Bounds().Dy())
			if image1.Bounds().Dx() > 1000 {
				smaller := resize.Resize(uint(image1.Bounds().Dx()/2.0), 0, image1, resize.Lanczos2)
				var b bytes.Buffer
				w := bufio.NewWriter(&b)
				err := jpeg.Encode(w, smaller, nil)
				images = append(images, CompressedImage{Image: smaller, ObjNr: objNr})
				if err != nil {
					fmt.Println(err)
				}
				ctx.DeleteObject(objNr)
				buf := new(bytes.Buffer)
				err = jpeg.Encode(buf, smaller, nil)
				images[len(images)-1].Ref, _, _, _ = createImageResource(ctx.XRefTable, buf)
				//imageO.ImageDict.StreamLength
			}
			fmt.Println(imageO.ImageDict.FilterPipeline[0].Name)
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
	f, err := os.Create("fixed2.pdf")
	defer f.Close()
	wr := bufio.NewWriter(f)
	api.WriteContext(ctx, wr)
}
