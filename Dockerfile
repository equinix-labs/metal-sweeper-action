FROM golang as build
WORKDIR /code
COPY . .
RUN go build .

FROM scratch
COPY --from=build /code/packet-sweeper-action /bin/packet-sweeper-action
ENTRYPOINT ["packet-sweeper-action"]
