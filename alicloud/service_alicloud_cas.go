package alicloud

import (
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CasService struct {
	client *connectivity.AliyunClient
}

func (s *CasService) DescribeCas(id string) (*cas.Certificate, error) {
	certificate := &cas.Certificate{}
	request := cas.CreateDescribeUserCertificateListRequest()
	request.RegionId = s.client.RegionId
	request.ShowSize = requests.NewInteger(PageSizeLarge)
	request.CurrentPage = requests.NewInteger(1)

	for i := 1; ; i++ {
		request.CurrentPage = requests.NewInteger(i)

		raw, err := s.client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
			return casClient.DescribeUserCertificateList(request)
		})
		if err != nil {
			return certificate, WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		res, _ := raw.(*cas.DescribeUserCertificateListResponse)
		for _, v := range res.CertificateList {
			if id == strconv.FormatInt(v.Id, 10) {
				return &v, nil
			}
		}

		if len(res.CertificateList) < PageSizeLarge {
			break
		}
	}

	return certificate, WrapErrorf(Error(GetNotFoundMessage("Cas", id)), NotFoundMsg, ProviderERROR)
}

func (s *CasService) DescribeSslCertificatesServiceCertificate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCasClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeUserCertificateDetail"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"CertId":   id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-07-13"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Cas.Sertificate", id)), NotFoundWithResponse, response)
	}
	return object, nil
}
