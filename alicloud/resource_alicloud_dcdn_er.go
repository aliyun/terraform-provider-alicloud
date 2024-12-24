package alicloud

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDcdnEr() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDcdnErCreate,
		Read:   resourceAliCloudDcdnErRead,
		Update: resourceAliCloudDcdnErUpdate,
		Delete: resourceAliCloudDcdnErDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"er_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_conf": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"staging": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"production": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_anhui": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_beijing": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_chongqing": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_fujian": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_gansu": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_guangdong": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_guangxi": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_guizhou": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_hainan": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_hebei": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_heilongjiang": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_henan": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_hong_kong": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_hubei": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_hunan": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_jiangsu": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_jiangxi": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_jilin": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_liaoning": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_macau": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_neimenggu": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_ningxia": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_qinghai": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_shaanxi": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_shandong": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_shanghai": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_shanxi": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_sichuan": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_taiwan": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_tianjin": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_xinjiang": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_xizang": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_yunnan": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_zhejiang": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"preset_canary_overseas": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"5ms", "50ms", "100ms"}, false),
									},
									"code_rev": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_hosts": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudDcdnErCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRoutine"
	request := make(map[string]interface{})
	var err error

	request["Name"] = d.Get("er_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("env_conf"); ok {
		envConfMap := map[string]interface{}{}
		for _, envConfList := range v.([]interface{}) {
			envConfArg := envConfList.(map[string]interface{})

			if staging, ok := envConfArg["staging"]; ok {
				stagingMap := map[string]interface{}{}
				for _, stagingList := range staging.([]interface{}) {
					stagingArg := stagingList.(map[string]interface{})

					if specName, ok := stagingArg["spec_name"]; ok {
						stagingMap["SpecName"] = specName
					}

					if codeRev, ok := stagingArg["code_rev"]; ok {
						stagingMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := stagingArg["allowed_hosts"]; ok {
						stagingMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(stagingMap) > 0 {
					envConfMap["staging"] = stagingMap
				}
			}

			if production, ok := envConfArg["production"]; ok {
				productionMap := map[string]interface{}{}
				for _, productionList := range production.([]interface{}) {
					productionArg := productionList.(map[string]interface{})

					if specName, ok := productionArg["spec_name"]; ok {
						productionMap["SpecName"] = specName
					}

					if codeRev, ok := productionArg["code_rev"]; ok {
						productionMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := productionArg["allowed_hosts"]; ok {
						productionMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(productionMap) > 0 {
					envConfMap["production"] = productionMap
				}
			}

			if presetCanaryAnhui, ok := envConfArg["preset_canary_anhui"]; ok {
				presetCanaryAnhuiMap := map[string]interface{}{}
				for _, presetCanaryAnhuiList := range presetCanaryAnhui.([]interface{}) {
					presetCanaryAnhuiArg := presetCanaryAnhuiList.(map[string]interface{})

					if specName, ok := presetCanaryAnhuiArg["spec_name"]; ok {
						presetCanaryAnhuiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryAnhuiArg["code_rev"]; ok {
						presetCanaryAnhuiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryAnhuiArg["allowed_hosts"]; ok {
						presetCanaryAnhuiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryAnhuiMap) > 0 {
					envConfMap["presetCanaryAnhui"] = presetCanaryAnhuiMap
				}
			}

			if presetCanaryBeijing, ok := envConfArg["preset_canary_beijing"]; ok {
				presetCanaryBeijingMap := map[string]interface{}{}
				for _, presetCanaryBeijingList := range presetCanaryBeijing.([]interface{}) {
					presetCanaryBeijingArg := presetCanaryBeijingList.(map[string]interface{})

					if specName, ok := presetCanaryBeijingArg["spec_name"]; ok {
						presetCanaryBeijingMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryBeijingArg["code_rev"]; ok {
						presetCanaryBeijingMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryBeijingArg["allowed_hosts"]; ok {
						presetCanaryBeijingMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryBeijingMap) > 0 {
					envConfMap["presetCanaryBeijing"] = presetCanaryBeijingMap
				}
			}

			if presetCanaryChongqing, ok := envConfArg["preset_canary_chongqing"]; ok {
				presetCanaryChongqingMap := map[string]interface{}{}
				for _, presetCanaryChongqingList := range presetCanaryChongqing.([]interface{}) {
					presetCanaryChongqingArg := presetCanaryChongqingList.(map[string]interface{})

					if specName, ok := presetCanaryChongqingArg["spec_name"]; ok {
						presetCanaryChongqingMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryChongqingArg["code_rev"]; ok {
						presetCanaryChongqingMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryChongqingArg["allowed_hosts"]; ok {
						presetCanaryChongqingMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryChongqingMap) > 0 {
					envConfMap["presetCanaryChongqing"] = presetCanaryChongqingMap
				}
			}

			if presetCanaryFujian, ok := envConfArg["preset_canary_fujian"]; ok {
				presetCanaryFujianMap := map[string]interface{}{}
				for _, presetCanaryFujianList := range presetCanaryFujian.([]interface{}) {
					presetCanaryFujianArg := presetCanaryFujianList.(map[string]interface{})

					if specName, ok := presetCanaryFujianArg["spec_name"]; ok {
						presetCanaryFujianMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryFujianArg["code_rev"]; ok {
						presetCanaryFujianMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryFujianArg["allowed_hosts"]; ok {
						presetCanaryFujianMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryFujianMap) > 0 {
					envConfMap["presetCanaryFujian"] = presetCanaryFujianMap
				}
			}

			if presetCanaryGansu, ok := envConfArg["preset_canary_gansu"]; ok {
				presetCanaryGansuMap := map[string]interface{}{}
				for _, presetCanaryGansuList := range presetCanaryGansu.([]interface{}) {
					presetCanaryGansuArg := presetCanaryGansuList.(map[string]interface{})

					if specName, ok := presetCanaryGansuArg["spec_name"]; ok {
						presetCanaryGansuMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGansuArg["code_rev"]; ok {
						presetCanaryGansuMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGansuArg["allowed_hosts"]; ok {
						presetCanaryGansuMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGansuMap) > 0 {
					envConfMap["presetCanaryGansu"] = presetCanaryGansuMap
				}
			}

			if presetCanaryGuangdong, ok := envConfArg["preset_canary_guangdong"]; ok {
				presetCanaryGuangdongMap := map[string]interface{}{}
				for _, presetCanaryGuangdongList := range presetCanaryGuangdong.([]interface{}) {
					presetCanaryGuangdongArg := presetCanaryGuangdongList.(map[string]interface{})

					if specName, ok := presetCanaryGuangdongArg["spec_name"]; ok {
						presetCanaryGuangdongMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGuangdongArg["code_rev"]; ok {
						presetCanaryGuangdongMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGuangdongArg["allowed_hosts"]; ok {
						presetCanaryGuangdongMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGuangdongMap) > 0 {
					envConfMap["presetCanaryGuangdong"] = presetCanaryGuangdongMap
				}
			}

			if presetCanaryGuangxi, ok := envConfArg["preset_canary_guangxi"]; ok {
				presetCanaryGuangxiMap := map[string]interface{}{}
				for _, presetCanaryGuangxiList := range presetCanaryGuangxi.([]interface{}) {
					presetCanaryGuangxiArg := presetCanaryGuangxiList.(map[string]interface{})

					if specName, ok := presetCanaryGuangxiArg["spec_name"]; ok {
						presetCanaryGuangxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGuangxiArg["code_rev"]; ok {
						presetCanaryGuangxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGuangxiArg["allowed_hosts"]; ok {
						presetCanaryGuangxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGuangxiMap) > 0 {
					envConfMap["presetCanaryGuangxi"] = presetCanaryGuangxiMap
				}
			}

			if presetCanaryGuizhou, ok := envConfArg["preset_canary_guizhou"]; ok {
				presetCanaryGuizhouMap := map[string]interface{}{}
				for _, presetCanaryGuizhouList := range presetCanaryGuizhou.([]interface{}) {
					presetCanaryGuizhouArg := presetCanaryGuizhouList.(map[string]interface{})

					if specName, ok := presetCanaryGuizhouArg["spec_name"]; ok {
						presetCanaryGuizhouMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGuizhouArg["code_rev"]; ok {
						presetCanaryGuizhouMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGuizhouArg["allowed_hosts"]; ok {
						presetCanaryGuizhouMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGuizhouMap) > 0 {
					envConfMap["presetCanaryGuizhou"] = presetCanaryGuizhouMap
				}
			}

			if presetCanaryHainan, ok := envConfArg["preset_canary_hainan"]; ok {
				presetCanaryHainanMap := map[string]interface{}{}
				for _, presetCanaryHainanList := range presetCanaryHainan.([]interface{}) {
					presetCanaryHainanArg := presetCanaryHainanList.(map[string]interface{})

					if specName, ok := presetCanaryHainanArg["spec_name"]; ok {
						presetCanaryHainanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHainanArg["code_rev"]; ok {
						presetCanaryHainanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHainanArg["allowed_hosts"]; ok {
						presetCanaryHainanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHainanMap) > 0 {
					envConfMap["presetCanaryHainan"] = presetCanaryHainanMap
				}
			}

			if presetCanaryHebei, ok := envConfArg["preset_canary_hebei"]; ok {
				presetCanaryHebeiMap := map[string]interface{}{}
				for _, presetCanaryHebeiList := range presetCanaryHebei.([]interface{}) {
					presetCanaryHebeiArg := presetCanaryHebeiList.(map[string]interface{})

					if specName, ok := presetCanaryHebeiArg["spec_name"]; ok {
						presetCanaryHebeiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHebeiArg["code_rev"]; ok {
						presetCanaryHebeiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHebeiArg["allowed_hosts"]; ok {
						presetCanaryHebeiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHebeiMap) > 0 {
					envConfMap["presetCanaryHebei"] = presetCanaryHebeiMap
				}
			}

			if presetCanaryHeilongjiang, ok := envConfArg["preset_canary_heilongjiang"]; ok {
				presetCanaryHeilongjiangMap := map[string]interface{}{}
				for _, presetCanaryHeilongjiangList := range presetCanaryHeilongjiang.([]interface{}) {
					presetCanaryHeilongjiangArg := presetCanaryHeilongjiangList.(map[string]interface{})

					if specName, ok := presetCanaryHeilongjiangArg["spec_name"]; ok {
						presetCanaryHeilongjiangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHeilongjiangArg["code_rev"]; ok {
						presetCanaryHeilongjiangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHeilongjiangArg["allowed_hosts"]; ok {
						presetCanaryHeilongjiangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHeilongjiangMap) > 0 {
					envConfMap["presetCanaryHeilongjiang"] = presetCanaryHeilongjiangMap
				}
			}

			if presetCanaryHenan, ok := envConfArg["preset_canary_henan"]; ok {
				presetCanaryHenanMap := map[string]interface{}{}
				for _, presetCanaryHenanList := range presetCanaryHenan.([]interface{}) {
					presetCanaryHenanArg := presetCanaryHenanList.(map[string]interface{})

					if specName, ok := presetCanaryHenanArg["spec_name"]; ok {
						presetCanaryHenanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHenanArg["code_rev"]; ok {
						presetCanaryHenanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHenanArg["allowed_hosts"]; ok {
						presetCanaryHenanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHenanMap) > 0 {
					envConfMap["presetCanaryHenan"] = presetCanaryHenanMap
				}
			}

			if presetCanaryHongKong, ok := envConfArg["preset_canary_hong_kong"]; ok {
				presetCanaryHongKongMap := map[string]interface{}{}
				for _, presetCanaryHongKongList := range presetCanaryHongKong.([]interface{}) {
					presetCanaryHongKongArg := presetCanaryHongKongList.(map[string]interface{})

					if specName, ok := presetCanaryHongKongArg["spec_name"]; ok {
						presetCanaryHongKongMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHongKongArg["code_rev"]; ok {
						presetCanaryHongKongMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHongKongArg["allowed_hosts"]; ok {
						presetCanaryHongKongMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHongKongMap) > 0 {
					envConfMap["presetCanaryHongKong"] = presetCanaryHongKongMap
				}
			}

			if presetCanaryHubei, ok := envConfArg["preset_canary_hubei"]; ok {
				presetCanaryHubeiMap := map[string]interface{}{}
				for _, presetCanaryHubeiList := range presetCanaryHubei.([]interface{}) {
					presetCanaryHubeiArg := presetCanaryHubeiList.(map[string]interface{})

					if specName, ok := presetCanaryHubeiArg["spec_name"]; ok {
						presetCanaryHubeiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHubeiArg["code_rev"]; ok {
						presetCanaryHubeiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHubeiArg["allowed_hosts"]; ok {
						presetCanaryHubeiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHubeiMap) > 0 {
					envConfMap["presetCanaryHubei"] = presetCanaryHubeiMap
				}
			}

			if presetCanaryHunan, ok := envConfArg["preset_canary_hunan"]; ok {
				presetCanaryHunanMap := map[string]interface{}{}
				for _, presetCanaryHunanList := range presetCanaryHunan.([]interface{}) {
					presetCanaryHunanArg := presetCanaryHunanList.(map[string]interface{})

					if specName, ok := presetCanaryHunanArg["spec_name"]; ok {
						presetCanaryHunanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHunanArg["code_rev"]; ok {
						presetCanaryHunanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHunanArg["allowed_hosts"]; ok {
						presetCanaryHunanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHunanMap) > 0 {
					envConfMap["presetCanaryHunan"] = presetCanaryHunanMap
				}
			}

			if presetCanaryJiangsu, ok := envConfArg["preset_canary_jiangsu"]; ok {
				presetCanaryJiangsuMap := map[string]interface{}{}
				for _, presetCanaryJiangsuList := range presetCanaryJiangsu.([]interface{}) {
					presetCanaryJiangsuArg := presetCanaryJiangsuList.(map[string]interface{})

					if specName, ok := presetCanaryJiangsuArg["spec_name"]; ok {
						presetCanaryJiangsuMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryJiangsuArg["code_rev"]; ok {
						presetCanaryJiangsuMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryJiangsuArg["allowed_hosts"]; ok {
						presetCanaryJiangsuMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryJiangsuMap) > 0 {
					envConfMap["presetCanaryJiangsu"] = presetCanaryJiangsuMap
				}
			}

			if presetCanaryJiangxi, ok := envConfArg["preset_canary_jiangxi"]; ok {
				presetCanaryJiangxiMap := map[string]interface{}{}
				for _, presetCanaryJiangxiList := range presetCanaryJiangxi.([]interface{}) {
					presetCanaryJiangxiArg := presetCanaryJiangxiList.(map[string]interface{})

					if specName, ok := presetCanaryJiangxiArg["spec_name"]; ok {
						presetCanaryJiangxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryJiangxiArg["code_rev"]; ok {
						presetCanaryJiangxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryJiangxiArg["allowed_hosts"]; ok {
						presetCanaryJiangxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryJiangxiMap) > 0 {
					envConfMap["presetCanaryJiangxi"] = presetCanaryJiangxiMap
				}
			}

			if presetCanaryJilin, ok := envConfArg["preset_canary_jilin"]; ok {
				presetCanaryJilinMap := map[string]interface{}{}
				for _, presetCanaryJilinList := range presetCanaryJilin.([]interface{}) {
					presetCanaryJilinArg := presetCanaryJilinList.(map[string]interface{})

					if specName, ok := presetCanaryJilinArg["spec_name"]; ok {
						presetCanaryJilinMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryJilinArg["code_rev"]; ok {
						presetCanaryJilinMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryJilinArg["allowed_hosts"]; ok {
						presetCanaryJilinMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryJilinMap) > 0 {
					envConfMap["presetCanaryJilin"] = presetCanaryJilinMap
				}
			}

			if presetCanaryLiaoning, ok := envConfArg["preset_canary_liaoning"]; ok {
				presetCanaryLiaoningMap := map[string]interface{}{}
				for _, presetCanaryLiaoningList := range presetCanaryLiaoning.([]interface{}) {
					presetCanaryLiaoningArg := presetCanaryLiaoningList.(map[string]interface{})

					if specName, ok := presetCanaryLiaoningArg["spec_name"]; ok {
						presetCanaryLiaoningMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryLiaoningArg["code_rev"]; ok {
						presetCanaryLiaoningMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryLiaoningArg["allowed_hosts"]; ok {
						presetCanaryLiaoningMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryLiaoningMap) > 0 {
					envConfMap["presetCanaryLiaoning"] = presetCanaryLiaoningMap
				}
			}

			if presetCanaryMacau, ok := envConfArg["preset_canary_macau"]; ok {
				presetCanaryMacauMap := map[string]interface{}{}
				for _, presetCanaryMacauList := range presetCanaryMacau.([]interface{}) {
					presetCanaryMacauArg := presetCanaryMacauList.(map[string]interface{})

					if specName, ok := presetCanaryMacauArg["spec_name"]; ok {
						presetCanaryMacauMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryMacauArg["code_rev"]; ok {
						presetCanaryMacauMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryMacauArg["allowed_hosts"]; ok {
						presetCanaryMacauMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryMacauMap) > 0 {
					envConfMap["presetCanaryMacau"] = presetCanaryMacauMap
				}
			}

			if presetCanaryNeimenggu, ok := envConfArg["preset_canary_neimenggu"]; ok {
				presetCanaryNeimengguMap := map[string]interface{}{}
				for _, presetCanaryNeimengguList := range presetCanaryNeimenggu.([]interface{}) {
					presetCanaryNeimengguArg := presetCanaryNeimengguList.(map[string]interface{})

					if specName, ok := presetCanaryNeimengguArg["spec_name"]; ok {
						presetCanaryNeimengguMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryNeimengguArg["code_rev"]; ok {
						presetCanaryNeimengguMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryNeimengguArg["allowed_hosts"]; ok {
						presetCanaryNeimengguMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryNeimengguMap) > 0 {
					envConfMap["presetCanaryNeimenggu"] = presetCanaryNeimengguMap
				}
			}

			if presetCanaryNingxia, ok := envConfArg["preset_canary_ningxia"]; ok {
				presetCanaryNingxiaMap := map[string]interface{}{}
				for _, presetCanaryNingxiaList := range presetCanaryNingxia.([]interface{}) {
					presetCanaryNingxiaArg := presetCanaryNingxiaList.(map[string]interface{})

					if specName, ok := presetCanaryNingxiaArg["spec_name"]; ok {
						presetCanaryNingxiaMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryNingxiaArg["code_rev"]; ok {
						presetCanaryNingxiaMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryNingxiaArg["allowed_hosts"]; ok {
						presetCanaryNingxiaMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryNingxiaMap) > 0 {
					envConfMap["presetCanaryNingxia"] = presetCanaryNingxiaMap
				}
			}

			if presetCanaryQinghai, ok := envConfArg["preset_canary_qinghai"]; ok {
				presetCanaryQinghaiMap := map[string]interface{}{}
				for _, presetCanaryQinghaiList := range presetCanaryQinghai.([]interface{}) {
					presetCanaryQinghaiArg := presetCanaryQinghaiList.(map[string]interface{})

					if specName, ok := presetCanaryQinghaiArg["spec_name"]; ok {
						presetCanaryQinghaiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryQinghaiArg["code_rev"]; ok {
						presetCanaryQinghaiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryQinghaiArg["allowed_hosts"]; ok {
						presetCanaryQinghaiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryQinghaiMap) > 0 {
					envConfMap["presetCanaryQinghai"] = presetCanaryQinghaiMap
				}
			}

			if presetCanaryShaanxi, ok := envConfArg["preset_canary_shaanxi"]; ok {
				presetCanaryShaanxiMap := map[string]interface{}{}
				for _, presetCanaryShaanxiList := range presetCanaryShaanxi.([]interface{}) {
					presetCanaryShaanxiArg := presetCanaryShaanxiList.(map[string]interface{})

					if specName, ok := presetCanaryShaanxiArg["spec_name"]; ok {
						presetCanaryShaanxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShaanxiArg["code_rev"]; ok {
						presetCanaryShaanxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShaanxiArg["allowed_hosts"]; ok {
						presetCanaryShaanxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShaanxiMap) > 0 {
					envConfMap["presetCanaryShaanxi"] = presetCanaryShaanxiMap
				}
			}

			if presetCanaryShandong, ok := envConfArg["preset_canary_shandong"]; ok {
				presetCanaryShandongMap := map[string]interface{}{}
				for _, presetCanaryShandongList := range presetCanaryShandong.([]interface{}) {
					presetCanaryShandongArg := presetCanaryShandongList.(map[string]interface{})

					if specName, ok := presetCanaryShandongArg["spec_name"]; ok {
						presetCanaryShandongMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShandongArg["code_rev"]; ok {
						presetCanaryShandongMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShandongArg["allowed_hosts"]; ok {
						presetCanaryShandongMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShandongMap) > 0 {
					envConfMap["presetCanaryShandong"] = presetCanaryShandongMap
				}
			}

			if presetCanaryShanghai, ok := envConfArg["preset_canary_shanghai"]; ok {
				presetCanaryShanghaiMap := map[string]interface{}{}
				for _, presetCanaryShanghaiList := range presetCanaryShanghai.([]interface{}) {
					presetCanaryShanghaiArg := presetCanaryShanghaiList.(map[string]interface{})

					if specName, ok := presetCanaryShanghaiArg["spec_name"]; ok {
						presetCanaryShanghaiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShanghaiArg["code_rev"]; ok {
						presetCanaryShanghaiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShanghaiArg["allowed_hosts"]; ok {
						presetCanaryShanghaiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShanghaiMap) > 0 {
					envConfMap["presetCanaryShanghai"] = presetCanaryShanghaiMap
				}
			}

			if presetCanaryShanxi, ok := envConfArg["preset_canary_shanxi"]; ok {
				presetCanaryShanxiMap := map[string]interface{}{}
				for _, presetCanaryShanxiList := range presetCanaryShanxi.([]interface{}) {
					presetCanaryShanxiArg := presetCanaryShanxiList.(map[string]interface{})

					if specName, ok := presetCanaryShanxiArg["spec_name"]; ok {
						presetCanaryShanxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShanxiArg["code_rev"]; ok {
						presetCanaryShanxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShanxiArg["allowed_hosts"]; ok {
						presetCanaryShanxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShanxiMap) > 0 {
					envConfMap["presetCanaryShanxi"] = presetCanaryShanxiMap
				}
			}

			if presetCanarySichuan, ok := envConfArg["preset_canary_sichuan"]; ok {
				presetCanarySichuanMap := map[string]interface{}{}
				for _, presetCanarySichuanList := range presetCanarySichuan.([]interface{}) {
					presetCanarySichuanArg := presetCanarySichuanList.(map[string]interface{})

					if specName, ok := presetCanarySichuanArg["spec_name"]; ok {
						presetCanarySichuanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanarySichuanArg["code_rev"]; ok {
						presetCanarySichuanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanarySichuanArg["allowed_hosts"]; ok {
						presetCanarySichuanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanarySichuanMap) > 0 {
					envConfMap["presetCanarySichuan"] = presetCanarySichuanMap
				}
			}

			if presetCanaryTaiwan, ok := envConfArg["preset_canary_taiwan"]; ok {
				presetCanaryTaiwanMap := map[string]interface{}{}
				for _, presetCanaryTaiwanList := range presetCanaryTaiwan.([]interface{}) {
					presetCanaryTaiwanArg := presetCanaryTaiwanList.(map[string]interface{})

					if specName, ok := presetCanaryTaiwanArg["spec_name"]; ok {
						presetCanaryTaiwanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryTaiwanArg["code_rev"]; ok {
						presetCanaryTaiwanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryTaiwanArg["allowed_hosts"]; ok {
						presetCanaryTaiwanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryTaiwanMap) > 0 {
					envConfMap["presetCanaryTaiwan"] = presetCanaryTaiwanMap
				}
			}

			if presetCanaryTianjin, ok := envConfArg["preset_canary_tianjin"]; ok {
				presetCanaryTianjinMap := map[string]interface{}{}
				for _, presetCanaryTianjinList := range presetCanaryTianjin.([]interface{}) {
					presetCanaryTianjinArg := presetCanaryTianjinList.(map[string]interface{})

					if specName, ok := presetCanaryTianjinArg["spec_name"]; ok {
						presetCanaryTianjinMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryTianjinArg["code_rev"]; ok {
						presetCanaryTianjinMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryTianjinArg["allowed_hosts"]; ok {
						presetCanaryTianjinMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryTianjinMap) > 0 {
					envConfMap["presetCanaryTianjin"] = presetCanaryTianjinMap
				}
			}

			if presetCanaryXinjiang, ok := envConfArg["preset_canary_xinjiang"]; ok {
				presetCanaryXinjiangMap := map[string]interface{}{}
				for _, presetCanaryXinjiangList := range presetCanaryXinjiang.([]interface{}) {
					presetCanaryXinjiangArg := presetCanaryXinjiangList.(map[string]interface{})

					if specName, ok := presetCanaryXinjiangArg["spec_name"]; ok {
						presetCanaryXinjiangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryXinjiangArg["code_rev"]; ok {
						presetCanaryXinjiangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryXinjiangArg["allowed_hosts"]; ok {
						presetCanaryXinjiangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryXinjiangMap) > 0 {
					envConfMap["presetCanaryXinjiang"] = presetCanaryXinjiangMap
				}
			}

			if presetCanaryXizang, ok := envConfArg["preset_canary_xizang"]; ok {
				presetCanaryXizangMap := map[string]interface{}{}
				for _, presetCanaryXizangList := range presetCanaryXizang.([]interface{}) {
					presetCanaryXizangArg := presetCanaryXizangList.(map[string]interface{})

					if specName, ok := presetCanaryXizangArg["spec_name"]; ok {
						presetCanaryXizangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryXizangArg["code_rev"]; ok {
						presetCanaryXizangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryXizangArg["allowed_hosts"]; ok {
						presetCanaryXizangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryXizangMap) > 0 {
					envConfMap["presetCanaryXizang"] = presetCanaryXizangMap
				}
			}

			if presetCanaryYunnan, ok := envConfArg["preset_canary_yunnan"]; ok {
				presetCanaryYunnanMap := map[string]interface{}{}
				for _, presetCanaryYunnanList := range presetCanaryYunnan.([]interface{}) {
					presetCanaryYunnanArg := presetCanaryYunnanList.(map[string]interface{})

					if specName, ok := presetCanaryYunnanArg["spec_name"]; ok {
						presetCanaryYunnanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryYunnanArg["code_rev"]; ok {
						presetCanaryYunnanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryYunnanArg["allowed_hosts"]; ok {
						presetCanaryYunnanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryYunnanMap) > 0 {
					envConfMap["presetCanaryYunnan"] = presetCanaryYunnanMap
				}
			}

			if presetCanaryZhejiang, ok := envConfArg["preset_canary_zhejiang"]; ok {
				presetCanaryZhejiangMap := map[string]interface{}{}
				for _, presetCanaryZhejiangList := range presetCanaryZhejiang.([]interface{}) {
					presetCanaryZhejiangArg := presetCanaryZhejiangList.(map[string]interface{})

					if specName, ok := presetCanaryZhejiangArg["spec_name"]; ok {
						presetCanaryZhejiangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryZhejiangArg["code_rev"]; ok {
						presetCanaryZhejiangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryZhejiangArg["allowed_hosts"]; ok {
						presetCanaryZhejiangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryZhejiangMap) > 0 {
					envConfMap["presetCanaryZhejiang"] = presetCanaryZhejiangMap
				}
			}

			if presetCanaryOverseas, ok := envConfArg["preset_canary_overseas"]; ok {
				presetCanaryOverseasMap := map[string]interface{}{}
				for _, presetCanaryOverseasList := range presetCanaryOverseas.([]interface{}) {
					presetCanaryOverseasArg := presetCanaryOverseasList.(map[string]interface{})

					if specName, ok := presetCanaryOverseasArg["spec_name"]; ok {
						presetCanaryOverseasMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryOverseasArg["code_rev"]; ok {
						presetCanaryOverseasMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryOverseasArg["allowed_hosts"]; ok {
						presetCanaryOverseasMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryOverseasMap) > 0 {
					envConfMap["presetCanaryOverseas"] = presetCanaryOverseasMap
				}
			}
		}

		envConfJson, err := convertMaptoJsonString(envConfMap)
		if err != nil {
			return WrapError(err)
		}

		request["EnvConf"] = envConfJson
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_er", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Name"]))

	return resourceAliCloudDcdnErRead(d, meta)
}

func resourceAliCloudDcdnErRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}

	object, err := dcdnService.DescribeDcdnEr(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("er_name", d.Id())

	description, err := base64.StdEncoding.DecodeString(fmt.Sprint(object["Description"]))
	if err != nil {
		return WrapError(err)
	}

	d.Set("description", string(description))

	if envConf, ok := object["EnvConf"]; ok {
		envConfMaps := make([]map[string]interface{}, 0)
		envConfArg := envConf.(map[string]interface{})
		envConfMap := map[string]interface{}{}

		if staging, ok := envConfArg["staging"]; ok {
			stagingMaps := make([]map[string]interface{}, 0)
			stagingArg := staging.(map[string]interface{})

			if len(stagingArg) > 0 {
				stagingMap := map[string]interface{}{}

				if specName, ok := stagingArg["SpecName"]; ok {
					stagingMap["spec_name"] = specName
				}

				if codeRev, ok := stagingArg["CodeRev"]; ok {
					stagingMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := stagingArg["AllowedHosts"]; ok {
					stagingMap["allowed_hosts"] = allowedHosts
				}

				stagingMaps = append(stagingMaps, stagingMap)
				envConfMap["staging"] = stagingMaps
			}
		}

		if production, ok := envConfArg["production"]; ok {
			productionMaps := make([]map[string]interface{}, 0)
			productionArg := production.(map[string]interface{})

			if len(productionArg) > 0 {
				productionMap := map[string]interface{}{}

				if specName, ok := productionArg["SpecName"]; ok {
					productionMap["spec_name"] = specName
				}

				if codeRev, ok := productionArg["CodeRev"]; ok {
					productionMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := productionArg["AllowedHosts"]; ok {
					productionMap["allowed_hosts"] = allowedHosts
				}

				productionMaps = append(productionMaps, productionMap)
				envConfMap["production"] = productionMaps
			}
		}

		if presetCanaryAnhui, ok := envConfArg["presetCanaryAnhui"]; ok {
			presetCanaryAnhuiMaps := make([]map[string]interface{}, 0)
			presetCanaryAnhuiArg := presetCanaryAnhui.(map[string]interface{})

			if len(presetCanaryAnhuiArg) > 0 {
				presetCanaryAnhuiMap := map[string]interface{}{}

				if specName, ok := presetCanaryAnhuiArg["SpecName"]; ok {
					presetCanaryAnhuiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryAnhuiArg["CodeRev"]; ok {
					presetCanaryAnhuiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryAnhuiArg["AllowedHosts"]; ok {
					presetCanaryAnhuiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryAnhuiMaps = append(presetCanaryAnhuiMaps, presetCanaryAnhuiMap)
				envConfMap["preset_canary_anhui"] = presetCanaryAnhuiMaps
			}
		}

		if presetCanaryBeijing, ok := envConfArg["presetCanaryBeijing"]; ok {
			presetCanaryBeijingMaps := make([]map[string]interface{}, 0)
			presetCanaryBeijingArg := presetCanaryBeijing.(map[string]interface{})

			if len(presetCanaryBeijingArg) > 0 {
				presetCanaryBeijingMap := map[string]interface{}{}

				if specName, ok := presetCanaryBeijingArg["SpecName"]; ok {
					presetCanaryBeijingMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryBeijingArg["CodeRev"]; ok {
					presetCanaryBeijingMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryBeijingArg["AllowedHosts"]; ok {
					presetCanaryBeijingMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryBeijingMaps = append(presetCanaryBeijingMaps, presetCanaryBeijingMap)
				envConfMap["preset_canary_beijing"] = presetCanaryBeijingMaps
			}
		}

		if presetCanaryChongqing, ok := envConfArg["presetCanaryChongqing"]; ok {
			presetCanaryChongqingMaps := make([]map[string]interface{}, 0)
			presetCanaryChongqingArg := presetCanaryChongqing.(map[string]interface{})

			if len(presetCanaryChongqingArg) > 0 {
				presetCanaryChongqingMap := map[string]interface{}{}

				if specName, ok := presetCanaryChongqingArg["SpecName"]; ok {
					presetCanaryChongqingMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryChongqingArg["CodeRev"]; ok {
					presetCanaryChongqingMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryChongqingArg["AllowedHosts"]; ok {
					presetCanaryChongqingMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryChongqingMaps = append(presetCanaryChongqingMaps, presetCanaryChongqingMap)
				envConfMap["preset_canary_chongqing"] = presetCanaryChongqingMaps
			}
		}

		if presetCanaryFujian, ok := envConfArg["presetCanaryFujian"]; ok {
			presetCanaryFujianMaps := make([]map[string]interface{}, 0)
			presetCanaryFujianArg := presetCanaryFujian.(map[string]interface{})

			if len(presetCanaryFujianArg) > 0 {
				presetCanaryFujianMap := map[string]interface{}{}

				if specName, ok := presetCanaryFujianArg["SpecName"]; ok {
					presetCanaryFujianMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryFujianArg["CodeRev"]; ok {
					presetCanaryFujianMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryFujianArg["AllowedHosts"]; ok {
					presetCanaryFujianMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryFujianMaps = append(presetCanaryFujianMaps, presetCanaryFujianMap)
				envConfMap["preset_canary_fujian"] = presetCanaryFujianMaps
			}
		}

		if presetCanaryGansu, ok := envConfArg["presetCanaryGansu"]; ok {
			presetCanaryGansuMaps := make([]map[string]interface{}, 0)
			presetCanaryGansuArg := presetCanaryGansu.(map[string]interface{})

			if len(presetCanaryGansuArg) > 0 {
				presetCanaryGansuMap := map[string]interface{}{}

				if specName, ok := presetCanaryGansuArg["SpecName"]; ok {
					presetCanaryGansuMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryGansuArg["CodeRev"]; ok {
					presetCanaryGansuMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryGansuArg["AllowedHosts"]; ok {
					presetCanaryGansuMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryGansuMaps = append(presetCanaryGansuMaps, presetCanaryGansuMap)
				envConfMap["preset_canary_gansu"] = presetCanaryGansuMaps
			}
		}

		if presetCanaryGuangdong, ok := envConfArg["presetCanaryGuangdong"]; ok {
			presetCanaryGuangdongMaps := make([]map[string]interface{}, 0)
			presetCanaryGuangdongArg := presetCanaryGuangdong.(map[string]interface{})

			if len(presetCanaryGuangdongArg) > 0 {
				presetCanaryGuangdongMap := map[string]interface{}{}

				if specName, ok := presetCanaryGuangdongArg["SpecName"]; ok {
					presetCanaryGuangdongMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryGuangdongArg["CodeRev"]; ok {
					presetCanaryGuangdongMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryGuangdongArg["AllowedHosts"]; ok {
					presetCanaryGuangdongMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryGuangdongMaps = append(presetCanaryGuangdongMaps, presetCanaryGuangdongMap)
				envConfMap["preset_canary_guangdong"] = presetCanaryGuangdongMaps
			}
		}

		if presetCanaryGuangxi, ok := envConfArg["presetCanaryGuangxi"]; ok {
			presetCanaryGuangxiMaps := make([]map[string]interface{}, 0)
			presetCanaryGuangxiArg := presetCanaryGuangxi.(map[string]interface{})

			if len(presetCanaryGuangxiArg) > 0 {
				presetCanaryGuangxiMap := map[string]interface{}{}

				if specName, ok := presetCanaryGuangxiArg["SpecName"]; ok {
					presetCanaryGuangxiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryGuangxiArg["CodeRev"]; ok {
					presetCanaryGuangxiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryGuangxiArg["AllowedHosts"]; ok {
					presetCanaryGuangxiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryGuangxiMaps = append(presetCanaryGuangxiMaps, presetCanaryGuangxiMap)
				envConfMap["preset_canary_guangxi"] = presetCanaryGuangxiMaps
			}
		}

		if presetCanaryGuizhou, ok := envConfArg["presetCanaryGuizhou"]; ok {
			presetCanaryGuizhouMaps := make([]map[string]interface{}, 0)
			presetCanaryGuizhouArg := presetCanaryGuizhou.(map[string]interface{})

			if len(presetCanaryGuizhouArg) > 0 {
				presetCanaryGuizhouMap := map[string]interface{}{}

				if specName, ok := presetCanaryGuizhouArg["SpecName"]; ok {
					presetCanaryGuizhouMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryGuizhouArg["CodeRev"]; ok {
					presetCanaryGuizhouMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryGuizhouArg["AllowedHosts"]; ok {
					presetCanaryGuizhouMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryGuizhouMaps = append(presetCanaryGuizhouMaps, presetCanaryGuizhouMap)
				envConfMap["preset_canary_guizhou"] = presetCanaryGuizhouMaps
			}
		}

		if presetCanaryHainan, ok := envConfArg["presetCanaryHainan"]; ok {
			presetCanaryHainanMaps := make([]map[string]interface{}, 0)
			presetCanaryHainanArg := presetCanaryHainan.(map[string]interface{})

			if len(presetCanaryHainanArg) > 0 {
				presetCanaryHainanMap := map[string]interface{}{}

				if specName, ok := presetCanaryHainanArg["SpecName"]; ok {
					presetCanaryHainanMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHainanArg["CodeRev"]; ok {
					presetCanaryHainanMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHainanArg["AllowedHosts"]; ok {
					presetCanaryHainanMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHainanMaps = append(presetCanaryHainanMaps, presetCanaryHainanMap)
				envConfMap["preset_canary_hainan"] = presetCanaryHainanMaps
			}
		}

		if presetCanaryHebei, ok := envConfArg["presetCanaryHebei"]; ok {
			presetCanaryHebeiMaps := make([]map[string]interface{}, 0)
			presetCanaryHebeiArg := presetCanaryHebei.(map[string]interface{})

			if len(presetCanaryHebeiArg) > 0 {
				presetCanaryHebeiMap := map[string]interface{}{}

				if specName, ok := presetCanaryHebeiArg["SpecName"]; ok {
					presetCanaryHebeiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHebeiArg["CodeRev"]; ok {
					presetCanaryHebeiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHebeiArg["AllowedHosts"]; ok {
					presetCanaryHebeiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHebeiMaps = append(presetCanaryHebeiMaps, presetCanaryHebeiMap)
				envConfMap["preset_canary_hebei"] = presetCanaryHebeiMaps
			}
		}

		if presetCanaryHeilongjiang, ok := envConfArg["presetCanaryHeilongjiang"]; ok {
			presetCanaryHeilongjiangMaps := make([]map[string]interface{}, 0)
			presetCanaryHeilongjiangArg := presetCanaryHeilongjiang.(map[string]interface{})

			if len(presetCanaryHeilongjiangArg) > 0 {
				presetCanaryHeilongjiangMap := map[string]interface{}{}

				if specName, ok := presetCanaryHeilongjiangArg["SpecName"]; ok {
					presetCanaryHeilongjiangMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHeilongjiangArg["CodeRev"]; ok {
					presetCanaryHeilongjiangMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHeilongjiangArg["AllowedHosts"]; ok {
					presetCanaryHeilongjiangMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHeilongjiangMaps = append(presetCanaryHeilongjiangMaps, presetCanaryHeilongjiangMap)
				envConfMap["preset_canary_heilongjiang"] = presetCanaryHeilongjiangMaps
			}
		}

		if presetCanaryHenan, ok := envConfArg["presetCanaryHenan"]; ok {
			presetCanaryHenanMaps := make([]map[string]interface{}, 0)
			presetCanaryHenanArg := presetCanaryHenan.(map[string]interface{})

			if len(presetCanaryHenanArg) > 0 {
				presetCanaryHenanMap := map[string]interface{}{}

				if specName, ok := presetCanaryHenanArg["SpecName"]; ok {
					presetCanaryHenanMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHenanArg["CodeRev"]; ok {
					presetCanaryHenanMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHenanArg["AllowedHosts"]; ok {
					presetCanaryHenanMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHenanMaps = append(presetCanaryHenanMaps, presetCanaryHenanMap)
				envConfMap["preset_canary_henan"] = presetCanaryHenanMaps
			}
		}

		if presetCanaryHongKong, ok := envConfArg["presetCanaryHongKong"]; ok {
			presetCanaryHongKongMaps := make([]map[string]interface{}, 0)
			presetCanaryHongKongArg := presetCanaryHongKong.(map[string]interface{})

			if len(presetCanaryHongKongArg) > 0 {
				presetCanaryHongKongMap := map[string]interface{}{}

				if specName, ok := presetCanaryHongKongArg["SpecName"]; ok {
					presetCanaryHongKongMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHongKongArg["CodeRev"]; ok {
					presetCanaryHongKongMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHongKongArg["AllowedHosts"]; ok {
					presetCanaryHongKongMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHongKongMaps = append(presetCanaryHongKongMaps, presetCanaryHongKongMap)
				envConfMap["preset_canary_hong_kong"] = presetCanaryHongKongMaps
			}
		}

		if presetCanaryHubei, ok := envConfArg["presetCanaryHubei"]; ok {
			presetCanaryHubeiMaps := make([]map[string]interface{}, 0)
			presetCanaryHubeiArg := presetCanaryHubei.(map[string]interface{})

			if len(presetCanaryHubeiArg) > 0 {
				presetCanaryHubeiMap := map[string]interface{}{}

				if specName, ok := presetCanaryHubeiArg["SpecName"]; ok {
					presetCanaryHubeiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHubeiArg["CodeRev"]; ok {
					presetCanaryHubeiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHubeiArg["AllowedHosts"]; ok {
					presetCanaryHubeiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHubeiMaps = append(presetCanaryHubeiMaps, presetCanaryHubeiMap)
				envConfMap["preset_canary_hubei"] = presetCanaryHubeiMaps
			}
		}

		if presetCanaryHunan, ok := envConfArg["presetCanaryHunan"]; ok {
			presetCanaryHunanMaps := make([]map[string]interface{}, 0)
			presetCanaryHunanArg := presetCanaryHunan.(map[string]interface{})

			if len(presetCanaryHunanArg) > 0 {
				presetCanaryHunanMap := map[string]interface{}{}

				if specName, ok := presetCanaryHunanArg["SpecName"]; ok {
					presetCanaryHunanMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryHunanArg["CodeRev"]; ok {
					presetCanaryHunanMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryHunanArg["AllowedHosts"]; ok {
					presetCanaryHunanMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryHunanMaps = append(presetCanaryHunanMaps, presetCanaryHunanMap)
				envConfMap["preset_canary_hunan"] = presetCanaryHunanMaps
			}
		}

		if presetCanaryJiangsu, ok := envConfArg["presetCanaryJiangsu"]; ok {
			presetCanaryJiangsuMaps := make([]map[string]interface{}, 0)
			presetCanaryJiangsuArg := presetCanaryJiangsu.(map[string]interface{})

			if len(presetCanaryJiangsuArg) > 0 {
				presetCanaryJiangsuMap := map[string]interface{}{}

				if specName, ok := presetCanaryJiangsuArg["SpecName"]; ok {
					presetCanaryJiangsuMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryJiangsuArg["CodeRev"]; ok {
					presetCanaryJiangsuMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryJiangsuArg["AllowedHosts"]; ok {
					presetCanaryJiangsuMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryJiangsuMaps = append(presetCanaryJiangsuMaps, presetCanaryJiangsuMap)
				envConfMap["preset_canary_jiangsu"] = presetCanaryJiangsuMaps
			}
		}

		if presetCanaryJiangxi, ok := envConfArg["presetCanaryJiangxi"]; ok {
			presetCanaryJiangxiMaps := make([]map[string]interface{}, 0)
			presetCanaryJiangxiArg := presetCanaryJiangxi.(map[string]interface{})

			if len(presetCanaryJiangxiArg) > 0 {
				presetCanaryJiangxiMap := map[string]interface{}{}

				if specName, ok := presetCanaryJiangxiArg["SpecName"]; ok {
					presetCanaryJiangxiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryJiangxiArg["CodeRev"]; ok {
					presetCanaryJiangxiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryJiangxiArg["AllowedHosts"]; ok {
					presetCanaryJiangxiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryJiangxiMaps = append(presetCanaryJiangxiMaps, presetCanaryJiangxiMap)
				envConfMap["preset_canary_jiangxi"] = presetCanaryJiangxiMaps
			}
		}

		if presetCanaryJilin, ok := envConfArg["presetCanaryJilin"]; ok {
			presetCanaryJilinMaps := make([]map[string]interface{}, 0)
			presetCanaryJilinArg := presetCanaryJilin.(map[string]interface{})

			if len(presetCanaryJilinArg) > 0 {
				presetCanaryJilinMap := map[string]interface{}{}

				if specName, ok := presetCanaryJilinArg["SpecName"]; ok {
					presetCanaryJilinMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryJilinArg["CodeRev"]; ok {
					presetCanaryJilinMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryJilinArg["AllowedHosts"]; ok {
					presetCanaryJilinMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryJilinMaps = append(presetCanaryJilinMaps, presetCanaryJilinMap)
				envConfMap["preset_canary_jilin"] = presetCanaryJilinMaps
			}
		}

		if presetCanaryLiaoning, ok := envConfArg["presetCanaryLiaoning"]; ok {
			presetCanaryLiaoningMaps := make([]map[string]interface{}, 0)
			presetCanaryLiaoningArg := presetCanaryLiaoning.(map[string]interface{})

			if len(presetCanaryLiaoningArg) > 0 {
				presetCanaryLiaoningMap := map[string]interface{}{}

				if specName, ok := presetCanaryLiaoningArg["SpecName"]; ok {
					presetCanaryLiaoningMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryLiaoningArg["CodeRev"]; ok {
					presetCanaryLiaoningMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryLiaoningArg["AllowedHosts"]; ok {
					presetCanaryLiaoningMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryLiaoningMaps = append(presetCanaryLiaoningMaps, presetCanaryLiaoningMap)
				envConfMap["preset_canary_liaoning"] = presetCanaryLiaoningMaps
			}
		}

		if presetCanaryMacau, ok := envConfArg["presetCanaryMacau"]; ok {
			presetCanaryMacauMaps := make([]map[string]interface{}, 0)
			presetCanaryMacauArg := presetCanaryMacau.(map[string]interface{})

			if len(presetCanaryMacauArg) > 0 {
				presetCanaryMacauMap := map[string]interface{}{}

				if specName, ok := presetCanaryMacauArg["SpecName"]; ok {
					presetCanaryMacauMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryMacauArg["CodeRev"]; ok {
					presetCanaryMacauMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryMacauArg["AllowedHosts"]; ok {
					presetCanaryMacauMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryMacauMaps = append(presetCanaryMacauMaps, presetCanaryMacauMap)
				envConfMap["preset_canary_macau"] = presetCanaryMacauMaps
			}
		}

		if presetCanaryNeimenggu, ok := envConfArg["presetCanaryNeimenggu"]; ok {
			presetCanaryNeimengguMaps := make([]map[string]interface{}, 0)
			presetCanaryNeimengguArg := presetCanaryNeimenggu.(map[string]interface{})

			if len(presetCanaryNeimengguArg) > 0 {
				presetCanaryNeimengguMap := map[string]interface{}{}

				if specName, ok := presetCanaryNeimengguArg["SpecName"]; ok {
					presetCanaryNeimengguMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryNeimengguArg["CodeRev"]; ok {
					presetCanaryNeimengguMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryNeimengguArg["AllowedHosts"]; ok {
					presetCanaryNeimengguMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryNeimengguMaps = append(presetCanaryNeimengguMaps, presetCanaryNeimengguMap)
				envConfMap["preset_canary_neimenggu"] = presetCanaryNeimengguMaps
			}
		}

		if presetCanaryNingxia, ok := envConfArg["presetCanaryNingxia"]; ok {
			presetCanaryNingxiaMaps := make([]map[string]interface{}, 0)
			presetCanaryNingxiaArg := presetCanaryNingxia.(map[string]interface{})

			if len(presetCanaryNingxiaArg) > 0 {
				presetCanaryNingxiaMap := map[string]interface{}{}

				if specName, ok := presetCanaryNingxiaArg["SpecName"]; ok {
					presetCanaryNingxiaMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryNingxiaArg["CodeRev"]; ok {
					presetCanaryNingxiaMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryNingxiaArg["AllowedHosts"]; ok {
					presetCanaryNingxiaMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryNingxiaMaps = append(presetCanaryNingxiaMaps, presetCanaryNingxiaMap)
				envConfMap["preset_canary_ningxia"] = presetCanaryNingxiaMaps
			}
		}

		if presetCanaryQinghai, ok := envConfArg["presetCanaryQinghai"]; ok {
			presetCanaryQinghaiMaps := make([]map[string]interface{}, 0)
			presetCanaryQinghaiArg := presetCanaryQinghai.(map[string]interface{})

			if len(presetCanaryQinghaiArg) > 0 {
				presetCanaryQinghaiMap := map[string]interface{}{}

				if specName, ok := presetCanaryQinghaiArg["SpecName"]; ok {
					presetCanaryQinghaiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryQinghaiArg["CodeRev"]; ok {
					presetCanaryQinghaiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryQinghaiArg["AllowedHosts"]; ok {
					presetCanaryQinghaiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryQinghaiMaps = append(presetCanaryQinghaiMaps, presetCanaryQinghaiMap)
				envConfMap["preset_canary_qinghai"] = presetCanaryQinghaiMaps
			}
		}

		if presetCanaryShaanxi, ok := envConfArg["presetCanaryShaanxi"]; ok {
			presetCanaryShaanxiMaps := make([]map[string]interface{}, 0)
			presetCanaryShaanxiArg := presetCanaryShaanxi.(map[string]interface{})

			if len(presetCanaryShaanxiArg) > 0 {
				presetCanaryShaanxiMap := map[string]interface{}{}

				if specName, ok := presetCanaryShaanxiArg["SpecName"]; ok {
					presetCanaryShaanxiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryShaanxiArg["CodeRev"]; ok {
					presetCanaryShaanxiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryShaanxiArg["AllowedHosts"]; ok {
					presetCanaryShaanxiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryShaanxiMaps = append(presetCanaryShaanxiMaps, presetCanaryShaanxiMap)
				envConfMap["preset_canary_shaanxi"] = presetCanaryShaanxiMaps
			}
		}

		if presetCanaryShandong, ok := envConfArg["presetCanaryShandong"]; ok {
			presetCanaryShandongMaps := make([]map[string]interface{}, 0)
			presetCanaryShandongArg := presetCanaryShandong.(map[string]interface{})

			if len(presetCanaryShandongArg) > 0 {
				presetCanaryShandongMap := map[string]interface{}{}

				if specName, ok := presetCanaryShandongArg["SpecName"]; ok {
					presetCanaryShandongMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryShandongArg["CodeRev"]; ok {
					presetCanaryShandongMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryShandongArg["AllowedHosts"]; ok {
					presetCanaryShandongMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryShandongMaps = append(presetCanaryShandongMaps, presetCanaryShandongMap)
				envConfMap["preset_canary_shandong"] = presetCanaryShandongMaps
			}
		}

		if presetCanaryShanghai, ok := envConfArg["presetCanaryShanghai"]; ok {
			presetCanaryShanghaiMaps := make([]map[string]interface{}, 0)
			presetCanaryShanghaiArg := presetCanaryShanghai.(map[string]interface{})

			if len(presetCanaryShanghaiArg) > 0 {
				presetCanaryShanghaiMap := map[string]interface{}{}

				if specName, ok := presetCanaryShanghaiArg["SpecName"]; ok {
					presetCanaryShanghaiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryShanghaiArg["CodeRev"]; ok {
					presetCanaryShanghaiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryShanghaiArg["AllowedHosts"]; ok {
					presetCanaryShanghaiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryShanghaiMaps = append(presetCanaryShanghaiMaps, presetCanaryShanghaiMap)
				envConfMap["preset_canary_shanghai"] = presetCanaryShanghaiMaps
			}
		}

		if presetCanaryShanxi, ok := envConfArg["presetCanaryShanxi"]; ok {
			presetCanaryShanxiMaps := make([]map[string]interface{}, 0)
			presetCanaryShanxiArg := presetCanaryShanxi.(map[string]interface{})

			if len(presetCanaryShanxiArg) > 0 {
				presetCanaryShanxiMap := map[string]interface{}{}

				if specName, ok := presetCanaryShanxiArg["SpecName"]; ok {
					presetCanaryShanxiMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryShanxiArg["CodeRev"]; ok {
					presetCanaryShanxiMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryShanxiArg["AllowedHosts"]; ok {
					presetCanaryShanxiMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryShanxiMaps = append(presetCanaryShanxiMaps, presetCanaryShanxiMap)
				envConfMap["preset_canary_shanxi"] = presetCanaryShanxiMaps
			}
		}

		if presetCanarySichuan, ok := envConfArg["presetCanarySichuan"]; ok {
			presetCanarySichuanMaps := make([]map[string]interface{}, 0)
			presetCanarySichuanArg := presetCanarySichuan.(map[string]interface{})

			if len(presetCanarySichuanArg) > 0 {
				presetCanarySichuanMap := map[string]interface{}{}

				if specName, ok := presetCanarySichuanArg["SpecName"]; ok {
					presetCanarySichuanMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanarySichuanArg["CodeRev"]; ok {
					presetCanarySichuanMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanarySichuanArg["AllowedHosts"]; ok {
					presetCanarySichuanMap["allowed_hosts"] = allowedHosts
				}

				presetCanarySichuanMaps = append(presetCanarySichuanMaps, presetCanarySichuanMap)
				envConfMap["preset_canary_sichuan"] = presetCanarySichuanMaps
			}
		}

		if presetCanaryTaiwan, ok := envConfArg["presetCanaryTaiwan"]; ok {
			presetCanaryTaiwanMaps := make([]map[string]interface{}, 0)
			presetCanaryTaiwanArg := presetCanaryTaiwan.(map[string]interface{})

			if len(presetCanaryTaiwanArg) > 0 {
				presetCanaryTaiwanMap := map[string]interface{}{}

				if specName, ok := presetCanaryTaiwanArg["SpecName"]; ok {
					presetCanaryTaiwanMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryTaiwanArg["CodeRev"]; ok {
					presetCanaryTaiwanMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryTaiwanArg["AllowedHosts"]; ok {
					presetCanaryTaiwanMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryTaiwanMaps = append(presetCanaryTaiwanMaps, presetCanaryTaiwanMap)
				envConfMap["preset_canary_taiwan"] = presetCanaryTaiwanMaps
			}
		}

		if presetCanaryTianjin, ok := envConfArg["presetCanaryTianjin"]; ok {
			presetCanaryTianjinMaps := make([]map[string]interface{}, 0)
			presetCanaryTianjinArg := presetCanaryTianjin.(map[string]interface{})

			if len(presetCanaryTianjinArg) > 0 {
				presetCanaryTianjinMap := map[string]interface{}{}

				if specName, ok := presetCanaryTianjinArg["SpecName"]; ok {
					presetCanaryTianjinMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryTianjinArg["CodeRev"]; ok {
					presetCanaryTianjinMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryTianjinArg["AllowedHosts"]; ok {
					presetCanaryTianjinMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryTianjinMaps = append(presetCanaryTianjinMaps, presetCanaryTianjinMap)
				envConfMap["preset_canary_tianjin"] = presetCanaryTianjinMaps
			}
		}

		if presetCanaryXinjiang, ok := envConfArg["presetCanaryXinjiang"]; ok {
			presetCanaryXinjiangMaps := make([]map[string]interface{}, 0)
			presetCanaryXinjiangArg := presetCanaryXinjiang.(map[string]interface{})

			if len(presetCanaryXinjiangArg) > 0 {
				presetCanaryXinjiangMap := map[string]interface{}{}

				if specName, ok := presetCanaryXinjiangArg["SpecName"]; ok {
					presetCanaryXinjiangMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryXinjiangArg["CodeRev"]; ok {
					presetCanaryXinjiangMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryXinjiangArg["AllowedHosts"]; ok {
					presetCanaryXinjiangMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryXinjiangMaps = append(presetCanaryXinjiangMaps, presetCanaryXinjiangMap)
				envConfMap["preset_canary_xinjiang"] = presetCanaryXinjiangMaps
			}
		}

		if presetCanaryXizang, ok := envConfArg["presetCanaryXizang"]; ok {
			presetCanaryXizangMaps := make([]map[string]interface{}, 0)
			presetCanaryXizangArg := presetCanaryXizang.(map[string]interface{})

			if len(presetCanaryXizangArg) > 0 {
				presetCanaryXizangMap := map[string]interface{}{}

				if specName, ok := presetCanaryXizangArg["SpecName"]; ok {
					presetCanaryXizangMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryXizangArg["CodeRev"]; ok {
					presetCanaryXizangMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryXizangArg["AllowedHosts"]; ok {
					presetCanaryXizangMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryXizangMaps = append(presetCanaryXizangMaps, presetCanaryXizangMap)
				envConfMap["preset_canary_xizang"] = presetCanaryXizangMaps
			}
		}

		if presetCanaryYunnan, ok := envConfArg["presetCanaryYunnan"]; ok {
			presetCanaryYunnanMaps := make([]map[string]interface{}, 0)
			presetCanaryYunnanArg := presetCanaryYunnan.(map[string]interface{})

			if len(presetCanaryYunnanArg) > 0 {
				presetCanaryYunnanMap := map[string]interface{}{}

				if specName, ok := presetCanaryYunnanArg["SpecName"]; ok {
					presetCanaryYunnanMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryYunnanArg["CodeRev"]; ok {
					presetCanaryYunnanMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryYunnanArg["AllowedHosts"]; ok {
					presetCanaryYunnanMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryYunnanMaps = append(presetCanaryYunnanMaps, presetCanaryYunnanMap)
				envConfMap["preset_canary_yunnan"] = presetCanaryYunnanMaps
			}
		}

		if presetCanaryZhejiang, ok := envConfArg["presetCanaryZhejiang"]; ok {
			presetCanaryZhejiangMaps := make([]map[string]interface{}, 0)
			presetCanaryZhejiangArg := presetCanaryZhejiang.(map[string]interface{})

			if len(presetCanaryZhejiangArg) > 0 {
				presetCanaryZhejiangMap := map[string]interface{}{}

				if specName, ok := presetCanaryZhejiangArg["SpecName"]; ok {
					presetCanaryZhejiangMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryZhejiangArg["CodeRev"]; ok {
					presetCanaryZhejiangMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryZhejiangArg["AllowedHosts"]; ok {
					presetCanaryZhejiangMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryZhejiangMaps = append(presetCanaryZhejiangMaps, presetCanaryZhejiangMap)
				envConfMap["preset_canary_zhejiang"] = presetCanaryZhejiangMaps
			}
		}

		if presetCanaryOverseas, ok := envConfArg["presetCanaryOverseas"]; ok {
			presetCanaryOverseasMaps := make([]map[string]interface{}, 0)
			presetCanaryOverseasArg := presetCanaryOverseas.(map[string]interface{})

			if len(presetCanaryOverseasArg) > 0 {
				presetCanaryOverseasMap := map[string]interface{}{}

				if specName, ok := presetCanaryOverseasArg["SpecName"]; ok {
					presetCanaryOverseasMap["spec_name"] = specName
				}

				if codeRev, ok := presetCanaryOverseasArg["CodeRev"]; ok {
					presetCanaryOverseasMap["code_rev"] = codeRev
				}

				if allowedHosts, ok := presetCanaryOverseasArg["AllowedHosts"]; ok {
					presetCanaryOverseasMap["allowed_hosts"] = allowedHosts
				}

				presetCanaryOverseasMaps = append(presetCanaryOverseasMaps, presetCanaryOverseasMap)
				envConfMap["preset_canary_overseas"] = presetCanaryOverseasMaps
			}
		}

		envConfMaps = append(envConfMaps, envConfMap)

		d.Set("env_conf", envConfMaps)
	}

	return nil
}

func resourceAliCloudDcdnErUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"Name": d.Id(),
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if d.HasChange("env_conf") {
		update = true
	}
	if v, ok := d.GetOk("env_conf"); ok {
		envConfMap := map[string]interface{}{}
		for _, envConfList := range v.([]interface{}) {
			envConfArg := envConfList.(map[string]interface{})

			if staging, ok := envConfArg["staging"]; ok {
				stagingMap := map[string]interface{}{}
				for _, stagingList := range staging.([]interface{}) {
					stagingArg := stagingList.(map[string]interface{})

					if specName, ok := stagingArg["spec_name"]; ok {
						stagingMap["SpecName"] = specName
					}

					if codeRev, ok := stagingArg["code_rev"]; ok {
						stagingMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := stagingArg["allowed_hosts"]; ok {
						stagingMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(stagingMap) > 0 {
					envConfMap["staging"] = stagingMap
				}
			}

			if production, ok := envConfArg["production"]; ok {
				productionMap := map[string]interface{}{}
				for _, productionList := range production.([]interface{}) {
					productionArg := productionList.(map[string]interface{})

					if specName, ok := productionArg["spec_name"]; ok {
						productionMap["SpecName"] = specName
					}

					if codeRev, ok := productionArg["code_rev"]; ok {
						productionMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := productionArg["allowed_hosts"]; ok {
						productionMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(productionMap) > 0 {
					envConfMap["production"] = productionMap
				}
			}

			if presetCanaryAnhui, ok := envConfArg["preset_canary_anhui"]; ok {
				presetCanaryAnhuiMap := map[string]interface{}{}
				for _, presetCanaryAnhuiList := range presetCanaryAnhui.([]interface{}) {
					presetCanaryAnhuiArg := presetCanaryAnhuiList.(map[string]interface{})

					if specName, ok := presetCanaryAnhuiArg["spec_name"]; ok {
						presetCanaryAnhuiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryAnhuiArg["code_rev"]; ok {
						presetCanaryAnhuiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryAnhuiArg["allowed_hosts"]; ok {
						presetCanaryAnhuiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryAnhuiMap) > 0 {
					envConfMap["presetCanaryAnhui"] = presetCanaryAnhuiMap
				}
			}

			if presetCanaryBeijing, ok := envConfArg["preset_canary_beijing"]; ok {
				presetCanaryBeijingMap := map[string]interface{}{}
				for _, presetCanaryBeijingList := range presetCanaryBeijing.([]interface{}) {
					presetCanaryBeijingArg := presetCanaryBeijingList.(map[string]interface{})

					if specName, ok := presetCanaryBeijingArg["spec_name"]; ok {
						presetCanaryBeijingMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryBeijingArg["code_rev"]; ok {
						presetCanaryBeijingMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryBeijingArg["allowed_hosts"]; ok {
						presetCanaryBeijingMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryBeijingMap) > 0 {
					envConfMap["presetCanaryBeijing"] = presetCanaryBeijingMap
				}
			}

			if presetCanaryChongqing, ok := envConfArg["preset_canary_chongqing"]; ok {
				presetCanaryChongqingMap := map[string]interface{}{}
				for _, presetCanaryChongqingList := range presetCanaryChongqing.([]interface{}) {
					presetCanaryChongqingArg := presetCanaryChongqingList.(map[string]interface{})

					if specName, ok := presetCanaryChongqingArg["spec_name"]; ok {
						presetCanaryChongqingMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryChongqingArg["code_rev"]; ok {
						presetCanaryChongqingMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryChongqingArg["allowed_hosts"]; ok {
						presetCanaryChongqingMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryChongqingMap) > 0 {
					envConfMap["presetCanaryChongqing"] = presetCanaryChongqingMap
				}
			}

			if presetCanaryFujian, ok := envConfArg["preset_canary_fujian"]; ok {
				presetCanaryFujianMap := map[string]interface{}{}
				for _, presetCanaryFujianList := range presetCanaryFujian.([]interface{}) {
					presetCanaryFujianArg := presetCanaryFujianList.(map[string]interface{})

					if specName, ok := presetCanaryFujianArg["spec_name"]; ok {
						presetCanaryFujianMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryFujianArg["code_rev"]; ok {
						presetCanaryFujianMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryFujianArg["allowed_hosts"]; ok {
						presetCanaryFujianMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryFujianMap) > 0 {
					envConfMap["presetCanaryFujian"] = presetCanaryFujianMap
				}
			}

			if presetCanaryGansu, ok := envConfArg["preset_canary_gansu"]; ok {
				presetCanaryGansuMap := map[string]interface{}{}
				for _, presetCanaryGansuList := range presetCanaryGansu.([]interface{}) {
					presetCanaryGansuArg := presetCanaryGansuList.(map[string]interface{})

					if specName, ok := presetCanaryGansuArg["spec_name"]; ok {
						presetCanaryGansuMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGansuArg["code_rev"]; ok {
						presetCanaryGansuMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGansuArg["allowed_hosts"]; ok {
						presetCanaryGansuMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGansuMap) > 0 {
					envConfMap["presetCanaryGansu"] = presetCanaryGansuMap
				}
			}

			if presetCanaryGuangdong, ok := envConfArg["preset_canary_guangdong"]; ok {
				presetCanaryGuangdongMap := map[string]interface{}{}
				for _, presetCanaryGuangdongList := range presetCanaryGuangdong.([]interface{}) {
					presetCanaryGuangdongArg := presetCanaryGuangdongList.(map[string]interface{})

					if specName, ok := presetCanaryGuangdongArg["spec_name"]; ok {
						presetCanaryGuangdongMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGuangdongArg["code_rev"]; ok {
						presetCanaryGuangdongMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGuangdongArg["allowed_hosts"]; ok {
						presetCanaryGuangdongMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGuangdongMap) > 0 {
					envConfMap["presetCanaryGuangdong"] = presetCanaryGuangdongMap
				}
			}

			if presetCanaryGuangxi, ok := envConfArg["preset_canary_guangxi"]; ok {
				presetCanaryGuangxiMap := map[string]interface{}{}
				for _, presetCanaryGuangxiList := range presetCanaryGuangxi.([]interface{}) {
					presetCanaryGuangxiArg := presetCanaryGuangxiList.(map[string]interface{})

					if specName, ok := presetCanaryGuangxiArg["spec_name"]; ok {
						presetCanaryGuangxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGuangxiArg["code_rev"]; ok {
						presetCanaryGuangxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGuangxiArg["allowed_hosts"]; ok {
						presetCanaryGuangxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGuangxiMap) > 0 {
					envConfMap["presetCanaryGuangxi"] = presetCanaryGuangxiMap
				}
			}

			if presetCanaryGuizhou, ok := envConfArg["preset_canary_guizhou"]; ok {
				presetCanaryGuizhouMap := map[string]interface{}{}
				for _, presetCanaryGuizhouList := range presetCanaryGuizhou.([]interface{}) {
					presetCanaryGuizhouArg := presetCanaryGuizhouList.(map[string]interface{})

					if specName, ok := presetCanaryGuizhouArg["spec_name"]; ok {
						presetCanaryGuizhouMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryGuizhouArg["code_rev"]; ok {
						presetCanaryGuizhouMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryGuizhouArg["allowed_hosts"]; ok {
						presetCanaryGuizhouMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryGuizhouMap) > 0 {
					envConfMap["presetCanaryGuizhou"] = presetCanaryGuizhouMap
				}
			}

			if presetCanaryHainan, ok := envConfArg["preset_canary_hainan"]; ok {
				presetCanaryHainanMap := map[string]interface{}{}
				for _, presetCanaryHainanList := range presetCanaryHainan.([]interface{}) {
					presetCanaryHainanArg := presetCanaryHainanList.(map[string]interface{})

					if specName, ok := presetCanaryHainanArg["spec_name"]; ok {
						presetCanaryHainanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHainanArg["code_rev"]; ok {
						presetCanaryHainanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHainanArg["allowed_hosts"]; ok {
						presetCanaryHainanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHainanMap) > 0 {
					envConfMap["presetCanaryHainan"] = presetCanaryHainanMap
				}
			}

			if presetCanaryHebei, ok := envConfArg["preset_canary_hebei"]; ok {
				presetCanaryHebeiMap := map[string]interface{}{}
				for _, presetCanaryHebeiList := range presetCanaryHebei.([]interface{}) {
					presetCanaryHebeiArg := presetCanaryHebeiList.(map[string]interface{})

					if specName, ok := presetCanaryHebeiArg["spec_name"]; ok {
						presetCanaryHebeiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHebeiArg["code_rev"]; ok {
						presetCanaryHebeiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHebeiArg["allowed_hosts"]; ok {
						presetCanaryHebeiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHebeiMap) > 0 {
					envConfMap["presetCanaryHebei"] = presetCanaryHebeiMap
				}
			}

			if presetCanaryHeilongjiang, ok := envConfArg["preset_canary_heilongjiang"]; ok {
				presetCanaryHeilongjiangMap := map[string]interface{}{}
				for _, presetCanaryHeilongjiangList := range presetCanaryHeilongjiang.([]interface{}) {
					presetCanaryHeilongjiangArg := presetCanaryHeilongjiangList.(map[string]interface{})

					if specName, ok := presetCanaryHeilongjiangArg["spec_name"]; ok {
						presetCanaryHeilongjiangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHeilongjiangArg["code_rev"]; ok {
						presetCanaryHeilongjiangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHeilongjiangArg["allowed_hosts"]; ok {
						presetCanaryHeilongjiangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHeilongjiangMap) > 0 {
					envConfMap["presetCanaryHeilongjiang"] = presetCanaryHeilongjiangMap
				}
			}

			if presetCanaryHenan, ok := envConfArg["preset_canary_henan"]; ok {
				presetCanaryHenanMap := map[string]interface{}{}
				for _, presetCanaryHenanList := range presetCanaryHenan.([]interface{}) {
					presetCanaryHenanArg := presetCanaryHenanList.(map[string]interface{})

					if specName, ok := presetCanaryHenanArg["spec_name"]; ok {
						presetCanaryHenanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHenanArg["code_rev"]; ok {
						presetCanaryHenanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHenanArg["allowed_hosts"]; ok {
						presetCanaryHenanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHenanMap) > 0 {
					envConfMap["presetCanaryHenan"] = presetCanaryHenanMap
				}
			}

			if presetCanaryHongKong, ok := envConfArg["preset_canary_hong_kong"]; ok {
				presetCanaryHongKongMap := map[string]interface{}{}
				for _, presetCanaryHongKongList := range presetCanaryHongKong.([]interface{}) {
					presetCanaryHongKongArg := presetCanaryHongKongList.(map[string]interface{})

					if specName, ok := presetCanaryHongKongArg["spec_name"]; ok {
						presetCanaryHongKongMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHongKongArg["code_rev"]; ok {
						presetCanaryHongKongMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHongKongArg["allowed_hosts"]; ok {
						presetCanaryHongKongMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHongKongMap) > 0 {
					envConfMap["presetCanaryHongKong"] = presetCanaryHongKongMap
				}
			}

			if presetCanaryHubei, ok := envConfArg["preset_canary_hubei"]; ok {
				presetCanaryHubeiMap := map[string]interface{}{}
				for _, presetCanaryHubeiList := range presetCanaryHubei.([]interface{}) {
					presetCanaryHubeiArg := presetCanaryHubeiList.(map[string]interface{})

					if specName, ok := presetCanaryHubeiArg["spec_name"]; ok {
						presetCanaryHubeiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHubeiArg["code_rev"]; ok {
						presetCanaryHubeiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHubeiArg["allowed_hosts"]; ok {
						presetCanaryHubeiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHubeiMap) > 0 {
					envConfMap["presetCanaryHubei"] = presetCanaryHubeiMap
				}
			}

			if presetCanaryHunan, ok := envConfArg["preset_canary_hunan"]; ok {
				presetCanaryHunanMap := map[string]interface{}{}
				for _, presetCanaryHunanList := range presetCanaryHunan.([]interface{}) {
					presetCanaryHunanArg := presetCanaryHunanList.(map[string]interface{})

					if specName, ok := presetCanaryHunanArg["spec_name"]; ok {
						presetCanaryHunanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryHunanArg["code_rev"]; ok {
						presetCanaryHunanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryHunanArg["allowed_hosts"]; ok {
						presetCanaryHunanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryHunanMap) > 0 {
					envConfMap["presetCanaryHunan"] = presetCanaryHunanMap
				}
			}

			if presetCanaryJiangsu, ok := envConfArg["preset_canary_jiangsu"]; ok {
				presetCanaryJiangsuMap := map[string]interface{}{}
				for _, presetCanaryJiangsuList := range presetCanaryJiangsu.([]interface{}) {
					presetCanaryJiangsuArg := presetCanaryJiangsuList.(map[string]interface{})

					if specName, ok := presetCanaryJiangsuArg["spec_name"]; ok {
						presetCanaryJiangsuMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryJiangsuArg["code_rev"]; ok {
						presetCanaryJiangsuMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryJiangsuArg["allowed_hosts"]; ok {
						presetCanaryJiangsuMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryJiangsuMap) > 0 {
					envConfMap["presetCanaryJiangsu"] = presetCanaryJiangsuMap
				}
			}

			if presetCanaryJiangxi, ok := envConfArg["preset_canary_jiangxi"]; ok {
				presetCanaryJiangxiMap := map[string]interface{}{}
				for _, presetCanaryJiangxiList := range presetCanaryJiangxi.([]interface{}) {
					presetCanaryJiangxiArg := presetCanaryJiangxiList.(map[string]interface{})

					if specName, ok := presetCanaryJiangxiArg["spec_name"]; ok {
						presetCanaryJiangxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryJiangxiArg["code_rev"]; ok {
						presetCanaryJiangxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryJiangxiArg["allowed_hosts"]; ok {
						presetCanaryJiangxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryJiangxiMap) > 0 {
					envConfMap["presetCanaryJiangxi"] = presetCanaryJiangxiMap
				}
			}

			if presetCanaryJilin, ok := envConfArg["preset_canary_jilin"]; ok {
				presetCanaryJilinMap := map[string]interface{}{}
				for _, presetCanaryJilinList := range presetCanaryJilin.([]interface{}) {
					presetCanaryJilinArg := presetCanaryJilinList.(map[string]interface{})

					if specName, ok := presetCanaryJilinArg["spec_name"]; ok {
						presetCanaryJilinMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryJilinArg["code_rev"]; ok {
						presetCanaryJilinMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryJilinArg["allowed_hosts"]; ok {
						presetCanaryJilinMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryJilinMap) > 0 {
					envConfMap["presetCanaryJilin"] = presetCanaryJilinMap
				}
			}

			if presetCanaryLiaoning, ok := envConfArg["preset_canary_liaoning"]; ok {
				presetCanaryLiaoningMap := map[string]interface{}{}
				for _, presetCanaryLiaoningList := range presetCanaryLiaoning.([]interface{}) {
					presetCanaryLiaoningArg := presetCanaryLiaoningList.(map[string]interface{})

					if specName, ok := presetCanaryLiaoningArg["spec_name"]; ok {
						presetCanaryLiaoningMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryLiaoningArg["code_rev"]; ok {
						presetCanaryLiaoningMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryLiaoningArg["allowed_hosts"]; ok {
						presetCanaryLiaoningMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryLiaoningMap) > 0 {
					envConfMap["presetCanaryLiaoning"] = presetCanaryLiaoningMap
				}
			}

			if presetCanaryMacau, ok := envConfArg["preset_canary_macau"]; ok {
				presetCanaryMacauMap := map[string]interface{}{}
				for _, presetCanaryMacauList := range presetCanaryMacau.([]interface{}) {
					presetCanaryMacauArg := presetCanaryMacauList.(map[string]interface{})

					if specName, ok := presetCanaryMacauArg["spec_name"]; ok {
						presetCanaryMacauMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryMacauArg["code_rev"]; ok {
						presetCanaryMacauMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryMacauArg["allowed_hosts"]; ok {
						presetCanaryMacauMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryMacauMap) > 0 {
					envConfMap["presetCanaryMacau"] = presetCanaryMacauMap
				}
			}

			if presetCanaryNeimenggu, ok := envConfArg["preset_canary_neimenggu"]; ok {
				presetCanaryNeimengguMap := map[string]interface{}{}
				for _, presetCanaryNeimengguList := range presetCanaryNeimenggu.([]interface{}) {
					presetCanaryNeimengguArg := presetCanaryNeimengguList.(map[string]interface{})

					if specName, ok := presetCanaryNeimengguArg["spec_name"]; ok {
						presetCanaryNeimengguMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryNeimengguArg["code_rev"]; ok {
						presetCanaryNeimengguMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryNeimengguArg["allowed_hosts"]; ok {
						presetCanaryNeimengguMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryNeimengguMap) > 0 {
					envConfMap["presetCanaryNeimenggu"] = presetCanaryNeimengguMap
				}
			}

			if presetCanaryNingxia, ok := envConfArg["preset_canary_ningxia"]; ok {
				presetCanaryNingxiaMap := map[string]interface{}{}
				for _, presetCanaryNingxiaList := range presetCanaryNingxia.([]interface{}) {
					presetCanaryNingxiaArg := presetCanaryNingxiaList.(map[string]interface{})

					if specName, ok := presetCanaryNingxiaArg["spec_name"]; ok {
						presetCanaryNingxiaMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryNingxiaArg["code_rev"]; ok {
						presetCanaryNingxiaMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryNingxiaArg["allowed_hosts"]; ok {
						presetCanaryNingxiaMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryNingxiaMap) > 0 {
					envConfMap["presetCanaryNingxia"] = presetCanaryNingxiaMap
				}
			}

			if presetCanaryQinghai, ok := envConfArg["preset_canary_qinghai"]; ok {
				presetCanaryQinghaiMap := map[string]interface{}{}
				for _, presetCanaryQinghaiList := range presetCanaryQinghai.([]interface{}) {
					presetCanaryQinghaiArg := presetCanaryQinghaiList.(map[string]interface{})

					if specName, ok := presetCanaryQinghaiArg["spec_name"]; ok {
						presetCanaryQinghaiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryQinghaiArg["code_rev"]; ok {
						presetCanaryQinghaiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryQinghaiArg["allowed_hosts"]; ok {
						presetCanaryQinghaiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryQinghaiMap) > 0 {
					envConfMap["presetCanaryQinghai"] = presetCanaryQinghaiMap
				}
			}

			if presetCanaryShaanxi, ok := envConfArg["preset_canary_shaanxi"]; ok {
				presetCanaryShaanxiMap := map[string]interface{}{}
				for _, presetCanaryShaanxiList := range presetCanaryShaanxi.([]interface{}) {
					presetCanaryShaanxiArg := presetCanaryShaanxiList.(map[string]interface{})

					if specName, ok := presetCanaryShaanxiArg["spec_name"]; ok {
						presetCanaryShaanxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShaanxiArg["code_rev"]; ok {
						presetCanaryShaanxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShaanxiArg["allowed_hosts"]; ok {
						presetCanaryShaanxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShaanxiMap) > 0 {
					envConfMap["presetCanaryShaanxi"] = presetCanaryShaanxiMap
				}
			}

			if presetCanaryShandong, ok := envConfArg["preset_canary_shandong"]; ok {
				presetCanaryShandongMap := map[string]interface{}{}
				for _, presetCanaryShandongList := range presetCanaryShandong.([]interface{}) {
					presetCanaryShandongArg := presetCanaryShandongList.(map[string]interface{})

					if specName, ok := presetCanaryShandongArg["spec_name"]; ok {
						presetCanaryShandongMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShandongArg["code_rev"]; ok {
						presetCanaryShandongMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShandongArg["allowed_hosts"]; ok {
						presetCanaryShandongMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShandongMap) > 0 {
					envConfMap["presetCanaryShandong"] = presetCanaryShandongMap
				}
			}

			if presetCanaryShanghai, ok := envConfArg["preset_canary_shanghai"]; ok {
				presetCanaryShanghaiMap := map[string]interface{}{}
				for _, presetCanaryShanghaiList := range presetCanaryShanghai.([]interface{}) {
					presetCanaryShanghaiArg := presetCanaryShanghaiList.(map[string]interface{})

					if specName, ok := presetCanaryShanghaiArg["spec_name"]; ok {
						presetCanaryShanghaiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShanghaiArg["code_rev"]; ok {
						presetCanaryShanghaiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShanghaiArg["allowed_hosts"]; ok {
						presetCanaryShanghaiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShanghaiMap) > 0 {
					envConfMap["presetCanaryShanghai"] = presetCanaryShanghaiMap
				}
			}

			if presetCanaryShanxi, ok := envConfArg["preset_canary_shanxi"]; ok {
				presetCanaryShanxiMap := map[string]interface{}{}
				for _, presetCanaryShanxiList := range presetCanaryShanxi.([]interface{}) {
					presetCanaryShanxiArg := presetCanaryShanxiList.(map[string]interface{})

					if specName, ok := presetCanaryShanxiArg["spec_name"]; ok {
						presetCanaryShanxiMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryShanxiArg["code_rev"]; ok {
						presetCanaryShanxiMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryShanxiArg["allowed_hosts"]; ok {
						presetCanaryShanxiMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryShanxiMap) > 0 {
					envConfMap["presetCanaryShanxi"] = presetCanaryShanxiMap
				}
			}

			if presetCanarySichuan, ok := envConfArg["preset_canary_sichuan"]; ok {
				presetCanarySichuanMap := map[string]interface{}{}
				for _, presetCanarySichuanList := range presetCanarySichuan.([]interface{}) {
					presetCanarySichuanArg := presetCanarySichuanList.(map[string]interface{})

					if specName, ok := presetCanarySichuanArg["spec_name"]; ok {
						presetCanarySichuanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanarySichuanArg["code_rev"]; ok {
						presetCanarySichuanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanarySichuanArg["allowed_hosts"]; ok {
						presetCanarySichuanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanarySichuanMap) > 0 {
					envConfMap["presetCanarySichuan"] = presetCanarySichuanMap
				}
			}

			if presetCanaryTaiwan, ok := envConfArg["preset_canary_taiwan"]; ok {
				presetCanaryTaiwanMap := map[string]interface{}{}
				for _, presetCanaryTaiwanList := range presetCanaryTaiwan.([]interface{}) {
					presetCanaryTaiwanArg := presetCanaryTaiwanList.(map[string]interface{})

					if specName, ok := presetCanaryTaiwanArg["spec_name"]; ok {
						presetCanaryTaiwanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryTaiwanArg["code_rev"]; ok {
						presetCanaryTaiwanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryTaiwanArg["allowed_hosts"]; ok {
						presetCanaryTaiwanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryTaiwanMap) > 0 {
					envConfMap["presetCanaryTaiwan"] = presetCanaryTaiwanMap
				}
			}

			if presetCanaryTianjin, ok := envConfArg["preset_canary_tianjin"]; ok {
				presetCanaryTianjinMap := map[string]interface{}{}
				for _, presetCanaryTianjinList := range presetCanaryTianjin.([]interface{}) {
					presetCanaryTianjinArg := presetCanaryTianjinList.(map[string]interface{})

					if specName, ok := presetCanaryTianjinArg["spec_name"]; ok {
						presetCanaryTianjinMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryTianjinArg["code_rev"]; ok {
						presetCanaryTianjinMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryTianjinArg["allowed_hosts"]; ok {
						presetCanaryTianjinMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryTianjinMap) > 0 {
					envConfMap["presetCanaryTianjin"] = presetCanaryTianjinMap
				}
			}

			if presetCanaryXinjiang, ok := envConfArg["preset_canary_xinjiang"]; ok {
				presetCanaryXinjiangMap := map[string]interface{}{}
				for _, presetCanaryXinjiangList := range presetCanaryXinjiang.([]interface{}) {
					presetCanaryXinjiangArg := presetCanaryXinjiangList.(map[string]interface{})

					if specName, ok := presetCanaryXinjiangArg["spec_name"]; ok {
						presetCanaryXinjiangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryXinjiangArg["code_rev"]; ok {
						presetCanaryXinjiangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryXinjiangArg["allowed_hosts"]; ok {
						presetCanaryXinjiangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryXinjiangMap) > 0 {
					envConfMap["presetCanaryXinjiang"] = presetCanaryXinjiangMap
				}
			}

			if presetCanaryXizang, ok := envConfArg["preset_canary_xizang"]; ok {
				presetCanaryXizangMap := map[string]interface{}{}
				for _, presetCanaryXizangList := range presetCanaryXizang.([]interface{}) {
					presetCanaryXizangArg := presetCanaryXizangList.(map[string]interface{})

					if specName, ok := presetCanaryXizangArg["spec_name"]; ok {
						presetCanaryXizangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryXizangArg["code_rev"]; ok {
						presetCanaryXizangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryXizangArg["allowed_hosts"]; ok {
						presetCanaryXizangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryXizangMap) > 0 {
					envConfMap["presetCanaryXizang"] = presetCanaryXizangMap
				}
			}

			if presetCanaryYunnan, ok := envConfArg["preset_canary_yunnan"]; ok {
				presetCanaryYunnanMap := map[string]interface{}{}
				for _, presetCanaryYunnanList := range presetCanaryYunnan.([]interface{}) {
					presetCanaryYunnanArg := presetCanaryYunnanList.(map[string]interface{})

					if specName, ok := presetCanaryYunnanArg["spec_name"]; ok {
						presetCanaryYunnanMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryYunnanArg["code_rev"]; ok {
						presetCanaryYunnanMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryYunnanArg["allowed_hosts"]; ok {
						presetCanaryYunnanMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryYunnanMap) > 0 {
					envConfMap["presetCanaryYunnan"] = presetCanaryYunnanMap
				}
			}

			if presetCanaryZhejiang, ok := envConfArg["preset_canary_zhejiang"]; ok {
				presetCanaryZhejiangMap := map[string]interface{}{}
				for _, presetCanaryZhejiangList := range presetCanaryZhejiang.([]interface{}) {
					presetCanaryZhejiangArg := presetCanaryZhejiangList.(map[string]interface{})

					if specName, ok := presetCanaryZhejiangArg["spec_name"]; ok {
						presetCanaryZhejiangMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryZhejiangArg["code_rev"]; ok {
						presetCanaryZhejiangMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryZhejiangArg["allowed_hosts"]; ok {
						presetCanaryZhejiangMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryZhejiangMap) > 0 {
					envConfMap["presetCanaryZhejiang"] = presetCanaryZhejiangMap
				}
			}

			if presetCanaryOverseas, ok := envConfArg["preset_canary_overseas"]; ok {
				presetCanaryOverseasMap := map[string]interface{}{}
				for _, presetCanaryOverseasList := range presetCanaryOverseas.([]interface{}) {
					presetCanaryOverseasArg := presetCanaryOverseasList.(map[string]interface{})

					if specName, ok := presetCanaryOverseasArg["spec_name"]; ok {
						presetCanaryOverseasMap["SpecName"] = specName
					}

					if codeRev, ok := presetCanaryOverseasArg["code_rev"]; ok {
						presetCanaryOverseasMap["CodeRev"] = codeRev
					}

					if allowedHosts, ok := presetCanaryOverseasArg["allowed_hosts"]; ok {
						presetCanaryOverseasMap["AllowedHosts"] = allowedHosts
					}
				}

				if len(presetCanaryOverseasMap) > 0 {
					envConfMap["presetCanaryOverseas"] = presetCanaryOverseasMap
				}
			}
		}

		envConfJson, err := convertMaptoJsonString(envConfMap)
		if err != nil {
			return WrapError(err)
		}

		request["EnvConf"] = envConfJson
	}

	if update {
		action := "EditRoutineConf"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudDcdnErRead(d, meta)
}

func resourceAliCloudDcdnErDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRoutine"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"Name": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, false)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
