FROM golang as build
WORKDIR /code
COPY . .
ENV CGO_ENABLED=0
RUN go build .

FROM alpine
COPY --from=build /code/metal-sweeper-action /bin/metal-sweeper-action
CMD ["metal-sweeper-action"]
