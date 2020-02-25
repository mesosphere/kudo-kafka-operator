package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	jsoniter "github.com/json-iterator/go"
	archiver "github.com/mholt/archiver/v3"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Connector Struct
type Connector struct {
	Resources []string    `json:"resources,omitempty" yaml:"resources,omitempty"`
	Config    interface{} `json:"config,omitempty" yaml:"config,omitempty"`
}

// ConfigFile Struct
type ConfigFile struct {
	Connectors map[string]Connector `json:"connectors,omitempty" yaml:"connectors,omitempty"`
	Resources  []string             `json:"resources,omitempty" yaml:"resources,omitempty"`
}

var (
	app        = kingpin.New("kafka-connectors", "A command-line kafka connectors configuration parser and setup helper.")
	config     = app.Flag("config", "Path to the configuration file.").Required().File()
	download   = app.Command("download", "Downloads the resources needed for connectors.")
	register   = app.Command("register", "Register the connectors in kafka-connect.")
	endpoint   = register.Arg("endpoint", "Kafka Connect REST endpoint").Required().String()
	configFile = &ConfigFile{}
)

func main() {
	parsed, _ := app.Parse(os.Args[1:])
	if *config != nil {
		b, _ := ioutil.ReadAll(*config)
		log.Printf("Parsing configuration: %s", string(b))
		if err := yaml.Unmarshal(b, configFile); err != nil {
			panic(err)
		}
		switch parsed {
		case "register":
			log.Printf("Registering to endpoint %s", *endpoint)
			for name, config := range configFile.Connectors {
				log.Printf("Parsing connector: %s", name)
				log.Printf("Registering connector configuration: %v\n", *endpoint, config.Config)
				RegisterConnector(*endpoint, config.Config)
			}

		case "download":
			downloadDirectory, _ := os.Getwd()
			for name, config := range configFile.Connectors {
				log.Printf("Parsing connector: %s", name)
				for _, resource := range config.Resources {
					log.Printf("Downloading file: %s\n", resource)
					filename, err := DownloadFile(downloadDirectory, resource)
					if err != nil {
						panic(err)
					}
					log.Printf("Extracting file: %s\n", resource)
					ExtractFile(path.Join(downloadDirectory, filename), downloadDirectory)
				}
			}
			log.Printf("Parsing download only connector resources")
			for _, resource := range configFile.Resources {
				log.Printf("Downloading file: %s\n", resource)
				filename, err := DownloadFile(downloadDirectory, resource)
				if err != nil {
					panic(err)
				}
				log.Printf("Extracting file: %s\n", resource)
				ExtractFile(path.Join(downloadDirectory, filename), downloadDirectory)
			}
		}
	}
}

// DownloadFile Download files to directory
func DownloadFile(downloadDirectory, url string) (string, error) {
	filename := path.Base(url)
	filepath := path.Join(downloadDirectory, filename)

	resp, err := http.Get(url)
	if err != nil {
		return filename, err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return filename, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filename, err
}

// ExtractFile Extract archives using p7zip
func ExtractFile(filepath, destination string) {
	err := archiver.Unarchive(filepath, destination)
	if err != nil {
		panic(err)
	}
}

// RegisterConnector register connectors to kafka-connect endpoint
func RegisterConnector(endpoint string, data interface{}) {
	if data == nil {
		return
	}
	payloadBytes, err := jsoniter.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", endpoint, "connectors"), body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
