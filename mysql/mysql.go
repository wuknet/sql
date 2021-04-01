package mysql

import (
	//"fmt"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//统计行数
func Count(db *sql.DB, filed, table_where string) uint64 {
	var id uint64
	sql := "select count(" + filed + ") as znum from " + table_where
	idarr, err := QueryOne(db, sql)
	if err == nil {
		if len(idarr) > 0 {
			id, _ = strconv.ParseUint(idarr[0], 0, 64)
			return id
		} else {
			return 0
		}
	} else {
		return 0
	}
}

//判断表中某个int字段的最大值
func MaxID(db *sql.DB, filed string, table_where string) uint64 {
	var id uint64
	sql := "select max(" + filed + ") as maxid from " + table_where
	//fmt.Println(sql)
	idarr, err := QueryOne(db, sql)
	if err == nil {
		if len(idarr) > 0 {
			id, _ = strconv.ParseUint(idarr[0], 0, 64)
			return id
		} else {
			return 0
		}
	} else {
		return 0
	}
}

//判断表中某个int字段累加总值
func SumID(db *sql.DB, filed string, table_where string) uint64 {
	var id uint64
	sql := "select sum(" + filed + ") as maxid from " + table_where
	//fmt.Println(sql)
	idarr, err := QueryOne(db, sql)
	if err == nil {
		if len(idarr) > 0 {
			id, _ = strconv.ParseUint(idarr[0], 0, 64)
			return id
		} else {
			return 0
		}
	} else {
		return 0
	}
}

//判断单选数据是否存在
func IsValid(db *sql.DB, sql string) bool {
	arr, err := QueryOne(db, sql)
	if err == nil {
		if len(arr) > 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//查询单一行数据
func QueryOne(db *sql.DB, sql string) ([]string, error) {
	//t1, _ := strconv.Atoi(fmt.Sprintf("%v", time.Now().UnixNano()/1e6))
	var arr []string
	rows, err := db.Query(sql)
	defer rows.Close()

	if err != nil {
		fmt.Println(sql + ":" + fmt.Sprintf("%s", err))
		return arr, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return arr, err
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return arr, err
		}
		for _, col := range values {
			if col != nil {
				arr = append(arr, string(col.([]byte)))
			} else {
				arr = append(arr, "")
			}
		}
	}
	//t2, _ := strconv.Atoi(fmt.Sprintf("%v", time.Now().UnixNano()/1e6))
	//t3 := strconv.Itoa(t2 - t1)
	//db.Exec("insert into sql_log(sqlstr,speed)values('" + sql + "'," + t3 + ")")
	return arr, err
}

func Query(db *sql.DB, sql string) ([][]string, error) {
	//t1, _ := strconv.Atoi(fmt.Sprintf("%v", time.Now().UnixNano()/1e6))

	var arr [][]string

	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(sql + ":" + fmt.Sprintf("%s", err))
		return arr, err
	}
	defer rows.Close()
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			fmt.Println(sql + ":" + fmt.Sprintf("%s", err))
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return arr, err
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		_ = rows.Scan(scanArgs...)
		var tmpArr []string
		for _, col := range values {
			if col != nil {
				tmpArr = append(tmpArr, string(col.([]byte)))
			} else {
				tmpArr = append(tmpArr, "")
			}
		}
		arr = append(arr, tmpArr)
	}
	//t2, _ := strconv.Atoi(fmt.Sprintf("%v", time.Now().UnixNano()/1e6))
	//t3 := strconv.Itoa(t2 - t1)
	//db.Exec("insert into sql_log(sqlstr,speed)values('" + sql + "'," + t3 + ")")
	return arr, err
}

//返回最新插入的ID和行数
func Insert(db *sql.DB, sql string) (int64, error) {
	result, err := db.Exec(sql)
	if err != nil {
		return 0, err
	} else {
		lastInsertID, err := result.LastInsertId() //插入数据的主键id
		//rowsaffected,_ := result.RowsAffected()  //影响行数
		return lastInsertID, err
	}
}

//更新数据库
func Update(db *sql.DB, sql string) (int64, error) {
	//result,err := DB.Exec("UPDATE users set age=? where id=?","30",3)
	stmt, err := db.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	} else {
		result, err := stmt.Exec()
		if err != nil {
			return 0, err
		} else {
			rowsaffected, err := result.RowsAffected()
			if err != nil {
				return 0, err
			} else {
				return rowsaffected, err
			}
		}
	}
}

//删除数据
func Delete(db *sql.DB, sql string) (int64, error) {
	//result,err := db.Exec("delete from users where id=?",1)
	result, err := db.Exec(sql)
	if err != nil {
		return 0, err
	} else {
		rowsaffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		} else {
			return rowsaffected, err
		}
	}
}

//利用mysqldump备份
func BackupMySqlDb(host, port, user, password, databaseName, tableName, sqlPath string) (error, string) {
	var cmd *exec.Cmd

	//路径  /usr/bin/mysqldump
	if tableName == "" {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName)
	} else {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName, tableName)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err, ""
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err, ""
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	now := time.Now().Format("20060102150405")
	var backupPath string
	if tableName == "" {
		backupPath = sqlPath + databaseName + "_" + now + ".sql"
	} else {
		backupPath = sqlPath + databaseName + "_" + tableName + "_" + now + ".sql"
	}
	err = ioutil.WriteFile(backupPath, bytes, 0644)

	if err != nil {
		log.Println(err)
		return err, ""
	}
	return nil, backupPath
}

func RestartMysql() (bool, error) {
	cmd := exec.Command("systemctl", "restart", "mysqld")
	err := cmd.Run()
	if err != nil {
		fmt.Println("mysql重启失败")
		fmt.Println(err)
		return false, err
	} else {
		fmt.Println("mysql重启成功")
		return true, err
	}
}
