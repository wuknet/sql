package page

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"

	"html/template"

	"../fun"
	//"github.com/julienschmidt/httprouter"
)

//////////////////////////////////////////////////////
func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if fun.ReadCookie(w, r, "addr") != "" && fun.ReadCookie(w, r, "username") != "" && fun.ReadCookie(w, r, "username") != "" {
		//http.Redirect(w, r, "/datalist", http.StatusFound)
	}

	AppPath := fun.GetAppPath() //读取当前目录路径

	/////////////////////////////////////////////////////////////////////////////////////
	data := map[string]template.HTML{
		//"save":        template.HTML(save),
	}

	mb, err := template.ParseFiles(AppPath + "/templates/index.html")
	if err != nil {
		fun.Log(0, "index.go:"+fmt.Sprintf("%s", err))
	}

	err = mb.Execute(w, data)
	if err != nil {
		fun.Log(0, "index.go:"+fmt.Sprintf("%s", err))
	}
}
