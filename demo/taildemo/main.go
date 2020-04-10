package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

// tailf 的用法示例
func main() {
	filename := "./my.log"
	config := tail.Config{
		ReOpen:true,									// 重新打开
		Follow:true,									// 是否跟随
		Location:&tail.SeekInfo{Offset:0, Whence:2},	// 从文件的哪个地方开始读
		MustExist:false,								// 文件不存在是否报错
		Poll:true,
	}

	// 返回的内容
	tails, err := tail.TailFile(filename, config)

	if err != nil {
		fmt.Println("tail file failed, err :", err)
		return
	}
	var (
		line *tail.Line
		ok bool
	)
	for {
		line, ok = <- tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename: %s\n", tails.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("msg:", line.Text)
	}
}