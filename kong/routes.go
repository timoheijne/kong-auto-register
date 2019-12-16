package kong

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/timoheijne/kong-auto-register/config"
)

type Route struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Paths   []string `json:"paths"`
	Methods []string `json:"methods"`
}

type RouteResponse struct {
	Next string  `json:"next"`
	Data []Route `json:"data"`
}

var existingRoutes []Route

// First we get all routes because we might also want to delete the ones no longer relevant
func compareRoutes() {
	svc := config.GetConfig().Service
	log.Println("Comparing routes for service:", svc.Name)

	hasNext := true
	url := "/services/" + config.GetConfig().Service.Name + "/routes"
	for hasNext {
		data := getRoutes(url)
		existingRoutes = append(existingRoutes, data.Data...)

		if data.Next != "" {
			url = data.Next
		} else {
			hasNext = false
		}
	}

	// Loop through routes for compare.

	for _, r := range config.GetConfig().Routes {
		if _, ok := findKongRoute(r); ok == true { // Check if route exists in KONG already
			// It does. Validate to see if its same as config defined
			if !validateRoute(r) {
				// Its not up to date. Commence patching
				log.Println("Route:", r.Name, "changed... Patching gateway")
				patchRoute(r)
			} else {
				log.Println("Route:", r.Name, "already exists and is up to date")
			}
		} else {

			log.Println("Route:", r.Name, "does not exist. Creating")
			// Create Route
			createRoute(r)
		}
	}

	// Loop through existing routes to check if some might be deleted
	for _, r := range existingRoutes {

		// Check if the existing route exists in our config file
		if !isRouteInConfig(r) {
			log.Println("Route:", r.Name, "is removed, Deleting from gateway")
			deleteRoute(r.Id)
		}
	}

}

func getRoutes(url string) RouteResponse {
	resp, err := http.Get(config.GetConfig().AdminEndpoint + url)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		panic("Service does not exist")
	} else {
		data := RouteResponse{}
		json.Unmarshal(readBody(resp.Body), &data)
		return data
	}
}

// Check if route exists in KONG
func findKongRoute(route config.Route) (Route, bool) {
	for _, r := range existingRoutes {
		if r.Name == route.Name {
			return r, true
		}
	}

	return Route{}, false
}

// Compare kong route with our config
func isRouteInConfig(route Route) bool {
	for _, r := range config.GetConfig().Routes {
		if r.Name == route.Name {
			return true
		}
	}

	return false
}

// Check if kong route is same as config defined routes
func validateRoute(route config.Route) bool {
	r, _ := findKongRoute(route)

	methods := route.Methods
	if len(route.Methods) > 0 && route.Methods[0] == "" {
		methods = []string{}
	}

	if r.Name != route.Name || !reflect.DeepEqual(r.Paths, route.Paths) || !reflect.DeepEqual(r.Methods, methods) {
		return false
	}

	return true
}

func patchRoute(route config.Route) {
	hc := http.Client{}

	form := url.Values{}
	for _, method := range route.Methods {
		form.Add("methods[]", method)
	}

	for _, path := range route.Paths {
		form.Add("paths[]", path)
	}

	req, err := http.NewRequest("PATCH", config.GetConfig().AdminEndpoint+"/routes/"+route.Name, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	check(err)

	resp, err := hc.Do(req)
	check(err)

	defer resp.Body.Close()
}

func deleteRoute(id string) {
	hc := http.Client{}

	req, err := http.NewRequest("DELETE", config.GetConfig().AdminEndpoint+"/routes/"+id, strings.NewReader(""))
	check(err)

	resp, err := hc.Do(req)
	check(err)

	defer resp.Body.Close()
}

func createRoute(route config.Route) {
	hc := http.Client{}

	// TODO: Move validation to another step.
	if len(route.Methods) == 0 {
		panic("No method provided, Please give a HTTP method like GET, POST, PATCH, PUT, DELETE or any other valid")
	} else if len(route.Paths) == 0 {
		panic("No paths provided")
	}

	form := url.Values{}
	form.Add("name", route.Name)
	form.Add("service.name", config.GetConfig().Service.Name)

	for _, method := range route.Methods {
		form.Add("methods[]", method)
	}

	for _, path := range route.Paths {
		form.Add("paths[]", path)
	}

	req, err := http.NewRequest("POST", config.GetConfig().AdminEndpoint+"/routes/", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	check(err)

	resp, err := hc.Do(req)
	check(err)

	defer resp.Body.Close()
}
