FROM golang
WORKDIR /app
COPY . .
RUN go mod tidy  \
    && go test ./... \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/server


FROM alpine
WORKDIR /app
COPY --from=0 /app/app .
CMD ./app 
