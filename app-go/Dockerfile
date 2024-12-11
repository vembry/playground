# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23.2-alpine AS build-stage

# copy app-go
WORKDIR /home/app/app-go
COPY ./app-go ./

# copy sdk-go
WORKDIR /home/app/sdk-go
COPY ./sdk-go ./

# initialize go workspace
WORKDIR /home/app
RUN go work init ./app-go ./sdk-go && go mod download 

# build app-go's binary 
WORKDIR /home/app/app-go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app-go

# rebuild the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /home/app/app-go/bin/app-go /app-go
COPY --from=gcr.io/distroless/base:debug /busybox/wget /

EXPOSE 80
