package alicloud

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRamRolePolicyDocument() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamRolePolicyDocumentRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validation.StringInSlice([]string{"1"}, false),
			},
			"statement": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Allow",
							ValidateFunc: validation.StringInSlice([]string{"Allow"}, false),
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "sts:AssumeRole",
							ValidateFunc: validation.StringInSlice([]string{"sts:AssumeRole"}, false),
						},
						"principal": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entity": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"RAM", "Service", "Federated"}, false),
									},
									"identifiers": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"condition": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},
									"variable": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"document": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudRamRolePolicyDocumentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if v, ok := d.GetOk("statement"); ok {
		ramService := RamService{client}
		doc, err := ramService.AssembleDataSourceRolePolicyDocument(v.([]interface{}), d.Get("version").(string))
		if err != nil {
			return WrapError(err)
		}

		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), doc)
		}

		d.Set("document", doc)

		d.SetId(tea.ToString(hashcode.String(doc)))
	}

	return nil
}
