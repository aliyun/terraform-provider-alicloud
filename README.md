Terraform Provider For Alibaba Cloud
==================
<img align="right" src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="200px"><br>
<img align="right" src="https://www.datocms-assets.com/2885/1506527326-color.svg" width="200px">

- Website: https://www.terraform.io
- Documentation: https://www.terraform.io/docs/providers/alicloud/
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.10 (to build the provider plugin)
-   [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports):
    ```
    go get golang.org/x/tools/cmd/goimports
    ```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-alicloud`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-alicloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-alicloud
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.10+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

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
cd $GOPATH/src/github.com/alibaba/terraform-provider
export ALICLOUD_ACCESS_KEY=xxx
export ALICLOUD_SECRET_KEY=xxx
export ALICLOUD_REGION=xxx
export ALICLOUD_ACCOUNT_ID=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./alicloud -v -run=TestAccAlicloud -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```

-> **Note:** The last line is optional, it allows to convert test results into a XML format compatible with xUnit.

Because some features are not available in all regions, the following environment variables can be set in order to
skip tests that use these features:
* ALICLOUD_SKIP_TESTS_FOR_SLB_SPECIFICATION=true    - Server Load Balancer with guaranteed performance specifications (old implementation has only shared performance)
* ALICLOUD_SKIP_TESTS_FOR_SLB_PAY_BY_BANDWIDTH=true - Server Load Balancer with a "pay by bandwidth" billing method (mostly available in China)
* ALICLOUD_SKIP_TESTS_FOR_FUNCTION_COMPUTE=true     - Function Compute
* ALICLOUD_SKIP_TESTS_FOR_PVTZ_ZONE=true            - Private Zone
* ALICLOUD_SKIP_TESTS_FOR_RDS_MULTIAZ=true          - Apsara RDS with multiple availability zones
* ALICLOUD_SKIP_TESTS_FOR_CLASSIC_NETWORK=true      - Classic network configuration

## Refer

Alibaba Cloud Provider Development Repository [terraform-provider](https://github.com/alibaba/terraform-provider)

Alibaba Cloud Provider [Official Docs](https://www.terraform.io/docs/providers/alicloud/index.html)
