package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

const format = "2006-01-02"
const allYear = 3000
const dec = allYear / 10

var dayHasRecord = [allYear][13][32]bool{}
var monthHasRecord = [allYear][13]bool{}
var yearHasRecord = [allYear / 10][11]bool{}

func main() {
	files, _ := ioutil.ReadDir("./days/")
	for _, f := range files {
		t, ok := getTime(f.Name())
		if !ok {
			continue
		}
		// fmt.Println(t)
		year, month, day := t.Year(), t.Month(), t.Day()
		dayHasRecord[year][month][day] = true
		monthHasRecord[year][0] = true
		monthHasRecord[year][int(month)] = true
		d, y := year/10, year%10
		yearHasRecord[d][10] = true
		yearHasRecord[d][y] = true
	}

	// for year, has := range monthHasRecord {
	// 	if has[0] {
	// 		fmt.Printf("%d: ", year)
	// 		for i := 1; i <= 12; i++ {
	// 			if has[i] {
	// 				fmt.Printf(" %02d", i)
	// 			}
	// 		}
	// 		fmt.Println()
	// 	}
	// }

	//
	//制作年代表
	//
	fmt.Print("# 年代表\n\n")
	fmt.Println("|0|1|2|3|4|5|6|7|8|9|")
	fmt.Println("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|")
	for d, has := range yearHasRecord {
		if !has[10] {
			continue
		}
		fmt.Print("|")
		y := d * 10
		for i := 0; i < 10; i++ {
			year := y + i
			if has[i] {
				fmt.Printf("[%d](%s)", year, yearHeaderLink(year))
			}
			fmt.Print("|")
		}
		fmt.Println()
	}

	for year := allYear - 1; year >= 0; year-- {
		if !monthHasRecord[year][0] {
			continue
		}
		fmt.Printf("\n%s\n\n", yearHeader(year))
		for m := 12; m > 0; m-- {
			if monthHasRecord[year][m] {
				fmt.Printf(" [-%s-](%s)", fmtNumber(m), monthHeaderLink(year, m))
			}
		}
		fmt.Println()
		for m := 12; m > 0; m-- {
			if !monthHasRecord[year][m] {
				continue
			}
			fmt.Printf("\n%s\n\n", monthHeader(year, m))
			// 输出月历
			firstDay := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
			fmt.Println("|周|一|二|三|四|五|六|日|")
			fmt.Println("|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|")
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
				if dayHasRecord[year][m][d] {
					weekRecord[wd] = fmt.Sprintf("[[%s\\|%s]]", day.Format(format), fmtNumber(d))
				} else {
					weekRecord[wd] = fmt.Sprintf("%s", fmtNumber(d))
				}
				if wd == 7 {
					fmt.Printf("|%s|\n", strings.Join(weekRecord[:], "|"))
					weekRecord = [8]string{}
				}
			}
			if weekRecord != [8]string{} {
				fmt.Printf("|%s|\n", strings.Join(weekRecord[:], "|"))
			}

			fmt.Printf("\n> [年代表](#%%20年代表) - [%d](%s)\n", year, yearHeaderLink(year))

			// fmt.Println("|周|一|二|三|四|五|六|日|")
		}
	}
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

func fmtNumber(i int) string {
	if i < 10 {
		return fmt.Sprintf("   %d  ", i)
	}
	return fmt.Sprintf("  %d  ", i)
}

func getTime(filename string) (time.Time, bool) {
	if len(filename) < 10 {
		return time.Time{}, false
	}
	t, err := time.Parse(format, filename[:10])
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
