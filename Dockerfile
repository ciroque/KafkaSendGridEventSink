############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git pkgconfig bash build-base
RUN git clone https://github.com/edenhill/librdkafka.git && \
	cd librdkafka && \
	./configure --prefix /usr && \
	make && \
	make install

WORKDIR $GOPATH/src/kafka-sendgrid-event-sink/
COPY . .

# Fetch dependencies.
# Using go get.
WORKDIR $GOPATH/src/kafka-sendgrid-event-sink/cmd/main/
RUN go get -d -v

# Build the binary.
WORKDIR $GOPATH/src/kafka-sendgrid-event-sink/
RUN go build -tags static -o /go/bin/kafka-sendgrid-event-sink cmd/main/main.go

############################
# STEP 2 build a small image
############################
FROM alpine:3.10

# Copy our static executable.
COPY --from=builder /go/bin/kafka-sendgrid-event-sink /go/bin/kafka-sendgrid-event-sink

# Run the hello binary.
ENTRYPOINT ["/go/bin/kafka-sendgrid-event-sink"]
