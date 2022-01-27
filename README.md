# pdfcomprezzor

Simple Wasm-Program to compress pdffiles with pictures in browser, based on pdfcpu. The programm replaces all images that are bigger than 1240x1754 are resized to this size and converted to JPG. This is done using pdfcpu. We decode, compress, delete and insert each embedded image. This reuses the object-id of the original object and thus replaces the image.

## Build

```
GOOS=js GOARCH=wasm go build -o pdfcomprezzor.wasm 
```
### Docker
If you want to try this project using docker, run the container by following these steps:

```sh
sudo docker build -t dockerized_pdf_comprezzor .
sudo docker run -p 8081:80 dockerized_pdf_comprezzor
```
Now you can visit localhost:8081 to compress files.

## API
### using worker.js

See test.js

## Running the example
Firstly, you have to copy a version of wasm_exec.js to this project. 
`cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .` 

 Serve the files with a server, that supports mime-type `application/wasm`

Navigate to index.html and open a file, wait for compression...

If you select two files or more,  they will be merged.

License: Apache 2.0

(c) 2020-2022, Henrik Reinstädtler
