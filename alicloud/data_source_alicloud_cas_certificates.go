package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudCasCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCascertsRead,

		Schema: map[string]*schema.Schema{
			"lang": {
				Type:     schema.TypeString,
				Computed: true,
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
							Type:     schema.TypeString,
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
							Type:     schema.TypeString,
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
	client := meta.(*connectivity.AliyunClient)

	args := cas.CreateDescribeUserCertificateListRequest()
	if v, ok := d.GetOk("show_size"); ok {
		args.ShowSize = requests.NewInteger(v.(int))
	} else {
		args.ShowSize = requests.NewInteger(50)
	}

	if v, ok := d.GetOk("current_page"); ok {
		args.CurrentPage = requests.NewInteger(v.(int))
	} else {
		args.CurrentPage = requests.NewInteger(1)
	}

	var allcerts []cas.Certificate
	raw, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
		return casClient.DescribeUserCertificateList(args)
	})
	if err != nil {
		return WrapError(err)
	}
	res, _ := raw.(*cas.DescribeUserCertificateListResponse)
	allcerts = append(allcerts, res.CertificateList...)

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), allcerts)
	}

	return nil
}
