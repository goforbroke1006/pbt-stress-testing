# pbt-stress-testing


### How to run

Goto [releases page](https://github.com/goforbroke1006/pbt-stress-testing/releases)
and download last zip with binaries.
Run binary with args.

### Apps arg samples

    ./01-get-balance \
        -base-url http://<API_BaseUrl> \
        -username <SomeRegisteredUser> -password <S3cretPassword> \
        -concurrency 100 -attempts 10 -timeout 5000

### How to build from sources

    cd ${GOPATH}/src/github.com/goforbroke1006/pbt-stress-testing
    make
    ./build/Release/<SomeApp> <some args...> 