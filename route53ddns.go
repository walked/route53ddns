package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

//Config is simply for defining the configuration of the download for utilization within datimsync
type Config struct {
	ZoneID    string
	Record    string
	AccessKey string
	SecretKey string
}

func main() {

	configureMode := flag.Bool("configure", false, "Use this flag to configure AWS credentials")
	flag.Parse()
	//Parse config file to Config struct
	//conf.toml should live in the same directory as the compiled exe
	c, err := ioutil.ReadFile("conf.toml")
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

	ip, err := checkIP()
	createRecord(conf, ip)

}
