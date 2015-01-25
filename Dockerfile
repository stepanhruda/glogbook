FROM golang

WORKDIR /go/src/github.com/stepanhruda/glogbook

EXPOSE 3000

RUN ["go", "get", "github.com/tools/godep"]
RUN ["go", "install", "github.com/tools/godep"]

COPY . /go/src/github.com/stepanhruda/glogbook

RUN CGO_ENABLED=0 GOOS=linux godep go build -a -tags netgo -ldflags '-w' .

ENTRYPOINT ["/go/src/github.com/stepanhruda/glogbook"]
