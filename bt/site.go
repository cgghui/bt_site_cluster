package bt

import (
	"encoding/json"
	"github.com/cgghui/cgghui"
	"strings"
)

type Site struct {
	ID            uint            `json:"id"`
	AddTime       string          `json:"addtime"`
	EndTime       string          `json:"edate"`
	BackupCount   int             `json:"backup_count"`
	BindDomainNum int             `json:"domain"`
	Name          string          `json:"name"`
	Path          string          `json:"path"`
	PHPVersion    string          `json:"php_version"`
	SSL           json.RawMessage `json:"ssl"`
	Status        string          `json:"status"`
	PS            string          `json:"ps"`
}

// CheckStatus 如果站点为开启状态 返回true
func (s *Site) CheckStatus() bool {
	return s.Status == "1"
}

// CheckStatic 如果站点为纯静态 返回true
func (s *Site) CheckStatic() bool {
	return s.PHPVersion == "静态"
}

const (
	PHPVersion72 = "72"
	PHPVersion00 = "00"
	DBNameMySQL  = "MySQL"
	BoolFalse    = "false"
	ListenPort80 = "80"
)

type CreateSiteConfig struct {
	bindDomain BindDomain
	phpVersion string
	port       string
	dbSoftName string
}

type BindDomain struct {
	Domain     string   `json:"domain"`
	DomainList []string `json:"domainlist"`
	Count      int      `json:"count"`
}

func (c *CreateSiteConfig) SetBindDomain(domains []string) {
	if len(domains) == 0 {
		return
	}
	c.bindDomain = BindDomain{
		Domain:     domains[0],
		DomainList: domains[1:],
		Count:      len(domains) - 1,
	}
}

func (c *CreateSiteConfig) GetBindDomainJSON() string {
	ret, err := json.Marshal(c.bindDomain)
	if err != nil {
		return ""
	}
	return string(ret)
}

func (c *CreateSiteConfig) SetPHPVersion(v string) {
	if v == "" {
		c.phpVersion = PHPVersion72
	} else {
		c.phpVersion = v
	}
}

func (c *CreateSiteConfig) SetPort(v string) {
	c.port = v
}

func (c *CreateSiteConfig) EnableSQL(name string) {
	if name == "" {
		c.dbSoftName = BoolFalse
		return
	}
	c.dbSoftName = name
}

func (c *CreateSiteConfig) GetNote() string {
	return strings.ReplaceAll(c.bindDomain.Domain, ".", "_")
}

func (c *CreateSiteConfig) SitePath() string {
	return SiteRoot + "/" + c.GetNote()
}

func (c *CreateSiteConfig) GeneratePassword(length int) string {
	return GeneratePassword(length)
}

func (c *CreateSiteConfig) GetDatabaseUsername() string {
	alias := strings.ReplaceAll(c.GetNote(), ".", "_")
	s := len(alias)
	if s > 16 {
		return alias[:16]
	}
	return alias
}

func GeneratePassword(length int) string {
	return cgghui.RandomString(length, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@")
}
