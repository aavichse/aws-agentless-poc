
FROM python:3.9-alpine3.18 as pybase

ENV OPTDIR=/opt/guardicore

RUN apk add --no-cache gcc musl-dev linux-headers libffi-dev expat-dev

WORKDIR $OPTDIR

COPY ./requirements.txt .
COPY ./dist   ./dist

RUN pip3 install --no-cache-dir -r requirements.txt

COPY ./agent/py/reveal ./reveal



FROM golang:1.22.2-alpine3.18 as gobase

ENV OUTDIR=/opt/guardicore

WORKDIR $OUTDIR

COPY ./agent/go/infra ./infra
COPY ./agent/go/onboarding ./onboarding
COPY ./agent/go/inventory ./inventory

RUN set -eux; \
    cd infra && go mod tidy; \
    cd ../onboarding && go mod tidy; \
    cd ../inventory && go mod tidy; \
    cd ../onboarding && go build  -o ../gc-onboarding; \
    cd ../inventory && go build   -o ../gc-inventory


FROM alpine:3.18

RUN apk add --no-cache  go python3 py3-pip libffi libgcc expat

ENV OPTDIR=/opt/guardicore

RUN apk add --no-cache nginx

COPY ./daemons/nginx /etc/nginx

WORKDIR $OPTDIR

COPY --from=gobase $OPTDIR/gc-onboarding ./gc-onboarding
COPY --from=gobase $OPTDIR/gc-inventory  ./gc-inventory
COPY --from=pybase /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages

COPY --from=pybase /usr/local/lib /usr/local/lib
COPY --from=pybase /usr/local/bin /usr/local/bin

COPY --from=pybase $OPTDIR/reveal ./reveal


COPY .aws/credentials /root/.aws/credentials

COPY ./scripts/start_services.sh /usr/local/bin/start_services.sh

EXPOSE 8080
CMD ["/usr/local/bin/start_services.sh"]