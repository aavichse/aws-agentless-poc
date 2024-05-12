FROM golang:1.22.2-alpine3.18 as gobase

ENV GOPATH=/go

WORKDIR $GOPATH/src/agentless

COPY ./agent/go/infra/go.mod ./agent/go/infra/go.sum ./infra/
COPY ./agent/go/inventory/go.mod ./agent/go/inventory/go.sum ./inventory/

RUN cd infra && go mod tidy
RUN cd inventory && go mod tidy

COPY ./agent/go/infra ./infra
COPY ./agent/go/inventory ./inventory

RUN mkdir -p .aws 
COPY .aws/credentials .aws/

WORKDIR $GOPATH/src/agentless/inventory
RUN go build -o $GOPATH/bin/gc-inventory


FROM alpine:3.18

ENV GOPATH=/go

COPY --from=gobase $GOPATH/src/agentless/.aws /root/.aws
COPY --from=gobase $GOPATH/bin/gc-inventory /usr/local/bin/gc-inventory

ENTRYPOINT ["gc-inventory"]
EXPOSE 8080