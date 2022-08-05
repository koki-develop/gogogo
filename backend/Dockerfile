FROM golang:1.18 as build
WORKDIR /var/task

COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN GO111MODULE=on GOOS=linux go build -ldflags="-s -w" -o dist/api pkg/handlers/api/lambda/main.go

ENTRYPOINT ["/var/task/dist/api"]

# ---------------------------

FROM golang:1.18

COPY --from=build /var/task/dist/api /var/task/dist/api

ENTRYPOINT ["/var/task/dist/api"]
