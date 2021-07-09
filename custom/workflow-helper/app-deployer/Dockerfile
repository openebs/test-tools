#Build Stage
FROM golang:1.14 AS builder

LABEL maintainer="LitmusChaos"

ARG TARGETOS=linux
ARG TARGETARCH

ADD . /app-deployer
WORKDIR /app-deployer

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

RUN go env

RUN CGO_ENABLED=0 go build -o /output/deployer -v

#Deploy Stage
FROM alpine:latest
ARG TARGETARCH

RUN apk add curl

#Installing Kubectl
ENV KUBECTL_VERSION="v1.19.0"
#Installing kubectl
RUN curl -sLO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/${TARGETARCH}/kubectl" && chmod +x kubectl && mv kubectl /usr/bin/kubectl

# Copy application manifests
COPY ./app-manifest /var/run

COPY --from=builder /output/deployer /var/run

ENTRYPOINT ["/var/run/deployer"]
