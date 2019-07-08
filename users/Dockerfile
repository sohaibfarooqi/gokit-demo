FROM golang


RUN mkdir -p /go/src/github.com/sohaibfarooqi/fragbook
ADD . /go/src/github.com/sohaibfarooqi/fragbook
WORKDIR /go/src/github.com/sohaibfarooqi/fragbook/users

RUN mkdir /root/.ssh/
ADD users/id_rsa /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

RUN curl -sSL https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
RUN echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ bionic main" > /etc/apt/sources.list.d/migrate.list
RUN apt-get update && apt-get install -y migrate

RUN go get -t -v ./...

ENTRYPOINT ["go", "run", "./cmd/main.go", "-zipkin", "http://zipkin:9411/api/v1/spans"]
