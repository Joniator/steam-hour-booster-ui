FROM golang:1.21-alpine AS build

COPY . /go

RUN apk add build-base
RUN make build

FROM alpine:3
COPY --from=build /go/steam-hour-booster-ui /app/steam-hour-booster-ui

ENTRYPOINT /app/steam-hour-booster-ui 
