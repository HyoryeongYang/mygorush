package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"runtime"
	"strings"
)

var defaultConf = []byte(`
core:
  enabled: true # enable httpd server
  address: "" # ip address to bind (default: any)
  shutdown_timeout: 30 # default is 30 second
  port: "8080" # ignore this port number if auto_tls is enabled (listen 443).
  max_notification: 100
  worker_num: 0 # default worker number is runtime.NumCPU()
  queue_num: 0 # default queue number is 8192

api:
  push_uri: "/api/push"

android:
  enabled: true
  apikey: ""
  max_retry: 0 # resend fail notification, default value zero is disabled

log:
  format: "string" # string or json
  access_log: "stdout" # stdout: output to console, or define log path like "log/access_log"
  access_level: "debug"
  error_log: "stderr" # stderr: output to console, or define log path like "log/error_log"
  error_level: "error"
  hide_token: true

`)

type ConfYaml struct {
	Core    SectionCore    `yaml:"core"`
	API     SectionAPI     `yaml:"api"`
	Android SectionAndroid `yaml:"android"`
	Log     SectionLog     `yaml:"log"`
}

type SectionCore struct {
	Enabled         bool           `yaml:"enabled"`
	Address         string         `yaml:"address"`
	ShutdownTimeout int64          `yaml:"shutdown_timeout"`
	Port            string         `yaml:"port"`
	MaxNotification int64          `yaml:"max_notification"`
	WorkerNum       int64          `yaml:"worker_num"`
	QueueNum        int64          `yaml:"queue_num"`
}

type SectionAPI struct {
	PushURI  string  `yaml:"push_uri"`
}

type SectionAndroid struct {
	Enabled    bool     `yaml:"enabled"`
	APIKey    string   `yaml:"apikey"`
	MaxRetry  int      `yaml:"max_notification"`
}

type SectionLog struct {
	Format      string `yaml:"format"`
	AccessLog   string `yaml:"access_log"`
	AccessLevel string `yaml:"access_level"`
	ErrorLog    string `yaml:"error_log"`
	ErrorLevel  string `yaml:"error_level"`
	HideToken   bool   `yaml:"hide_token"`
}

func LoadConfig (ConfPath string) (ConfYaml, error) {
	var conf ConfYaml

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()         // read in environment variables that match
	viper.SetEnvPrefix("mygorush") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)
		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		// Search config in home directory with name ".gorush" (without extension).
		viper.AddConfigPath("/etc/gorush/")
		viper.AddConfigPath("$HOME/.gorush")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
			// load default config
			return conf, err
		}
	}

	// Core
	conf.Core.Address = viper.GetString("core.address")
	conf.Core.Port = viper.GetString("core.port")
	conf.Core.ShutdownTimeout = int64(viper.GetInt("core.shutdown_timeout"))
	conf.Core.Enabled = viper.GetBool("core.enabled")
	conf.Core.MaxNotification = int64(viper.GetInt("core.max_notification"))

	// Api
	conf.API.PushURI = viper.GetString("api.push_uri")

	// Android
	conf.Android.Enabled = viper.GetBool("android.enabled")
	conf.Android.APIKey = viper.GetString("android.apikey")
	conf.Android.MaxRetry = viper.GetInt("android.max_retry")

	// log
	conf.Log.Format = viper.GetString("log.format")
	conf.Log.AccessLog = viper.GetString("log.access_log")
	conf.Log.AccessLevel = viper.GetString("log.access_level")
	conf.Log.ErrorLog = viper.GetString("log.error_log")
	conf.Log.ErrorLevel = viper.GetString("log.error_level")
	conf.Log.HideToken = viper.GetBool("log.hide_token")

	if conf.Core.WorkerNum == int64(0) {
		conf.Core.WorkerNum = int64(runtime.NumCPU())
	}

	if conf.Core.QueueNum == int64(0) {
		conf.Core.QueueNum = int64(8192)
	}

	return conf, nil
}



