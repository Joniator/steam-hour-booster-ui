FROM golang:1.21-alpine AS build

RUN apk add build-base nodejs

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY . /app
RUN make build

FROM alpine:3
COPY --from=build /go/steam-hour-booster-ui /steam-hour-booster-ui

EXPOSE 8080

CMD [ "/steam-hour-booster-ui" ]
