
class HtmlToMarkdownConverter {

    handleFileSelect(callback, accept = '.html,.htm') {
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
                let text = e.target.result;
                let markdown = this.convertToMarkdown(text);
                callback(markdown);
            };
            reader.readAsText(file);
        };
        input.click();
    }


    convertToMarkdown(html) {
        let turndownService = new TurndownService()
        return turndownService.turndown(html)
    }
}