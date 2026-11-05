package oss

import (
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// OssService *connectivity.AliyunClient
type OssService struct {
	client *connectivity.AliyunClient
}

func NewOssService(client *connectivity.AliyunClient) *OssService { return &OssService{client: client} }

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {
	request := map[string]string{"bucketName": id}
	var requestInfo *oss.Client
	raw, err := s.client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketInfo(request["bucketName"])
	})
	if err != nil {
		if ossNotFoundError(err) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunOssGoSdk)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetBucketInfo", errs.AliyunOssGoSdk)
	}

	logs.AddDebug("GetBucketInfo", raw, requestInfo, request)
	response, _ = raw.(oss.GetBucketInfoResult)
	return
}

func (s *OssService) WaitForOssBucket(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeOssBucket(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
				// for delete bucket replication
			} else if status == Deleted && errs.IsExpectedErrors(err, []string{"AccessDenied"}) {
				return nil
			} else {
				return errs.WrapError(err)
			}
		}

		if object.BucketInfo.Name != "" && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.BucketInfo.Name, status, errs.ProviderERROR)
		}
	}
}

func (s *OssService) WaitForOssBucketObject(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		exist, err := bucket.IsObjectExist(id)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "IsObjectExist", errs.AliyunOssGoSdk)
		}
		logs.AddDebug("IsObjectExist", exist)

		if !exist {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, strconv.FormatBool(exist), status, errs.ProviderERROR)
		}
	}
}

func (s *OssService) DescribeOssBucketReplication(id string) (response string, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return response, errs.WrapError(err)
	}
	bucket := parts[0]
	ruleId := parts[1]

	request := map[string]string{"bucketName": bucket, "ruleId": ruleId}
	var requestInfo *oss.Client
	raw, err := s.client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketReplication(bucket)
	})
	if err != nil {
		if ossNotFoundError(err) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunOssGoSdk)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetBucketReplication", errs.AliyunOssGoSdk)
	}

	logs.AddDebug("GetBucketReplication", raw, requestInfo, request)
	response, _ = raw.(string)
	return
}
