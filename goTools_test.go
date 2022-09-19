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
