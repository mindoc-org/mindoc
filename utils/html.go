package utils

import (
	"regexp"
	"strings"
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
