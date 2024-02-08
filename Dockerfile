# syntax=docker/dockerfile:1

FROM golang:1.22 AS build
WORKDIR /app
COPY ./src/go.mod ./src/go.sum .
RUN go mod download
COPY ./src/ ./
RUN mkdir -p ./build
RUN CGO_ENABLED=0 GOOS=linux go build -o build/olap-srv ./cmd

FROM debian:12
WORKDIR /
COPY --from=build /app/build/olap-srv ./olap-srv
EXPOSE 8080
CMD ["./olap-srv"]
