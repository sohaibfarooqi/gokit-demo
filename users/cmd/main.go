package main

import (
  "flag"
  "fmt"
  "net"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  "github.com/sohaibfarooqi/fragbook/users/pkg"
  "github.com/go-kit/kit/log"

  "github.com/prometheus/client_golang/prometheus"
  endpoint "github.com/go-kit/kit/endpoint"
  kitprom "github.com/go-kit/kit/metrics/prometheus"
  opentracing "github.com/opentracing/opentracing-go"
  zipkin "github.com/openzipkin/zipkin-go-opentracing"
  promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

)

var httpAddr = flag.String("http-addr", ":8081", "HTTP listen address")
var debugAddr = flag.String("debug-addr", ":8080", "HTTP listen address")
var zipkinURL = flag.String("zipkin", "", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")

func main(){

  flag.Parse()

  // Create a single logger, which we'll use and give to other components.
  var logger log.Logger
  logger = log.NewLogfmtLogger(os.Stderr)
  logger = log.With(logger, "ts", log.DefaultTimestampUTC)
  logger = log.With(logger, "caller", log.DefaultCaller)

  logger.Log("Zipkin URL", *zipkinURL)
  collector, err := zipkin.NewHTTPCollector(*zipkinURL)
  if err != nil {
    logger.Log("zipkin http collector err", err)
    os.Exit(1)
  }
  defer collector.Close()

  recorder := zipkin.NewRecorder(collector, false, "http://users:8081", "users")
  tracer, err := zipkin.NewTracer(
    recorder,
    zipkin.ClientServerSameSpan(true),
    zipkin.TraceID128Bit(true),
  )
  if err != nil {
    logger.Log("unable to create zipking tracer", err)
    os.Exit(1)
  }

  opentracing.InitGlobalTracer(tracer)
  HttpSummaryMiddleware(logger)

  var s pkg.UsersService
  {
    s = pkg.NewPGService()
  }

  var h http.Handler
  {
    h = pkg.MakeHttpHandler(s, log.With(logger, "component", "HTTP"), tracer)
  }

  errs := make(chan error)
  go func() {
    http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
    debugListener, err := net.Listen("tcp", *debugAddr)
    if err != nil {
      logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
    }
    errs <- http.Serve(debugListener, http.DefaultServeMux)
  }()

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

func HttpSummaryMiddleware(logger log.Logger) (mw map[string][]endpoint.Middleware){
  mw = map[string][]endpoint.Middleware{}
  duration := kitprom.NewSummaryFrom(prometheus.SummaryOpts{
    Help:      "Request duration in seconds.",
    Name:      "request_duration_seconds",
    Namespace: "example",
    Subsystem: "users",
  }, []string{"method", "success"})
  mw["Create"] = []endpoint.Middleware{pkg.LoggingMiddleware(log.With(logger, "method", "Create")), pkg.InstrumentingMiddleware(duration.With("method", "Create"))}
  return mw
}
