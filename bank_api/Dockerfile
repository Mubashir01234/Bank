# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY . /app
#RUN go mod download

WORKDIR /app/cmd

RUN go build -o /community

# Run the tests in the container
FROM build-stage AS test-stage

COPY --from=build-stage /app /app
WORKDIR /app
RUN go test -v ./..

# Deploy the application binary into a lean image
FROM quay.io/centos/centos:stream9-minimal AS release-stage

WORKDIR /

COPY --from=build-stage /community /app
COPY --from=build-stage /app/config /config
COPY --from=build-stage /app/internal/kafka/schemas/*.avsc /internal/kafka/schemas/

EXPOSE 8080

ENTRYPOINT ["/app"]
