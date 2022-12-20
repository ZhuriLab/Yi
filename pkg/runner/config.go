package runner

import (
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"bytes"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

/**
  @author: yhy
  @since: 2022/12/12
  @desc: //TODO
**/

type QLFile struct {
	GoQL   []string `mapstructure:"go_ql"`
	JavaQL []string `mapstructure:"java_ql"`
}

var QLFiles *QLFile

// HotConf 使用 viper 对配置热加载
func HotConf() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logging.Logger.Fatalf("cmd.HotConf, fail to get current path: %v", err)
	}
	// 配置文件路径 当前文件夹 + config.yaml
	configFile := path.Join(dir, ConfigFileName)
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	// watch 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		logging.Logger.Infoln("config file changed: ", e.Name)
		oldQls := QLFiles
		ReadYamlConfig(configFile)
		newQls := QLFiles
		// 规则更新时，从数据库中获取项目，用新规则跑一遍
		NewRules(oldQls, newQls)
	})
}

// Init 加载配置
func Init() {
	//配置文件路径 当前文件夹 + config.yaml
	configFile := path.Join(Pwd, ConfigFileName)

	// 检测配置文件是否存在
	if !utils.Exists(configFile) {
		WriteYamlConfig(configFile)
		logging.Logger.Infof("%s not find, Generate profile.", configFile)
	} else {
		logging.Logger.Infoln("Load profile ", configFile)
	}
	ReadYamlConfig(configFile)

}

func ReadYamlConfig(configFile string) {
	// 加载config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to read 'config.yaml': %+v", err)
	}
	err = viper.Unmarshal(&QLFiles)
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to parse 'config.yaml', check format: %v", err)
	}
}

func WriteYamlConfig(configFile string) {
	// 生成默认config
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultYamlByte))
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to read default config bytes: %v", err)
	}
	// 写文件
	err = viper.SafeWriteConfigAs(configFile)
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to write 'config.yaml': %v", err)
	}
}
