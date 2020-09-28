downloadBlob = function (data, fileName, mimeType) {
    var blob, url;
    blob = new Blob([data], {
        type: mimeType
    });
    url = window.URL.createObjectURL(blob);
    downloadURL(url, fileName);
    setTimeout(function () {
        return window.URL.revokeObjectURL(url);
    }, 1000);
};

downloadURL = function (data, fileName) {
    var a;
    a = document.createElement('a');
    a.href = data;
    a.download = fileName;
    document.body.appendChild(a);
    a.style = 'display: none';
    a.click();
    a.remove();
};

function fn() {
    document.querySelector('input').addEventListener('change', function () {
        var files = [];
        var self = this;
        var readFileIntoFiles = function () {
            if (this.result) {
                var arrayBuffer = this.result;
                array = new Uint8Array(arrayBuffer);
                console.log(files.push(array));
                reader.onload = readFileIntoFiles;
                if (files.length < self.files.length)
                    reader.readAsArrayBuffer(self.files[files.length]);
                else {
                    var worker = new Worker('worker.js');

                    worker.addEventListener('message', function (e) {
                        console.log('Worker said: ', e);
                        if (e.data.type == "log") {
                            let div = document.createElement("div");
                            div.textContent = e.data.message;
                            document.querySelector("body").appendChild(div);
                        } else if (e.data.type == "result") {
                            alert(`TOOK: ${e.data.time}`)
                            downloadBlob(e.data.result, "smaller.pdf", "application/pdf");
                        }
                    }, false);
                    let action = files.length > 1 ? "merge" : "compress";
                    let r =
                        worker.postMessage({
                            array: files,
                            action
                        });
                }
            }

        }
        var reader = new FileReader();
        reader.onload = readFileIntoFiles;
        reader.readAsArrayBuffer(this.files[0]);

    }, false);
}
document.addEventListener('DOMContentLoaded', fn, false);