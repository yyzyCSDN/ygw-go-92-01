FROM golang:1.23.12 AS build
WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o /regenbrake ./cmd/regenbrake

FROM golang:1.23.12
ENV GOPROXY=off GOSUMDB=off
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN mkdir -p /app/data
COPY --from=build /regenbrake /usr/local/bin/regenbrake
EXPOSE 8090
CMD ["/usr/local/bin/regenbrake", "-addr", "0.0.0.0:8090", "-dir", "/app/data"]
