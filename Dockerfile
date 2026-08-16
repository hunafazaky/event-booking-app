# --- Stage 2: Runner ---
FROM golang:1.26.5 AS runner
WORKDIR /app

RUN go install github.com/air-verse/air@latest

# Dev container: run as root so bind-mounted host files (owned by your host
# user) can be written to by Air (tmp/main) and go build. Re-enable a
# non-root user in a separate production-stage Dockerfile instead.
# RUN useradd -m -u 1001 gouser
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
# USER gouser

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
