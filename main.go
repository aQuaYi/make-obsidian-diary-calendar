package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const format = "2006-01-02"
const maxYear = 3000
const maxDecade = maxYear / 10

var days = [maxYear][13][32]bool{}
var monthes = [maxYear][13]bool{}
var decades = [maxDecade]bool{}

func main() {
	counter()
	content := makeContent()
	create("new.md", content)
}

func create(fileName, content string) {
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
	content += yearTable()
	// 制作年段落
	for year := maxYear - 1; year >= 0; year-- {
		if !yearHasRecord(year) {
			continue
		}
		content += yearSection(year)
	}
	return content
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

func getDate(filename string) (time.Time, bool) {
	if len(filename) < 10 {
		return time.Time{}, false
	}
	t, err := time.Parse(format, filename[:10])
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func set(date time.Time) {
	year, month, day := date.Year(), date.Month(), date.Day()
	days[year][month][day] = true
	monthes[year][0] = true
	monthes[year][int(month)] = true
	dec := year / 10
	decades[dec] = true
}

func thisDecadeHasRecord(dec [11]bool) bool {
	return dec[10]
}

func counter() {
	files, _ := ioutil.ReadDir("./Diary/")
	for _, f := range files {
		date, ok := getDate(f.Name())
		if !ok {
			continue
		}
		set(date)
	}
}

func yearTable() string {
	content := fmt.Sprint("# 年份\n\n")
	content += fmt.Sprintln("|0|1|2|3|4|5|6|7|8|9|")
	content += fmt.Sprintln("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|")
	for d := 0; d < maxDecade; d++ {
		if !decades[d] {
			continue
		}
		content += fmt.Sprint("|")
		y := d * 10
		for i := 0; i < 10; i++ {
			year := y + i
			if yearHasRecord(year) {
				content += fmt.Sprintf("**[%d](%s)**", year, yearHeaderLink(year))
			}
			content += fmt.Sprint("|")
		}
		content += fmt.Sprintln()
	}
	return content
}

func yearHasRecord(year int) bool {
	return monthes[year][0]
}

func yearSection(year int) string {
	content := ""
	// add year header
	content += fmt.Sprintf("\n%s\n\n", yearHeader(year))
	//add month line
	content += monthesLine(year)
	//
	content += fmt.Sprintln()

	for m := 12; m > 0; m-- {
		if !monthes[year][m] {
			continue
		}
		content += monthView(year, m)
	}
	return content
}

func monthesLine(year int) string {
	content := ""
	for m := 12; m > 0; m-- {
		if monthes[year][m] {
			content += fmt.Sprintf(" [-%s-](%s)", fmtNum(m), monthHeaderLink(year, m))
		}
	}
	return content
}

func monthView(year, month int) string {
	content := ""
	content += fmt.Sprintf("\n%s\n\n", monthHeader(year, month))
	// 输出月历
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	content += fmt.Sprintln("|周|一|二|三|四|五|六|日|")
	content += fmt.Sprintln("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|")
	week := 0
	weekRecord := [8]string{}
	for day := firstDay; day.Month() == firstDay.Month(); day = day.Add(time.Hour * 24) {
		_, w := day.ISOWeek()
		if week != w {
			week = w
			weekRecord[0] = fmt.Sprintf("**%d**", week)
		}
		d := day.Day()
		wd := int(day.Weekday())
		if wd == 0 {
			wd = 7 //星期天的 weekday 是 0
		}
		if days[year][month][d] {
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
	content += fmt.Sprintf("\n> [年份](#%%20年份) - [%d](%s)\n", year, yearHeaderLink(year))
	return content
}
