FROM golang


RUN mkdir -p /go/src/github.com/sohaibfarooqi/fragbook
ADD . /go/src/github.com/sohaibfarooqi/fragbook
WORKDIR /go/src/github.com/sohaibfarooqi/fragbook/users

RUN mkdir /root/.ssh/
ADD users/id_rsa /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

RUN go get -t -v ./...
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT watcher -run github.com/sohaibfarooqi/fragbook/users/cmd -watch github.com/sohaibfarooqi/fragbook/users
