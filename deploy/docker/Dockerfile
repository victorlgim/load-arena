FROM golang:1.24 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/load-arena ./cmd/api

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/load-arena /load-arena
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/load-arena"]
