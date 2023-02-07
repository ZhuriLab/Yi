package runner

import (
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/**
  @author: yhy
  @since: 2022/10/13
  @desc: //TODO
**/

type Options struct {
	Target   string
	Targets  string
	Token    string
	Proxy    string
	UserName string
	Pwd      string
	Path     string
	Port     string
	Thread   int
	Session  *utils.Session
}

var Option Options

type DirName struct {
	ZipDir    string
	DbDir     string
	GithubDir string
	ResDir    string
}

var DirNames DirName

var Pwd string

var ProgressBar map[string]float32

// Languages Codeql支持的语言，这里要根据机器上的配置，若要支持其他语言，请自行安装语言，以及指定对应语言的 codeql 规则
var Languages = []string{"Go", "Java"}

func ParseArguments() {

	flag.StringVar(&Option.Target, "t", "", "target")
	flag.StringVar(&Option.Targets, "f", "", "target file")
	flag.StringVar(&Option.Token, "token", "", "github personal access token")
	flag.StringVar(&Option.Proxy, "p", "", "http(s) proxy. eg: http://127.0.0.1:8080")
	flag.StringVar(&Option.UserName, "user", "yhy", "http username")
	flag.StringVar(&Option.Pwd, "pwd", "", "http pwd")
	flag.StringVar(&Option.Path, "path", "", "codeql path")
	flag.StringVar(&Option.Port, "port", "8888", "web port(default: 8888)")
	flag.IntVar(&Option.Thread, "thread", 5, "thread")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if Option.Path == "" {
		logging.Logger.Errorln("Please specify the rule base address of codeql.\neg: -path /Users/yhy/CodeQL/codeql")
		os.Exit(1)
	}

	if !strings.HasSuffix(Option.Path, "/") && !strings.HasSuffix(Option.Path, "\\") {
		if os.Getenv("GOOS") == "windows" {
			Option.Path += "\\"
		} else {
			Option.Path += "/"
		}
	}

	if Option.Pwd == "" {
		Option.Pwd = utils.RandStr()
		logging.Logger.Infof("Web access account password:%s/%s", Option.UserName, Option.Pwd)
	}

	Pwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	today := time.Now().Format("2006-01-02") + "/"
	DirNames = DirName{
		ZipDir:    Pwd + "/db/zip/" + today,
		ResDir:    Pwd + "/db/results/" + today,
		DbDir:     Pwd + "/db/database/" + today,
		GithubDir: Pwd + "/github/" + today,
	}

	os.MkdirAll(DirNames.ZipDir, 0755)
	os.MkdirAll(DirNames.ResDir, 0755)
	os.MkdirAll(DirNames.DbDir, 0755)
	os.MkdirAll(DirNames.GithubDir, 0755)

	Option.Session = utils.NewSession(Option.Proxy)

	// 生成配置文件，并监控更改
	Init()
	HotConf()

	ProgressBar = make(map[string]float32)
}
