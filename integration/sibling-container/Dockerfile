FROM golang:1.26
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "test", "-count=1", "-run", "^TestSiblingContainerUpstreamProxy$", "-v", "./pkg/config"]
