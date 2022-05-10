package kernel

import (
	"bytes"
	"github.com/cgghui/bt_site_cluster/bt"
	"github.com/cgghui/bt_site_cluster/tmpl"
	"github.com/cgghui/cgghui"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type EditFile struct {
	Local  string `json:"local"`
	Server string `json:"server"`
	tpl    *template.Template
}

type SiteConfig struct {
	bpt                *BtPanelThread
	BindDomain         []string   `json:"bind_domain"`
	PHPVersion         string     `json:"php_version"`
	Path               string     `json:"path"`
	Program            string     `json:"program"`
	ProgramName        string     `json:"program_name"`
	SqlFilePath        string     `json:"sql_file_path"`
	RewriteFilePath    string     `json:"rewrite_file_path"`
	EditFiles          []EditFile `json:"edit_files"`
	ResetAdminPassword bool       `json:"reset_admin_password"`
}

func (s *SiteConfig) SetBtPanelThread(bpt *BtPanelThread) {
	s.bpt = bpt
}

func (s *SiteConfig) Run() {
	if len(s.BindDomain) == 0 {
		s.bpt.log(s.bpt.logF, "至少需要绑定一个域名，请检查未设置域名的站点")
		return
	}
	// 数据库
	var sqlTmpl *template.Template
	var err error
	if s.SqlFilePath != "" {
		sqlTmpl = template.New(filepath.Base(s.SqlFilePath))
		sqlTmpl, err = tmpl.RegisterTemplateFuncMap(sqlTmpl).ParseFiles(s.SqlFilePath)
		if err != nil {
			s.bpt.log(s.bpt.logF, "解析数据库备份文件失败 Err: %v", err)
			return
		}
	}
	// 伪静态
	var rewrite *bt.ResponseGetFileBody
	if s.RewriteFilePath != "" {
		var text []byte
		text, err = ioutil.ReadFile(s.RewriteFilePath)
		if err != nil {
			s.bpt.log(s.bpt.logF, "读取伪静态文件失败 Err: %v", err)
			return
		}
		rewrite = &bt.ResponseGetFileBody{Data: string(text)}
	}
	// 修改程序文件
	if len(s.EditFiles) > 0 {
		var e error
		for i, f := range s.EditFiles {
			t := template.New(filepath.Base(f.Local))
			t, e = tmpl.RegisterTemplateFuncMap(t).ParseFiles(f.Local)
			if e != nil {
				s.bpt.log(s.bpt.logF, "【%s】修改程序文件之模板载入失败 Err: %v", s.BindDomain[0], err)
				return
			}
			s.EditFiles[i].tpl = t
		}
	}
	//////////////////////////////////////////////////////////////
	// 创建站点
	conf := bt.CreateSiteConfig{}
	conf.SetBindDomain(s.BindDomain)
	conf.SetPHPVersion(s.PHPVersion)
	conf.SetPort(bt.ListenPort80)
	if s.SqlFilePath != "" {
		conf.EnableSQL(bt.DBNameMySQL)
	}
	siteRootPath := conf.SitePath() // 站点根目录
	var cs *bt.ResponseCreateSite
	if cs, err = s.bpt.s.CreateSiteWithTimeout(s10, &conf); err != nil {
		s.bpt.log(s.bpt.logF, "【%s】站点创建失败 Err: %v", s.BindDomain[0], err)
		return
	}
	if !cs.SiteStatus {
		s.bpt.log(s.bpt.logF, "【%s】站点创建失败 Err: %v", s.BindDomain[0], cs.Msg)
		return
	}
	// 创建数据库 （创建数据库失败时）
	if sqlTmpl != nil && !cs.DatabaseStatus {
		dbName := conf.GetDatabaseUsername()
		dbPass := bt.GeneratePassword(27)
		if err = s.bpt.s.CreateDatabaseWithTimeout(s10, dbName, dbName, bt.DBNameMySQL, dbPass); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】站点数据库创建失败 Err: %v", s.BindDomain[0], err)
			return
		}
		cs.DatabaseStatus = true
		cs.DatabaseUser = dbName
		cs.DatabasePass = dbPass
	}
	// 全局级别的模板变量
	TplVar := map[string]string{
		"NameDB": conf.GetDatabaseUsername(),
		"Domain": s.BindDomain[0],
		"DbUser": cs.DatabaseUser,
		"DbPass": cs.DatabasePass,
	}
	// 配置伪静态
	if rewrite != nil {
		if err = s.bpt.s.SetSiteFileContentWithTimeout(s10, bt.RewriteRoot+"/"+s.BindDomain[0]+".conf", rewrite); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】站点伪静态创建失败 Err: %v", s.BindDomain[0], err)
			return
		}
	}
	_, _ = s.bpt.s.DeleteFileWithTimeout(s10, []string{"404.html", "index.html", ".htaccess"}, siteRootPath)
	// 解压缩程序包
	if err = s.bpt.s.UnZipWithTimeout(time.Minute, s.Program, siteRootPath); err != nil {
		s.bpt.log(s.bpt.logF, "【%s】站点程序压缩包解压失败 Err: %v", s.BindDomain[0], err)
		return
	}
	// 上传数据库
	if sqlTmpl != nil {
		var tmp *os.File
		var storeName = conf.GetNote() + ".sql"
		var storePath = "./temp/" + storeName
		if tmp, err = os.Create(storePath); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】临时SQL文件创建失败 Err: %v", s.BindDomain[0], err) // 解析成功的SQL文件
			return
		}
		if s.ResetAdminPassword {
			min := int64(10)
			max := int64(20)
			TplVar["new_password"] = cgghui.RandomString(int(cgghui.RangeRand(min, max)))
		}
		if err = sqlTmpl.Execute(tmp, TplVar); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】解析SQL文件失败 Err: %v", s.BindDomain[0], err)
			return
		}
		if _, err = tmp.Seek(0, io.SeekStart); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】临时SQL文件指针移至内容开头失败 Err: %v", s.BindDomain[0], err)
			return
		}
		err = s.bpt.s.UploadWithTimeout(m10, storePath, bt.DBackupRoot+"/", tmp, true)
		if err != nil {
			s.bpt.log(s.bpt.logF, "【%s】站点上传sql备份文件失败 Err: %v", s.BindDomain[0], err)
			return
		}
		if err = s.bpt.s.RestoreBackupDBWithTimeout(m10, TplVar["NameDB"], bt.DBackupRoot+"/"+storeName); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】站点还原数据库备份失败 Err: %v", s.BindDomain[0], err)
			return
		}
	}
	// 检测程序包是否解压完成
	st := 120
	tt := time.NewTimer(time.Second)
	defer tt.Stop()
	for range tt.C {
		if st <= 0 {
			break
		}
		st--
		unzipOK := true
		for _, task := range s.bpt.btTaskList {
			to, fail := task.ParseOther()
			if fail != nil {
				continue
			}
			if to.DFile == siteRootPath { // 解压任务仍在进行
				unzipOK = false
				break
			}
		}
		if unzipOK { // 解压缩完成
			break
		}
	}
	// 修改程序文件
	buf := bytes.NewBuffer(nil)
	for _, sf := range s.EditFiles {
		buf.Reset()
		if err = sf.tpl.Execute(buf, TplVar); err != nil {
			s.bpt.log(s.bpt.logF, "【%s】解析模板文件失败 %s->%s Err: %v", s.BindDomain[0], sf.Local, sf.Server, err)
			continue
		}
		err = s.bpt.s.SetSiteFileContentWithTimeout(s10, siteRootPath+sf.Server, &bt.ResponseGetFileBody{Data: buf.String()})
		if err != nil {
			s.bpt.log(s.bpt.logF, "【%s】修改程序文件失败 %s->%s Err: %v", s.BindDomain[0], sf.Local, sf.Server, err)
			continue
		}
	}
	if _, ok := TplVar["new_password"]; ok {
		s.bpt.log(s.bpt.logS, "【%s】id=%d db_user=%s db_pwd=%s admin_pwd=%s", s.BindDomain[0], cs.SiteId, cs.DatabaseUser, cs.DatabasePass, TplVar["new_password"])
	} else {
		s.bpt.log(s.bpt.logS, "【%s】id=%d db_user=%s db_pwd=%s", s.BindDomain[0], cs.SiteId, cs.DatabaseUser, cs.DatabasePass)
	}
}
