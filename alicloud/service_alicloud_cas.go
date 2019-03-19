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
	request := cas.CreateDescribeUserCertificateListRequest()
	request.ShowSize = requests.NewInteger(PageSizeLarge)
	request.CurrentPage = requests.NewInteger(1)

	for i := 1; ; i++ {
		request.CurrentPage = requests.NewInteger(i)

		raw, err := s.client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
			return casClient.DescribeUserCertificateList(request)
		})
		if err != nil {
			return nil, WrapError(err)
		}

		res, _ := raw.(*cas.DescribeUserCertificateListResponse)
		for _, v := range res.CertificateList {
			if id == strconv.Itoa(v.Id) {
				return &v, nil
			}
		}

		if len(res.CertificateList) < PageSizeLarge {
			break
		}
	}

	return nil, WrapErrorf(Error(GetNotFoundMessage("Cas", id)), NotFoundMsg, ProviderERROR)
}
