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
		fmt.Fprintf(w, "Database link lost!")
		return
	} else {
		//sql := "select table_name from information_schema.tables where table_schema='" + database_name + "' and table_type='base table';"
		sql := "select `table_name`,`table_type`,`engine`,`data_length`,AUTO_INCREMENT,TABLE_COLLATION from information_schema.tables where table_schema='" + database_name + "';"
		rows, _ := mysql.Query(db, sql)
		tableNum := len(rows)
		tablelist := `<table class="tablelist" border="1" width="100%">`
		tablelist += "<tr style=\"background:#eeeeee;\">"
		tablelist += "<td>表名</td>"
		tablelist += "<td>类型</td>"
		tablelist += "<td>驱动</td>"
		tablelist += "<td>大小</td>"
		tablelist += "<td>自增</td>"
		tablelist += "<td>编码</td>"
		tablelist += "</tr>"
		for i := 0; i < tableNum; i++ {
			tablelist += "<tr>"
			tablelist += `<td><a href="/rowslist/?database_name=` + database_name + `&tablename=` + rows[i][0] + `">` + rows[i][0] + `</a></td>`
			tablelist += `<td>` + rows[i][1] + `</td>`
			tablelist += `<td>` + rows[i][2] + `</td>`
			tablelist += `<td>` + rows[i][3] + `</td>`
			tablelist += `<td>` + rows[i][4] + `</td>`
			tablelist += `<td>` + rows[i][5] + `</td>`
			tablelist += "</tr>"
		}
		tablelist += "</table>"

		/////////////////////////////////////////////////////////////////////////////////////
		data := map[string]template.HTML{
			"tableNum":      template.HTML(strconv.Itoa(tableNum)),
			"database_name": template.HTML(database_name),
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
