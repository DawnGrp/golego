package config

import (
	"flag"
	"fmt"
	"golego/modules/helper"
	"io/ioutil"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	default_Config_Path = "./config.json"
)

var config map[string]gjson.Result

func GetInfo() helper.Info {
	return helper.Info{
		Name: "config",
	}
}

func init() {
	LoadConfig()
}

func Add(name string, value gjson.Result) {
	config[name] = value
}

func Get(name string) (value gjson.Result, ok bool) {
	value, ok = config[name]
	return
}

func LoadConfig() {

	configPath := flag.String("c", "", "custom config file path")
	flag.Parse()
	if configPath == nil || *configPath == "" {
		*configPath = default_Config_Path
	}

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	config = gjson.ParseBytes(configData).Map()

}

func SaveConfig() {

	configArray := []string{}

	for name, value := range config {
		configArray = append(configArray, fmt.Sprintf("\t\"%s\": %s", name, value.String()))
	}

	configString := fmt.Sprintf("{\n%s\n}", strings.Join(configArray, ",\n"))

	err := ioutil.WriteFile(default_Config_Path, []byte(configString), 0555)
	if err != nil {
		panic(err)
	}

	fmt.Println("write to config over")
}
