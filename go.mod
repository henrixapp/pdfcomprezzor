module github.com/henrixapp/pdfcomprezzor

go 1.16

//replace github.com/pdfcpu/pdfcpu => ./vendor/github.com/pdfcpu/pdfcpu

require (
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pdfcpu/pdfcpu v0.4.0
	github.com/rivo/uniseg v0.4.4 // indirect
)

replace github.com/pdfcpu/pdfcpu => github.com/christian-heusel/pdfcpu v0.4.0
