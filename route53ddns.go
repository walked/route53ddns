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

//Config is simply for defining the configuration the access / service info
type Config struct {
	ZoneID    string
	Record    string
	AccessKey string
	SecretKey string
}

//Configuration - global configuration
var Configuration Config

// Service Stuff
var logger service.Logger

type program struct{}

var stopControl bool

func (p *program) Start(s service.Service) error {
	stopControl = false
	go p.run()
	return nil
}
func (p *program) run() {

	currentIP, err := checkIP()
	createRecord(Configuration, currentIP)
	var ip string

	for {
		ip, err = checkIP()
		if currentIP != ip {
			createRecord(Configuration, ip)
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
		time.Sleep(15 * time.Minute)
	}
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	stopControl = true
	return nil
}

func main() {

	configureMode := flag.Bool("configure", false, "Use this flag to configure AWS credentials")
	installMode := flag.Bool("install", false, "Use this flag to install the Windows Service")
	uninstallMode := flag.Bool("uninstall", false, "Use this flag to uninstall the Windows Service")
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
	//var conf Config
	if _, err := toml.Decode(string(c), &Configuration); err != nil {
		fmt.Print("Cannot decode conf.toml")
		os.Exit(1)
	}

	// Decodes config file access key and sets temp environment variable for S3 actions
	accessKey, _ := base64.StdEncoding.DecodeString(Configuration.AccessKey)
	secretKey, _ := base64.StdEncoding.DecodeString(Configuration.SecretKey)
	os.Setenv("AWS_ACCESS_KEY_ID", string(accessKey))
	os.Setenv("AWS_SECRET_ACCESS_KEY", string(secretKey))

	if *configureMode {
		fmt.Println("Configure Mode")
		configure(Configuration)
		os.Exit(0)
	}

	svcConfig := &service.Config{
		Name:        "R53DDNSSRV",
		DisplayName: "Route53DDNS Service",
		Description: "Service for the Purpose of Performing DDNS to Route53",
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

	if *installMode {
		fmt.Println("Install Mode")
		s.Install()
		os.Exit(0)
	}
	if *uninstallMode {
		fmt.Println("Uninstall Mode")
		s.Stop()
		time.Sleep(5 * time.Second)
		s.Uninstall()
		logger.Info("uninstalled")
		os.Exit(0)
	}

	//s.Install() to install
	//s.Uninstall() to uninstall
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
