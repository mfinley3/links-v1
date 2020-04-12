package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/influx"
	"github.com/go-pg/pg/v9"
	influxdb "github.com/influxdata/influxdb1-client/v2"

	linkrepository "github.com/mfinley3/links-v1/internal/links/postgres"
	links "github.com/mfinley3/links-v1/internal/links/service"
	transportHTTP "github.com/mfinley3/links-v1/internal/links/transport/http"
)

type bufWriter struct {
	buf bytes.Buffer
}

func (w *bufWriter) Write(bp influxdb.BatchPoints) error {
	for _, p := range bp.Points() {
		fmt.Fprintf(&w.buf, p.String()+"\n")
	}
	return nil
}

func extractAndPrintMessage(expected []string, msg string) error {
	for _, pattern := range expected {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(msg)
		if len(match) != 2 {
			return fmt.Errorf("pattern not found! {%s} [%s]: %v\n", pattern, msg, match)
		}
		fmt.Println(match[1])
	}
	return nil
}

func main() {

	opts, err := pg.ParseURL("postgres://postgres:postgres@localhost:5432/links?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db := pg.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	influx := influx.New(map[string]string{}, influxdb.BatchPointsConfig{
		Database: "db0",
	}, kitlog.NewNopLogger())

	ctx, cancel := context.WithCancel(context.Background())
	startMetricWriter(ctx, influx)
	redirectCounter := influx.NewCounter("redirect_counter")

	lr := linkrepository.New(db)
	ls := links.New(lr)

	errs := make(chan error, 1)

	go startServer(ls, redirectCounter, errs)

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-sig)
	}()

	err = <-errs
	cancel()
	log.Fatal("links-v1 terminated: " + err.Error())
}

func startServer(linksvc links.Service, redirectCounter *influx.Counter, errs chan error) {
	addr := fmt.Sprintf("127.0.0.1:8080")

	mux := chi.NewRouter()
	mux.Mount("/", transportHTTP.Handler(linksvc, redirectCounter))
	log.Printf("listening on: %s\n", addr)
	errs <- http.ListenAndServe(addr, mux)
}

func startMetricWriter(ctx context.Context, influx *influx.Influx) error {
	httpClient, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return err
	}

	go influx.WriteLoop(ctx, time.NewTicker(1*time.Second).C, httpClient)
	return nil
}
