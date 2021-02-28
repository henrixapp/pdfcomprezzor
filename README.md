# pdfcomprezzor

Simple Wasm-Program to compress pdffiles with pictures in browser, based on pdfcpu. The programm replaces all images that are bigger than 1240x1754 are resized to this size and converted to JPG. This is done using pdfcpu. We decode, compress, delete and insert each embedded image. This reuses the object-id of the original object and thus replaces the image.

## Build
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
