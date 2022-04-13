package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cgghui/bt_site_cluster/bt"
	"io"
	"log"
	"os"
	"time"
)

func LoadBtPanelConfig(ret interface{}) error {
	fp, err := os.Open("./bt_panel_config.json")
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()
	if err = json.NewDecoder(fp).Decode(ret); err != nil {
		return err
	}
	return nil
}

func loadSiteConfig(ip string) ([]SiteConfig, error) {
	fp, err := os.Open("./site." + ip + ".json")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fp.Close()
	}()
	var ret []SiteConfig
	if err = json.NewDecoder(fp).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type BtPanelThread struct {
	logWrite   io.Writer
	logS       *os.File
	logF       *os.File
	o          *bt.Option
	s          *bt.Session
	ip         string
	btTaskList []bt.ResponseTask
}

func NewBtPanelThread(ctx context.Context, opt *bt.Option, logS, logF *os.File) {
	fp, err := os.OpenFile("./log/"+opt.GetAddress()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("无法生成日志文件 Err: %v", err)
		return
	}
	thread := &BtPanelThread{
		logWrite: io.MultiWriter(fp, os.Stdout),
		logS:     logS,
		logF:     logF,
		o:        opt,
		ip:       opt.GetAddress(),
	}
	// 加载站点
	var siteList []SiteConfig
	if siteList, err = loadSiteConfig(thread.ip); err != nil {
		thread.log(nil, "加载站点列表失败 Err: %v", err)
		return
	}
	// 登录宝塔
	if err = thread.loginBtPanel(ctx); err != nil {
		return
	}
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				thread.btTaskList, _ = thread.s.GetTaskListsWithTimeout(s03)
				t.Reset(time.Second)
			case <-ctx.Done():
				return
			}
		}
	}()
	for _, site := range siteList {
		cstSend(&site, thread)
	}
}

var retry = []time.Duration{
	time.Second,
	10 * time.Second,
	20 * time.Second,
	30 * time.Second,
	40 * time.Second,
	50 * time.Second,
	time.Minute,
	3 * time.Minute,
	6 * time.Minute,
	9 * time.Minute,
}

func (p *BtPanelThread) loginBtPanel(ctx context.Context) error {
	var err error
	t := time.NewTimer(0)
	defer t.Stop()
	s := len(retry)
	n := 0
	i := 0
	for range t.C {
		if n >= 3 {
			i++
			n = 0
		}
		if i >= s {
			n = 0
			i = 0
		}
		ct, cancel := context.WithTimeout(ctx, s10)
		p.s, err = p.o.Login(ct)
		cancel()
		if err == nil {
			break
		}
		n++
		p.log(nil, "%s后重试 %v", retry[i], err)
		t.Reset(retry[i])
	}
	return err
}

func (p *BtPanelThread) log(to io.Writer, s string, arg ...interface{}) {
	if len(arg) == 0 {
		s = fmt.Sprintf("[%s] [%s] "+s+"\n", time.Now().Format("2006-01-02 15:04:05"), p.ip)
	} else {
		g := make([]interface{}, 2, len(arg)+2)
		g[0] = time.Now().Format("2006-01-02 15:04:05")
		g[1] = p.ip
		g = append(g, arg...)
		s = fmt.Sprintf("[%s] [%s] "+s+"\n", g...)
	}
	sByte := []byte(s)
	_, _ = p.logWrite.Write(sByte)
	if to != nil {
		_, _ = to.Write(sByte)
	}
}

const (
	s03 = 3 * time.Second
	s10 = 10 * time.Second
	m10 = 10 * time.Minute
)
