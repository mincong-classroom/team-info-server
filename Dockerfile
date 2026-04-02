FROM golang:1.24-alpine AS builder
WORKDIR /workspace
COPY go.mod ./
COPY . ./
RUN go build -o bin/server ./main.go

FROM scratch
COPY --from=builder /workspace/bin/server /server
EXPOSE 8090
ENTRYPOINT ["/server"]
