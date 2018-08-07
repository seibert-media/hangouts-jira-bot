FROM golang:1.10 as builder

ARG GIT_HOST
ARG REPO
ARG NAME

ADD ./ /go/src/${GIT_HOST}/${REPO}/${NAME}
WORKDIR /go/src/${GIT_HOST}/${REPO}/${NAME}/

RUN make buildgo

CMD ["/bin/bash"]

FROM scratch

LABEL maintainer //SEIBERT/MEDIA GmbH <docker@seibert-media.net>
LABEL type "public"
LABEL versioning "simple"

ARG VERSION
ARG GIT_HOST
ARG REPO
ARG NAME

COPY --from=builder /go/src/${GIT_HOST}/${REPO}/${NAME}/app /
ADD ./files/go-cloud-debug /
ADD ./files/source-context.json /
COPY files/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["./app"]
