package pagination

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

//Pagination 分页器
type Pagination struct {
	Request *http.Request
	Total   int
	Pernum  int
	BaseUrl string
}

//NewPagination 新建分页器
func NewPagination(req *http.Request, total int, pernum int, baseUrl string) *Pagination {
	return &Pagination{
		Request: req,
		Total:   total,
		Pernum:  pernum,
		BaseUrl: baseUrl,
	}
}

func (p *Pagination) HtmlPages() template.HTML {
	return template.HTML(p.Pages())
}

//Pages 渲染生成html分页标签
func (p *Pagination) Pages() string {
	queryParams := p.Request.URL.Query()
	//从当前请求中获取page
	page := queryParams.Get("page")
	if page == "" {
		page = "1"
	}
	//将页码转换成整型，以便计算
	pagenum, _ := strconv.Atoi(page)
	if pagenum == 0 {
		return ""
	}

	//计算总页数
	var totalPageNum = int(math.Ceil(float64(p.Total) / float64(p.Pernum)))

	lang := p.getLang()

	//首页链接
	var firstLink string
	//上一页链接
	var prevLink string
	//下一页链接
	var nextLink string
	//末页链接
	var lastLink string
	//中间页码链接
	var pageLinks []string

	//首页和上一页链接
	if pagenum > 1 {
		firstLink = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, p.pageURL("1"), i18n.Tr(lang, "page.first"))
		prevLink = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, p.pageURL(strconv.Itoa(pagenum-1)), i18n.Tr(lang, "page.prev"))
	} else {
		firstLink = fmt.Sprintf(`<li class="disabled"><a href="#">%s</a></li>`, i18n.Tr(lang, "page.first"))
		prevLink = fmt.Sprintf(`<li class="disabled"><a href="#">%s</a></li>`, i18n.Tr(lang, "page.prev"))
	}

	//末页和下一页
	if pagenum < totalPageNum {
		lastLink = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, p.pageURL(strconv.Itoa(totalPageNum)), i18n.Tr(lang, "page.last"))
		nextLink = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, p.pageURL(strconv.Itoa(pagenum+1)), i18n.Tr(lang, "page.next"))
	} else {
		lastLink = fmt.Sprintf(`<li class="disabled"><a href="#">%s</a></li>`, i18n.Tr(lang, "page.last"))
		nextLink = fmt.Sprintf(`<li class="disabled"><a href="#">%s</a></li>`, i18n.Tr(lang, "page.next"))
	}

	//生成中间页码链接
	pageLinks = make([]string, 0, 10)
	startPos := pagenum - 3
	endPos := pagenum + 3
	if startPos < 1 {
		endPos = endPos + int(math.Abs(float64(startPos))) + 1
		startPos = 1
	}
	if endPos > totalPageNum {
		endPos = totalPageNum
	}
	for i := startPos; i <= endPos; i++ {
		var s string
		if i == pagenum {
			s = fmt.Sprintf(`<li class="active"><a href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
		} else {
			s = fmt.Sprintf(`<li><a href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
		}
		pageLinks = append(pageLinks, s)
	}

	return fmt.Sprintf(`<ul class="pagination">%s%s%s%s%s</ul>`, firstLink, prevLink, strings.Join(pageLinks, ""), nextLink, lastLink)
}

//pageURL 生成分页url
func (p *Pagination) pageURL(page string) string {
	//基于当前url新建一个url对象
	u, _ := url.Parse(p.BaseUrl + p.Request.URL.String())
	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()
	return u.String()
}

func (p *Pagination) getLang() string {
	lang, _ := web.AppConfig.String("default_lang")
	ulang := p.Request.FormValue("lang")
	if len(ulang) == 0 {
		clang, err := p.Request.Cookie("lang")
		if err != nil {
			return lang
		}
		ulang = clang.Value
	}
	if !i18n.IsExist(ulang) {
		return lang
	}
	return ulang
}

type Page struct {
	PageNo		int         `json:"PageNo"`
	PageSize	int         `json:"PageSize"`
	TotalPage	int         `json:"TotalPage"`
	TotalCount	int         `json:"TotalCount"`
	FirstPage	bool        `json:"FirstPage"`
	LastPage	bool        `json:"LastPage"`
	List		interface{} `json:"List"`
}

func PageUtil(count int, pageNo int, pageSize int, list interface{}) Page {
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}
	return Page {
		PageNo: pageNo,
		PageSize: pageSize,
		TotalPage: tp,
		TotalCount: count,
		FirstPage: pageNo == 1,
		LastPage: pageNo == tp,
		List: list,
	}
}
