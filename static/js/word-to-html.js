
class WordToHtmlConverter {

    handleFileSelect(callback, accept = '.doc,.docx') {
        let input = document.createElement('input');
        input.type = 'file';
        input.accept = accept;
        input.onchange = (e)=>{
            let file = e.target.files[0];
            if (!file) {
                return;
            }
            let reader = new FileReader();
            reader.onload = (e)=>{
                let arrayBuffer = e.target.result;
                this.convertToHtml(arrayBuffer, (html)=>{
                    callback(html);
                });
            };
            reader.readAsArrayBuffer(file);
        };
        input.click();
    }



    convertToHtml(arrayBuffer, callback) {
        try {
            mammoth.convertToHtml({arrayBuffer: arrayBuffer}).then(callback, function(error) {
                layer.msg('Error: ' + error);
                return
            });
        } catch (error) {
            layer.msg('Error: ' + error);
            return
        }
    }

    replaceHtmlBase64(html) {
        let regex = /<img\s+src="data:image\/[^;]*;base64,([^"]*)"/g;
        let matches = [];
        let match;

        while ((match = regex.exec(html)) !== null) {
            matches.push(match[1]);
        }

        if (matches.length === 0) {
            return new Promise((resolve, reject) => {
                resolve(html);
            });
        }

        // 将base64转为blob
        let promises = matches.map((base64)=>{
            return new Promise((resolve, reject)=>{
                let blob = this.base64ToBlob(base64);
                let reader = new FileReader();
                reader.onload = (e)=>{
                    resolve({base64, blob, url: e.target.result});
                };
                reader.readAsDataURL(blob);
            });
        });

        return Promise.all(promises).then((results)=>{
            let htmlCopy = html;
            return Promise.all(results.map((result) => {
                return this.uploadFile(result.blob).then((data) => {
                    htmlCopy = htmlCopy.replace(`data:image/png;base64,${result.base64}`, data.url);
                });
            })).then(() => {
                return htmlCopy;
            });
        });
    }

    uploadFile(blob) {
        let file = new File([blob], 'image.jpg', { type: 'image/jpeg' });
        let formData = new FormData();
        formData.append('editormd-file-file', file);
        return fetch(window.fileUploadURL, {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {return data})
            .catch(error => {return error});
    }


    base64ToBlob(base64, type) {
        let binary = atob(base64);
        let array = [];
        for (let i = 0; i < binary.length; i++) {
            array.push(binary.charCodeAt(i));
        }
        return new Blob([new Uint8Array(array)], {type: type});
    }
}