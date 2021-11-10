# Start from the latest golang base image
FROM golang:latest as builder

## Add Maintainer Info
LABEL maintainer="Yogesh <katreddyveey@vmware.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

RUN mkdir data
RUN chmod 755 /data

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

EXPOSE 8082

# Command to run the executable
CMD ["./main"]
