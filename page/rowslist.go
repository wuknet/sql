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
func Rowslist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

		sqlstr := "select column_name,DATA_TYPE,COLUMN_TYPE,COLUMN_KEY,EXTRA from information_schema.columns where table_name = '" + tablename + "' and table_schema='" + database_name + "';"
		rows, err := mysql.Query(db, sqlstr)
		if err != nil {
			rowslist += fmt.Sprintf("%v", err)
		} else {
			if len(rows) > 0 {
				//columns := strings.Split(rows[0][0], ",")
				columnNum := len(rows)
				rowslist += `<table border="1" class="rowslist" width="100%">
			<tr style="background:#eeeeee;">`
				for i := 0; i < columnNum; i++ {
					column_key := "-" //判断是否有主键
					if rows[i][3] != "" {
						column_key = "(" + rows[i][3] + ")"
					}
					auto_increment := ""
					if rows[i][4] == "auto_increment" {
						auto_increment = `&nbsp;<font style="color:red;">&uarr;</font>`
					}

					rowslist += `
				<td align="center"><b>` + rows[i][0] + `</b>` + auto_increment + `<br />` + rows[i][2] + `<br />` + column_key + `</td>`
				}
				rowslist += "</tr>"

				sqlstr = "select * from " + database_name + "." + tablename + " limit 30"
				rows, err = mysql.Query(db, sqlstr)
				if err != nil {
					rowslist += "<tr><td colspan=\"" + strconv.Itoa(columnNum) + "\" align=\"center\">" + fmt.Sprintf("%v", err) + "</td></tr>"
				} else {
					if len(rows) > 0 {
						for i := 0; i < len(rows); i++ {
							rowslist += "<tr>"
							for j := 0; j < len(rows[i]); j++ {
								rowslist += `
				            <td>` + rows[i][j] + `</td>`
							}
							rowslist += "</tr>"
						}
					} else {
						rowslist += "<tr><td colspan=\"" + strconv.Itoa(columnNum) + "\" align=\"center\">没有任何记录！</td></tr>"
					}
				}
				rowslist += "</table>"
			} else {
				rowslist += "<div style=\"color:red;\">没有找到这个表！</div>"
			}
		}

		/////////////////////////////////////////////////////////////////////////////////////
		data := map[string]template.HTML{
			"database_name": template.HTML(database_name),
			"tablename":     template.HTML(tablename),
			"rowslist":      template.HTML(rowslist),
		}

		mb, err := template.ParseFiles(AppPath + "/templates/rowslist.html")
		if err != nil {
			fun.Log(0, "rowslist.go:"+fmt.Sprintf("%s", err))
		}

		err = mb.Execute(w, data)
		if err != nil {
			fun.Log(0, "rowslist.go:"+fmt.Sprintf("%s", err))
		}
	}
}
