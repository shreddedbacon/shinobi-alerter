package shinobiclient

// https://shinobi.video/docs/api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ShinobiClient interface {
	TriggerMotion(string) (string, error)
	GetMonitors(string) (string, error)
	GetVideos(string) (string, error)
	GetVideosById(string, string, string) (string, error)
	GetStartedMonitors(string) (string, error)
	RunRequest(string, string, string, string) (string, error)
}

type Shinobi struct {
	netClient     *http.Client
	shinobiConfig ShinobiConfig
}

type ShinobiConfig struct {
	Server  string          `json:"server"`
	Apikey  string          `json:"apikey"`
	Cameras []ShinobiCamera `json:"cameras"`
}

type ShinobiCamera struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Group  string `json:"group"`
	Region string `json:"region"`
}

func New(config string) ShinobiClient {
	shinobiConfig := ShinobiConfig{}
	netClientTimeout := 10
	var netClient = &http.Client{
		Timeout: time.Second * time.Duration(netClientTimeout),
	}
	// ignore ssl cert warnings
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	json.Unmarshal([]byte(config), &shinobiConfig)
	return &Shinobi{
		netClient:     netClient,
		shinobiConfig: shinobiConfig,
	}
}

func (sa *Shinobi) GetVideos(cameraGroup string) (string, error) {
	str, err := sa.RunRequest(
		"GET",
		sa.shinobiConfig.Server,
		"/"+sa.shinobiConfig.Apikey+"/videos/"+cameraGroup+"/",
		"",
	)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (sa *Shinobi) GetVideosById(cameraGroup string, cameraID string, query string) (string, error) {
	if query != "" {
		query = "?" + query
	}
	str, err := sa.RunRequest(
		"GET",
		sa.shinobiConfig.Server,
		"/"+sa.shinobiConfig.Apikey+"/videos/"+cameraGroup+"/"+cameraID+query,
		"",
	)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (sa *Shinobi) GetMonitors(cameraGroup string) (string, error) {
	str, err := sa.RunRequest(
		"GET",
		sa.shinobiConfig.Server,
		"/"+sa.shinobiConfig.Apikey+"/monitor/"+cameraGroup+"/",
		"",
	)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (sa *Shinobi) GetStartedMonitors(cameraGroup string) (string, error) {
	str, err := sa.RunRequest(
		"GET",
		sa.shinobiConfig.Server,
		"/"+sa.shinobiConfig.Apikey+"/smonitor/"+cameraGroup+"/",
		"",
	)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (sa *Shinobi) TriggerMotion(host string) (string, error) {
	for _, camera := range sa.shinobiConfig.Cameras {
		if camera.IP == host {
			str, err := sa.RunRequest(
				"GET",
				sa.shinobiConfig.Server,
				"/"+sa.shinobiConfig.Apikey+"/motion/"+camera.Group+"/"+camera.Name+"?data={\"plug\":\""+camera.Name+"\",\"name\":\""+camera.Region+"\",\"reason\":\"motion\",\"confidence\":200}",
				"",
			)
			if err != nil {
				return "", err
			}
			return str, nil
		}
	}
	return "", nil
}

func (sa *Shinobi) RunRequest(httpMethod string, URL string, apiPath string, queryParam string) (string, error) {
	apiRequest, apiRequestErr := http.NewRequest(httpMethod, URL+apiPath, bytes.NewBuffer([]byte(queryParam)))
	if apiRequestErr != nil {
		return "", apiRequestErr
	}
	fmt.Println(apiRequest)
	apiResponse, apiResponseErr := sa.netClient.Do(apiRequest)
	if apiResponseErr != nil {
		return "", apiResponseErr
	}
	defer apiResponse.Body.Close()
	apiResponseBody, _ := ioutil.ReadAll(apiResponse.Body)
	apiReturnBody := string(apiResponseBody)
	if apiResponse.StatusCode != 200 {
		return apiReturnBody, errors.New("error performing check or connecting to shinobi")
	}
	return apiReturnBody, nil
}
