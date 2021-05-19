package main

import (
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

var (
	apis = make([]*url.URL, 0)
	port = flag.Int("port", 8888, "listen http port")
)

func init() {
	martianLog.SetLevel(martianLog.Error)
	flag.Parse()
}

func main() {
	for _, s := range flag.Args() {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		apis = append(apis, u)
	}

	p := martian.NewProxy()
	defer p.Close()

	ca, privateKey, _ := mitm.NewAuthority("name", "org", 24*365*time.Hour)
	conf, _ := mitm.NewConfig(ca, privateKey)
	p.SetMITM(conf)

	//proxy, _ := url.Parse("http://localhost:8080")
	//p.SetDownstreamProxy(proxy)

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("starting listen on %s", l.Addr().String())

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
