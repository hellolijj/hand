package dd

import (
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var DdConf DDConf

type DDConf struct {
	AppKey         string `yaml:"AppKey"`
	AppSecret string `yaml:"AppSecret"`
	DDServerAddress           string `yaml:"DDServerAddress"`
}

func init() {
	yamlFile, err := ioutil.ReadFile("./conf/config.yaml")
	if err != nil {
		log.Fatal("init conf/dingding.yaml发生错误:", err)
	}
	err = yaml.Unmarshal(yamlFile, &DdConf)
	if err != nil {
		log.Fatal("init conf/dingding.yaml发生错误:", err)
	}
	return
}