FROM golang:1.18-alpine as builder

# RUN go install golang.org/dl/go1.18@latest \
#     && go1.18 download

WORKDIR /build

COPY . .

RUN go mod tidy

WORKDIR /build/cmd
RUN go build -o ./app .

ENTRYPOINT ["/build/cmd/app"]