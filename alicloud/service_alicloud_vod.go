package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VodService struct {
	client *connectivity.AliyunClient
}

func (s *VodService) DescribeVodDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVodDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *VodService) VodStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeVodDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["DomainStatus"].(string) == failState {
				return object, object["DomainStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["DomainStatus"].(string)))
			}
		}
		return object, object["DomainStatus"].(string), nil
	}
}

func (s *VodService) DescribeVodDomainCertificateInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVodDomainCertificateInfo"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
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
	v, e := jsonpath.Get("$.CertInfos.CertInfo", response)
	if e != nil || v == nil {
		return object, WrapErrorf(e, FailedGetAttributeMsg, id, "$.CertInfos.CertInfo", response)
	}
	certs, ok := v.([]interface{})
	if !ok || len(certs) == 0 {
		return nil, nil
	}
	object = certs[0].(map[string]interface{})
	return object, nil
}

func (s *VodService) SetVodDomainSSLCertificate(d *schema.ResourceData) error {
	client := s.client
	var response map[string]interface{}
	var err error
	action := "SetVodDomainSSLCertificate"
	request := map[string]interface{}{
		"DomainName":  d.Id(),
		"SSLProtocol": d.Get("ssl_protocol"),
	}
	if v, ok := d.GetOk("cert_name"); ok {
		request["CertName"] = v
	}
	if v, ok := d.GetOk("cert_type"); ok {
		request["CertType"] = v
	}
	if v, ok := d.GetOk("env"); ok {
		request["Env"] = v
	}
	certType := d.Get("cert_type").(string)
	if certType == "upload" {
		sslPub, pubOk := d.GetOk("ssl_pub")
		sslPri, priOk := d.GetOk("ssl_pri")
		if d.Get("ssl_protocol").(string) != "off" {
			if !pubOk {
				return WrapError(fmt.Errorf("ssl_pub must be set when cert_type is \"upload\" and ssl_protocol is not \"off\""))
			}
			if !priOk {
				return WrapError(fmt.Errorf("ssl_pri must be set when cert_type is \"upload\" and ssl_protocol is not \"off\""))
			}
		}
		if pubOk {
			request["SSLPub"] = sslPub
		}
		if priOk {
			request["SSLPri"] = sslPri
		}
	} else if certType == "cas" {
		if v, ok := d.GetOk("cert_id"); ok {
			request["CertId"] = v
		}
		if v, ok := d.GetOk("cert_region"); ok {
			request["CertRegion"] = v
		}
	} else {
		if v, ok := d.GetOk("cert_id"); ok {
			request["CertId"] = v
		}
		if v, ok := d.GetOk("cert_region"); ok {
			request["CertRegion"] = v
		}
		if v, ok := d.GetOk("ssl_pub"); ok {
			request["SSLPub"] = v
		}
		if v, ok := d.GetOk("ssl_pri"); ok {
			request["SSLPri"] = v
		}
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDomain.notOnline", "CreateCertificateFailed"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *VodService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagVodResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("vod", "2017-03-21", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagVodResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("vod", "2017-03-21", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *VodService) DescribeVodEditingProject(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetEditingProject"
	request := map[string]interface{}{
		"ProjectId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("vod", "2017-03-21", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidEditingProject.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("VOD:EditingProject", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Project", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Project", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
