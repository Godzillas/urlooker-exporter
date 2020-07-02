package g

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/toolkits/file"
)

type WebConf struct {
	ReportAddr string `json:"reportaddrs"`
	Index      string `json:"index"`
	Type       string `json:"type"`
	Interval   int    `json:"interval"`
	Timeout    int    `json:"timeout"`
}

type GlobalConfig struct {
	Debug          bool     `json:"debug"`
	InstanceRegion string   `json:"instanceRegion"`
	GetUrlAddr     string   `json:"geturladdrs"`
	Worker         int      `json:"worker"`
	ItemAddr       string   `json:"itemAddr"`
	CacheFile      string   `json:"cacheFile"`
	CacheTimeSec   string   `json:"cacheTimeSec"`
	DNSServer      string   `json:"dnsServer"`
	ListenPort     string   `json:"listenPort"`
	Web            *WebConf `json:"web"`
}

var (
	Config *GlobalConfig
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	Config = &c

	log.Println("load configuration file", cfg, "successfully")
	log.Println("listen at ", Config.ListenPort)
	return nil
}
