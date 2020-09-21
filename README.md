# pdfcomprezzor

Simple WASM-Program to compress pdffiles with pictures in browser, based on pdfcpu. The programm replaces all images that are bigger than 1240x1754 are resized to this size and converted to JPG. This is done using pdfcpu. We decode, compress, delete and insert each embedded image. This reuses the object of the original object and thus replaces the image.

## Build
```
go mod vendor
```
 Deactivate the init function in `vendor/github.com/pdfcpu/pdfcpu/pkg/font/metrics.go`, otherwise the program can not start in WASM. This is caused initialy by a call to `User.Dir()`.

```
GOOS=js GOARCH=wasm go build -o pdfcomprezzor.wasm 
```
## API
### using worker.js
 
```js
var l={l:0};
var worker = new Worker('worker.js');

worker.addEventListener('message', function(e) {
  console.log('Worker said: ', e);
  if(e.data.type=="log"){
  let div = document.createElement( "div");
	div.textContent =e.data.message;
    document.querySelector("body").appendChild(div);
  } else if (e.data.type=="result"){
   console.log(l);
   alert(`TOOK: ${e.data.time}`)
    downloadBlob(e.data.result,"smaller.pdf","application/pdf");
    }
}, false);
worker.postMessage({array,l});
```
where array is a Uint8Array containing the pdf file. You could load it, via 

```js
var reader = new FileReader();
reader.onload = function() {
    var arrayBuffer = this.result;
    array = new Uint8Array(arrayBuffer);
    // code from above....
  };
reader.readAsArrayBuffer(this.files[0]);
```

## Running the example

 Serve the files with a server, that supports mime-type `application/wasm`

Navigate to index.html and open a file, wait for compression...

License: Apache 2.0 
(c) 2020, Henrik Reinstädtler