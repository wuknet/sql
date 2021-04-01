package page

import (
	"fmt"
	"html/template"
	"net/http"

	//"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	//"context"
	//"github.com/go-session/session"

	"../fun"
	"../mysql"
)

//////////////////////////////////////////////////////
func Column(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//t := time.Now() //开始计算执行时间
	AppPath := fun.GetAppPath() //读取当前目录路径

	addr := fun.ReadCookie(w, r, "addr")
	username := fun.ReadCookie(w, r, "username")
	password := fun.ReadCookie(w, r, "password")
	database_name := fun.Get_str(r, "database_name")
	tablename := fun.Get_str(r, "tablename")
	if tablename == "" {
		tablename = fun.ReadCookie(w, r, "tablename")
	}

	db, linkzt, _ := mysql.Conndb(addr, username, password, database_name)
	defer db.Close()
	if linkzt == false {
		fmt.Fprintf(w, "Database link lost!")
		return
	} else {
		rowslist := ""
		//select column_name,columnName, data_type dataType, column_comment columnComment, column_key columnKey, extra from information_schema.columns

		sqlstr := "select column_name,DATA_TYPE,COLUMN_TYPE,COLUMN_KEY,EXTRA,COLUMN_DEFAULT,IS_NULLABLE,COLLATION_NAME from information_schema.columns where table_name = '" + tablename + "' and table_schema='" + database_name + "';"
		rows, _ := mysql.Query(db, sqlstr)
		if len(rows) > 0 {
			//columns := strings.Split(rows[0][0], ",")
			columnNum := len(rows)
			rowslist += `<table border="1" class="rowslist" width="100%">`
			rowslist += `<tr style="background:#eeeeee;">`
			rowslist += `<td>字段名称</td>`
			rowslist += `<td>类型</td>`
			rowslist += `<td>主键</td>`
			rowslist += `<td>默认值</td>`
			rowslist += `<td>支持NULL</td>`
			rowslist += `<td>编码</td>`
			rowslist += `</tr>`
			for i := 0; i < columnNum; i++ {
				column_key := "-" //判断是否有主键
				if rows[i][3] != "" {
					column_key = "(" + rows[i][3] + ")"
				}
				auto_increment := ""
				if rows[i][4] == "auto_increment" {
					auto_increment = `<font style="color:red;">&uarr;</font>`
				}
				rowslist += `<tr>`
				rowslist += `<td><b>` + rows[i][0] + `</b>` + auto_increment + `</td>`
				rowslist += `<td>` + rows[i][2] + `</td>`
				rowslist += `<td>` + column_key + `</td>`
				rowslist += `<td>` + rows[i][5] + `</td>`
				rowslist += `<td>` + rows[i][6] + `</td>`
				rowslist += `<td>` + rows[i][7] + `</td>`
				rowslist += "</tr>"
			}
			rowslist += "</table>"
		}

		/////////////////////////////////////////////////////////////////////////////////////
		data := map[string]template.HTML{
			"database_name": template.HTML(database_name),
			"tablename":     template.HTML(tablename),
			"rowslist":      template.HTML(rowslist),
		}

		mb, err := template.ParseFiles(AppPath + "/templates/column.html")
		if err != nil {
			fun.Log(0, "column.go:"+fmt.Sprintf("%s", err))
		}

		err = mb.Execute(w, data)
		if err != nil {
			fun.Log(0, "column.go:"+fmt.Sprintf("%s", err))
		}
	}
}
