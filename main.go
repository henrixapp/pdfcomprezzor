package main

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"

	pdfcomprezzor "github.com/henrixapp/pdfcomprezzor/pkg"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

var done = make(chan struct{})

// var body js.Value
// var doc js.Value
var LogCallback js.Value

func Log(a ...interface{}) {

	if !LogCallback.IsUndefined() {
		LogCallback.Invoke(js.Null(), fmt.Sprint(a...))
	}
}

var Bytes []byte
var onCompress = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	LogCallback = args[len(args)-1:][0]
	Log("pdfcomprezzor/compress")
	Log("File size in Bytes:", args[0].Get("length").Int())
	array := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(array, args[0])
	res, err := pdfcomprezzor.Compress(array, Log, 1000)
	if err != nil {
		return 0
	}
	Bytes = res
	Log(len(res), " Bytes")
	return len(res)
})

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
	Bytes = pdfcomprezzor.Merge(seekers, Log)
	return len(Bytes)
})

// onReadBack, reads back value to JS.
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
