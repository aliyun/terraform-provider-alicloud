module github.com/aliyun/terraform-provider-alicloud

require (
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/Sirupsen/logrus v0.0.0-20181010200618-458213699411 // indirect
	github.com/alibabacloud-go/tea-rpc v1.1.6
	github.com/alibabacloud-go/tea-utils v1.3.4
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.623
	github.com/aliyun/aliyun-datahub-sdk-go v0.0.0-20180929121038-c1c85baca7c0
	github.com/aliyun/aliyun-log-go-sdk v0.1.13
	github.com/aliyun/aliyun-mns-go-sdk v0.0.0-20191115025756-088ba95470f4
	github.com/aliyun/aliyun-oss-go-sdk v2.1.2+incompatible
	github.com/aliyun/aliyun-tablestore-go-sdk v4.1.3+incompatible
	github.com/aliyun/credentials-go v1.1.2
	github.com/aliyun/fc-go-sdk v0.0.0-20200925033337-c013428cbe21
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/denverdino/aliyungo v0.0.0-20201111050017-ecb1ab1b9105
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/gogap/errors v0.0.0-20160523102334-149c546090d0 // indirect
	github.com/gogap/stack v0.0.0-20150131034635-fef68dddd4f8 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.4.0
	github.com/hashicorp/vault v0.10.4
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/keybase/go-crypto v0.0.0-20190416182011-b785b22cc757 // indirect
	github.com/klauspost/compress v0.0.0-20180801095237-b50017755d44 // indirect
	github.com/klauspost/cpuid v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/valyala/bytebufferpool v0.0.0-20180905182247-cdfbe9377474 // indirect
	github.com/valyala/fasthttp v0.0.0-20180927122258-761788a34bb6 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20190221042446-c2654d5206da // indirect
)

replace github.com/Sirupsen/logrus v0.0.0-20181010200618-458213699411 => github.com/sirupsen/logrus v0.0.0-20181010200618-458213699411

go 1.13
