package tmpl

import (
	"bytes"
	"github.com/cgghui/cgghui"
	"github.com/mozillazg/go-pinyin"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	wordList map[string]*[]string
	kvList   map[string]map[string]string
	kvKeys   map[string]*[]string
	mutex    = &sync.Mutex{}
	pyArg    = pinyin.NewArgs()
)

const (
	pathSep = string(os.PathSeparator)
	ms      = 10 * time.Millisecond
)

// loadFileList 处理目录下的文件
// wg 以等待进度完成处理
// folder 目录路径
// ext 文件后缀
// RunCall 执行处理
func loadFileList(wg *sync.WaitGroup, folder, ext string, RunCall func(filePath string) error) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	t := time.NewTimer(ms)
	defer t.Stop()
	for _, file := range files {
		fp := folder + pathSep + file.Name()
		if file.IsDir() {
			for range t.C {
				t.Reset(ms)
				if err = loadFileList(wg, fp, ext, RunCall); err == nil {
					break
				}
			}
		} else {
			if ext != "" {
				if path.Ext(fp) == ext {
					handleFile(wg, pathReal(fp), RunCall)
				}
			} else {
				handleFile(wg, pathReal(fp), RunCall)
			}
		}
	}
	return nil
}

func handleFile(wg *sync.WaitGroup, fp string, RunCall func(filePath string) error) {
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTimer(ms)
		defer t.Stop()
		for range t.C {
			if err = RunCall(fp); err != nil {
				eMsg := err.Error()
				if strings.Contains(eMsg, "too many open files") || strings.Contains(eMsg, "Too many connections") {
					t.Reset(ms)
					continue
				}
				log.Printf("Call Func Error: %s %v", fp, err)
			}
			break
		}
	}()
}

func pathReal(fp string) string {
	fp = strings.ReplaceAll(fp, "\\", pathSep)
	fp = strings.ReplaceAll(fp, "/", pathSep)
	return fp
}

func LoadTmplData() error {
	var err error
	// kv
	kvList = make(map[string]map[string]string)
	kvKeys = make(map[string]*[]string)
	wg := sync.WaitGroup{}
	sep := []byte("=>")
	err = loadFileList(&wg, "./data/kv", ".txt", func(filePath string) error {
		ret := make(map[string]string)
		key := make([]string, 0)
		err = cgghui.LoadFileLine(filePath, func(line []byte) bool {
			r := bytes.SplitN(line, sep, 2)
			k := string(bytes.Trim(r[0], " "))
			v := string(bytes.Trim(r[1], " "))
			ret[k] = v
			key = append(key, k)
			return true
		})
		mutex.Lock()
		name := strings.Replace(filepath.Base(filePath), ".txt", "", 1)
		kvList[name] = ret
		kvKeys[name] = &key
		mutex.Unlock()
		return err
	})
	wg.Wait()
	if err != nil {
		return err
	}
	// word
	wordList = make(map[string]*[]string)
	err = loadFileList(&wg, "./data/word", ".txt", func(filePath string) error {
		ret := make([]string, 0)
		err = cgghui.LoadFileLine(filePath, func(line []byte) bool {
			ret = append(ret, string(line))
			return true
		})
		mutex.Lock()
		wordList[strings.Replace(filepath.Base(filePath), ".txt", "", 1)] = &ret
		mutex.Unlock()
		return err
	})
	wg.Wait()
	return err
}

func RegisterTemplateFuncMap(t *template.Template) *template.Template {
	t.Funcs(template.FuncMap{
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"rand_bool": func() bool {
			return cgghui.RandomSliceString(&[]string{"0", "1"}) == "1"
		},
		"md5": func(str string) string {
			return cgghui.MD5(str)
		},
		"split": func(str string, sep string) []string {
			return strings.Split(str, sep)
		},
		"join": func(sep string, str ...string) string {
			return strings.Join(str, sep)
		},
		"now_timestamp": func() int64 {
			return time.Now().Unix()
		},
		"now": func() time.Time {
			return time.Now()
		},
		"substr": func(str string, l int, add ...string) string {
			r := []rune(str)
			if len(r) <= l {
				return str
			}
			str = string(r[0:l])
			if len(add) == 0 {
				return str
			}
			return str + add[0]
		},
		"replace": func(str string, old, new string, n int) string {
			return strings.Replace(str, old, new, n)
		},
		"get_word": func(name string, n int) []string {
			var r []string
			for i := 0; i < n; i++ {
				r = append(r, RandWord(name))
			}
			return r
		},
		"get_kv": func(name string, n int) []map[string]string {
			var r []map[string]string
			for i := 0; i < n; i++ {
				if v := RandKv(name); v != nil {
					r = append(r, v)
				}
			}
			return r
		},
		"rand_kv":   RandKv,
		"rand_word": RandWord,
		"rand_string": func(l, max int) string {
			if max != l {
				l = int(cgghui.RangeRand(int64(l), int64(max)))
			}
			return cgghui.RandomString(l, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		},
		"py": func(s string) string {
			return strings.Join(pinyin.LazyPinyin(s, pyArg), "")
		},
		"py_first": func(s string) string {
			first := ""
			for _, py := range pinyin.LazyPinyin(s, pyArg) {
				first += string(py[0])
			}
			return first
		},
	})
	return t
}

func RandWord(name string) string {
	if _, ok := wordList[name]; !ok {
		return ""
	}
	return cgghui.RandomSliceString(wordList[name])
}

func RandKv(name string) map[string]string {
	if _, ok := kvList[name]; !ok {
		return nil
	}
	k := cgghui.RandomSliceString(kvKeys[name])
	return map[string]string{"key": k, "val": kvList[name][k]}
}
