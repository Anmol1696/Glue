ARG GO_VERSION=1.13

FROM golang:${GO_VERSION}-alpine as builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk --no-cache --update upgrade && apk add --no-cache ca-certificates git curl

# Update certs
RUN update-ca-certificates

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Copy the local package files to the container's workspace.
COPY ./ ./

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

# Final stage: the running container.
FROM scratch AS final

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /app

USER nobody:nobody

ENTRYPOINT ["/app"]
