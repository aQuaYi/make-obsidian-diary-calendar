// 本文件用于随机创建日记文件数据
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	format := "2006-01-02"
	staStr := "2009-05-09"
	endStr := "2020-10-10"
	start, _ := time.Parse(format, staStr)
	end, _ := time.Parse(format, endStr)
	for day := start; day.Before(end); day = day.Add(time.Hour * 24) {
		if randOK() {
			create(fileName(day, format))
		}
	}
}

func randOK() bool {
	return rand.Intn(100) > 50
}

func fileName(day time.Time, format string) string {
	return day.Format(format) + ".md"
}

func create(fileName string) {
	//1.创建文件file.txt
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

	//写入数据
	str := "[[Calendar]]"
	file.WriteString(str)
}
