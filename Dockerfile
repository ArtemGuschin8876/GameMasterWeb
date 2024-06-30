FROM --platform=linux/amd64 golang:latest

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
ADD go.mod go.sum .
RUN go mod download
ADD . .
RUN go build ./cmd/api

CMD ./api
