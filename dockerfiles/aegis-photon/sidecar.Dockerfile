
# .-'_.---._'-.
# ||####|(__)||   Protect your secrets, protect your business.
#   \\()|##//       Secure your sensitive data with Aegis.
#    \\ |#//                    <aegis.ist>
#     .\_/.
#

# builder image
FROM golang:1.20.1-alpine3.17 as builder
RUN mkdir /build
COPY app /build/app
COPY core /build/core
COPY sdk /build/sdk
COPY vendor /build/vendor
COPY go.mod /build/go.mod
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -o aegis-sidecar ./app/sidecar/cmd/main.go

# generate clean, final image for end users
FROM photon:5.0

LABEL "maintainers"="Volkan Özçelik <volkan@aegis.ist>"
LABEL "version"="0.18.0"
LABEL "website"="https://aegis.ist/"
LABEL "repo"="https://github.com/shieldworks/aegis"
LABEL "documentation"="https://aegis.ist/docs/"
LABEL "contact"="https://aegis.ist/contact/"
LABEL "community"="https://aegis.ist/contact/#community"
LABEL "changelog"="https://aegis.ist/changelog"

COPY --from=builder /build/aegis-sidecar .

# executable
ENTRYPOINT [ "./aegis-sidecar" ]
CMD [ "" ]
