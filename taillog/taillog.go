package taillog

import (
	"github.com/hpcloud/tail"
)

var(
	tailObj *tail.Tail
)

// 专门从日志文件收集日志的模块
func Init(filename string) (err error) {
	config := tail.Config{
		ReOpen:true,									// 重新打开
		Follow:true,									// 是否跟随
		//Location:&tail.SeekInfo{Offset:0, Whence:2},	// 从文件的哪个地方开始读
		MustExist:false,								// 文件不存在是否保存
		Poll:true,
	}

	// 返回的内容
	tailObj, err = tail.TailFile(filename, config)
	return
}

func ReadChan() (chan *tail.Line){
	return tailObj.Lines
}