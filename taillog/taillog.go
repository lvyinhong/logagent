package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"llvvlv00.org/logagent/kafka"
)


// 日志收集的任务
type TailTask struct {
	Path string
	Topic string
	Instance *tail.Tail
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		Path:path,
		Topic:topic,
	}
	tailObj.Init()
	return
}

func (t *TailTask) Init() {
	config := tail.Config{
		ReOpen:true,									// 重新打开
		Follow:true,									// 是否跟随
		//Location:&tail.SeekInfo{Offset:0, Whence:2},	// 从文件的哪个地方开始读
		MustExist:false,								// 文件不存在是否保存
		Poll:true,
	}
	var err error
	t.Instance,err= tail.TailFile(t.Path, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
	}
	go t.run()	//直接读取日志发往kafka
}

func (t *TailTask) run () {
	for {
		select {
		case line:= <-t.Instance.Lines: //从tailObj 的通道中读取日志数据
			//发往kafka
			//kafka.SendToKafka(t.Topic, line.Text)

			// 先把日志数据投递到通道中
			kafka.SendToChan(t.Topic, line.Text)
			// 在kafka包中有单独的goroutine 从通道中取数据发往kafka
		}
	}
}