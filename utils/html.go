package utils

import (
	"regexp"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"github.com/lifei6671/mindoc/conf"
)

func StripTags(s string) string  {

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src := re.ReplaceAllStringFunc(s, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return src
}
//自动提取文章摘要
func AutoSummary(body string,l int) string {

	//匹配图片，如果图片语法是在代码块中，这里同样会处理
	re := regexp.MustCompile(`<p>(.*?)</p>`)

	contents := re.FindAllString(body, -1)

	if len(contents) <= 0 {
		return  ""
	}
	content := ""
	for _,s := range contents {
		b := strings.Replace(StripTags(s),"\n","", -1)

		if l <= 0 {
			break
		}
		l = l - len([]rune(b))

		content += b

	}
	return content
}

//安全处理HTML文档，过滤危险标签和属性.
func SafetyProcessor(html string) string {

	//安全过滤，移除危险标签和属性
	if docQuery, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html)); err == nil {
		docQuery.Find("script").Remove()
		docQuery.Find("form").Remove()
		docQuery.Find("link").Remove()
		docQuery.Find("applet").Remove()
		docQuery.Find("frame").Remove()
		docQuery.Find("meta").Remove()
		docQuery.Find("iframe").Remove()
		docQuery.Find("*").Each(func(i int, selection *goquery.Selection) {

			if href, ok := selection.Attr("href"); ok && strings.HasPrefix(href, "javascript:") {
				selection.SetAttr("href", "#")
			}
			if src, ok := selection.Attr("src"); ok && strings.HasPrefix(src, "javascript:") {
				selection.SetAttr("src", "#")
			}

			selection.RemoveAttr("onafterprint").
				RemoveAttr("onbeforeprint").
				RemoveAttr("onbeforeunload").
				RemoveAttr("onload").
				RemoveAttr("onclick").
				RemoveAttr("onkeydown").
				RemoveAttr("onkeypress").
				RemoveAttr("onkeyup").
				RemoveAttr("ondblclick").
				RemoveAttr("onmousedown").
				RemoveAttr("onmousemove").
				RemoveAttr("onmouseout").
				RemoveAttr("onmouseover").
				RemoveAttr("onmouseup")
		})

		//处理外链
		docQuery.Find("a").Each(func(i int, contentSelection *goquery.Selection) {
			if src, ok := contentSelection.Attr("href"); ok {
				if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
					if conf.BaseUrl != "" && !strings.HasPrefix(src, conf.BaseUrl) {
						contentSelection.SetAttr("target", "_blank")
					}
				}
			}
		})
		//添加文档标签包裹
		if selector := docQuery.Find("article.markdown-article-inner").First(); selector.Size() <= 0 {
			docQuery.Children().WrapAllHtml("<article class=\"markdown-article-inner\"></article>")
		}
		//解决文档内容缺少包裹标签的问题
		if selector := docQuery.Find("div.markdown-article").First(); selector.Size() <= 0 {
			if selector := docQuery.Find("div.markdown-toc").First(); selector.Size() > 0 {
				docQuery.Find("div.markdown-toc").NextAll().WrapAllHtml("<div class=\"markdown-article\"></div>")
			}
		}


		if html, err := docQuery.Html(); err == nil {
			return  strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(html), "<html><head></head><body>"), "</body></html>")
		}
	}
	return html
}