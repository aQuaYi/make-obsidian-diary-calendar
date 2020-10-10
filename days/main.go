// 本文件用于随机创建日记文件数据
package main

import (
	"fmt"
	"os"
)

func main() {

	return
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
	str := "[[📅 Calendar]]"
	file.WriteString(str)
}
