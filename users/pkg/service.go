package pkg

import (
  "context"
  "errors"
  "sync"

  opentracing "github.com/opentracing/opentracing-go"
)

// CRUD interface for users
type UsersService interface {
  Create(ctx context.Context, u User) (User, error)
}

type User struct {
  FirstName   string
  LastName    string
  Email       string
  Password    string
}

var (
  ErrAlreadyExists   = errors.New("user already exists")
  ErrNotFound        = errors.New("user not found")
)

type InMemService struct {
  mtx sync.RWMutex
  m map[string]User
}

// Inmemory service implementation
func NewInMemService() UsersService {
  return &InMemService{
    m: map[string] User{},
  }
}

func (s *InMemService) Create(ctx context.Context, u User) (User, error) {
  span := opentracing.SpanFromContext(ctx)
  defer span.Finish()
  s.mtx.Lock()
  defer s.mtx.Unlock()
  if _, ok := s.m[u.Email]; ok {
    span.SetTag("error", ErrAlreadyExists)
    return User{}, ErrAlreadyExists // POST = create, don't overwrite
  }
  s.m[u.Email] = u
  span.SetTag("user", u.Email)
  return u, nil
}

