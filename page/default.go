package page

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"../fun"
	"github.com/julienschmidt/httprouter"
)

func Default(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	Hostname := r.Host
	Hostname = strings.Split(Hostname, ":")[0]

	urlname, _ := url.PathUnescape(r.RequestURI) //解码url
	urlname = strings.ToLower(urlname)           //全部转小写
	urls := strings.Split(urlname, "/")          //转成数组

	switch urls[1] {
	case "":
		{
			Index(w, r, ps)
		}
	case "login":
		{
			Login(w, r, ps)
		}
	case "run":
		{
			Run(w, r, ps)
		}
	case "datalist":
		{
			Datalist(w, r, ps)
		}
	case "tablelist":
		{
			Tablelist(w, r, ps)
		}
	case "rowslist":
		{
			Rowslist(w, r, ps)
		}
	case "column":
		{
			Column(w, r, ps)
		}
	case "exit":
		{
			fun.WriteCookie(w, r, "username", "", 0)
			fun.WriteCookie(w, r, "password", "", 0)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	default:
		{
			AppPath := fun.GetAppPath() //读取当前目录路径
			PathFile := AppPath + "/www/" + urlname
			str, err := ioutil.ReadFile(PathFile)
			if err != nil {
				str, err = ioutil.ReadFile(AppPath + "/www/error/wuknet_404.html")
				if err != nil {
					fmt.Fprintf(w, Hostname+"：Read Error!"+PathFile)
				} else {
					w.Header().Add("Content-Type", GetContentType(urlname))
					w.Write(str)
				}
			} else {
				w.Header().Add("Content-Type", GetContentType(urlname))
				w.Write(str)
			}
		}
	}
}

//判断文件类型
func GetContentType(filename string) string {
	var contentType string
	if strings.HasSuffix(filename, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(filename, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(filename, ".txt") {
		contentType = "text/html"
	} else if strings.HasSuffix(filename, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".svg") {
		contentType = "image/svg+xml"
	} else {
		contentType = "text/plain"
	}
	return contentType
}
