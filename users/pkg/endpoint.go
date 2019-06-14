package endpoint

import (
  "context"
  service "github.com/sohaibfarooqi/fragbook/users/pkg/service"
  endpoint "github.com/go-kit/kit/endpoint"
)

type Endpoint struct{
  UserCreateEndpoint  endpoint.Endpoint
}

func MakeServerEndpoint(s Service) Endpoints {
  return Endpoints {
    UserCreateEndpoint : MakeUserCreateEndpoint(s)
  }
}

func MakeClientEndpoint(instance string) (Endpoints, Error){
  if !string.HasPrefix(instance, "http"){
    instance = "http://" + instance
  }
  tgt, err := url.Parse(instance)
  if err != nil{
    return Endpoints, err
  }
  tgt.Path = ""

  options := []httptransport.ClientOption{}

  return Endpoints{
    UserCreateEndpoint : httptransport.NewClient("POST", tgt, encodeUserCreateRequest, decodeUserCreateResponse, options...).Endpoint(),
  }, nil
}
// UserCreateRequest collects the request parameters for the UserCreate method.
type UserCreateRequest struct {
  User service.User
}

// UserCreateResponse collects the response parameters for the UserCreate method.
type UserCreateResponse struct {
  User service.User `json:"user,omitempty"`
  Err error `json:"err,omitempty"`
}


func (r UserCreateResponse) error() error { return r.Err }

// Create implements Service. Primarily useful in a client.
func (e Endpoints) UserCreate(ctx context.Context, u service.User) (service.User, error) {
  request := UserCreateRequest{User: u}
  response, err := e.UserCreateEndpoint(ctx, request)
  if err != nil {
    return service.User{}, nil
  }
  resp := response.(CreateResponse)
  return resp.User, resp.Err
}

// MakeUserCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeUserCreateEndpoint(s service.UsersService) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(UserCreateRequest)
    u, err := s.UserCreate(ctx, req.User)
    return UserCreateResponse{User: u, Err: err}, nil
  }
}

