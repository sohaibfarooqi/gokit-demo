FROM golang

RUN mkdir -p /root/.ssh
COPY users/id_rsa /root/.ssh/
RUN chmod 600 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n\tIdentityFile ~/.ssh/id_rsa" >> /root/.ssh/config

RUN export GIT_TERMINAL_PROMPT=1
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

RUN mkdir -p /go/src/github.com/sohaibfarooqi/fragbook
ADD . /go/src/github.com/sohaibfarooqi/fragbook
WORKDIR /go/src/github.com/sohaibfarooqi/fragbook/users

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run /go/src/github.com/sohaibfarooqi/fragbook/users/cmd  -watch /go/src/github.com/sohaibfarooqi/fragbook/users
