package pkg

import (
  "context"
  "net/url"
  "strings"

  "github.com/go-kit/kit/endpoint"
  httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
  CreateUserEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s UsersService) Endpoints {
  return Endpoints {
    CreateUserEndpoint : MakeUserCreateEndpoint(s),
  }
}

func MakeClientEndpoints(instance string) (Endpoints, error){
  if !strings.HasPrefix(instance, "http"){
    instance = "http://" + instance
  }
  tgt, err := url.Parse(instance)
  if err != nil{
    return Endpoints{}, err
  }
  tgt.Path = ""

  options := []httptransport.ClientOption{}

  return Endpoints{
    CreateUserEndpoint : httptransport.NewClient("POST", tgt, encodeCreateUserRequest, decodeCreateUserResponse, options...).Endpoint(),
  }, nil
}
// UserCreateRequest collects the request parameters for the UserCreate method.
type UserCreateRequest struct {
  User User
}

// UserCreateResponse collects the response parameters for the UserCreate method.
type UserCreateResponse struct {
  User User `json:"user,omitempty"`
  Err error `json:"err,omitempty"`
}


func (r UserCreateResponse) error() error { return r.Err }

// Create implements Service. Primarily useful in a client.
func (e Endpoints) UserCreate(ctx context.Context, u User) (User, error) {
  request := UserCreateRequest{User: u}
  response, err := e.CreateUserEndpoint(ctx, request)
  if err != nil {
    return User{}, nil
  }
  resp := response.(UserCreateResponse)
  return resp.User, resp.Err
}

// MakeUserCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeUserCreateEndpoint(s UsersService) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(UserCreateRequest)
    u, err := s.Create(ctx, req.User)
    return UserCreateResponse{User: u, Err: err}, nil
  }
}

