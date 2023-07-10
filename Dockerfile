FROM golang:1.20 AS builder

WORKDIR /build

COPY . .
RUN go build -o ./program cmd/program.go 

FROM ubuntu:latest
WORKDIR /app 
COPY --from=builder /build/program /app/

CMD [ "/app/program" ]


