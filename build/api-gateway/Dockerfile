FROM golang:1.22

WORKDIR /api-gateway/

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./api-gateway

EXPOSE 42069

CMD ["./api-gateway"]