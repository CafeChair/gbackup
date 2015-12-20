package plugin

import (
	"encoding/json"
	"io/ioutil"
	"log"
	// "os"
	"strings"
	"sync"
)

type ScriptConfig struct {
	Dir       string
	LogDir    string
	RedisAddr string
}

type RedisConfig struct {
	Addr string
	Port string
}

type GlobalConfig struct {
	Script *ScriptConfig
	Redis  *RedisConfig
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ToString(filename string) (string, error) {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(str)), nil
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify config file")
	}
	ConfigFile = cfg
	configContent, err := ToString(cfg)
	if err != nil {
		log.Fatalln("read config file: ", cfg, "fail: ", err)
	}
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file: ", cfg, "fail: ", err)
	}
	lock.Lock()
	defer lock.Unlock()
	config = &c
}

// func IsExist(fp string) bool {
// 	_, err := os.Stat(fp)
// 	return err == nil || os.IsExist(err)
// }

// func IsFile(fp string) bool {
// 	f, e := os.Stat(fp)
// 	if e != nil {
// 		return false
// 	}
// 	return !f.IsDir()
// }
