package goSmallTools

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestHtU(t *testing.T) {
	t.Log(HumpToUnderline("ID"))
	t.Log(HumpToUnderline("userID"))
	t.Log(HumpToUnderline("SuperAdmin"))
}
func TestUtH(t *testing.T) {
	t.Log(UnderlineToHump("id"))
	t.Log(UnderlineToHump("mail_address"))
}
func TestParameter(t *testing.T) {
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Error("时区文件加载失败:", err)
		cstSh = time.FixedZone("CST", 8*3600)
	}
	timeFormat := "2006-01-02 15:04:05"
	str := `[{"key":"id","value":"1"},{"key":"userName","value":"test","isor":true,"isFuzzy":true},{"key":"time","value":[1662369862,1662369877]},{"key":"enable","isNull":true}]`
	pStr := ParameterToWhere(str, cstSh, timeFormat)
	t.Log(pStr)
}
func TestCheckEmail(t *testing.T) {
	pStr := CheckEmail("test@test.com")
	t.Log(pStr)
}
func TestCheckPhone(t *testing.T) {
	pStr := CheckPhone("86", "13433335555")
	t.Log(pStr)
}

func TestTimeStampToTimeStr(t *testing.T) {
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Error("时区文件加载失败:", err)
		cstSh = time.FixedZone("CST", 8*3600)
	}
	timeFormat := "2006-01-02 15:04:05"
	tn := time.Now()
	fmt.Println(tn.Unix(), tn.In(cstSh).Format(timeFormat))
	timeStr := TimeStampToTimeStr(tn.Unix(), cstSh, timeFormat)
	fmt.Println(timeStr)
}

func TestTimeStampStrToTimeStr(t *testing.T) {
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Error("时区文件加载失败:", err)
		cstSh = time.FixedZone("CST", 8*3600)
	}
	timeFormat := "2006-01-02 15:04:05"
	tn := time.Now()
	fmt.Println(tn.Unix(), tn.In(cstSh).Format(timeFormat))
	timeStr, err := TimeStampStrToTimeStr(strconv.Itoa(int(tn.Unix())), cstSh, timeFormat)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	fmt.Println(timeStr)
}

func TestCheckString(t *testing.T) {
	a := "测a试jsdij/[]$%^&*()_+{}|:<>?~`select 1234567890-=\\';lkjhgfdsazxcvbnm,./!@#$%^&*()_+{}|:<>?~`"
	b := "\\x dfas '%dfs%' -- sad-ds"
	CheckString("\xAA")
	CheckString("\x12b")
	CheckString("\x12")
	CheckString("\x9A")
	CheckString("\x12")
	CheckString("\xAA\x12b\x12\x9A\x12")
	CheckString("asdfghjkl; ")
	CheckString("'sa'")
	CheckString("\\sa'")
	CheckString(a)
	fmt.Println(">>>", ReplaceAllString(a, ""))
	CheckString(b)
	fmt.Println(">>>", ReplaceAllString(b, ""))
	SetRegexpString(`(?:a)|(?:')|(?:\\)|(?:--)`)
	if errArr := CheckStringReturnList(a); len(errArr) > 0 {
		fmt.Println(a, "非法字符:", errArr)
	}
	if errArr := CheckStringReturnList(b); len(errArr) > 0 {
		fmt.Println(b, "非法字符:", errArr)
	}
}

func TestSetRegexpType(t *testing.T) {
	SetRegexpType(Int)
	a := "123"
	errArr := CheckRegexpType(a)
	fmt.Println("Int:", a, "=>", errArr)
	SetRegexpType(Float)
	a = "123.456"
	errArr = CheckRegexpType(a)
	fmt.Println("Float:", a, "=>", errArr)
	SetRegexpType(OnlyChinese)
	a = "测试"
	errArr = CheckRegexpType(a)
	fmt.Println("OnlyChinese:", a, "=>", errArr)
	SetRegexpType(OnlyLetter)
	a = "asdf"
	errArr = CheckRegexpType(a)
	fmt.Println("OnlyLetter:", a, "=>", errArr)
	SetRegexpType(Email)
	a = "test@test.com"
	errArr = CheckRegexpType(a)
	fmt.Println("Email:", a, "=>", errArr)
	SetRegexpType(Url)
	a = "https://www.test.com"
	errArr = CheckRegexpType(a)
	fmt.Println("Url:", a, "=>", errArr)
	SetRegexpType(Account)
	a = "Dv2_d"
	errArr = CheckRegexpType(a)
	fmt.Println("Account:", a, "=>", errArr)
	SetRegexpType(Password)
	a = "%xa6tGKkKd9Sz&b!8*6GjL!NP$tk^2DdSxeBE#jCF@5rH9VwB6QExCTt%oUwskbq"
	errArr = CheckRegexpType(a)
	fmt.Println("Password:", a, "=>", errArr)
}
