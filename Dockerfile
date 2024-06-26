FROM golang:1.20

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /usr/local/bin/app cmd/api/main.go
EXPOSE 8080
CMD ["/usr/local/bin/app", "/usr/src/app/config/local.yml"]
