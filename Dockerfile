FROM golang:1.20 AS build-pdfcomprezzor
WORKDIR /go/src
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN GOOS=js GOARCH=wasm go build -o pdfcomprezzor.wasm
RUN cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
FROM nginx:alpine
COPY --from=build-pdfcomprezzor /go/src/pdfcomprezzor.wasm /usr/share/nginx/html/
COPY --from=build-pdfcomprezzor /go/src/wasm_exec.js /usr/share/nginx/html/
COPY index.html worker.js test.js /usr/share/nginx/html/