FROM golang:1.21-alpine AS build

RUN apk add build-base npm
WORKDIR /app

# Install go mods early for caching
COPY go.mod go.sum .
RUN go mod download

COPY . /app
RUN make build

FROM alpine:3
COPY --from=build /app/steam-hour-booster-ui /steam-hour-booster-ui

EXPOSE 8080

CMD [ "/steam-hour-booster-ui" ]
