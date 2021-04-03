package page

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	//"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	//"context"
	//"github.com/go-session/session"

	"../fun"
	"../mysql"
)

//////////////////////////////////////////////////////
func Run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//t := time.Now() //开始计算执行时间
	AppPath := fun.GetAppPath() //读取当前目录路径

	addr := fun.ReadCookie(w, r, "addr")
	username := fun.ReadCookie(w, r, "username")
	password := fun.ReadCookie(w, r, "password")
	database_name := fun.Get_str(r, "database_name")
	if database_name == "" {
		database_name = "未知"
	}
	tablename := fun.Get_str(r, "tablename")
	sqlstr := fun.Get_str(r, "sqlstr")
	if tablename == "" {
		tablename = "表名"
	}
	isrun := fun.Get_str(r, "isrun")

	runResult := ""
	if isrun == "1" {
		db, linkzt, _ := mysql.Conndb(addr, username, password, database_name)
		defer db.Close()
		if linkzt == false {
			fmt.Fprintf(w, "Database link lost!")
			return
		} else {
			sqlstr = strings.Trim(sqlstr, " ") //清除二边空格
			sqls := strings.Split(sqlstr, " ") //分解

			commstr := strings.ToLower(sqls[0])
			if commstr == "select" {
				if strings.Contains(strings.ToLower(sqlstr), " limit ") == false {
					if strings.HasSuffix(sqlstr, ";") {
						sqlstr = strings.TrimRight(sqlstr, ";")
					}
					sqlstr = sqlstr + " limit 30;"
				}
			}
			runResult = "状态：<br />"
			if commstr == "drop" || commstr == "delete" || commstr == "create" {
				_, err := db.Exec(sqlstr)
				if err != nil {
					runResult += "<div style=\"color:red;\">" + fmt.Sprintf("%v", err) + "</div>"
				} else {
					runResult += "<div style=\"color:blue;\">Run successfully！</div>"
				}
			} else {
				rows, err := mysql.Query(db, sqlstr)
				if err != nil {
					runResult += "查询错误：" + fmt.Sprintf("%v", err)
				} else {
					if len(rows) > 0 {
						runResult += `<table border="1" class="runResult">`
						for i := 0; i < len(rows); i++ {
							runResult += `<tr>`
							for j := 0; j < len(rows[i]); j++ {
								runResult += "<td>&nbsp;" + rows[i][j] + "</td>"
							}
							runResult += `</tr>`
						}
						runResult += `</table>`
					} else {
						runResult += "没有任何记录！"
					}
				}
			}
		}
	} else { //运行样例
		runComm := fun.Get_str(r, "runcomm")
		switch runComm {
		case "create_table":
			{
				sqlstr = `
CREATE TABLE IF NOT EXISTS ` + tablename + ` (
	id INT UNSIGNED AUTO_INCREMENT,
	title VARCHAR(100) NOT NULL,
	content VARCHAR(40) NOT NULL,
	create_date DATETIME,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
#建表语法样例`
			}

		case "delete_table":
			{
				sqlstr = `
DROP TABLE ` + tablename + `
#删除表样例`
			}

		case "alter_table":
			{
				sqlstr = `
ALTER TABLE ` + tablename + ` ADD COLUMN_Name VARCHAR(100) NOT NULL
ALTER TABLE ` + tablename + ` MODIFY COLUMN_Name VARCHAR(100) NOT NULL
ALTER TABLE ` + tablename + ` DROP COLUMN_Name
#字段操作样例
`
			}
		}
		runResult = "请输入SQL语句！"
	}

	/////////////////////////////////////////////////////////////////////////////////////
	data := map[string]template.HTML{
		"database_name": template.HTML(database_name),
		"tablename":     template.HTML(tablename),
		"sqlstr":        template.HTML(sqlstr),
		"runResult":     template.HTML(runResult),
	}

	mb, err := template.ParseFiles(AppPath + "/templates/run.html")
	if err != nil {
		fun.Log(0, "run.go:"+fmt.Sprintf("%s", err))
	}

	err = mb.Execute(w, data)
	if err != nil {
		fun.Log(0, "run.go:"+fmt.Sprintf("%s", err))
	}
}
