FROM golang:1.21 as build

WORKDIR /app
COPY . .

RUN go mod download
RUN go vet -v

#must be empty to keep using secret from keychain
#exec: "dbus-launch": executable file not found in $PATH >> need to communicate with the host
#works fine with distroless if we pass the token as env var.
ENV GITHUB_TOKEN=""
ENV KEYCHAIN_APP_SERVICE="github-forkrefresh"
ENV KEYCHAIN_USERNAME="dmore"

RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/repos_repo.json /
COPY --from=build /go/bin/app /
CMD ["/app"]
