# Start from golang base latest image
FROM golang:latest

# Set maintainer
LABEL maintainer="Setia Budi"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy all resource into container
COPY . .

# Download and install dependencies
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o main .

ENTRYPOINT ./main