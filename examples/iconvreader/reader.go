package main

import (
	"fmt"
	"gopkg.in/iconv.v1"
	"io"
	"os"
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
