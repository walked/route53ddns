package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

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
