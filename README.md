
github-forkrefresh-docker
=========================

A variant of the github-forkrefresh go mini project I did before. 

This time dockerised so it can be triggered from within github workflow events...So, the target for this is to not need to use the github secret locally and create a workflow event that handles that within github... So next I will try to set github events to handle the build and run. That way there is no need to a GITHUB_TOKEN locally.. nor the need to use OS/Keychain will be needed either.

Overall a github fork refresher run on the original project branches (from your public forks) so they are updated. 

To run this locally you will need :
====================================
- GITHUB_TOKEN
- Needs a list of the public repos you want to keep updated from your original projects. (I put a few of mine in there as a sample)

- This variant if used locally works with the Github token, (not the OS/Keychain currently). As is containerised, container needs 
access to the host and don't want to expose either. 

Neeeext
=======
So next I will try to set github events to handle the build and run. That way there is no need to a GITHUB_TOKEN locally.. nor the need to use OS/Keychain will be needed either.

What does it do:
===============

    It calls github api fork refresh so your public forks are up-to-date with its source.
    there is a repos_repo.json json array file. make sure your forking public repos are there.
    That is your forks, not the originals.

    repos_repo.json
    [
       "yourgithubuser/yourpublicfork",
       "yourgithubuser/yourpublicfork2"
    ]

    if works ok if parent repos use master and main branchs. forking from develop should also work. 

    tells github to refresh the fork from the original so your public forks are refreshed from the source.


Distroless and Ubuntu latest:
==============================

    - uses distroless from google. I like it =) Thin as good...(and safer)
    - also setup a Dockerfile with go onto ubuntu-latest as github expects it.


Docker build:
============

No vulnerabilities on the current versions.

```bash
docker build -t app .
....
=> [build 7/7] RUN CGO_ENABLED=0 go build -o /go/bin/app                                                                        11.5s 
 => [stage-1 2/3] COPY --from=build /app/repos_repo.json /                                                                        0.1s 
 => [stage-1 3/3] COPY --from=build /go/bin/app /                                                                                 0.1s 
 => exporting to image                                                                                                            0.1s 
 => => exporting layers                                                                                                           0.1s 
 => => writing image sha256:cd975626feeafc46c32475b63430b8aeea27199e2b532f76ea955164b26c2331                                      0.0s 
 => => naming to docker.io/library/app                                                                                            0.0s

What's Next?
  View a summary of image vulnerabilities and recommendations → docker scout quickview
$ docker scout quickview
    ✓ SBOM of image already cached, 8 packages indexed

  Your image  app:latest                 │    0C     0H     0M     0L   
  Base image  distroless/static:nonroot  │    0C     0H     0M     0L   

```
Docker run:
==========

Uses env list. Will change that on the event side of things to get this built and run.

docker run -d -t -i --env-file .env_list --name githubforkrefresh docker.io/library/app

```bash
    docker build --build-arg some_variable_name=a_value
    rm -rf packages*
    cd github-forkrefresh/app
    go build packages/app
```
Local build
============

```bash
#export GITHUB_TOKEN='leave out for now'
export KEYCHAIN_APP_SERVICE="github-forkrefresh"
export KEYCHAIN_USERNAME="dmore"
go build -o packages/app
cd packages 
cp ../repos_repo.json .
./app
cat app.log
```

Docker compose: 
===============
```bash
docker compose up -d
docker compose down
```

How to build it:
================
```go
    rm -rf packages*
    cd github-forkrefresh/app
    go build packages/app
```
How to run it:
==============
    cd github-forkrefresh/app
    go run main.go

    this is the core of it, if you just wanna know

needs an env list like so

```go
myenvfile

KEYCHAIN_APP_SERVICE=github-forkrefresh
GITHUB_TOKEN=yertoken
#GITHUB_TOKEN=
KEYCHAIN_USERNAME=dmore
```

```go
func fork_refresh_call(branch string, reponame string, method string) (string, error) {
    absPath, _ := filepath.Abs("../"+ branch + ".json")
    f, err := os.Open(absPath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    httpposturl := "https://api.github.com/repos/" + reponame + "/merge-upstream"
    fmt.Println("url: %v", httpposturl)
    request, err := http.NewRequest("POST", httpposturl, f)
    if err != nil {
        log.Fatal(err)
    }
    request.Header.Set("Content-Type", "application/json; charset=UTF-8")
    request.Header.Set("Accept", "application/vnd.github.v3+json")
    request.Header.Set("Authorization", "token " + token_variable)

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
    //fmt.Println("response :", response.Errorf)
    fmt.Println("response Status:", response.Status)
    b, err := io.ReadAll(response.Body)
    // b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
    if err != nil {
        log.Fatalln(err)
        return "nil", err
    }
    return string(b), nil
    //return fmt.Println(string(b))
}


```

Dependencies:
=============
    - it holds dependencies to : 
    -   Depends on zalando/go-keyring to retrieve and pull secrets. Currently using version 0.2.3.
