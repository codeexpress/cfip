# Is a given IP a Cloudflare IP address?
`cfip` is a script that checks if a given IP address is a Cloudflare IP address. It can be run as:
- a standlone tool
- a HTTP API

## Cloud
`cfip` runs as a API endpoint on api.code.express. To check IPv4 and IPv6 addresses respectively. Invoke it as follows following:

https://api.code.express/cfip/104.31.122.34

https://api.code.express/cfip/2606:4700:30::681b:804b

Please do not use this for production use cases or cause undue load on the server. If you need a more stable api, download and run a local copy by following the instructions below.

## Download

### Latest Releases
`cfip` is available for 32/64 bit linux, OS X and Windows systems.
Latest versions can be downloaded from the
[Release](https://github.com/codeexpress/cfip/releases) tab above.

### Using go install
If you have golang tools installed, you can install `cfip` as follows:
```
$ go install github.com/codeexpress/cfip@latest
```

### Build from source
This is a golang project with no dependencies. Assuming you have golang compiler installed,
the following will build the binary from scratch
```
$ git clone https://github.com/codeexpress/cfip
$ cd cfif
$ go build 
```

## Usage

### As a HTTP API
`cfip` is best used as a HTTP api listening on localhost. One downloaded, run this to start the server:
```
$ ./cfip -s 8080
```

Then, to check an IP address, use curl or browser to open a URL eg. `curl http://localhost:8080/cfip/<ip_address>`
Eg.
```sh
$ curl http://localhost:8080/cfip/104.31.122.34
{"ip": "104.31.122.34", "cloudflare_ip": "true"}
```

Also works for IPv6, eg.

![cfip server on localhost checking IPv6 address](https://user-images.githubusercontent.com/14211134/59569036-c5a4f000-9038-11e9-9cd0-03a053398cc7.png)

### As a standalone binary
`cfip` can be used as a standalone binary as well. Simply invoke as follows:
```sh
$ ./cfip 104.31.122.34
{"ip": "104.31.122.34", "cloudflare_ip": "true"}
```
