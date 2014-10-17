package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	var dbFlag = flag.String("db", "", "mmdb gz file")
	var ipFlag = flag.String("ip", "", "ip address")

	flag.Parse()

	if *dbFlag == "" || *ipFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	ip := net.ParseIP(*ipFlag)

	db, err := openDB(*dbFlag)
	if err != nil {
		fmt.Println("openDB failed:", err)
		os.Exit(1)
	}
	defer db.Close()

	city, err := db.City(ip)
	if err != nil {
		fmt.Println("db.City failed:", err)
		os.Exit(1)
	}

	fmt.Printf("%#v", city)
}

func openDB(file string) (db *geoip2.Reader, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return geoip2.FromBytes(data)
}
