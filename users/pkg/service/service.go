package service

import (
  "context"
  "errors"
  "sync"
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

var ErrAlreadyExists   = errors.New("user already exists")

type InMemService struct {
  mtx sync.RWMutex
  m map[string]User
}

// Inmemory service implementation
func NewInMemService() UsersService {
  return &InMemService{
    m: map[string]User{},
  }
}

func (s *InMemService) Create(ctx context.Context, u User) (User, error) {
  s.mtx.Lock()
  defer s.mtx.Unlock()
  if _, ok := s.m[u.Email]; ok {
    return User{}, ErrAlreadyExists // POST = create, don't overwrite
  }
  s.m[u.Email] = u
	return u, nil
}

// New returns a UsersService with all of the expected middleware wired in.
func New(middleware []Middleware) UsersService {
	var svc UsersService = NewInMemService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
