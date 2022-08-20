# We need a golang build environment first 
FROM golang:1.18-alpine as builder

WORKDIR /build

COPY . .

RUN go mod tidy

WORKDIR /build/cmd
RUN go build -o ./app .

# We use a Docker multi-stage build here in order that we only take the compiled go executable
FROM alpine:latest
 
COPY --from=0 "/build/cmd/app" app
 
ENTRYPOINT ./app