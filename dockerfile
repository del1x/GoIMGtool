FROM golang:1.24.6 AS builder


RUN apt-get update && apt-get install -y gcc libwebp-dev libgl1-mesa-dev libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libxxf86vm-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY cmd/ ./cmd/
COPY gui/ ./gui/
COPY fileio/ ./fileio/
COPY processor/ ./processor/
COPY config/ ./config/


RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o goimgtool ./cmd


FROM debian:bookworm-slim


RUN apt-get update && apt-get install -y libwebp7 libx11-6 libxrandr2 libxinerama1 libxcursor1 libxi6 libgl1-mesa-glx libxxf86vm1 && rm -rf /var/lib/apt/lists/*


COPY --from=builder /app/goimgtool /usr/local/bin/goimgtool


WORKDIR /app


ENTRYPOINT ["/usr/local/bin/goimgtool"]