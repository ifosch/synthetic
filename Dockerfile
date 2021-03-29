FROM golang:alpine AS builder

RUN apk add --update \
	git

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

# Add all the source code (except what's ignored
# under `.dockerignore`) to the build context.
COPY . /code

# Build final artifact
RUN set -ex && \
	CGO_ENABLED=0 go build -o synthetic cmd/main.go

FROM scratch

# Retrieve trusted CAs' certificates from previous stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Retrieve the binary from the previous stage
COPY --from=builder /code/synthetic /app/synthetic

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "/app/synthetic" ]
