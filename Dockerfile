FROM golang:1.21-alpine AS build

RUN apk add build-base npm
WORKDIR /app

# Setup early for caching
COPY --parents Makefile go.mod go.sum web/package.json web/package-lock.json .
RUN make setup

COPY . /app
RUN make build

FROM alpine:3
COPY --from=build /app/build/steam-hour-booster-ui /steam-hour-booster-ui

EXPOSE 8080

CMD [ "/steam-hour-booster-ui" ]
