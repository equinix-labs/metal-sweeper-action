# syntax=docker/dockerfile:1
FROM golang:1 as build
WORKDIR /code
COPY . .
# Need to manually bump this before each release
ENV ACTION_VERSION=0.5.0
ENV CGO_ENABLED=0
RUN go build -ldflags "-X 'main.version=${ACTION_VERSION}'"

FROM alpine:3
COPY --from=build /code/metal-sweeper-action /bin/metal-sweeper-action
CMD ["metal-sweeper-action"]