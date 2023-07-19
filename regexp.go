package goSmallTools

import (
	"fmt"
	"regexp"
)

// 正则表达式
//
// Regular expression
var regexpString string = ``

// 默认正则表达式
//
// Default regular expression
const defaultRegexpString string = `([\x00-\x08\x0B-\x0C\x0E-\x1F\x7F]+)|(?:')|(?:\\)|(?:--)|(\b(\\x|�|select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`

const (
	Default     = iota //	默认
	Int                //	整数
	Float              //	浮点数
	OnlyLetter         //	只有字母
	OnlyChinese        //	只有中文
	Email              //	邮箱
	Url                //	网址
	Account            //	账号,首字母为字母, 5-16位
	Password           //	密码, 6-64位, 可以包含字母, 数字, 特殊字符(!@#$%^&*)
)

// 根据类型设置正则表达式
//
//	regexpType	int	正则表达式类型
func SetRegexpType(regexpType int) {
	switch regexpType {
	case Default:
		regexpString = ""
	case Int:
		regexpString = `^-?\d+$`
	case Float:
		regexpString = `^(-?\d+)(\.\d+)?$`
	case OnlyLetter:
		regexpString = `^[a-zA-Z]+$`
	case OnlyChinese:
		regexpString = "^[\u4e00-\u9fa5]+$"
	case Email:
		regexpString = `^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+`
	case Url:
		regexpString = `^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`
	case Account:
		regexpString = `^[a-zA-Z][a-zA-Z0-9_]{4,15}$`
	case Password:
		regexpString = `^[a-zA-Z0-9!@#$%^&*]{6,64}$`
	}
}

// 检查字符串是否满足正则表达式
//
//	str	string	需要检查的字符串
//	return	bool	是否满足正则表达式
func CheckRegexpType(str string) bool {
	var re *regexp.Regexp
	if regexpString == "" {
		re = regexp.MustCompile(defaultRegexpString)
	} else {
		re = regexp.MustCompile(regexpString)
	}
	return re.MatchString(str)
}

// 设置正则表达式
//
//	regexpStr	string	正则表达式
//							""时修改为默认值
func SetRegexpString(regexpStr string) {
	regexpString = regexpStr
}

// 检查字符串中是否有linux下的非法字符和SQL注入
//
//	str	string	需要检查的字符串
//	return	bool	是否有非法字符
func CheckString(str string) bool {
	var re *regexp.Regexp
	if regexpString == "" {
		re = regexp.MustCompile(defaultRegexpString)
	} else {
		re = regexp.MustCompile(regexpString)
	}
	strList := re.FindAllString(str, -1)
	if len(strList) > 0 {
		fmt.Println("   ", str, "=>", strList, len(strList))
	}
	return len(strList) > 0
}

// 检查字符串中是否有linux下的非法字符和SQL注入并返回非法字符数组
//
//	str	string		需要检查的字符串
//	return	[]string	非法字符数组
func CheckStringReturnList(str string) []string {
	var re *regexp.Regexp
	if regexpString == "" {
		re = regexp.MustCompile(defaultRegexpString)
	} else {
		re = regexp.MustCompile(regexpString)
	}
	fmt.Println("Regexp:", re)
	strList := re.FindAllString(str, -1)
	return strList
}

// 替换字符串中的非法字符
//
//	str	string	需要替换的字符串
//	repl	string	替换的字符
//	return	string	替换后的字符串
func ReplaceAllString(str string, repl string) string {
	var re *regexp.Regexp
	if regexpString == "" {
		re = regexp.MustCompile(defaultRegexpString)
	} else {
		re = regexp.MustCompile(regexpString)
	}
	reStr := re.ReplaceAllString(str, repl)
	return reStr
}

// 检查电子邮件地址
//
//	mail	string	电子邮件地址
//	return	error 	错误信息
func CheckEmail(mail string) error {
	// regexpStr := `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	// regexpStr := `^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	regexpStr := `^([\w\.\_\-]{2,245})@(\w{1,}).([a-z]{2,4})$`
	re, err := regexp.Compile(regexpStr)
	if err != nil {
		return err
	}
	if !re.MatchString(mail) {
		return fmt.Errorf("err mail path")
	}
	return nil
}

// 验证手机号码
//
//	code	string	国家码
//	phone	string	电话号码,注:美国:001，加拿大:1
//	return	error 	错误信息
func CheckPhone(code string, phone string) error {
	regexpStr := phoneReg[code]
	re, err := regexp.Compile(regexpStr)
	if err != nil {
		return err
	}
	if !re.MatchString(phone) {
		return fmt.Errorf("err phone path")
	}
	return nil

}
