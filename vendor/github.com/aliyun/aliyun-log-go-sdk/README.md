This is a Golang SDK for Alibaba Cloud [Log Service](https://sls.console.aliyun.com/).

API Reference :

* [Chinese](https://help.aliyun.com/document_detail/29007.html)
* [English](https://intl.aliyun.com/help/doc-detail/29007.htm)

[![Build Status](https://travis-ci.org/galaxydi/go-loghub.svg?branch=master)](https://travis-ci.org/galaxydi/go-loghub)
[![Coverage Status](https://coveralls.io/repos/github/galaxydi/go-loghub/badge.svg?branch=master&foo=bar)](https://coveralls.io/github/galaxydi/go-loghub?branch=master&foo=bar)


# Install Instruction

### LogHub Golang SDK

```
go get github.com/aliyun/aliyun-log-go-sdk
```

# Example 

### Write and Read LogHub

[loghub_sample.go](example/loghub/loghub_sample.go)

### Use Index on LogHub (SLS)

[index_sample.go](example/index/index_sample.go)

### Create Config for Logtail

[log_config_sample.go](example/config/log_config_sample.go)

### Create Machine Group for Logtail

[machine_group_sample.go](example/machine_group/machine_group_sample.go)

# For developer
### Update log protobuf
`protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gofast_out=. log.proto`