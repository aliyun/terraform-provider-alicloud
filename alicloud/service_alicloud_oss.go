package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type OssService struct {
	client *connectivity.AliyunClient
}

func (s *OssService) QueryOssBucketById(id string) (info *oss.BucketInfo, err error) {
	raw, err := s.client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketInfo(id)
	})
	if err != nil {
		return nil, err
	}
	bucket, _ := raw.(oss.GetBucketInfoResult)
	return &bucket.BucketInfo, nil
}

func (s *OssService) WaitForOssBucket(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		exist, err := bucket.IsObjectExist(id)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, "IsObjectExist", AliyunOssGoSdk)
		}
		addDebug("IsObjectExist", exist)

		if !exist {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(exist), status, ProviderERROR)
		}
	}
}
