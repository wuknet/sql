package fun

import (
	//"fmt"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	//github.com/gorilla/mux
	//生成唯一ID
	//生成唯一ID
	//编码转换

	//"github.com/alexedwards/scs"
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

//温度转换成度，开尔文转摄氏度
func WeatherTo(temp string) string {
	temp_f, _ := strconv.ParseFloat(temp, 64)
	temp_f -= 273.15
	return strconv.FormatFloat(temp_f, 'f', 2, 64)
}

//写入文件
func WriteFile(filepath string, content string) bool {
	AppPath := GetAppPath() //读取当前目录路径
	filepath = AppPath + "/" + filepath
	f, err := os.Create(filepath)
	defer f.Close()
	if err != nil {
		return false
	} else {
		_, err = f.Write([]byte(content))
		if err != nil {
			return false
		} else {
			return true
		}
	}
}

//读取文件内容
func ReadFile(filepath string) (string, error) {
	AppPath := GetAppPath() //读取当前目录路径
	c, err := ioutil.ReadFile(AppPath + "/" + filepath)
	if err != nil {
		panic(err)
	}
	return string(c), err
}

//中文转Unicode
func Unicode(sText string) string {
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	return textUnquoted
}

//获取当前url
func Geturl(r *http.Request) string {
	scheme := "https://"
	///if r.TLS != nil {
	//	scheme = "https://"
	//}
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

//生产二维图形的url,可以用 <img src=url>显示
func Qrcode(w http.ResponseWriter, r *http.Request, url string) string {
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err == nil {
		dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(png))
		return dataURI
	} else {
		return ""
	}
}

//通过域名返回IP
func Domainip(domain string) string {
	ns, err := net.LookupHost(domain)
	if err != nil {
		return ""
	} else {
		for _, ip := range ns {
			return ip
		}
		return ""
	}
}

//html原码输出
func Htmlencode(mystring string) string {
	mystring = strings.Replace(mystring, "\n", "<br />", -1) //换行
	//mystring = strings.Replace(mystring, "\r" , "<br />", -1)//回车
	return mystring
}

//文字过滤
func FilterWords(mystring string, IllegalWords string) string {
	wordlist := strings.Split(IllegalWords, "|")
	for i := 0; i < len(wordlist); i++ {
		wordlist2 := strings.Split(wordlist[i], ",")
		mystring = strings.Replace(mystring, wordlist2[0], wordlist2[1], -1)
	}
	return mystring
}

//获取任意长度的随机数
func Sjm(len int) string {
	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

//按数字建立二层目录
func Numdir(num string) (string, string) {
	i, _ := strconv.Atoi(num)
	s := 100000000 / 1000
	x := 0
	y := 0
	ss := i / s
	if ss <= 0 {
		x = 0
	} else {
		x = ss
	}
	y = i / 255
	return strconv.Itoa(x), strconv.Itoa(y)
	//fmt.Println(i + x + y)
	//fmt.Println("i=" + strconv.Itoa(i) + " s=" + strconv.Itoa(s) + " ss=" + strconv.Itoa(ss) + " x=" + strconv.Itoa(x) + " y=" + strconv.Itoa(y))
}

//判断浏览器版本是不是中文
func IsZh(r *http.Request) bool {
	if strings.Contains(r.Header.Get("Accept-Language"), "zh") {
		return true
	} else {
		return false
	}
}

// 判断访问的是不是手机设备
func Isphone(r *http.Request) bool {
	useragent := r.UserAgent()
	keywords := []string{"Android", "iPhone", "Windows Phone"}
	for i := range keywords {
		if strings.Contains(useragent, keywords[i]) {
			return true
		}
	}
	return false
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFile(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//判断文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

//判断字符串是否全是数字
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func Get_int(r *http.Request, varname string) int {
	// 解析客户端请求的信息
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	revar := r.Form.Get(varname)
	if revar != "" {
		revar2, err := strconv.Atoi(revar)
		if err != nil {
			return 0
		} else {
			return revar2
		}
	} else {
		return 0
	}
}

func Get_str(r *http.Request, varname string) string {
	// 解析客户端请求的信息
	if err := r.ParseForm(); err != nil {
		return ""
	}
	revar := r.Form.Get(varname)
	//revar = html.EscapeString(revar)
	return revar
}

//是不是电子邮箱格式
func IsMail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//是不是手机号
func IsMobile(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
//       负数 - 在从字符串结尾的指定位置开始
//       0 - 在字符串中的第一个字符处开始
//length:正数 - 从 start 参数所在的位置返回
//       负数 - 从字符串末端返回
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//分隔符的多内容替换
func ReplaceX(LabelContel string, replaceKeys string) string {
	if replaceKeys != "" {
		replaceKeysArr := strings.Split(replaceKeys, "#")
		for x := 0; x < len(replaceKeysArr); x++ {
			replaceKeysArr2 := strings.Split(replaceKeysArr[x], "$")
			LabelContel = strings.Replace(LabelContel, replaceKeysArr2[0], replaceKeysArr2[1], -1)
		}
	}
	return LabelContel
}

func Md532(str string) string {
	md5 := md5.New()
	md5.Write([]byte(str))
	MD5Str := hex.EncodeToString(md5.Sum(nil))
	return MD5Str
}

func GetAppPath() string { //获取当前实际路径，在服务状态下
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

//判断是否为utf-8编码
func isUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

//删除数组复重内容
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

//获取字符串中间的内容
func Between(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

func ReadCookie(w http.ResponseWriter, req *http.Request, cookiename string) string {
	// read cookie
	cookie, err := req.Cookie(cookiename)
	if err == nil {
		cookievalue, _ := url.QueryUnescape(cookie.Value) //解码
		return cookievalue
	} else {
		return ""
	}
}

func WriteCookie(w http.ResponseWriter, req *http.Request, cookiename string, cookievalue string, MaxAgeValue int) {
	if MaxAgeValue == 0 {
		MaxAgeValue = 86400 * 365 * 10
	}
	cookievalue = url.QueryEscape(cookievalue) //加码
	cookie := http.Cookie{
		Name:  cookiename,
		Value: cookievalue,
		//Domain:   "887d.com",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   MaxAgeValue,
	}
	http.SetCookie(w, &cookie)
}

func DeleteCookie(w http.ResponseWriter, req *http.Request, cookiename string) {
	cookie := http.Cookie{Name: cookiename, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
}

//设置为某种语言
func SetLang(w http.ResponseWriter, req *http.Request, lang string) {
	WriteCookie(w, req, "lang", lang, 0)
}

//读取使用什么语言
func Lang(w http.ResponseWriter, req *http.Request) string {
	/*
		lang := ReadCookie(w, req, "lang")
		if lang == "" {
			lang = "zh-CN"
		}
	*/
	if strings.Contains(req.Host, "en.") || strings.Contains(req.Host, ":8988") {
		return "en"
	} else {
		return "zh-CN"
	}
}

//中间用##号隔开，前面是中文输出，后面是英文输出
func Fprintf(w http.ResponseWriter, req *http.Request, text string) {
	if strings.Contains(text, "{{") {
		texts := strings.Split(text, "{{")
		for i := 1; i < len(texts); i++ {
			texts2 := strings.Split(texts[i], "}}")
			retext := "{{" + texts2[0] + "}}"
			texts3 := strings.Split(texts2[0], "##")

			if Lang(w, req) == "zh-CN" {
				text = strings.Replace(text, retext, texts3[0], -1)
			} else {
				if len(texts) == 1 {
					text = strings.Replace(text, retext, texts3[0], -1)
				} else {
					text = strings.Replace(text, retext, texts3[1], -1)
				}
			}
		}
		fmt.Fprintf(w, text)
	} else {
		texts := strings.Split(text, "##")
		if Lang(w, req) == "zh-CN" {
			fmt.Fprintf(w, texts[0])
		} else {
			if len(texts) == 1 {
				fmt.Fprintf(w, texts[0])
			} else {
				fmt.Fprintf(w, texts[1])
			}
		}
	}
}
