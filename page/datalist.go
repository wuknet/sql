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
func Datalist(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//t := time.Now() //开始计算执行时间
	AppPath := fun.GetAppPath() //读取当前目录路径

	addr := fun.ReadCookie(w, r, "addr")
	username := fun.ReadCookie(w, r, "username")
	password := fun.ReadCookie(w, r, "password")
	database_name := fun.ReadCookie(w, r, "database_name")
	db, linkzt, _ := mysql.Conndb(addr, username, password, database_name)
	defer db.Close()
	if linkzt == false {
		fmt.Fprintf(w, "Database link lost!")
		return
	} else {
		datalist := ""
		rows, _ := mysql.Query(db, "SHOW DATABASES")
		for i := 0; i < len(rows); i++ {
			datalist += `
			<a href="/tablelist/?database_name=` + rows[i][0] + `" class="databox">` + rows[i][0] + `</a>`
		}

		/////////////////////////////////////////////////////////////////////////////////////
		data := map[string]template.HTML{
			"datalist": template.HTML(datalist),
		}

		mb, err := template.ParseFiles(AppPath + "/templates/datalist.html")
		if err != nil {
			//fun.Log(0, "login.go:"+fmt.Sprintf("%s", err))
		}

		err = mb.Execute(w, data)
		if err != nil {
			//fun.Log(0, "login.go:"+fmt.Sprintf("%s", err))
		}
	}
}
