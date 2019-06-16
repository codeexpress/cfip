package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Version = "1.0.1"
)

var (
	serverPtr = flag.Int("s", 0,
		"Create a `HTTP API server`")
	Urls = [...]string{
		"https://www.cloudflare.com/ips-v4",
		"https://www.cloudflare.com/ips-v6",
	}
)

type CFIP struct {
	t     time.Time
	cidrs []net.IPNet
}

var cfip = CFIP{
	t:     time.Now().AddDate(0, 0, -1), // yesterday's date, purposefully stale
	cidrs: []net.IPNet{},
}

func main() {
	initFlags()

	if *serverPtr != 0 {
		setupServer()
	} else {
		fmt.Println(checkIP(os.Args[1]))
	}
}

// provide a handler and start listening on the given port
func setupServer() {
	http.HandleFunc("/cfip/", func(w http.ResponseWriter, r *http.Request) {
		input := strings.Split(r.URL.Path, "/")[2]

		//log request on stdout
		fmt.Print(time.Now().Format("Mon Jan _2 15:04:05 2006"))
		fmt.Printf(" => %s %s %s\n", r.Header.Get("X-Forwarded-For"), r.Method, r.URL)

		if input == "ips-v4.gz" {
			// Read static file and send it
			// A cron job creates this file every 6 hours
			// $ curl -s https://www.cloudflare.com/ips-v4 | ./cidr2ip | gzip > ips-v4.gz
			ips, _ := ioutil.ReadFile("ips-v4.gz")
			w.Write(ips)
		} else { //an ip address, presumably
			fmt.Fprintf(w, checkIP(input)+"\n")
		}
	})
	port_str := strconv.Itoa(*serverPtr)
	fmt.Println("Serving on port " + port_str + ". To query, either:")
	fmt.Println("Use curl. Eg. 'curl http://localhost:" + port_str + "/cfip/104.32.122.34' or")
	fmt.Println("open the URL in the browser")
	http.ListenAndServe("127.0.0.1:"+port_str, nil)
}

// takes in a ip address and checks if it belongs to
// CF ranges
func checkIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "Invalid IP address provided"
	}

	cidrs := cfcidrs()
	result := `{"ip": "` + ip.String() + `", "cloudflare_ip": "`
	for i := 0; i < len(cidrs); i++ {
		if cidrs[i].Contains(ip) {
			result += `true"}`
			return result
		}
	}
	result += `false"}`
	return result
}

// check cache
// if not stale, return from cache
// if stale, update cache and return
func cfcidrs() []net.IPNet {
	// refresh cache if it is stale (> 6 hrs. old)
	timeDiff := time.Now().Sub(cfip.t).Hours()
	if timeDiff > 6 {
		fmt.Fprintf(os.Stderr, "Cached IPs are "+strconv.FormatFloat(timeDiff, 'G', 4, 64)+" hours stale, refetching\n")
		updateCache()
	}
	return cfip.cidrs
}

// fetch and save CIDRs from Cloudflare's website
func updateCache() {
	var cidrs []net.IPNet

	for _, url := range Urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't GET CF's ip ranges from their website\n")
			// don't update cache, use previous values
			return
		}
		defer resp.Body.Close()
		//fmt.Println("Response status:", resp.Status)

		if resp.Status != "200 OK" {
			fmt.Fprintf(os.Stderr, "non 200 HTTP code from CF\n")
			// don't update cache, use previous values
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan(); i++ {
			//fmt.Println(net.ParseCIDR(scanner.Text()))
			_, cidr, err := net.ParseCIDR(scanner.Text())
			if err != nil {
				panic(err)
			}
			cidrs = append(cidrs, *cidr)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
	cfip.cidrs = cidrs
	cfip.t = time.Now()
}

func initFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Cloudflare IP cfip %s\n", Version)
		fmt.Fprintf(os.Stderr, "Usage:   $ cfip <IP address to check> \n")
		fmt.Fprintf(os.Stderr, "Usage:   $ cfip -s <port> \n")
		fmt.Fprintf(os.Stderr, "Example: $ cfip 172.64.0.10\n")
		fmt.Fprintf(os.Stderr, "Example: $ cfip -s 8080\n")
		fmt.Fprintf(os.Stderr, "------------\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *serverPtr == 0 && len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Neither test IP address was provided nor the program is run in server mode.\n\n")
		flag.Usage()
		os.Exit(0)
	}

	if *serverPtr != 0 && len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "Either provide a test IP or run in server mode, not both.\n\n")
		flag.Usage()
		os.Exit(0)
	}
}
