FROM golang:1.12
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build --ldflags="-w -s" -mod=readonly

FROM scratch
COPY --from=0 /build/k8s-csrapprove /
ENTRYPOINT ["/k8s-csrapprove"]
