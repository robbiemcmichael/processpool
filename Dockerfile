FROM golang:1.17 as build

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/ ./cmd/...

####

FROM scratch

COPY --from=build /src/bin/processpool /bin/
