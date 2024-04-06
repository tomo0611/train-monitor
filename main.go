package main

import (
	"github.com/tomo0611/train-monitor/downloader"
	"github.com/tomo0611/train-monitor/parser"
)

/*
	1 特急
	2 快速急行
	3 急行
	4 準急
	5 区間準急
	6 普通
	7 区間快速
	8 区間急行
	9 区間快速急行
	10  一般貸切
	11  特急貸切
*/

func main() {

	//【A】近鉄奈良線 (大阪難波～近鉄奈良)
	data, err := downloader.GetData()
	if err != nil {
		panic(err)
	}
	parser.ParseData(data)
}
