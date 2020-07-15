# Debian buster (10) based image
FROM golang:1.14.4-buster

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build 

# ----- DONE BUILDING ----
# TODO: create steps

WORKDIR /dist

CMD curl -fLo air https://git.io/linux_air

RUN cp /build/hack-a-bot .

# Sensitive information in this file!!!!
# COPY .air.conf /.air.conf

# EXPOSE 3000

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["./air"]
