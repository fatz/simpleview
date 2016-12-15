package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/jsonq"
	"github.com/namsral/flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Simpleview app
type Simpleview struct {
	IcingaEndpoint string
	IcingaUsername string
	IcingaPassword string
	ProjectsList   string
	Projects       []string
	Cache          Overview
	CacheTimeout   time.Duration
	CacheRequest   time.Time
	Debug          bool
}

type Host struct {
	Name string `json:"name"`
}

func (simpleview *Simpleview) readJSON(jsonstring string) (iresp *jsonq.JsonQuery, err error) {

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(jsonstring))
	err = dec.Decode(&data)
	iresp = jsonq.NewQuery(data)
	return
}

// workaround for the new format of sending an error object
// we simply build our own empty results
func (simpleview *Simpleview) buildEmptyResults() (iresp *jsonq.JsonQuery) {
	iresp, _ = simpleview.readJSON(`{"results": []}`)
	return
}

func (simpleview *Simpleview) getRequest(urlString string, filter string) (iresp *jsonq.JsonQuery, err error) {
	// encodedFilter, err := url.Parse(filter)
	if err != nil {
		return
	}
	geturl := fmt.Sprintf("%s/%s?filter=%s", simpleview.IcingaEndpoint, urlString, filter)
	myurl, err := url.Parse(geturl)
	myurl.User = url.UserPassword(simpleview.IcingaUsername, simpleview.IcingaPassword)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(myurl.String())

	if simpleview.Debug {
		fmt.Printf("%s\n", geturl)
	}

	if err != nil {
		err = fmt.Errorf("Request error: %s", err)
		return
	}

	jsonstring, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("ReadAll error: %s", err)
		return
	}
	iresp, err = simpleview.readJSON(string(jsonstring))

	if err != nil {
		err = fmt.Errorf("Read error: %s\n\n %s", err, jsonstring)
		return
	}

	// so lets check if we just get an error object
	if resp.StatusCode == http.StatusNotFound {
		_, err1 := iresp.Int("error")
		_, err2 := iresp.String("status")

		if err1 == nil && err2 == nil {
			//we found an expected error object lets just return empty results
			iresp = simpleview.buildEmptyResults()
			return
		}
	}

	return
}

func (simpleview *Simpleview) getHosts() (jq []map[string]interface{}, err error) {
	url := "v1/objects/hosts"
	// filter := "host.last_hard_state>0&host.last_in_downtime==false&host.acknowledgement<2"
	// filter := "host.last_hard_state>0"

	filter := `host.last_hard_state%3E0%26%26host.last_in_downtime==false%26%26host.acknowledgement%3C2`

	iresp, err := simpleview.getRequest(url, filter)
	if err != nil {

		err = fmt.Errorf("request error %s", err)
		return
	}

	jq, err = iresp.ArrayOfObjects("results")

	if err != nil {
		err = fmt.Errorf("object error %s", err)
		return
	}

	return

}

func (simpleview *Simpleview) getServices() (jq []map[string]interface{}, err error) {
	url := "v1/objects/services"
	// filter := "host.last_hard_state>0&host.last_in_downtime==false&host.acknowledgement<2"
	// filter := "host.last_hard_state>0"

	filter := `host.state==HostUp%26%26service.last_hard_state%3E0%26%26service.acknowledgement%3C2%26%26service.last_in_downtime==false`

	iresp, err := simpleview.getRequest(url, filter)
	if err != nil {
		err = fmt.Errorf("request error %s", err)
		return
	}

	jq, err = iresp.ArrayOfObjects("results")

	// fmt.Printf("%s", iresp)

	if err != nil {
		err = fmt.Errorf("object error %s", err)
		return
	}

	return

}

type Overview struct {
	Hosts    []map[string]interface{} `json:"hosts"`
	Services []map[string]interface{} `json:"services"`
}

func (simpleview *Simpleview) getOverview(c *gin.Context) {
	if simpleview.CacheRequest.Sub(time.Now())*-1 > simpleview.CacheTimeout {
		var overview Overview
		hosts, err := simpleview.getHosts()
		if simpleview.Debug {
			fmt.Print(err)
		}

		if err != nil {
			c.JSON(500, fmt.Sprintf("error: %v", err))
			return
		}

		overview.Hosts = hosts
		services, err := simpleview.getServices()

		if err != nil {
			c.JSON(500, fmt.Sprintf("error: %v", err))
			return
		}
		overview.Services = services
		simpleview.CacheRequest = time.Now()
		simpleview.Cache = overview
	}

	c.JSON(200, simpleview.Cache)
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}

func main() {
	var simpleview Simpleview
	f := flag.NewFlagSetWithEnvPrefix(os.Args[0], "SIMPLEVIEW", 0)
	f.StringVar(&simpleview.IcingaEndpoint, "icingaendpoint", "http://api-icinga2.mycompany.example.com:5665", "the url to the icinga2 api")
	f.StringVar(&simpleview.IcingaUsername, "icingausername", "simpleview", "username to authenticate against icinga")
	f.StringVar(&simpleview.IcingaPassword, "icingapassword", "password", "password")
	f.StringVar(&simpleview.ProjectsList, "projects", "", "filter on this comma-sep. list of projects")
	f.BoolVar(&simpleview.Debug, "debug", false, "turn on debug output")

	simpleview.CacheTimeout = 30 * time.Second
	simpleview.CacheRequest = time.Now().Add(-simpleview.CacheTimeout)

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		f.PrintDefaults()
		os.Exit(1)
	}
	f.Parse(os.Args[1:])

	var port string

	port = os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": "OK",
		})
	})

	r.GET("/v1/overview", simpleview.getOverview)

	// r.GET("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "public"}))
	r.Use(static.Serve("/", BinaryFileSystem("public")))

	listen := fmt.Sprintf("0.0.0.0:%s", port)

	r.Run(listen)

	fmt.Printf("listening on: %s", listen)
}
