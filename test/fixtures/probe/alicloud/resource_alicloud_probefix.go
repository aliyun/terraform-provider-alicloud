package alicloud

// Fixture resource for tier-0 source-schema parser hermetic tests.
// 刻意布置五类 doc↔source gap;不对应真实 provider 资源,仅供 test/probe_test.sh。

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudProbeFix() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudProbeFixCreate,
		Read:   resourceAliCloudProbeFixRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flag_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forcenew_field": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deprecated_field": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'deprecated_field' has been deprecated.",
			},
			"undocumented_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nested_block": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inner_field": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"inner_forcenew": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudProbeFixCreate(d *schema.ResourceData, meta interface{}) error {
	action := "CreateProbeFix"
	_ = action
	return nil
}

func resourceAliCloudProbeFixRead(d *schema.ResourceData, meta interface{}) error {
	action := "DescribeProbeFix"
	_ = action
	return nil
}
