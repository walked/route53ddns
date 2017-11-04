package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kardianos/service"
)

//Config is simply for defining the configuration of the download for utilization within datimsync
type Config struct {
	ZoneID    string
	Record    string
	AccessKey string
	SecretKey string
}

var logger service.Logger

type program struct{}

var stopControl bool

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	stopControl = false
	go p.run()
	return nil
}
func (p *program) run() {
	configureMode := flag.Bool("configure", false, "Use this flag to configure AWS credentials")
	flag.Parse()
	//Parse config file to Config struct
	//conf.toml should live in the same directory as the compiled exe
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)

	c, err := ioutil.ReadFile(exePath + "/conf.toml")
	if err != nil {
		fmt.Println("Error Reading Config")
		os.Exit(1)
	}
	var conf Config
	if _, err := toml.Decode(string(c), &conf); err != nil {
		fmt.Print("Cannot decode conf.toml")
		os.Exit(1)
	}

	// Decodes config file access key and sets temp environment variable for S3 actions
	accessKey, _ := base64.StdEncoding.DecodeString(conf.AccessKey)
	secretKey, _ := base64.StdEncoding.DecodeString(conf.SecretKey)
	os.Setenv("AWS_ACCESS_KEY_ID", string(accessKey))
	os.Setenv("AWS_SECRET_ACCESS_KEY", string(secretKey))

	if *configureMode {
		fmt.Println("Configure Mode")
		configure(conf)
		os.Exit(0)
	}

	currentIP, err := checkIP()
	createRecord(conf, currentIP)
	var ip string

	for {
		ip, err = checkIP()
		if currentIP != ip {
			createRecord(conf, ip)
			currentIP = ip
			fmt.Println("Updated IP")

		} else {
			fmt.Println("No update")
		}
		if err != nil {
			fmt.Println(err)
		}
		if stopControl == true {
			break
		}
		time.Sleep(15 * time.Second)
	}
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	stopControl = true
	return nil
}

func main() {

	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
