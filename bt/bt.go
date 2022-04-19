package bt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cgghui/cgghui"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"

type Bt struct {
	Link        string
	reqLogin    *http.Request
	reqGetToken *http.Request
}

var SiteRoot = "/www/wwwroot"
var RewriteRoot = "/www/server/panel/vhost/rewrite"
var DBackupRoot = "/www/backup/database"

func New(link string) *Bt {
	return &Bt{
		Link: strings.Trim(strings.TrimSpace(link), "/") + "/",
	}
}

type Option struct {
	Link     string   `json:"url"`      //后面地址路径
	Username string   `json:"username"` // 账号
	Password string   `json:"password"` // 密码
	Code     string   `json:"code"`     //
	Node     string   `json:"node"`     // 备注
	s        *Session // 共享会话
}

// LoadOption 从文件加载option
func LoadOption(filePath string) (*Option, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var opt Option
	if err = json.Unmarshal(content, &opt); err != nil {
		return nil, err
	}
	return &opt, err
}

func (o *Option) SetLoginSession(s *Session) *Session {
	o.s = s
	return s
}

func (o *Option) GetLoginSession() *Session {
	return o.s
}

func (o *Option) Login(ctx context.Context) (*Session, error) {
	return New(o.Link).Login(ctx, o.Code, o.Username, o.Password)
}

func (o *Option) GetAddress() string {
	l, _ := url.Parse(o.Link)
	return strings.Split(l.Host, ":")[0]
}

type Session struct {
	bt          *Bt
	btToken     string
	cookieStr   string
	cookieToken string
	code        string
}

func (b *Bt) Login(ctx context.Context, code, username, password string) (*Session, error) {
	var err error
	if b.reqLogin == nil {
		u := cgghui.MD5(username)
		p := cgghui.MD5(cgghui.MD5(password) + "_bt.cn")
		s := strings.NewReader("username=" + u + "&password=" + p + "&code=")
		b.reqLogin, err = http.NewRequest(http.MethodPost, b.Link+"login", s)
		if err != nil {
			return nil, fmt.Errorf("登录宝塔失败 Error[1]: %v", err)
		}
		b.reqLogin.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		b.reqLogin.Header.Add("Origin", b.Link)
		b.reqLogin.Header.Add("Referer", b.Link+code+"/")
		b.reqLogin.Header.Add("User-Agent", UserAgent)
	}

	resp, err := http.DefaultClient.Do(b.reqLogin.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("登录宝塔失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("登录宝塔失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("登录宝塔失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if !ret.Status {
		return nil, fmt.Errorf("登录宝塔失败 Error[5]: %s", ret.Msg)
	}
	ckArr := make(map[string]string, 0)
	ckStr := make([]string, 0)
	for _, cookie := range resp.Cookies() {
		ckArr[cookie.Name] = cookie.Value
		ckStr = append(ckStr, cookie.Name+"="+cookie.Value)
	}
	ckArr["sites_path"] = SiteRoot
	ckArr["serverType"] = "nginx"
	ckArr["site_type"] = "-1"
	ckArr["order"] = "id%20desc"
	ckStr = append(ckStr, "sites_path="+SiteRoot)
	ckStr = append(ckStr, "serverType=nginx")
	ckStr = append(ckStr, "site_type=-1")
	ckStr = append(ckStr, "order=id%20desc")
	obj := &Session{
		bt:          b,
		cookieStr:   strings.Join(ckStr, "; "),
		cookieToken: ckArr["request_token"],
		btToken:     "",
		code:        code,
	}
	if err = obj.getToken(); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *Session) GetAddress() string {
	l, _ := url.Parse(s.bt.Link)
	return strings.Split(l.Host, ":")[0]
}

func (s *Session) getToken() error {
	var err error
	if s.bt.reqGetToken == nil {
		s.bt.reqGetToken, err = http.NewRequest(http.MethodGet, s.bt.Link, nil)
		if err != nil {
			return fmt.Errorf("获取TOKEN失败 Error[1]: %v", err)
		}
		s.bt.reqGetToken.Header.Add("Referer", s.bt.Link+s.code+"/")
		s.bt.reqGetToken.Header.Add("Cookie", s.cookieStr)
		s.bt.reqGetToken.Header.Add("User-Agent", UserAgent)
	}
	resp, err := http.DefaultClient.Do(s.bt.reqGetToken)
	if err != nil {
		return fmt.Errorf("获取TOKEN失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("获取TOKEN失败 Error[3]: http code is %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("TK宝塔失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	s.btToken = doc.Find("#request_token_head").AttrOr("token", "")
	return nil
}

func (s *Session) DownFile(downFileURL, savePath, saveName string) (*ResponseBase, error) {
	body := url.Values{}
	body.Add("url", downFileURL)
	body.Add("path", savePath)
	body.Add("filename", saveName)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=DownloadFile", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("下载文件失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载文件失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载文件失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("下载文件失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) BackupSite(id uint) (*ResponseBase, error) {
	body := url.Values{}
	body.Add("id", strconv.FormatUint(uint64(id), 10))
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"site?action=ToBackup", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("备份站点失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("备份站点失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("备份站点失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("备份站点失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) BackupList(id uint) (*ResponseBackupList, error) {
	body := url.Values{}
	body.Add("table", "backup")
	body.Add("type", "0")
	body.Add("limit", "100")
	body.Add("p", "1")
	body.Add("search", strconv.FormatUint(uint64(id), 10))
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"data?action=getData", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("备份列表失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("备份列表失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("备份列表失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBackupList
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("备份列表失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) GetDir(path string) (*ResponseDirs, error) {
	body := url.Values{}
	body.Add("p", "1")
	body.Add("showRow", "2000")
	body.Add("path", path)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=GetDir", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("获取目录列表 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取目录列表 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取目录列表 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseDirs
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取目录列表 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) DownFileLocal(target string, fp *os.File) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return 0, fmt.Errorf("下载文件 Error[1]: %v", err)
	}
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("下载文件 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("下载文件 Error[3]: http code is %d", resp.StatusCode)
	}
	sv, err := io.Copy(fp, resp.Body)
	defer resp.Body.Close()
	return sv, err
}

func (s *Session) DeleteFileWithTimeout(out time.Duration, files []string, rootPath string) (*ResponseBase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.DeleteFileWithContext(ctx, files, rootPath)
}

func (s *Session) DeleteFileWithContext(ctx context.Context, files []string, rootPath string) (*ResponseBase, error) {
	body := url.Values{}
	data, _ := json.Marshal(files)
	body.Add("data", string(data))
	body.Add("type", "4")
	body.Add("path", rootPath)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=SetBatchData", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("删除文件 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("删除文件 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("删除文件 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("删除文件 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return &ret, nil
}

func (s *Session) ClearRecycle() (*ResponseBase, error) {
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=Close_Recycle_bin", nil)
	if err != nil {
		return nil, fmt.Errorf("清空回收站失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("清空回收站失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("清空回收站失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("清空回收站失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) GetPlugin(ctx context.Context) ([]*ResponsePlugin, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.bt.Link+"plugin?action=get_index_list", nil)
	if err != nil {
		return nil, fmt.Errorf("获取服务器软件失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取服务器软件失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取服务器软件失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret []*ResponsePlugin
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取服务器软件失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return ret, nil
}

func (s *Session) GetStatus(ctx context.Context) (*ResponseNetWork, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.bt.Link+"system?action=GetNetWork", nil)
	if err != nil {
		return nil, fmt.Errorf("获取服务器状态失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取服务器状态失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取服务器状态失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseNetWork
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取服务器状态失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) CreateSiteWithTimeout(out time.Duration, conf *CreateSiteConfig) (*ResponseCreateSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.CreateSiteWithContext(ctx, conf)
}

func (s *Session) CreateSiteWithContext(ctx context.Context, conf *CreateSiteConfig) (*ResponseCreateSite, error) {
	ps := conf.GetNote()
	pwd := conf.GeneratePassword(27)
	body := url.Values{}
	body.Add("webname", conf.GetBindDomainJSON())
	body.Add("type", "PHP")
	body.Add("port", conf.port)
	body.Add("ps", ps)
	body.Add("path", SiteRoot+"/"+ps)
	body.Add("type_id", "0")
	body.Add("version", conf.phpVersion)
	body.Add("ftp", "false")
	body.Add("sql", conf.dbSoftName)
	if conf.dbSoftName != "false" {
		body.Add("datauser", conf.GetDatabaseUsername())
		body.Add("datapassword", pwd)
	}
	body.Add("codeing", "utf8")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"site?action=AddSite", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建站点失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("创建站点失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("创建站点失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseCreateSite
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("创建站点失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return &ret, nil
}

func (s *Session) CreateDatabaseWithTimeout(out time.Duration, dbName, dbUser, dbType, dbPass string) error {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.CreateDatabaseWithContext(ctx, dbName, dbUser, dbType, dbPass)
}

func (s *Session) CreateDatabaseWithContext(ctx context.Context, dbName, dbUser, dbType, dbPass string) error {
	body := url.Values{}
	body.Add("name", dbName)
	body.Add("db_user", dbUser)
	body.Add("password", dbPass)
	body.Add("dtype", dbType)
	body.Add("dataAccess", "%")
	body.Add("address", "%")
	body.Add("ps", dbName)
	body.Add("codeing", "utf8")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"database?action=AddDatabase", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("创建数据库失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return fmt.Errorf("创建数据库失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("创建数据库失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("创建数据库失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if !ret.Status {
		return fmt.Errorf("创建数据库失败 Error[5]: %s", ret.Msg)
	}
	return nil
}

func (s *Session) GetDatabaseList() (*ResponseTableDatabases, error) {
	body := url.Values{}
	body.Add("table", "databases")
	body.Add("tojs", "database.get_list")
	body.Add("limit", "10000")
	body.Add("p", "1")
	body.Add("search", "")
	body.Add("order", "id desc")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"data?action=getData", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("获取数据库列表失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取数据库列表失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取数据库列表失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseTableDatabases
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取数据库列表失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return &ret, nil
}

// GetSiteListWithTimeout 战点列表
func (s *Session) GetSiteListWithTimeout(out time.Duration) (*ResponseTableSites, error) {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.GetSiteListWithContext(ctx)
}

// GetSiteListWithContext 战点列表
func (s *Session) GetSiteListWithContext(ctx context.Context) (*ResponseTableSites, error) {
	body := url.Values{}
	body.Add("table", "sites")
	body.Add("type", "-1")
	body.Add("limit", "10000")
	body.Add("p", "1")
	body.Add("search", "")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"data?action=getData", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("获取站点列表失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("获取站点列表失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取站点列表失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseTableSites
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取站点列表失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

func (s *Session) GetSiteBindDomain(site *Site) ([]ResponseTableDomain, error) {
	body := url.Values{}
	body.Add("table", "domain")
	body.Add("list", "True")
	body.Add("search", strconv.FormatUint(uint64(site.ID), 10))
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"data?action=getData", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("获取站点绑定域名失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取站点绑定域名失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取站点绑定域名失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret []ResponseTableDomain
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取站点绑定域名失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return ret, nil
}

func (s *Session) GetSiteFileContent(path string) (*ResponseGetFileBody, error) {
	body := url.Values{}
	body.Add("path", path)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=GetFileBody", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("获取文件内容失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取文件内容失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取文件内容失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseGetFileBody
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取文件内容失败 Error[4]: %v", err)
	}
	defer resp.Body.Close()
	return &ret, nil
}

// SetSiteFileContentWithTimeout 保存修改文件
// path 为文件的完整路径
func (s *Session) SetSiteFileContentWithTimeout(out time.Duration, path string, r *ResponseGetFileBody) error {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.SetSiteFileContentWithContext(ctx, path, r)
}

// SetSiteFileContentWithContext 保存修改文件
// path 为文件的完整路径
func (s *Session) SetSiteFileContentWithContext(ctx context.Context, path string, r *ResponseGetFileBody) error {
	body := url.Values{}
	body.Add("data", r.Data)
	body.Add("encoding", "utf-8")
	body.Add("path", path)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=SaveFileBody", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("写入文件内容失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return fmt.Errorf("写入文件内容失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("写入文件内容失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("获取文件内容失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if ret.Status == false {
		return fmt.Errorf("获取文件内容失败 Error[5]: %s", ret.Msg)
	}
	return nil
}

// UnZipWithTimeout 解压文件
// src 压缩文件路径 dst 解压到的文件夹
func (s *Session) UnZipWithTimeout(out time.Duration, src, dst string) error {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.UnZipWithContext(ctx, src, dst)
}

// UnZipWithContext 解压文件
// src 压缩文件路径 dst 解压到的文件夹
func (s *Session) UnZipWithContext(ctx context.Context, src, dst string) error {
	body := url.Values{}
	body.Add("sfile", src)
	body.Add("dfile", dst)
	body.Add("type", "zip")
	body.Add("coding", "UTF-8")
	body.Add("password", "")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=UnZip", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("解压文件内容失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return fmt.Errorf("解压文件内容失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("解压文件内容失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("解压文件内容失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if !ret.Status {
		return fmt.Errorf("解压文件内容失败 Error[5]: %s", ret.Msg)
	}
	return nil
}

// GetTaskListsWithTimeout 获取任务列表
func (s *Session) GetTaskListsWithTimeout(out time.Duration) ([]ResponseTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.GetTaskListsWithContext(ctx)
}

// GetTaskListsWithContext 获取任务列表
func (s *Session) GetTaskListsWithContext(ctx context.Context) ([]ResponseTask, error) {
	body := url.Values{}
	body.Add("status", "-3")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"task?action=get_task_lists", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("获取任务列表失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("获取任务列表失败 Error[2]: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取任务列表失败 Error[3]: http code is %d", resp.StatusCode)
	}
	var ret []ResponseTask
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("获取任务列表失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return ret, nil
}

// UploadWithTimeout 上传文件
func (s *Session) UploadWithTimeout(out time.Duration, localFilePath, saveFilePath string, fp *os.File, isRemoveFile bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.UploadWithContext(ctx, localFilePath, saveFilePath, fp, isRemoveFile)
}

// UploadWithContext 上传文件
func (s *Session) UploadWithContext(ctx context.Context, localFilePath, saveFilePath string, fp *os.File, isRemoveFile bool) error {
	if isRemoveFile {
		defer func() {
			_ = fp.Close()
			_ = os.Remove(localFilePath)
		}()
	}
	var (
		body bytes.Buffer
		ct   string
		err  error
	)
	if err = upload(filepath.Base(localFilePath), saveFilePath, fp, &body, &ct); err != nil {
		return fmt.Errorf("上传文件失败 Error[1]: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=upload", &body)
	if err != nil {
		return fmt.Errorf("上传文件失败 Error[2]: %v", err)
	}
	req.Header.Add("Content-Type", ct)
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return fmt.Errorf("上传文件失败 Error[3]: %v", err)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("上传文件失败 Error[4]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if !ret.Status {
		return fmt.Errorf("上传文件失败 Error[5]: %v", ret.Msg)
	}
	return nil
}

func upload(name, saveFilePath string, fp *os.File, ret *bytes.Buffer, ct *string) error {
	fi, _ := fp.Stat()
	sendWriter := multipart.NewWriter(ret)
	_ = sendWriter.WriteField("f_path", saveFilePath)
	_ = sendWriter.WriteField("f_name", name)
	_ = sendWriter.WriteField("f_size", strconv.FormatInt(fi.Size(), 10))
	_ = sendWriter.WriteField("f_start", "0")
	fileWriter, err := sendWriter.CreateFormFile("blob", "blob")
	if err != nil {
		return err
	}
	if _, err = io.Copy(fileWriter, fp); err != nil {
		return err
	}
	*ct = sendWriter.FormDataContentType()
	_ = sendWriter.Close()
	return nil
}

// RestoreBackupDBWithTimeout 还原SQL备份
func (s *Session) RestoreBackupDBWithTimeout(out time.Duration, dbName, backupFilePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), out)
	defer cancel()
	return s.RestoreBackupDBWithContext(ctx, dbName, backupFilePath)
}

// RestoreBackupDBWithContext 还原SQL备份
func (s *Session) RestoreBackupDBWithContext(ctx context.Context, dbName, backupFilePath string) error {
	data := url.Values{}
	data.Add("file", backupFilePath)
	data.Add("name", dbName)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"database?action=InputSql", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("还原数据库备份文件失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response

	if resp, err = http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		return fmt.Errorf("还原数据库备份文件失败 Error[2]: %v", err)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("还原数据库备份文件失败 Error[3]: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if !ret.Status {
		return fmt.Errorf("还原数据库备份文件失败 Error[5]: %s", ret.Msg)
	}
	return nil
}

// SiteDatabasePermission 将站点数据库更变为可远程操作
func (s *Session) SiteDatabasePermission(dbName string) error {
	body := url.Values{}
	body.Add("name", dbName)
	body.Add("dataAccess", "%")
	body.Add("access", "%")
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"database?action=SetDatabaseAccess", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("修改数据库权限失败 %s Error[1]: %v", dbName, err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("修改数据库权限失败 %s Error[2]: %v", dbName, err)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("修改数据库权限失败 %s Error[3]: %v", dbName, err)
	}
	defer resp.Body.Close()
	if !ret.Status {
		return fmt.Errorf("修改数据库权限失败 %s Error[4]: %v", dbName, err)
	}
	return nil
}

func (s *Session) SiteDelete(id []string) error {
	body := url.Values{}
	body.Add("database", "1")
	body.Add("path", "1")
	body.Add("ftp", "1")
	body.Add("sites_id", strings.Join(id, ","))
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"site?action=delete_website_multiple", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("删除站点失败 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除站点失败 Error[2]: %v", err)
	}
	var ret ResponseDeleteSite
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("删除站点失败 Error[3]: %v", err)
	}
	defer resp.Body.Close()
	if !ret.Status {
		return fmt.Errorf("删除站点失败 Error[4]: %v", err)
	}
	return nil
}

func (s *Session) CreateFile(savePath string) error {
	body := url.Values{}
	body.Add("path", savePath)
	req, err := http.NewRequest(http.MethodPost, s.bt.Link+"files?action=CreateFile", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("创建文件 Error[1]: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("x-http-token", s.btToken)
	req.Header.Add("x-cookie-token", s.cookieToken)
	req.Header.Add("Cookie", s.cookieStr)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("创建文件 Error[2]: %v", err)
	}
	var ret ResponseBase
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return fmt.Errorf("创建文件 Error[3]: %v", err)
	}
	defer resp.Body.Close()
	if !ret.Status {
		return fmt.Errorf("创建文件 Error[4]: %v", ret.Msg)
	}
	return nil
}
