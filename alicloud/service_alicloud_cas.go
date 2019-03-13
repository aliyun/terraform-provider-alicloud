package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CasService struct {
	client *connectivity.AliyunClient
}

func (s *CasService) DescribeCas(id string) (*cas.Certificate, error) {
	var cert *cas.Certificate

	request := cas.CreateDescribeUserCertificateListRequest()
	request.ShowSize = requests.NewInteger(PageSizeLarge)
	request.CurrentPage = requests.NewInteger(1)

	raw, err := s.client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
		return casClient.DescribeUserCertificateList(request)
	})
	if err != nil {
		return nil, WrapError(err)
	}
	res, _ := raw.(*cas.DescribeUserCertificateListResponse)

	if len(res.CertificateList) == 0 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Cas", id)), NotFoundMsg, ProviderERROR)
	}

	for _, v := range res.CertificateList {
		if id == strconv.Itoa(v.Id) {
			cert = &v
			break
		}
	}

	if cert == nil {
		return cert, WrapErrorf(Error(GetNotFoundMessage("Cas", id)), NotFoundMsg, ProviderERROR)
	}

	return cert, nil
}
