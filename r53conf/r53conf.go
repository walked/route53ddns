package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	configureMode := flag.Bool("configure", true, "Use this flag to configure AWS credentials")
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
}

func configure(c Config) {
	fmt.Print("AWS Access Key ID: ")
	var keyID, secret, zone, record string
	fmt.Scanln(&keyID)
	fmt.Print("AWS Secret Key: ")
	fmt.Scanln(&secret)
	fmt.Print("Zone ID: ")
	fmt.Scanln(&zone)
	fmt.Print("DNS Record: ")
	fmt.Scanln(&record)

	//Configure mapping for toml file
	// Note we encode access/secret keys to base64 to minimixe visibility.
	//**--!! This is not encryption do not treat it as such; it minimizes prying eyes and thats it!!--**
	var configuration = map[string]interface{}{
		"ZoneID":    zone,
		"AccessKey": base64.StdEncoding.EncodeToString([]byte(keyID)),
		"SecretKey": base64.StdEncoding.EncodeToString([]byte(secret)),
		"Record":    record,
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(configuration); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
	err := ioutil.WriteFile("conf.toml", buf.Bytes(), 0755)
	if err != nil {
		fmt.Println("Error Writing file")
		os.Exit(1)
	}

}
