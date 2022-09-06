package goTools

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
