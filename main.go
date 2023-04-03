package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
)

const CONFIG_FILE = "/etc/config/configmap-microservice-demo.yaml"

// const CONFIG_FILE = "configmap-microservice-demo.yaml"
const BIND = "0.0.0.0:8084"

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

/*
This is the struct that holds our application's configuration
*/
type Config struct {
	Message string `yaml:"message"`
}

/*
Simple Yaml Config file loader
*/
func loadConfig(configFile string) *Config {
	conf := &Config{}
	configData, err := ioutil.ReadFile(configFile)
	check(err)

	err = yaml.Unmarshal(configData, conf)
	check(err)
	log.Println(conf)
	return conf
}

func main() {
	confManager := NewMutexConfigManager(loadConfig(CONFIG_FILE))
	//confManager := NewChannelConfigManager(loadConfig(CONFIG_FILE))

	// Create a single GET Handler to print out our simple config message
	router := httprouter.New()
	router.GET("/", func(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
		conf := confManager.Get()
		fmt.Fprintf(resp, "%s", conf.Message)
	})

	// Watch the file for modification and update the config manager with the new config when it's available
	watcher, err := WatchFile(CONFIG_FILE, time.Second, func() {
		fmt.Printf("Configfile Updated\n")
		conf := loadConfig(CONFIG_FILE)
		confManager.Set(conf)
	})
	check(err)

	// Clean up
	defer func() {
		watcher.Close()
		confManager.Close()
	}()

	fmt.Printf("Listening on '%s'....\n", BIND)
	err = http.ListenAndServe(BIND, router)
	check(err)
}
