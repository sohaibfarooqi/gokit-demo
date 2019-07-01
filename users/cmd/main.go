package main

import (
  "flag"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  opentracing "github.com/opentracing/opentracing-go"
  zipkin "github.com/openzipkin/zipkin-go-opentracing"

  "github.com/sohaibfarooqi/fragbook/users/pkg"
  "github.com/go-kit/kit/log"
)
var fs = flag.NewFlagSet("users", flag.ExitOnError)
var httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
var zipkinURL = fs.String("zipkin-url", "http://zipkin:9411/api/v1/spans", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")

func main(){

  fs.Parse(os.Args)

  // Create a single logger, which we'll use and give to other components.
  var logger log.Logger
  logger = log.NewLogfmtLogger(os.Stderr)
  logger = log.With(logger, "ts", log.DefaultTimestampUTC)
  logger = log.With(logger, "caller", log.DefaultCaller)

  logger.Log("tracer", "Zipkin", "URL", *zipkinURL)
  collector, err := zipkin.NewHTTPCollector(*zipkinURL)
  if err != nil {
    logger.Log("zipkin http collector err", err)
    os.Exit(1)
  }
  defer collector.Close()

  recorder := zipkin.NewRecorder(collector, false, "http://users:8081", "users")
  tracer, err := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(true), zipkin.TraceID128Bit(true))
  if err != nil {
    logger.Log("unable to create zipking tracer", err)
    os.Exit(1)
  }
  opentracing.InitGlobalTracer(tracer)

  var s pkg.UsersService
  {
    s = pkg.NewInMemService()
  }

  var h http.Handler
  {
    h = pkg.MakeHttpHandler(s, log.With(logger, "component", "HTTP"), tracer)
  }

  errs := make(chan error)
  go func() {
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
    errs <- fmt.Errorf("%s", <-c)
  }()

  go func() {
    logger.Log("transport", "HTTP", "addr", *httpAddr)
    errs <- http.ListenAndServe(*httpAddr, h)
  }()

  logger.Log("exit", <-errs)
}

