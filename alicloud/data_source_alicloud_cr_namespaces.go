package alicloud

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"regexp"
)

func dataSourceAlicloudCRNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCRNamespacesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_create": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_visibility": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudCRNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	invoker := NewInvoker()

	var resp *cr.GetNamespaceListResponse

	if err := invoker.Run(func() error {
		raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			req := cr.CreateGetNamespaceListRequest()
			return crClient.GetNamespaceList(req)
		})
		resp, _ = raw.(*cr.GetNamespaceListResponse)
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_namespaces", "GetNamespaceList", AlibabaCloudSdkGoERROR)
	}

	var crResp crDescribeNamespaceListResponse
	err := json.Unmarshal(resp.GetHttpContentBytes(), &crResp)
	if err != nil {
		return WrapError(err)
	}

	var ids []string
	var s []map[string]interface{}

	for _, ns := range crResp.Data.Namespace {
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(ns.Namespace) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"name": ns.Namespace,
		}

		if detailedEnabled, ok := d.GetOk("enable_details"); ok && !detailedEnabled.(bool) {
			ids = append(ids, ns.Namespace)
			s = append(s, mapping)
			continue
		}

		raw, err := crService.DescribeNamespace(ns.Namespace)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		var resp crDescribeNamespaceResponse
		err = json.Unmarshal(raw.GetHttpContentBytes(), &resp)
		if err != nil {
			return WrapError(err)
		}

		mapping["auto_create"] = resp.Data.Namespace.AutoCreate
		mapping["default_visibility"] = resp.Data.Namespace.DefaultVisibility

		ids = append(ids, ns.Namespace)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("namespaces", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
