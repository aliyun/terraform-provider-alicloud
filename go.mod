module github.com/terraform-providers/terraform-provider-alicloud

require (
	github.com/Sirupsen/logrus v0.0.0-20181010200618-458213699411 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.60.339
	github.com/aliyun/aliyun-datahub-sdk-go v0.0.0-20180929121038-c1c85baca7c0
	github.com/aliyun/aliyun-log-go-sdk v0.1.6-0.20191218073949-14a8ed32b27b
	github.com/aliyun/aliyun-mns-go-sdk v0.0.0-20191115025756-088ba95470f4
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190528142024-f8d6d645dc4b
	github.com/aliyun/aliyun-tablestore-go-sdk v4.1.3+incompatible
	github.com/aliyun/fc-go-sdk v0.0.0-20190326033901-db3e654c23d6
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/denverdino/aliyungo v0.0.0-20191128015008-acd8035bbb1d
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/gogap/errors v0.0.0-20160523102334-149c546090d0 // indirect
	github.com/gogap/stack v0.0.0-20150131034635-fef68dddd4f8 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform v0.12.9
	github.com/hashicorp/terraform-plugin-sdk v1.4.0
	github.com/hashicorp/vault v0.10.4
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/keybase/go-crypto v0.0.0-20190416182011-b785b22cc757 // indirect
	github.com/klauspost/compress v0.0.0-20180801095237-b50017755d44 // indirect
	github.com/klauspost/cpuid v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/valyala/bytebufferpool v0.0.0-20180905182247-cdfbe9377474 // indirect
	github.com/valyala/fasthttp v0.0.0-20180927122258-761788a34bb6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20190221042446-c2654d5206da // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace github.com/Sirupsen/logrus v0.0.0-20181010200618-458213699411 => github.com/sirupsen/logrus v0.0.0-20181010200618-458213699411

go 1.13
