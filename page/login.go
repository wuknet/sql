package page

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"html/template"

	"../fun"
	"../mysql"
	//"github.com/julienschmidt/httprouter"
)

//////////////////////////////////////////////////////
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if fun.ReadCookie(w, r, "addr") != "" && fun.ReadCookie(w, r, "username") != "" && fun.ReadCookie(w, r, "username") != "" {
		http.Redirect(w, r, "/datalist", http.StatusFound)
	}

	AppPath := fun.GetAppPath() //读取当前目录路径

	save := fun.Get_str(r, "save")
	addr := fun.Get_str(r, "addr")
	username := fun.Get_str(r, "username")
	username = strings.Trim(username, "	")
	password := fun.Get_str(r, "password")
	database_name := fun.Get_str(r, "database_name")

	if save == "1" {
		db, linkzt, sqlstr := mysql.Conndb(addr, username, password, database_name)
		defer db.Close()
		if linkzt == false {
			fmt.Fprintf(w, "c|"+username+"|")
			fun.Fprintf(w, r, "数据库连接失败##Database link lost!")
			fmt.Fprintf(w, sqlstr)
		} else {
			fun.WriteCookie(w, r, "addr", addr, 0)
			fun.WriteCookie(w, r, "username", username, 0)
			fun.WriteCookie(w, r, "password", password, 0)
			fun.WriteCookie(w, r, "database_name", database_name, 0)
			http.Redirect(w, r, "/datalist", http.StatusFound)
		}
	} else {

		/////////////////////////////////////////////////////////////////////////////////////
		data := map[string]template.HTML{
			//"save":        template.HTML(save),
		}

		mb, err := template.ParseFiles(AppPath + "/templates/login.html")
		if err != nil {
			fun.Log(0, "login.go:"+fmt.Sprintf("%s", err))
		}

		err = mb.Execute(w, data)
		if err != nil {
			fun.Log(0, "login.go:"+fmt.Sprintf("%s", err))
		}
	}
}
