FROM golang:1.23.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
COPY ./static /static
RUN go build -o bin/app ./cmd

FROM alpine:latest AS final
RUN apk update && apk add --no-cache ca-certificates
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /app/bin/app .
COPY --from=build /static .
EXPOSE 8090
ENTRYPOINT [ "./app" ]