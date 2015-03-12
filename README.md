iconv: libiconv for go
======

[![Build Status](https://travis-ci.org/qiniu/iconv.png?branch=master)](https://travis-ci.org/qiniu/iconv) [![Build Status](https://drone.io/github.com/qiniu/iconv/status.png)](https://drone.io/github.com/qiniu/iconv/latest)

![logo](http://qiniutek.com/images/logo-2.png)

iconv is a libiconv wrapper for go. libiconv Convert string to requested character encoding.

iconv project's homepage is: https://github.com/qiniu/iconv

# Document

See http://godoc.org/github.com/qiniu/iconv

Note: Open returns a conversion descriptor cd, cd contains a conversion state and can not be used in multiple threads simultaneously.

# Install

```
go get github.com/qiniu/iconv
```

# Example

## Convert string

```go
package main

import (
	"fmt"
	"github.com/qiniu/iconv"
)

func main() {

	cd, err := iconv.Open("gbk", "utf-8") // convert utf-8 to gbk
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()

	gbk := cd.ConvString("你好，世界！")

	fmt.Println(gbk)
}
```

## Output to io.Writer

```go
package main

import (
	"fmt"
	"github.com/qiniu/iconv"
)

func main() {

	cd, err := iconv.Open("gbk", "utf-8") // convert utf-8 to gbk
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()

	output := ... // eg. output := os.Stdout || ouput, err := os.Create(file)
	autoSync := false // buffered or not
	bufSize := 0 // default if zero
	w := iconv.NewWriter(cd, output, bufSize, autoSync)

	fmt.Fprintln(w, "你好，世界！")

	w.Sync() // if autoSync = false, you need call Sync() by yourself
}
```

## Input from io.Reader

```go
package main

import (
	"fmt"
	"io"
	"os"
	"github.com/qiniu/iconv"
)

func main() {

	cd, err := iconv.Open("utf-8", "gbk") // convert gbk to utf8
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()
	
	input := ... // eg. input := os.Stdin || input, err := os.Open(file)
	bufSize := 0 // default if zero
	r := iconv.NewReader(cd, input, bufSize)
	
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		fmt.Println("\nio.Copy failed:", err)
		return
	}
}
```

