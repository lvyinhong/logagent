package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"llvvlv00.org/logagent/kafka"
)


// 日志收集的任务
type TailTask struct {
	Path string
	Topic string
	Instance *tail.Tail

	// 实现退出t.Run()
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		Path:path,
		Topic:topic,
		ctx:ctx,
		cancelFunc:cancel,
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
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s 结束了", t.Path , t.Topic)
			return

		case line:= <-t.Instance.Lines: //从tailObj 的通道中读取日志数据
			//发往kafka
			//kafka.SendToKafka(t.Topic, line.Text)

			// 先把日志数据投递到通道中
			kafka.SendToChan(t.Topic, line.Text)
			// 在kafka包中有单独的goroutine 从通道中取数据发往kafka
		}
	}
}