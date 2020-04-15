FROM golang:1.13 as builder

WORKDIR /src/hangouts-jira-bot

# Copy rest of the application source code
COPY . ./

# Compile the application to /app.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o /app ./cmd/hangouts-jira-bot

FROM scratch

LABEL maintainer Team Google <team-google@seibert-media.net>

COPY --from=builder /app /
COPY files/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/app"]
