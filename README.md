Terraform Provider For Alibaba Cloud
==================

- Tutorials: [learn.hashicorp.com](https://learn.hashicorp.com/terraform?track=getting-started#getting-started)
- Documentation: https://www.terraform.io/docs/providers/alicloud/index.html
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Supported Versions
------------------

| Terraform version | minimum provider version |maximum provider version
| ---- | ---- | ----| 
| >= 0.11.x	| 1.0.0	| latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)
-   [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports):
    ```
    go get golang.org/x/tools/cmd/goimports
    ```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/aliyun/terraform-provider-alicloud`

```sh
$ mkdir -p $GOPATH/src/github.com/aliyun; cd $GOPATH/src/github.com/aliyun
$ git clone git@github.com:aliyun/terraform-provider-alicloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/aliyun/terraform-provider-alicloud
$ make build
```

Using the provider
----------------------
Please see [instructions](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs#authentication) on how to configure the Alibaba Cloud Provider.


## Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-alicloud
...
```

Running `make dev` or `make devlinux` or `devwin` will only build the specified developing provider which matchs the local system.
And then, it will unarchive the provider binary and then replace the local provider plugin.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Acceptance Testing
Before making a release, the resources and data sources are tested automatically with acceptance tests (the tests are located in the alicloud/*_test.go files).
You can run them by entering the following instructions in a terminal:
```
cd $GOPATH/src/github.com/aliyun/terraform-provider-alicloud
export ALICLOUD_ACCESS_KEY=xxx
export ALICLOUD_SECRET_KEY=xxx
export ALICLOUD_REGION=xxx
export ALICLOUD_ACCOUNT_ID=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./alicloud -v -run=TestAccAlicloud -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```

-> **Note:** The last line is optional, it allows converting test results into an XML format compatible with xUnit.


-> **Note:** Most test cases will create PayAsYouGo resources when running above test command. However, currently not all
 account site type support create PayAsYouGo resources, so you need set your account site type before running the command:
```
# If your account belongs to domestic site
export ALICLOUD_ACCOUNT_SITE=Domestic

# If your account belongs to international site
export ALICLOUD_ACCOUNT_SITE=International
```
The setting of account site type can skip some unsupported cases automatically.
