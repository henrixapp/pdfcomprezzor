# pdfcomprezzor

Simple Wasm-Program to compress pdffiles with pictures in browser, based on pdfcpu. The programm replaces all images that are bigger than 1240x1754 are resized to this size and converted to JPG. This is done using pdfcpu. We decode, compress, delete and insert each embedded image. This reuses the object-id of the original object and thus replaces the image.

## Build
```
go mod vendor
sed -i 's/init(/init2(/g' vendor/github.com/pdfcpu/pdfcpu/pkg/font/metrics.go
```
 The last code deactivates the init function in `vendor/github.com/pdfcpu/pdfcpu/pkg/font/metrics.go`, otherwise the program can not start in Wasm. This is caused initialy by a call to `User.Dir()` and a subsequent `os.Exit(1)`.

```
GOOS=js GOARCH=wasm go build -o pdfcomprezzor.wasm 
```
## API
### using worker.js

See test.js

## Running the example

 Serve the files with a server, that supports mime-type `application/wasm`

Navigate to index.html and open a file, wait for compression...

If you select two files or more,  they will be merged.

License: Apache 2.0

(c) 2020, Henrik Reinstädtler
