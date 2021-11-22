FROM golang:1.17.2-alpine AS go_build

WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o oapi3gen .

FROM alpine

COPY --from=go_build /build/oapi3gen /oapi3gen
COPY base.tmpl /base.tmpl
COPY ./echo/server.tmpl /echo/server.tmpl

ENTRYPOINT ["/oapi3gen"]