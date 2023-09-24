FROM ubuntu:latest as build

#ARG GO_VERSION
ENV GO_VERSION="1.21.1"

RUN apt-get update
RUN apt-get install -y wget git gcc

RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

#WORKDIR $GOPATH
WORKDIR /app
COPY . .

RUN go mod download
RUN go vet -v

ENV GITHUB_TOKEN=""
ENV KEYCHAIN_APP_SERVICE="github-forkrefresh"
ENV KEYCHAIN_USERNAME="dmore"
ENV REPOS_GIST="https://gist.githubusercontent.com/dmore/5c26c5c2484aa13736f22d80e8bf4e7e/raw/88ebe3b0d641fc7e8715bfe4056625ac2532953b/repos_repo.json"

RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/repos_repo.json /
COPY --from=build /go/bin/app /
CMD ["/app"]
