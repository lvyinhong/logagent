package taillog

import (
	"fmt"
	"llvvlv00.org/logagent/etcd"
	"time"
)

var tskMgr *taillogMgr

type taillogMgr struct {
	logEntry [] *etcd.LogEntry
	tskMap map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}


// 初始化一个taillogMgr
func Init(logEntrys []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntry: logEntrys,	// 将当前的日志收集项任务信息都保存
		tskMap:make(map[string]*TailTask, 16),
		newConfChan:make(chan [] *etcd.LogEntry),
	}
	for _, logEntry := range logEntrys {
		// 初始化时起了多少个tailtask都要记录下来，方便接下来的管理
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)//路径和topic 组成唯一的配置键
		tskMgr.tskMap[mk] = tailObj
	}

	// 监听etcd配置的更新，并管理taillogMgr
	go tskMgr.run()
}

// 监听newConfChan。有了新的配置过来后做对应的处理
func (t *taillogMgr) run() {
	for {
		select {
			case newConf := <-t.newConfChan:
				//1、配置新增
				for _, conf := range newConf {
					mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)//路径和topic 组成唯一的配置键
					if _, ok:=t.tskMap[conf.Path];!ok{
						//如果配置中没有，则新增的有效
						tailObj := NewTailTask(conf.Path, conf.Topic)
						t.tskMap[mk] = tailObj	//存入管理映射集合中
					}
				}
				//2、配置删除// 找出原来 t.tskMap中有，新的newConf中没有要删除
				for _, c1 := range t.logEntry {
					isDelete := true
					for _,c2 := range newConf {
						if c1.Path == c2.Path && c1.Topic == c2.Topic {
							isDelete = false
							continue
						}
					}
					if isDelete {
						//把c1对应的tailObj停止掉
						mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)//路径和topic 组成唯一的配置键
						t.tskMap[mk].cancelFunc()
					}
				}

				//3、配置变更:
				fmt.Println("有新的配置变更:", newConf)
			default:
			time.Sleep(time.Second)
		}

	}
}

// 向外暴露tskMgr的newConfChan
func NewConfCh()chan<-[]*etcd.LogEntry{
	return tskMgr.newConfChan
}
