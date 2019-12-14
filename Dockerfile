############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git tree
WORKDIR $GOPATH/src/KafkaSendGridEventSink/
COPY . .
RUN pwd && tree

# Fetch dependencies.
# Using go get.
WORKDIR $GOPATH/src/KafkaSendGridEventSink/cmd/main/
RUN go get -d -v

# Build the binary.
WORKDIR $GOPATH/src/KafkaSendGridEventSink/
RUN go build -o /go/bin/KafkaSendGridEventSink cmd/main/main.go

############################
# STEP 2 build a small image
############################
FROM alpine:3.10

# Copy our static executable.
COPY --from=builder /go/bin/KafkaSendGridEventSink /go/bin/KafkaSendGridEventSink

# Run the hello binary.
ENTRYPOINT ["/go/bin/KafkaSendGridEventSink"]
