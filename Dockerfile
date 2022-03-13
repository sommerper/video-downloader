FROM golang:1.17 as dev
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN sh install.sh
RUN go build -v main.go
CMD ["/usr/src/app/main"]
# ENTRYPOINT ["tail", "-f", "/dev/null"]