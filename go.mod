module github.com/henrixapp/pdfcomprezzor

go 1.14

//replace github.com/pdfcpu/pdfcpu => ./vendor/github.com/pdfcpu/pdfcpu

require (
	github.com/hhrutter/lzw v0.0.0-20190829144645-6f07a24e8650
	github.com/hhrutter/tiff v0.0.0-20190829141212-736cae8d0bc7
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pdfcpu/pdfcpu v0.3.5
	github.com/pkg/errors v0.9.1
	golang.org/x/image v0.0.0-20200618115811-c13761719519
)
