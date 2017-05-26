package bergamot

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func init() {
	// fmt.Println("init")
	// configName := "app"
	// if os.Getenv("CONFIG") != "" {
	// 	configName = os.Getenv("CONFIG")
	// }
	// fmt.Println(configName)
	// viper.SetConfigFile(configName)
	// viper.AddConfigPath("$HOME")
	// viper.AddConfigPath("$WORKDIR")
	// viper.AddConfigPath(".")
	// if os.Getenv("WORKDIR") != "" {
	// 	fmt.Println("has workdir", os.Getenv("WORKDIR"))
	// 	viper.AddConfigPath(os.Getenv("WORKDIR"))
	// }

}

type Datasource interface{}
type Server interface{}

type Parser interface {
	GetDatasources(*viper.Viper) map[string]Datasource
}

// App app name
type App struct {
	Datasources map[string]Datasource
	Servers     map[string]Server
	Component   string
	Config      *viper.Viper
}

// New bootstrap an app with a provided the configuration
func New(config *viper.Viper, parser Parser) error {
	err := config.ReadInConfig()
	if err != nil {
		return err
	}

	keys := config.AllKeys()
	fmt.Println(keys)
	for _, k := range keys {
		fmt.Println(k)
	}

	// fmt.Println("all stttings")
	for k, v := range config.AllSettings() {
		fmt.Println(k, ":", v)
	}

	data := parser.GetDatasources(config)
	fmt.Println(data)
	// fmt.Println(config.AllSettings())

	return err
}

// Auto start and bootstrap service
func Auto() error {
	v := viper.New()
	// setting configuration file name
	configName := "config"
	if os.Getenv("CONFIG") != "" {
		configName = os.Getenv("CONFIG")
	}
	v.SetConfigName(configName)
	// setting configuration file paths
	v.AddConfigPath("$HOME")
	v.AddConfigPath(".")
	if os.Getenv("WORKDIR") != "" {
		v.AddConfigPath("$WORKDIR")
	}
	v.AutomaticEnv()
	return New(v, DefaultParser{})
}

// DefaultParser default parser
type DefaultParser struct{}

func (DefaultParser) GetDatasources(config *viper.Viper) map[string]Datasource {
	ds := make(map[string]Datasource)
	var settings interface{}
	for k, v := range config.AllSettings() {
		if k == "datasources" {
			settings = v
			break
		}
	}
	fmt.Println("got this")
	fmt.Println(settings)

	return ds
}
