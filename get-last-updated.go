package main 

import (
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"os"
)

type Latest struct { 
	Last_updated string
}


func main() {
    argLen := len(os.Args)
    if argLen < 2 {
	    showUsage()
	    os.Exit(1)
    }

    url := os.Args[1]   // https://hub.docker.com/v2/repositories/selenium/standalone-chrome/tags/latest/
    resp, err := http.Get(url)
    if err != nil {
	log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
	    log.Fatalln(err)
    }

    var jsonResp Latest
    sb := string(body)
    json.Unmarshal([]byte(sb), &jsonResp)
    fmt.Printf(jsonResp.Last_updated)
    //fmt.Println(resp.StatusCode)
}

func showUsage() {
	fmt.Println(`Usage: 
    get-last-updated TAG_URL

    TAG_URL -> URL for a container image which includes its last updated time (Required)

    Example Usage:
    $ get-last-updated https://hub.docker.com/v2/repositories/selenium/standalone-chrome/tags/latest/
    `)
}