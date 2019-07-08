CREATE TABLE users (
  id           BIGSERIAL       PRIMARY KEY,
  FirstName    varchar(20),
  LastName     varchar(20),
  Email        varchar(20),
  Password     varchar(40)
);
