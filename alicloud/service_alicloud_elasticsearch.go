package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/denverdino/aliyungo/common"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type ElasticsearchService struct {
	client *connectivity.AliyunClient
}

func (s *ElasticsearchService) DescribeInstance(instanceId string) (v elasticsearch.DescribeInstanceResponse, err error) {
	request := elasticsearch.CreateDescribeInstanceRequest()
	request.InstanceId = instanceId
	request.SetContentType("application/json")

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DescribeInstance(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{InstanceNotFound}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("elasticsearch", instanceId))
			}

			return err
		}

		resp, _ := raw.(*elasticsearch.DescribeInstanceResponse)
		if resp == nil || resp.Result.InstanceId != instanceId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("elasticsearch", instanceId))
		}

		v = *resp
		return nil
	})

	return
}

func (s *ElasticsearchService) WaitForElasticsearchInstance(instanceId string, status ElasticsearchStatus, timeout int) error {
	for {
		resp, err := s.DescribeInstance(instanceId)
		if err != nil {
			//
		} else if resp.Result.Status == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalLong
		time.Sleep(DefaultIntervalLong * time.Second)
	}
	return nil
}
