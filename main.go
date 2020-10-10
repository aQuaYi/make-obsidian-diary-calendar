package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const calendarName = "Calendar.md"
const format = "2006-01-02"
const path = "./Diary/"

const maxYear = 3000
const maxDecade = maxYear / 10

var days = [maxYear][13][32]bool{}
var monthes = [maxYear][13]bool{}
var decades = [maxDecade]bool{}

func main() {
	counter()
	content := makeContent()
	create(calendarName, content)
}

func counter() {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		date, ok := dateOf(f.Name())
		if !ok {
			continue
		}
		record(date)
	}
}

func dateOf(filename string) (time.Time, bool) {
	if len(filename) < 10 {
		return time.Time{}, false
	}
	t, err := time.Parse(format, filename[:10])
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func record(date time.Time) {
	year, month, day := date.Year(), date.Month(), date.Day()
	// 天有记录
	days[year][month][day] = true
	// 月有记录
	monthes[year][int(month)] = true
	// 年有记录
	monthes[year][0] = true
	dec := year / 10
	// 年代有记录
	decades[dec] = true
}

func dayHasRecord(year, month, day int) bool {
	return days[year][month][day]
}

func monthHasRecord(year, month int) bool {
	return monthes[year][month]
}

func yearHasRecord(year int) bool {
	return monthes[year][0]
}

func decadeHasRecord(decade int) bool {
	return decades[decade]
}

func create(fileName, content string) {
	os.Remove(fileName)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("文件打开/创建失败,原因是:", err)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println("文件关闭失败,原因是:", err)
		}
	}()
	file.WriteString(content)
}

func makeContent() string {
	content := ""
	//制作年份表
	content += yearsTable()
	// 制作年段落
	for year := maxYear - 1; year >= 0; year-- {
		if !yearHasRecord(year) {
			continue
		}
		content += yearSection(year)
	}
	return content
}

func yearsTable() string {
	content := fmt.Sprint("# 年份\n\n")
	content += fmt.Sprintln("|0|1|2|3|4|5|6|7|8|9|")
	content += fmt.Sprintln("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|")
	for d := 0; d < maxDecade; d++ {
		if !decadeHasRecord(d) {
			continue
		}
		content += fmt.Sprint("|")
		dec := d * 10
		for y := 0; y < 10; y++ {
			year := dec + y
			if yearHasRecord(year) {
				content += fmt.Sprintf("**[%d](%s)**|", year, yearHeaderLink(year))
			} else {
				content += fmt.Sprint("|")
			}
		}
		content += fmt.Sprintln()
	}
	return content
}

func yearSection(year int) string {
	content := ""
	// add year header
	content += fmt.Sprintf("\n%s\n\n", yearHeader(year))
	//add month line
	content += monthesLine(year)
	//
	content += fmt.Sprintln()
	// add month view
	for month := 12; month > 0; month-- {
		if !monthHasRecord(year, month) {
			continue
		}
		content += monthView(year, month)
	}
	return content
}

func monthesLine(year int) string {
	content := ""
	for m := 12; m > 0; m-- {
		if monthHasRecord(year, m) {
			content += fmt.Sprintf(" [-%s-](%s)", fmtNum(m), monthHeaderLink(year, m))
		}
	}
	return content
}

func monthView(year, month int) string {
	content := ""
	// add month header
	content += fmt.Sprintf("\n%s\n\n", monthHeader(year, month))
	// 输出月历
	content += fmt.Sprintln("|周|一|二|三|四|五|六|日|")
	content += fmt.Sprintln("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|")
	//
	thisWeek := 0
	weekRecord := [8]string{}
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	for day := first; isSameMonth(day, first); day = next(day) {
		_, weekNum := day.ISOWeek()
		if thisWeek != weekNum {
			thisWeek = weekNum
			weekRecord[0] = fmt.Sprintf("**%d**", thisWeek)
		}
		d := day.Day()
		wd := int(day.Weekday())
		if wd == 0 {
			wd = 7 //星期天的 weekday 是 0
		}
		if dayHasRecord(year, month, d) {
			weekRecord[wd] = fmt.Sprintf("[[%s\\|%s]]", day.Format(format), fmtNum(d))
		} else {
			weekRecord[wd] = fmt.Sprintf("%s", fmtNum(d))
		}
		if wd == 7 {
			content += fmt.Sprintf("|%s|\n", strings.Join(weekRecord[:], "|"))
			weekRecord = [8]string{}
		}
	}
	if weekRecord != [8]string{} {
		content += fmt.Sprintf("|%s|\n", strings.Join(weekRecord[:], "|"))
	}
	// 添加一个小尾巴，方便跳转
	content += fmt.Sprintf("\n> [年份](#%%20年份) - [%d](%s)\n", year, yearHeaderLink(year))
	return content
}

func isSameMonth(day, first time.Time) bool {
	return day.Month() == first.Month()
}

func next(day time.Time) time.Time {
	return day.AddDate(0, 0, 1)
}

func yearHeader(year int) string {
	return fmt.Sprintf("## %d", year)
}

func yearHeaderLink(year int) string {
	return fmt.Sprintf("##%%20%d", year)
}

func monthHeader(year, month int) string {
	return fmt.Sprintf("### %d-%02d", year, month)
}

func monthHeaderLink(year, month int) string {
	return fmt.Sprintf("###%%20%d-%02d", year, month)
}

func fmtNum(i int) string {
	return fmt.Sprintf("%d", i)
}
