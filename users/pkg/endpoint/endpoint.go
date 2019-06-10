package endpoint

import (
	"context"
	service "github.com/sohaibfarooqi/fragbook/users/pkg/service"
	endpoint "github.com/go-kit/kit/endpoint"
)

// CreateRequest collects the request parameters for the Create method.
type CreateRequest struct {
	User service.User
}

// CreateResponse collects the response parameters for the Create method.
type CreateResponse struct {
	User service.User `json:"user,omitempty"`
	Err error `json:"err,omitempty"`
}

// MakeCreateEndpoint returns an endpoint that invokes Create on the service.
func MakeCreateEndpoint(s service.UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		u, err := s.Create(ctx, req.User)
		return CreateResponse{User: u, Err: err}, nil
	}
}

// Failed implements Failer.
func (r CreateResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Create implements Service. Primarily useful in a client.
func (e Endpoints) Create(ctx context.Context, u service.User) (service.User, error) {
	request := CreateRequest{User: u}
	response, err := e.CreateEndpoint(ctx, request)
	if err != nil {
		return service.User{}, nil
	}
	resp := response.(CreateResponse)
	return resp.User, resp.Err
}
