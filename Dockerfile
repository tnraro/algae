FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -o=build/algae ./cmd/algae

FROM docker:27.1-cli
ENV GIN_MODE=release
ENV DATA_DIR=data
WORKDIR /app
COPY --from=build /app/build/algae ./

ENTRYPOINT ["./algae"]