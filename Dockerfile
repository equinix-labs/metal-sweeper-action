FROM golang as build
WORKDIR /code
COPY . .
ENV CGO_ENABLED=0
RUN go build .

FROM alpine
COPY --from=build /code/packet-sweeper-action /bin/packet-sweeper-action
CMD ["packet-sweeper-action"]
