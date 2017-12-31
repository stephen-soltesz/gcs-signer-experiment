package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

var (
	keyfile = flag.String("key", "", "Path fo name of PRM file.")
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
	url, err := storage.SignedURL("soltesz-urlsign-mlab-sandbox", "sample.txt", &storage.SignedURLOptions{
		GoogleAccessID: data["client_id"],
		PrivateKey:     []byte(data["private_key"]),
		Method:         "PUT",
		Expires:        time.Now().Add(48 * time.Hour),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
}
