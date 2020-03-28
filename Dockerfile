# Build environment
# -----------------
FROM golang:1.14-alpine as build-env
WORKDIR /ex8j

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app


# Deployment environment
# ----------------------
FROM alpine

COPY --from=build-env /ex8j/bin/app /ex8j/

CMD ["/ex8j/app"]