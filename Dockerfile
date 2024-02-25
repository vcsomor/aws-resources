# Build Stage
FROM vcsomor/alpine-golang-buildimage:1.21.7 AS build-stage

LABEL app="build-aws-resources"
LABEL REPO="https://github.com/vcsomor/aws-resources"

ENV PROJPATH=/go/src/github.com/vcsomor/aws-resources

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/vcsomor/aws-resources
WORKDIR /go/src/github.com/vcsomor/aws-resources

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/vcsomor/aws-resources"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/aws-resources/bin

WORKDIR /opt/aws-resources/bin

COPY --from=build-stage /go/src/github.com/vcsomor/aws-resources/bin/aws-resources /opt/aws-resources/bin/
RUN chmod +x /opt/aws-resources/bin/aws-resources

# Create appuser
RUN adduser -D -g '' aws-resources
USER aws-resources

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/aws-resources/bin/aws-resources"]
