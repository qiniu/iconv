package main

import (
	"fmt"
	"io"
	"os"
	"github.com/qiniu/iconv"
)

func main() {

	cd, err := iconv.Open("utf-8", "gbk") // gbk => utf8
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()
	
	r := iconv.NewReader(cd, os.Stdin, 0)
	
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Println("\nio.Copy failed:", err)
		return
	}
}

