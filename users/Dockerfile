FROM golang

RUN mkdir -p /go/src/github.com/sohaibfarooqi/fragbook
ADD . /go/src/github.com/sohaibfarooqi/fragbook
WORKDIR /go/src/github.com/sohaibfarooqi/fragbook/users

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT watcher -run github.com/sohaibfarooqi/fragbook/users/cmd -watch github.com/sohaibfarooqi/fragbook/users
