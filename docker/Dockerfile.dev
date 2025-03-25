FROM golang:1.24-alpine

WORKDIR /app

ENV APP_ENV=production

COPY . .

RUN go mod download

RUN go build -o main ./main.go

EXPOSE 8080

CMD [ "./main" ]
