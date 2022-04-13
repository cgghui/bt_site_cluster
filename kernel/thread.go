package kernel

import (
	"context"
	"sort"
)

var cst = make([]*CreateSiteThread, 0)

func cstSend(site *SiteConfig, bpt *BtPanelThread) {
	length := make([]int, 0)
	for _, c := range cst {
		length = append(length, c.n)
	}
	sort.Slice(length, func(i, j int) bool {
		return length[i] < length[j]
	})
	for _, c := range cst {
		if c.n == length[0] {
			site.SetBtPanelThread(bpt)
			go c.send(site)
			break
		}
	}
}

type CreateSiteThread struct {
	name    string
	n       int
	channel chan *SiteConfig // 不在执行时，可以往该通道发送数据
}

func NewCreateSiteThread(name string) *CreateSiteThread {
	obj := &CreateSiteThread{
		name:    name,
		n:       0,
		channel: make(chan *SiteConfig),
	}
	cst = append(cst, obj)
	return obj
}

func (c *CreateSiteThread) send(site *SiteConfig) {
	c.n++
	c.channel <- site
}

func (c *CreateSiteThread) Listen(ctx context.Context) {
	for {
		select {
		case s := <-c.channel:
			c.n--
			if s != nil {
				s.Run()
				s.bpt.log(nil, "[thread:%v] 创建站点【%s】完成", c.name, s.BindDomain[0])
			}
		case <-ctx.Done():
			return
		}
	}
}
