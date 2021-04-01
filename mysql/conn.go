package mysql

import (
	//"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conndb(addr string, username string, password string, database_name string) (*sql.DB, bool, string) {
	//DB, _ := sql.Open("mysql", "bbsgood:bbsgood!#%&(x@tcp(localhost:3306)/bbsgood?charset=utf8")
	sqlstr := username + ":" + password + "@tcp(" + addr + ")/" + database_name + "?charset=utf8"
	DB, _ := sql.Open("mysql", sqlstr)
	//设置数据库最大超时时间
	//DB.SetConnMaxLifetime(300 * time.Second)
	//打开的最大连接数
	//DB.SetMaxOpenConns(0)
	//设置上数据库最大闲置连接数
	//DB.SetMaxIdleConns(100)

	//验证连接
	if err := DB.Ping(); err != nil {
		return DB, false, sqlstr
	} else {
		return DB, true, sqlstr
	}
}
