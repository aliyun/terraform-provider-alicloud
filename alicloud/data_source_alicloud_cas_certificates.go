package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudCasCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCascertsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"province": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sans": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"buy_in_aliyun": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCascertsRead(d *schema.ResourceData, meta interface{}) error {
	var allcerts []cas.Certificate
	client := meta.(*connectivity.AliyunClient)

	args := cas.CreateDescribeUserCertificateListRequest()
	args.ShowSize = requests.NewInteger(PageSizeLarge)

	for i := 1; ; i++ {
		args.CurrentPage = requests.NewInteger(i)

		raw, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
			return casClient.DescribeUserCertificateList(args)
		})
		if err != nil {
			return WrapError(err)
		}

		res, _ := raw.(*cas.DescribeUserCertificateListResponse)
		allcerts = append(allcerts, res.CertificateList...)

		if len(res.CertificateList) < PageSizeLarge {
			break
		}
	}

	var s []map[string]interface{}
	var ids []string
	for _, cert := range allcerts {
		if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(cert.Name) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":            cert.Id,
			"name":          cert.Name,
			"common":        cert.Common,
			"finger_print":  cert.Fingerprint,
			"issuer":        cert.Issuer,
			"org_name":      cert.OrgName,
			"province":      cert.Province,
			"city":          cert.City,
			"country":       cert.Country,
			"start_date":    cert.StartDate,
			"end_date":      cert.EndDate,
			"sans":          cert.Sans,
			"expired":       cert.Expired,
			"buy_in_aliyun": cert.BuyInAliyun,
		}
		s = append(s, mapping)
		ids = append(ids, strconv.Itoa(cert.Id))
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("certificates", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
