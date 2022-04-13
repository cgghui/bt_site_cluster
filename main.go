package main

import (
	"context"
	"github.com/cgghui/bt_site_cluster/bt"
	"github.com/cgghui/bt_site_cluster/kernel"
	"github.com/cgghui/bt_site_cluster/tmpl"
	"log"
	"os"
	"strconv"
	"sync"
)

const FLAG = os.O_RDWR | os.O_CREATE | os.O_APPEND

func main() {
	var (
		opt []bt.Option
		err error
	)
	var fpS, fpF *os.File
	if fpS, err = os.OpenFile("./log/success.log", FLAG, 0666); err != nil {
		log.Panicf("无法创建日志文件 Err: %v", err)
	}
	if fpF, err = os.OpenFile("./log/failure.log", FLAG, 0666); err != nil {
		log.Panicf("无法创建日志文件 Err: %v", err)
	}
	if err = tmpl.LoadTmplData(); err != nil {
		panic(err)
	}
	if err = kernel.LoadBtPanelConfig(&opt); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := strconv.Itoa(i)
			kernel.NewCreateSiteThread(k).Listen(ctx)
		}(i + 1)
	}
	for i := range opt {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			kernel.NewBtPanelThread(ctx, &opt[i], fpS, fpF)
		}(i)
	}
	wg.Wait()
}
