package goSmallTools

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// 驼峰转下划线
func HumpToUnderline(str string) string {
	if str == "ID" {
		return "id"
	}
	a := str
	alist := strings.Split(a, "ID")
	a = strings.Join(alist, "_id")
	reg := regexp.MustCompile(`[A-Z]`)
	result := reg.FindAllStringSubmatch(a, -1)
	for i := 0; i < len(result); i++ {
		if len(result[i]) > 0 {
			b := strings.Split(a, result[i][0])
			a = strings.Join(b, "_"+strings.ToLower(result[i][0]))
		}
	}
	return a
}

// 下划线转驼峰
func UnderlineToHump(str string) string {
	if str == "id" {
		return "ID"
	}
	a := str
	temp := strings.Split(a, "_")
	var upperStr string
	for i := 0; i < len(temp); i++ {
		b := []rune(temp[i])
		if i != 0 {
			for j := 0; j < len(b); j++ {
				if j == 0 {
					b[j] -= 32
					upperStr += string(b[j])
				} else {
					if string(b[j-1]) == "I" && string(b[j]) == "d" {
						b[j] -= 32
					}
					upperStr += string(b[j])
				}
			}
		} else {
			upperStr += string(temp[i])
		}
	}
	return upperStr
}

// 处理字典数组转为MySQL where语句
func ParameterToWhere(str string, cst *time.Location, timeFormat string) string {
	plist := gjson.Parse(str)
	where := ""
	if plist.Exists() && plist.Type.String() == "JSON" {
		for _, v := range plist.Array() {
			key := gjson.Get(v.String(), "key")
			isor := gjson.Get(v.String(), "isor")
			if where != "" {
				if isor.Exists() && isor.Bool() {
					where += " OR "
				} else {
					where += " AND "
				}
			}
			where += "`" + HumpToUnderline(key.String()) + "`"
			isNull := gjson.Get(v.String(), "isNull")
			if isNull.Exists() {
				if isNull.Bool() {
					where += " IS NULL"
				} else {
					where += " IS NOT NULL"
				}
				continue
			}
			value := gjson.Get(v.String(), "value")
			switch value.Type.String() {
			case "Number", "String":
				isFuzzy := gjson.Get(v.String(), "isFuzzy")
				if isFuzzy.Exists() && isFuzzy.Bool() {
					where += " LIKE '%" + value.String() + "%'"
				} else {
					val := value.String()
					if key.String() == "language" {
						val = strings.ToLower(val)
						vs := strings.Split(val, "-")
						if len(vs) > 1 {
							val = strings.Join(vs, "_")
						} else {
							vs = strings.Split(val, "_")
							if len(vs) == 1 && val == "zh" {
								val = val + "_cn"
							}
						}
					}
					if key.String() == "timeZone" {
						vs := strings.Split(val, ":")
						if len(vs) > 1 {
							hour, err := strconv.ParseFloat(vs[0], 64)
							if err == nil {
								min, err := strconv.ParseFloat(vs[1], 64)
								if err == nil {
									min = min / 60
									time := hour + min
									val = fmt.Sprintf("%f", time)
								}
							}
						}
					}
					where += "='" + val + "'"
				}
			case "JSON":
				if key.String() == "freeTime" {
					where += "='" + value.String() + "'"
					continue
				}
				var strs string
				for _, s := range value.String() {
					strs = fmt.Sprintf("%c", s)
					break
				}
				if strs == "[" {
					if len(value.Array()) == 2 {
						t1 := time.Unix(0, value.Array()[0].Int()*1000000).In(cst).Format(timeFormat)
						t2 := time.Unix(0, value.Array()[1].Int()*1000000).In(cst).Format(timeFormat)
						where += " BETWEEN '" + t1 + "' AND '" + t2 + "'"
					}
				}
			}
		}
	}
	return where
}

// 检查电子邮件地址
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
//	code	国家码
//	phone	电话号码
//
// 美国:001，加拿大:1
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

// 时间戳转时间字符串
//
//	ts		int64		时间戳10-19位
//	cst		*time.Location	时区
//	timeFormat	string		时间格式'2006-01-02 15:04:05'
func TimeStampToTimeStr(ts int64, cst *time.Location, timeFormat string) string {
	timeStamp := ts
	tsStr := strconv.Itoa(int(ts))
	return tsTostr(int64(timeStamp), len(tsStr), cst, timeFormat)
}

// 时间戳字符串转时间字符串
//
//	ts		string		时间戳字符串10-19位
//	cst		*time.Location	时区
//	timeFormat	string		时间格式'2006-01-02 15:04:05'
func TimeStampStrToTimeStr(tsStr string, cst *time.Location, timeFormat string) (string, error) {
	timeStamp, err := strconv.Atoi(tsStr)
	if err != nil {
		return "", err
	}
	return tsTostr(int64(timeStamp), len(tsStr), cst, timeFormat), nil
}

func tsTostr(ts int64, tsLen int, cst *time.Location, timeFormat string) string {
	timeStamp := ts
	if tsLen < 19 {
		for i := 0; i < 19-tsLen; i++ {
			timeStamp *= 10
		}
	}
	return time.Unix(0, int64(timeStamp)).In(cst).Format(timeFormat)
}
