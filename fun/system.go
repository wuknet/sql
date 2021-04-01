package fun

import (
	//"fmt"

	//"log"

	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	//"log"
	//"html/template"
	//"github.com/julienschmidt/httprouter"
)

//PageLinkNum表示一页显示多少按钮,PageNum表示有多少页,CurrentPage表示当前页
//url表示连接表达式，如 ///ranking/?page=1   page参数是变量{page}
func Fpage(PageLinkNum int, PageNum int, CurrentPage int, url string) string {
	str := ""
	if CurrentPage > PageNum { //如果当前页大于总数则
		CurrentPage = PageNum
	}
	wz := CurrentPage - PageLinkNum/2 //翻页开始位置是当前页 减掉 页数平均值
	if wz < 1 {
		wz = 1
	}
	if wz > 1 { //如果翻页大于1
		str += `<a href="` + strings.Replace(url, "{page}", "1", -1) + `" class="pageinput">1..</a> `
	}
	if PageNum < PageLinkNum { //如果实际按钮小于总页数，按实际分页数来
		PageLinkNum = PageNum
	}
	for i := wz; i < wz+PageLinkNum; i++ {
		if i > PageNum { //如果循环的页码大于总页数，跳出循环
			break
		}
		if i == CurrentPage {
			str += `<a href="` + strings.Replace(url, "{page}", strconv.Itoa(i), -1) + `" class="pageinput2">` + strconv.Itoa(i) + `</a> `
		} else {
			str += `<a href="` + strings.Replace(url, "{page}", strconv.Itoa(i), -1) + `" class="pageinput">` + strconv.Itoa(i) + `</a> `
		}
	}
	if PageNum > wz+PageLinkNum { //如果总数大于当前总翻页数，表示还可以往后翻
		str += `<a href="` + strings.Replace(url, "{page}", strconv.Itoa(wz+PageLinkNum), -1) + `" class="pageinput">..</a> `
	}
	return str
}
