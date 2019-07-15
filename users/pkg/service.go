package pkg

import (
  "context"
  "errors"
  "os"

  "github.com/go-pg/pg"
  opentracing "github.com/opentracing/opentracing-go"
)

// CRUD interface for users
type UsersService interface {
  Create(ctx context.Context, u User) (User, error)
}

type User struct {
  tableName struct{}  `sql:"users,alias:user"`

  FirstName   string  `sql:"firstname"`
  LastName    string  `sql:"lastname"`
  Email       string  `sql:"email"`
  Password    string  `sql:"password"`
}

var (
  ErrAlreadyExists   = errors.New("user already exists")
  ErrNotFound        = errors.New("user not found")
)

type PGService struct {
  db *pg.DB
}

func NewPGService() UsersService {
  return &PGService{
    db: pg.Connect(&pg.Options{
      Addr: os.Getenv("PG_HOST") + ":" + os.Getenv("PG_PORT"),
      User: os.Getenv("PG_USER"),
      Password: os.Getenv("PG_PASSWORD"),
      Database: os.Getenv("PG_DB"),
    }),
  }
}

func (s *PGService) Create(ctx context.Context, u User) (User, error) {
  span := opentracing.SpanFromContext(ctx)
  defer span.Finish()

  err := s.db.Insert(&u)

  if err != nil {
    return User{}, err
  }

  span.SetTag("user", u.Email)
  return u, nil
}

