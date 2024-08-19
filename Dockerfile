FROM --platform=${BUILDPLATFORM} golang:1.21-alpine AS builder
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

RUN apk add make curl git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN make GOOS=$TARGETOS GOARCH=$TARGETARCH build

FROM alpine:3.18.4
LABEL maintainer="Tom Helander <thomas.helander@gmail.com>"

WORKDIR /app

COPY --from=builder /src/output/sdtd_exporter .

EXPOSE 9810

ENTRYPOINT ["/app/sdtd_exporter"]
