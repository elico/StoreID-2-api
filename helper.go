package main

// Resources:
// - http://www.jokecamp.com/blog/examples-of-creating-base64-hashes-using-hmac-sha256-in-different-languages/
// - https://gobyexample.com/base64-encoding

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	//"net"
	// Net libs can be used to parse IP addresses
	"net/http"
	"net/url"
	"os"
	"strings"
)

var debug *bool
var storeIdApiUrl *string
var user *string
var pass *string
var tlsConfig *tls.Config
var tlsCert *string
var pemCert []byte
var dontVerifyTls *bool
var useOsTls *bool
var http_version *string

var err error

func process_request(line string) {
	answer := "ERR"
	lparts := strings.Split(strings.TrimRight(line, "\n"), " ")
	if len(lparts[0]) > 0 {
		if *debug {
			fmt.Fprintln(os.Stderr, "ERRlog: Proccessing request => \""+strings.TrimRight(line, "\n")+"\"")
		}
	}

	var client *http.Client
	switch *http_version {
	case "2":
		client = &http.Client{
			Transport: &http2.Transport{TLSClientConfig: tlsConfig},
		}
	default:
		client = &http.Client{}
	}

	testurl, _ := url.Parse(*storeIdApiUrl)
	testurlVals := url.Values{}
	testurlVals.Set("url", base64.StdEncoding.EncodeToString([]byte(lparts[1])))
	//testurlVals.Set("other", port)
	//testurlVals.Set("otherother", srcip)
	testurl.RawQuery = testurlVals.Encode()

	request, err := http.NewRequest("GET", testurl.String(), nil)
	request.Close = true
	request.SetBasicAuth(*user, *pass)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERRlog: reporting a http connection error1 => \""+err.Error()+"\"")
		fmt.Println(lparts[0] + " " + answer)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERRlog: reporting a http connection error2 => \""+err.Error()+"\"")
		fmt.Println(lparts[0] + " " + answer)
		return
	}

	if body != nil {
		response := string(body)
		if strings.HasPrefix(response, "OK ") {
			fmt.Println(lparts[0] + " " + response)
		}
	}

	fmt.Println(lparts[0] + " " + answer)
}

func init() {

	fmt.Fprintln(os.Stderr, "ERRlog: Starting Fake helper")

	debug = flag.Bool("d", false, "Debug mode can be \"yes\" or something else for no")
	storeIdApiUrl = flag.String("api-url", "http://ngtech.co.il/fake-storeid-url/", "The url of the api")
	user = flag.String("api-user", "admin", "Basic auth username for server authentication")
	pass = flag.String("api-pass", "admin", "Basic auth password for server authentication")
	http_version = flag.String("api-httpv", "1", "http client version: 1\\2")
	tlsCert = flag.String("tlscert", "cert.pem", "tls certificate")
	dontVerifyTls = flag.Bool("skiptls", false, "Verify tls certificate, use \"1\" to enable")
	useOsTls = flag.Bool("ostls", false, "Use OS tls certificates, use \"1\" to enable")

	flag.Parse()
	flagsMap := make(map[string]interface{})
	flagsMap["debug"] = *debug
	flagsMap["api-url"] = *storeIdApiUrl
	flagsMap["api_user"] = *user
	flagsMap["api_pass"] = *pass
	flagsMap["api-httpv"] = *http_version
	flagsMap["tlscert"] = *tlsCert
	flagsMap["skiptls"] = *dontVerifyTls
	flagsMap["skiptls"] = *useOsTls

	if *debug {
		fmt.Fprintln(os.Stderr, "ERRlog: Config Variables:")
		for k, v := range flagsMap {
			fmt.Fprintf(os.Stderr, "ERRlog:\t%v =>  %v\n", k, v)
		}
	}
}

func main() {
	if *http_version == "2" && strings.HasPrefix(*storeIdApiUrl, "http://") {
		fmt.Fprintf(os.Stderr, "ERRlog: ### The http2 library doesn't support \"https://\" scheme, you are using => %v\n", *storeIdApiUrl)
		return
	}

	switch {
	case *http_version == "2" && *useOsTls:

	case (*http_version == "2") && !*dontVerifyTls:
		tlsConfig = &tls.Config{RootCAs: x509.NewCertPool()}

		var err error
		pemCert, err = ioutil.ReadFile(*tlsCert)
		if err != nil {
			fmt.Println(err)
		}
		ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemCert)
		if !ok {
			panic("Couldn't load PEM data")
		}
	case (*http_version == "2" && (*dontVerifyTls)):
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	default:

	}

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// You may check here if err == io.EOF
			break
		}
		if strings.HasPrefix(line, "q") || strings.HasPrefix(line, "Q") {
			fmt.Fprintln(os.Stderr, "ERRlog: Exiting cleanly")
			break
		}

		go process_request(line)

	}
}
