# UltiPkg

Ultipkg helps make code `go get`-able. Ultipkg does not hold any package code; instead it will redirect the go tool to the respective Git repository.

It supports SSL so that we can support the latest version of Go which requires a `-insecure` flag before you can use non-ssl.

## Hosting Ultipkg

To host Ultipkg, you'll just need a copy of the binary for your operating system and architecture. With that, you can start it just by executing the binary.

There are some option you can pass as environment variables. The below is an example for a linux server:

    $ DOMAIN=pkg.example.com ADDR=0.0.0.0:80 ADDR_TLS=0.0.0.0:443 \
    SSL_CERTIFICATE=/path/to/cert SSL_PRIVATEKEY=/path/to/key \
    LOG_LEVEL=warn ultipkg

    {"addr":"0.0.0.0:443","level":"info","msg":"starting https server","time":"2016-02-14T02:00:31-05:00","version":"XXXXXX"}
    {"addr":"0.0.0.0.0:8080","level":"info","msg":"starting http server","time":"2016-02-14T02:00:31-05:00","version":"XXXXXX"}


## Using Ultipkg

Issue the following command:

    go get pkg.ulti.io/{project}/{repository}/{optional_path}

This will clone the repo into the following path:

    $GOPATH/src/pkg.ulti.io/{project}/{repository}

### Caveats

Make sure the following is true (if not, make it so):

* You have access to the git repo you are trying to import
* You have configured your SSH keys with your repo host
