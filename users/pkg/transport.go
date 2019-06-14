package pkg

import (
  "context"
  "bytes"
  "io/ioutil"
  "encoding/json"
  "errors"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/go-kit/kit/log"
  "github.com/go-kit/kit/transport"
  httptransport "github.com/go-kit/kit/transport/http"
)

var(
  ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHttpHandler(s UsersService, logger log.Logger) http.Handler {
  r := mux.NewRouter()
  e := MakeServerEndpoints(s)
  options := []httptransport.ServerOption{
    httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
    httptransport.ServerErrorEncoder(encodeError),
  }

  r.Methods("POST").Path("/users/").Handler(httptransport.NewServer(
    e.CreateUserEndpoint,
    decodeCreateUserRequest,
    encodeResponse,
    options...,
  ))

  return r
}

func decodeCreateUserRequest(_ context.Context, r *http.Request)(request interface{}, err error){
  var req UserCreateRequest
  if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
    return nil, e
  }
  return req, nil
}

func encodeCreateUserRequest(ctx context.Context, req *http.Request, request interface{}) error{
  req.URL.Path = "/users/"
  return encodeRequest(ctx, req, request)
}

func decodeCreateUserResponse(_ context.Context, resp *http.Response)(interface{}, error){
  var response UserCreateResponse
  err := json.NewDecoder(resp.Body).Decode(&response)
  return response, err
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
  var buf bytes.Buffer
  err := json.NewEncoder(&buf).Encode(request)
  if err != nil {
    return err
  }
  req.Body = ioutil.NopCloser(&buf)
  return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
  if e, ok := response.(errorer); ok && e.error() != nil {
    // service specific error
    encodeError(ctx, e.error(), w)
    return nil
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  return json.NewEncoder(w).Encode(response)
}

type errorer interface {error() error}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
  if err == nil {
    panic("encodeError with nil error")
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(codeFrom(err))
  json.NewEncoder(w).Encode(map[string]interface{}{
    "error": err.Error(),
  })
}

func codeFrom(err error) int {
  switch err {

    case ErrNotFound:
      return http.StatusNotFound

    case ErrAlreadyExists:
      return http.StatusBadRequest

    default:
      return http.StatusInternalServerError
  }
}
