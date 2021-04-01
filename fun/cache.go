package fun

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
	//"github.com/tidwall/gjson"
)

//////////////////////////////////////////////////////
type Datajg struct {
	title      string
	content    interface{}
	expire     int
	hits       int //访问次数
	createtime string
}

var (
	CacheVal    []Datajg
	SinCacheval Datajg
)

//////////////////////////////////////////////////////
//检测缓存是否存在
func CheckCache(title string) (int, bool) {
	i := 0
	Isexist := false
	for i = 0; i < len(CacheVal); i++ {
		if CacheVal[i].title == title {
			Isexist = true
			break
		}
	}
	return i, Isexist
}

//添加或设置缓存
func SetCache(title string, content interface{}, expire int) {
	SinCacheval.title = title
	SinCacheval.content = content
	SinCacheval.expire = expire

	index, Isexist := CheckCache(title) //检测缓存是否存在，如果存在返回索引数据
	if Isexist == false {               //如果不存在
		SinCacheval.title = title
		SinCacheval.content = content
		if expire != 0 { //如果是0不过期
			expire_unix, _ := strconv.Atoi(fmt.Sprintf("%v", time.Now().Unix()))
			expire = expire_unix + expire
		}
		SinCacheval.expire = expire
		SinCacheval.hits = 1
		SinCacheval.createtime = Gettime()
		CacheVal = append(CacheVal, SinCacheval)
	} else { //存在则编辑
		CacheVal[index].content = content
		CacheVal[index].expire = expire
	}
}

//获取缓存内容
func GetCache(title string) (interface{}, bool) {
	DelExpireCache() //删除过期的缓存
	index, Isexist := CheckCache(title)
	//fmt.Println(len(CacheVal))
	if Isexist == false {
		return "", false
	} else {
		CacheVal[index].hits += 1
		return CacheVal[index].content, true
	}
}

//删除缓存项目
func DelCache(title string) {
	for i := 0; i < len(CacheVal); i++ {
		if CacheVal[i].title == title {
			CacheVal = append(CacheVal[:i], CacheVal[i+1:]...)
			break
		}
	}
}

//删除所有缓存
func DelAllCache() {
	CacheVal = CacheVal[0:0]
}

//删除所有过期缓存项目
func DelExpireCache() {
	NowTime, _ := strconv.Atoi(fmt.Sprintf("%v", time.Now().Unix()))
	for i := 0; i < len(CacheVal); i++ {
		if NowTime-CacheVal[i].expire >= 0 && CacheVal[i].expire != 0 {
			CacheVal = append(CacheVal[:i], CacheVal[i+1:]...)
			i--
		}
	}
}

//获取单个缓存大小
func CacheSize(title string) int {
	index, Isexist := CheckCache(title)
	if Isexist == false {
		return 0
	} else {
		return int(unsafe.Sizeof(CacheVal[index].content))
	}
}

//获取总缓存大小
func AllCacheSize() int {
	znum := 0
	for i := 0; i < len(CacheVal); i++ {
		znum = znum + int(unsafe.Sizeof(CacheVal[i]))
	}
	return znum
}

func CacheList() string {
	list := `<table><tr><td>缓存名称</td><td>访问量</td><td>到期时间</td><td>创建时间</td><td>操作</td></tr>`
	for i := 0; i < len(CacheVal); i++ {
		list += `<tr><td>` + CacheVal[i].title + `</td><td>` + strconv.Itoa(CacheVal[i].hits) + `</td><td>` + strconv.Itoa(CacheVal[i].expire) + `</td><td>` + CacheVal[i].createtime + `</td><td><a href="?cls=` + CacheVal[i].title + `">清除</a></td></tr>`
	}
	list += `</table>`
	return list
}
