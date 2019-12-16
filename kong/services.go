package kong

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/timoheijne/kong-auto-register/config"
)

type Service struct {
	Name string `json:"name"`
	Port int    `json:"port"`
	Host string `json:"host"`
	Path string `json:"path"`
}

// compareService will check if the service exists in kong if not create
func compareService() {
	svc := config.GetConfig().Service
	log.Println("Comparing service:", svc.Name)

	resp, err := http.Get(config.GetConfig().AdminEndpoint + "/services/" + svc.Name)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		createService(svc)
	} else {
		data := Service{}
		json.Unmarshal(readBody(resp.Body), &data)

		runServiceCompare(data, svc)
	}
}

// runServiceCompare compares an existing service in kong with our config file defined
func runServiceCompare(kong Service, configSrv config.Service) {
	if kong.Name != configSrv.Name {
		return
	}

	p, _ := strconv.Atoi(configSrv.Port)
	if kong.Port != p || kong.Host != configSrv.Host || kong.Path != configSrv.Path {
		log.Println("Service", configSrv.Name, "changed... Patching gateway")
		patchService(configSrv)
	} else {
		log.Println("Service", configSrv.Name, "already exists and up to date")
	}
}

// patchService will update the values of the already existing service
func patchService(svc config.Service) {
	log.Println("Patching service:", svc.Name)
	hc := http.Client{}

	form := url.Values{}
	form.Add("host", svc.Host)
	form.Add("port", svc.Port)
	form.Add("path", svc.Path)

	req, err := http.NewRequest("PATCH", config.GetConfig().AdminEndpoint+"/services/"+svc.Name, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	check(err)

	resp, err := hc.Do(req)
	check(err)

	defer resp.Body.Close()
}

// createService... Name says all
func createService(srv config.Service) {
	log.Println("Creating service:", srv.Name)
	hc := http.Client{}

	form := url.Values{}
	form.Add("name", srv.Name)
	form.Add("host", srv.Host)
	form.Add("port", srv.Port)
	form.Add("path", srv.Path)

	req, err := http.NewRequest("POST", config.GetConfig().AdminEndpoint+"/services", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	check(err)

	resp, err := hc.Do(req)
	check(err)

	defer resp.Body.Close()
}
