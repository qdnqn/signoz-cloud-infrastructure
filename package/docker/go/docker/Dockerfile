FROM golang:1.19-alpine as build

WORKDIR /service

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o app

FROM alpine:3.11.3

WORKDIR /service
COPY --from=build /service/app /service/app
COPY configs/ /service/configs

EXPOSE 8080

CMD [ "/service/app" ]
