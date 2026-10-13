package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CasService struct {
	client *connectivity.AliyunClient
}

func (s *CasService) DescribeSslCertificatesServiceCertificate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeUserCertificateDetail"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"CertId":   id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2018-07-13", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}

	object = v.(map[string]interface{})

	if _, idExist := response["Id"]; !idExist {
		return object, WrapErrorf(NotFoundErr("Cas:Certificate", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CasService) SslCertificatesServiceCertificateStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSslCertificatesServiceCertificate(d.Id())
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		return object, fmt.Sprint(object["Id"]), nil
	}
}
