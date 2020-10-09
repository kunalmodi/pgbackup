FROM golang:1.14-alpine AS builder

RUN apk update
RUN apk add --no-cache git curl

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/pgbackup

FROM alpine

RUN apk update
RUN apk add --no-cache postgresql-client

COPY --from=builder /go/bin/pgbackup /go/bin/pgbackup

CMD ["/go/bin/pgbackup"]
