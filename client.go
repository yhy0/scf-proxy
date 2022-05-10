package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/martian/v3"
	martianLog "github.com/google/martian/v3/log"
	"github.com/google/martian/v3/mitm"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	Port    int
	ApiFile string
}

var (
	apis = make([]*url.URL, 0)
	o    = Options{}
)

func init() {
	martianLog.SetLevel(martianLog.Error)

	flag.IntVar(&o.Port, "port", 8888, "listen http port")
	flag.StringVar(&o.ApiFile, "aF", "", "scf api file ")

	flag.Parse()
}

func open(path string) (lines []string, Error error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func main() {
	var (
		list []string
		err  error
	)

	if o.ApiFile == "" {
		list = flag.Args()
	} else {
		list, err = open(o.ApiFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, s := range list {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		apis = append(apis, u)
	}

	l, err := net.Listen("tcp", fmt.Sprintf("[::]:%d", o.Port))
	if err != nil {
		log.Fatalf("net.Listen(): got %v, want no error", err)
	}

	log.Infof("starting listen on %s", l.Addr().String())

	p := martian.NewProxy()
	defer p.Close()

	p.SetTimeout(15 * time.Second)

	// Test TLS server.
	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", 24*365*time.Hour)
	if err != nil {
		log.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}
	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		log.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}

	p.SetMITM(mc)

	p.SetRequestModifier(new(T))

	go p.Serve(l)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

type T struct {
	martian.RequestModifier
}

func (T) ModifyRequest(req *http.Request) error {
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Error(err)
		return err
	}

	api := apis[rand.Intn(len(apis))]
	newReq, _ := http.NewRequest(http.MethodPost, api.String(), strings.NewReader(base64.StdEncoding.EncodeToString(b)))
	newReq.URL = api
	newReq.Header.Set("Url", req.URL.String())
	*req = *newReq
	return nil
}
