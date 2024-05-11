FROM golang:1.21.5
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/server -o /build/server
CMD ["/build/server"]
