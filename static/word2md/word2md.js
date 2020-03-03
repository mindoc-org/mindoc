(function () {
    var options = {
        convertImage: mammoth.images.imgElement(function (image) {
            var result = "";
            return image.read("base64").then(function (imageBuffer, $callback) {
                var fileName = Date.parse(new Date());

                switch (image.contentType) {
                case "image/png":
                    fileName += ".png";
                    break;
                case "image/jpg":
                    fileName += ".jpg";
                    break
                case "image/jpeg":
                    fileName += ".jpeg";
                    break;
                case "image/gif":
                    fileName += ".gif";
                    break;
                default:
                    layer.msg("不支持的图片格式");
                    return;
                }
                var form = new FormData();
                form.append('editormd-image-file', base64ToBlob(imageBuffer, image.contentType), fileName);

                var layerIndex = 0;

                return {
                    src: _ajax(window.imageUploadURL, form, function (ret) {
                        if (ret.success == 1) {
                            //return ret.url;
                            //cm.replaceSelection("![](" + ret.url  + ")");
                        }
                        console.log(ret.message);
                    })
                };

            });
        })
    };
    var _ajax = function (url, data, callback) {
        $.ajax({
            "type": 'post',
            "cache": false,
            "url": url,
            "data": data,
            "dateType": "json",
            "processData": false,
            "contentType": false,
            "mimeType": "multipart/form-data",
            async: false,
            success: function (ret) {
                callback(JSON.parse(ret));
                result = JSON.parse(ret).url;
            },
            error: function (err) {
                console.log('请求失败')
            }
        })
        return result;
    };
    function base64ToBlob(base64, mime) {
        mime = mime || "";
        const sliceSize = 1024;
        const byteChars = window.atob(base64);
        const byteArrays = [];
        for (
            let offset = 0, len = byteChars.length;
            offset < len;
            offset += sliceSize) {
            const slice = byteChars.slice(offset, offset + sliceSize);

            const byteNumbers = new Array(slice.length);
            for (let i = 0; i < slice.length; i++) {
                byteNumbers[i] = slice.charCodeAt(i);
            }

            const byteArray = new Uint8Array(byteNumbers);

            byteArrays.push(byteArray);
        }

        return new Blob(byteArrays, {
            type: mime
        });
    }
    document.getElementById("document")
    .addEventListener("change", handleFileSelect, false);

    function handleFileSelect(event) {
        readFileInputEventAsArrayBuffer(event, function (arrayBuffer) {
            mammoth.convertToHtml({
                arrayBuffer: arrayBuffer
            }, options)
            .then(displayResult)
            .done();
        });
    }

    function displayResult(result) {
        document.getElementById("output").innerHTML = result.value;

        var messageHtml = result.messages.map(function (message) {
                return '<li class="' + message.type + '">' + escapeHtml(message.message) + "</li>";
            }).join("");

        document.getElementById("messages").innerHTML = "<ul>" + messageHtml + "</ul>";
    }

    function readFileInputEventAsArrayBuffer(event, callback) {
        var file = event.target.files[0];

        var reader = new FileReader();

        reader.onload = function (loadEvent) {
            var arrayBuffer = loadEvent.target.result;
            callback(arrayBuffer);
        };

        reader.readAsArrayBuffer(file);
    }

    function escapeHtml(value) {
        return value
        .replace(/&/g, '&amp;')
        .replace(/"/g, '&quot;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
    }
})();
