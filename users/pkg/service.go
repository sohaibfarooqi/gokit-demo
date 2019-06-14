package pkg

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

