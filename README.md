# pbt-stress-testing


### How to run

    make
    ./build/Release/01-get-balance \
        -base-url http://<API_BaseUrl> \
        -username <SomeRegisteredUser> -password <S3cretPassword> \
        -concurrency 100 -attempts 10 -timeout 5000 