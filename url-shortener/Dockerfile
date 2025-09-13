FROM golang:alpine AS build

# Set the working directory inside the container
WORKDIR /app
COPY . .
RUN apk add ca-certificates && go mod download
RUN go build  -o /main .

FROM build AS development
RUN apk add ca-certificates make curl git postgresql-client redis \
  && go install github.com/cweill/gotests/gotests@latest \
  && go install github.com/fatih/gomodifytags@latest \
  && go install github.com/josharian/impl@latest \
  && go install github.com/haya14busa/goplay/cmd/goplay@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install honnef.co/go/tools/cmd/staticcheck@latest \
  && go install github.com/ramya-rao-a/go-outline@latest \
  && go install golang.org/x/tools/gopls@latest \
  && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
COPY .env.sample .env
CMD ["go", "run", "main.go"]

# Use 'scratch' image for super-mini build.
# .dockerignore may not be required because we copy only needed
FROM scratch AS production
COPY --from=build /main .
COPY --from=build /app/.env.sample ./.env
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# check .env file for this port

EXPOSE 8080
CMD ["./main"]