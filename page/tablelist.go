package page

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	//"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	//"context"
	//"github.com/go-session/session"

	"../fun"
	"../mysql"
)

//////////////////////////////////////////////////////
func Tablelist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//t := time.Now() //开始计算执行时间
	AppPath := fun.GetAppPath() //读取当前目录路径

	addr := fun.ReadCookie(w, r, "addr")
	username := fun.ReadCookie(w, r, "username")
	password := fun.ReadCookie(w, r, "password")
	database_name := fun.Get_str(r, "database_name")

	//fun.WriteCookie(w, r, "database_name", database_name, 0)
	db, linkzt, _ := mysql.Conndb(addr, username, password, database_name)
	defer db.Close()
	if linkzt == false {
		http.Redirect(w, r, "/datalist", http.StatusFound)
		fmt.Fprintf(w, "Database link lost!")
		return
	} else {
		tablelist := ""
		tableNum := 0 //表数量
		//sql := "select table_name from information_schema.tables where table_schema='" + database_name + "' and table_type='base table';"
		sql := "select `table_name`,`table_type`,`engine`,`data_length`,AUTO_INCREMENT,TABLE_COLLATION from information_schema.tables where table_schema='" + database_name + "';"
		rows, err := mysql.Query(db, sql)
		if err != nil {
			tablelist += fmt.Sprintf("%v", err)
		} else {
			tableNum = len(rows)
			tablelist += `<table class="tablelist" border="1" width="100%">`
			tablelist += "<tr style=\"background:#eeeeee;\">"
			tablelist += "<td>表名</td>"
			tablelist += "<td>类型</td>"
			tablelist += "<td>驱动</td>"
			tablelist += "<td>大小</td>"
			tablelist += "<td>自增</td>"
			tablelist += "<td>编码</td>"
			tablelist += "</tr>"
			if tableNum > 0 {
				for i := 0; i < tableNum; i++ {
					tablelist += "<tr>"
					tablelist += `<td><a href="/rowslist/?database_name=` + database_name + `&tablename=` + rows[i][0] + `">` + rows[i][0] + `</a></td>`
					tablelist += `<td>` + rows[i][1] + `</td>`
					tablelist += `<td>` + rows[i][2] + `</td>`
					datasize, _ := strconv.ParseInt(rows[i][3], 10, 64)
					tablelist += `<td>` + formatFileSize(datasize) + `</td>`
					tablelist += `<td>` + rows[i][4] + `</td>`
					tablelist += `<td>` + rows[i][5] + `</td>`
					tablelist += "</tr>"
				}
			} else {
				tablelist += "<tr>"
				tablelist += `<td colspan="6" align="center">还没有任何表！</td>`
				tablelist += "</tr>"
			}
			tablelist += "</table>"
		}
		////////////////////
		database_sizes, _ := mysql.QueryOne(db, "select sum(DATA_LENGTH) as data from information_schema.TABLES where table_schema='"+database_name+"';")
		database_size, _ := strconv.ParseInt(database_sizes[0], 10, 64)
		database_sizes[0] = formatFileSize(database_size)
		/////////////////////////////////////////////////////////////////////////////////////
		data := map[string]template.HTML{
			"tableNum":      template.HTML(strconv.Itoa(tableNum)),
			"database_name": template.HTML(database_name),
			"database_size": template.HTML(database_sizes[0]),
			"tablelist":     template.HTML(tablelist),
		}

		mb, err := template.ParseFiles(AppPath + "/templates/tablelist.html")
		if err != nil {
			fun.Log(0, "login.go:"+fmt.Sprintf("%s", err))
		}

		err = mb.Execute(w, data)
		if err != nil {
			fun.Log(0, "login.go:"+fmt.Sprintf("%s", err))
		}
	}
}

func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
