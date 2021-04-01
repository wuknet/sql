package fun

import (
	//"fmt"
	"bytes"
	"math"
	"strconv"
	"time"
)

//获取相差分钟时间
func GetMoDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 60
		return hour
	} else {
		return hour
	}
}

//按时间生成一个唯一数字字符串
//按时间到秒加上6位随机数
func GetUniqueTimeString() string {
	return FormatTime(Gettime(), "20060102150405") + Sjm(6)
}

//获取按日期为目录的地址，按年/月-日/分配
func GetDataPath(mytime string) string {
	if mytime == "" {
		mytime = Gettime()
	}
	return FormatTime(mytime, "2006/01-02")
}

func Gettime() string {
	mytime := time.Now().Format("2006-01-02 15:04:05")
	return mytime
}

//设置时间格式，只有日期
func SettimeDate(mytime string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", mytime)
	return t.Format("2006-01-02")
}

//自定义返回日期格式 如fromattime参数为2006/01-02
func FormatTime(mytime string, fromattime string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", mytime)
	return t.Format(fromattime)
}

//把时间转换成时间戳
func TimetoTimeC(mytime string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")                      //设置时区
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", mytime, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix()
}

//把时间戳转成时间
func TimeCtoTime(TimeC int64) string {
	tm := time.Unix(TimeC, 0)
	mytime := tm.Format("2006-01-02 15:04:05")
	return mytime
}

//获取把时间变成几分钟前，几天前
func StrTime(atime int64) string {
	var byTime = []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"年前", "天前", "小时前", "分钟前", "秒钟前"}
	now := time.Now().Unix()
	ct := now - atime
	if ct < 0 {
		return "刚刚"
	}
	var res string
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i]
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = MergeString(tempStr, unit[i]) //此处调用了一个我自己封装的字符串拼接的函数（你也可以自己实现）
		}
		break //我想要的形式是精确到最大单位，即："2天前"这种形式，如果想要"2天12小时36分钟48秒前"这种形式，把此处break去掉，然后把字符串拼接调整下即可（别问我怎么调整，这如果都不会我也是无语）
	}
	return res
}
func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}
	return buffer.String()
}

//
