package main

import (
	"fmt"
	"github.com/qiniu/iconv"
)

func main() {

	cd, err := iconv.Open("gbk", "utf-8")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()

	gbk := cd.ConvString(
`		你好，世界！你好，世界！你好，世界！你好，世界！
		你好，世界！你好，世界！你好，世界！你好，世界！`)
	fmt.Println(gbk)
}

