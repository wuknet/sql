package fun

import (
	"context"
	"net/http"
	"github.com/go-session/session"
)

//此函数放在使用页面输出的最前面
func Session_start(w http.ResponseWriter, r *http.Request) session.Store {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return nil
	} else {
	    return store
	}
}

func Session_save(store session.Store,savename string,content string) bool {
    store.Set(savename, content)
	err := store.Save()
	if err != nil {
		return false
	}
	return true
}

func Session_read(store session.Store,readname string) string {
	foo, ok := store.Get(readname)
	if ok {
		return foo.(string)
	} else {
		return ""
	}
}