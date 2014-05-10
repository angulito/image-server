package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// ServerConfiguration struct
// Most of this configuration comes from json config
type ServerConfiguration struct {
	ServerPort            string   `json:"server_port"`
	StatusPort            string   `json:"status_port"`
	SourceDomain          string   `json:"source_domain"`
	WhitelistedExtensions []string `json:"whitelisted_extensions"`
	MaximumWidth          int      `json:"maximum_width"`
	LocalBasePath         string   `json:"local_base_path"`
	MantaBasePath         string   `json:"manta_base_path"`
	MantaConcurrency      int      `json:"manta_concurrency"`
	DefaultQuality        uint     `json:"default_quality"`
	GraphiteEnabled       bool     `json:"graphite_enabled"`
	GraphiteHost          string   `json:"graphite_host"`
	GraphitePort          int      `json:"graphite_port"`
	Environment           string   `json:"environment"`
	NamespaceMappings     map[string]string
	Events                *EventChannels
	Adapters              *Adapters
}

// EventChannels struct
// Available image processing/downloading events
type EventChannels struct {
	ImageProcessed              chan *ImageConfiguration
	ImageProcessedWithErrors    chan *ImageConfiguration
	OriginalDownloaded          chan *ImageConfiguration
	OriginalDownloadUnavailable chan *ImageConfiguration
}

type NamespaceMapping struct {
	Namespace string
	Source    string
}

func LoadServerConfiguration(path string) (*ServerConfiguration, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
		return nil, fmt.Errorf("configuration error: %v\n", err)
	}

	var config *ServerConfiguration
	json.Unmarshal(configFile, &config)
	config.Events = &EventChannels{
		ImageProcessed:     make(chan *ImageConfiguration),
		OriginalDownloaded: make(chan *ImageConfiguration),
	}
	return config, nil
}
