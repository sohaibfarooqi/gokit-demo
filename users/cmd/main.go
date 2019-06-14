package main

import (
  "flag"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  zipkingoopentracing "github.com/openzipkin/zipkin-go-opentracing"

  "github.com/sohaibfarooqi/fragbook/users/pkg"
  "github.com/go-kit/kit/log"
)

var logger log.Logger

var(
  fs         = flag.NewFlagSet("users", flag.ExitOnError)
  debugAddr  = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
  httpAddr   = fs.String("http-addr", ":8081", "HTTP listen address")
  zipkinURL  = fs.String("zipkin-url", "", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")
)

func main(){
  // Create a single logger, which we'll use and give to other components.
  logger = log.NewLogfmtLogger(os.Stderr)
  logger = log.With(logger, "ts", log.DefaultTimestampUTC)
  logger = log.With(logger, "caller", log.DefaultCaller)

  logger.Log("tracer", "Zipkin", "URL", *zipkinURL)
  collector, err := zipkingoopentracing.NewHTTPCollector(*zipkinURL)
  if err != nil {
    logger.Log("err", err)
    os.Exit(1)
  }
  defer collector.Close()

  var s pkg.UsersService
  {
    s = pkg.NewInMemService()
  }

  var h http.Handler
  {
    h = pkg.MakeHttpHandler(s, log.With(logger, "component", "HTTP"))
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

