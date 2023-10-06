FROM golang:1.20 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o /apiserver cmd/apiserver/main.go 


FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /apiserver /apiserver

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/apiserver"]
