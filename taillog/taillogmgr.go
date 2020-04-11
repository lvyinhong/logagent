package taillog

import "llvvlv00.org/logagent/etcd"

type taillogMgr struct {
	logEntry [] *etcd.LogEntry
	tskMap map[string]*TailTask
}
var tskMgr *taillogMgr

// 初始化一个taillogMgr
func Init(logEntrys []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntry: logEntrys,	// 将当前的日志收集项任务信息都保存
	}
	for _, logEntry := range logEntrys {
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
}
