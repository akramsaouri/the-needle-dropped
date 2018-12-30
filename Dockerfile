FROM golang:1.9 as builder
# Install dep tool
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep
WORKDIR /go/src/github.com/akramsaouri/gocker/
COPY Gopkg.toml Gopkg.lock ./
# Install dependencies without checking for go code
RUN dep ensure -vendor-only
COPY src ./src
# Build binary file
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o app ./src

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Use build artifact from build stage
COPY --from=builder /go/src/github.com/akramsaouri/gocker/app . 

# Set port var and pass it to the http server
ENV PORT 80 
EXPOSE $PORT

CMD ["./app"]
