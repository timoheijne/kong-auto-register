package kong

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.liquidstudio.nl/deeplens/kong-auto-registration/config"
)

func checkKong() {
	resp, err := http.Get(config.GetConfig().AdminEndpoint)
	check(err)

	if resp.StatusCode != 200 {
		log.Fatalln("Kong ADMIN API not found")
	}
}

func Run() {
	checkKong()

	compareService()
	compareRoutes()
}

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func readBody(body io.ReadCloser) []byte {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes
}
