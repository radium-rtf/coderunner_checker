FROM golang:1.22.5-alpine3.20 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22.5-alpine3.20 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/checker

# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
CMD ["/app"]