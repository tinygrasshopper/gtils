# (G)o-U(tils)
Arguably useful & reusable packages for our golang projects.

[![wercker status](https://app.wercker.com/status/f80b31a3ddb734d6327e3fd9e250dec3/m "wercker status")](https://app.wercker.com/project/bykey/f80b31a3ddb734d6327e3fd9e250dec3)

[![GoDoc](http://godoc.org/github.com/pivotalservices/gtils?status.png)](http://godoc.org/github.com/pivotalservices/gtils)


## Running tests / build pipeline locally

```

# install the wercker cli
$ curl -L https://install.wercker.com | sh

# make sure a docker host is running
$ boot2docker up && $(boot2docker shellinit)

# run the build pipeline locally, to test your code locally
$ ./testrunner

```

