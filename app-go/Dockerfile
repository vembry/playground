# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22.6-alpine AS build-stage

WORKDIR /home/app/app-go
COPY ./app-go ./
# RUN go mod download

WORKDIR /home/app/sdk-go
COPY ./sdk-go ./
# RUN go mod download

WORKDIR /home/app
RUN go work init ./app-go ./sdk-go && go mod download 

WORKDIR /home/app/app-go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app-go


# # Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /home/app/app-go/bin/app-go /app-go
COPY --from=gcr.io/distroless/base:debug /busybox/wget /

EXPOSE 80
