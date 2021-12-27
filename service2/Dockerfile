FROM golang
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/hash-server


FROM alpine
WORKDIR /app
COPY --from=0 /app/app .
CMD ./app --port 8080 --host '0.0.0.0'
