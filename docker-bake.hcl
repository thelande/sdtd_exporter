DATE = formatdate("YYYY.MM.DD", timestamp())
APP = "sdtd_exporter"
SOURCE = "https://github.com/thelande/sdtd_exporter"
variable "GIT_SHA" {}

variable "VERSION" {
  default = "0.2.4"
}

group "default" {
    targets = ["image-local"]
}

target "image" {
    inherits = ["docker-metadata-action"]
    args = {}
    labels = {
        "org.opencontainers.image.vendor" = "thelande"
        "org.opencontainers.image.source" = "https://github.com/thelande/sdtd_exporter"
        "org.opencontainers.image.created" = "${DATE}"
        "org.opencontainers.image.revision" = "${GIT_SHA}"
        "org.opencontainers.image.title" = "${APP}"
        "org.opencontainers.image.url" = "${SOURCE}"
        "org.opencontainers.image.version" = "${VERSION}"
    }
    no-cache = true
}

target "image-local" {
    inherits = ["image"]
    output = ["type=docker"]
    tags = ["${APP}:${DATE}"]
}

target "image-all" {
    inherits = ["image"]
    platforms = [
        "linux/amd64",
        "linux/arm64"
    ]
    tags = [
        "docker.io/thelande/${APP}:rolling",
        "docker.io/thelande/${APP}:sha-${GIT_SHA}",
        "docker.io/thelande/${APP}:v${VERSION}"
    ]
}

target "docker-metadata-action" {}
