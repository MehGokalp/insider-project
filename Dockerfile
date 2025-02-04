FROM golang:1.23-alpine AS build

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go mod download

# Build the Go app
RUN go build -o ./sbin/app .

FROM alpine:edge

COPY --from=build /src /root/go/src/github.com/mehgokalp/insider-project

EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/root/go/src/github.com/mehgokalp/insider-project/sbin/app", "server"]