# Build the golang runtime environment using an alias：builder
FROM golang:1.20 as builder

# environment variables
ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

# Set the working directory - where all our files live in the working directory 
# now：COPY go.mod go.sum ./ && COPY . .
# In the docker environment of golang, there will be：
## /app/build/Dockerfile
## /app/cmd/demo/main.go
## /app/go.mod
## /app/README.MD
WORKDIR /app
COPY . .
# download dependencies
RUN go mod download

# compile app
RUN go build -v -a -installsuffix cgo -o sever cmd/sever/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

# set workspace
WORKDIR /bin/

# Copy the binaries compiled by the previous container to the working directory
# That is: copy golang environment/working directory/demo alpine environment/working directory
COPY --from=builder /app/sever .

# exec command：/bin/demo
ENTRYPOINT ["/bin/sever"]