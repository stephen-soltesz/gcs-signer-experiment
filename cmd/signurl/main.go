package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

var (
	keyfile  = flag.String("key", "", "Path fo name of JSON key file.")
	datafile = flag.String("data", "", "Path fo name of data file.")
	bucket   = flag.String("bucket", "soltesz-urlsign-mlab-sandbox",
		"Name of GCS bucket (without gs:// prefix), the given key must "+
			"have write privileges for this bucket.")
)

// Store key in datastore.
// Fetch key and extract private_key during start up.
// Load whitelist from a set of machine inventory.
// Accept requests from whitelisted machines.
// Generate and return Signed URLs
// That's it.

func main() {
	flag.Parse()

	pkey, err := ioutil.ReadFile(*keyfile)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]string{}
	err = json.Unmarshal(pkey, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	url, err := storage.SignedURL(*bucket, *datafile, &storage.SignedURLOptions{
		GoogleAccessID: data["client_id"],
		PrivateKey:     []byte(data["private_key"]),
		Method:         "PUT",
		Expires:        time.Now().Add(48 * time.Hour),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
	c := &http.Client{}
	in, err := os.Open(*datafile)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", url, in)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", resp)
}
