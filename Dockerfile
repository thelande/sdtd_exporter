FROM --platform=${BUILDPLATFORM} golang:1.26-alpine AS builder
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

RUN apk add --no-cache make curl git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN make GOOS=$TARGETOS GOARCH=$TARGETARCH build

FROM alpine:3
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

WORKDIR /app

COPY --from=builder /src/output/sdtd_exporter .

EXPOSE      9816
USER        nobody
ENTRYPOINT  ["/app/sdtd_exporter"]
