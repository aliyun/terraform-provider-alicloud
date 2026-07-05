package alicloud

// Fixture resource for tier-0 MECHANICAL diff hermetic tests (F3 T0-mech).
// 刻意布置多类 TF↔API 陷阱,供 test/probe_test.sh 断言机械 diff 的 finding vs queue 路由。
// 不对应真实 provider 资源。用 RpcPost 风格,使 (product,version,action) 三元组可抽取。

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudProbeMech() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudProbeMechCreate,
		Update: resourceAliCloudProbeMechUpdate,
		Read:   resourceAliCloudProbeMechRead,
		Schema: map[string]*schema.Schema{
			// 控制项:文档=源码=API 一致,不应产 finding(零误报守门)
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// enum 超集:TF 放行 Standard,API 枚举无 → api_gap_enum_superset
			"storage_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IA",
				ValidateFunc: StringInSlice([]string{"Standard", "IA", "Archive"}, false),
			},
			// required:TF Optional 无 Default,API required=true → api_gap_required
			"required_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// range:TF IntBetween(0,4095),API max 255 → api_gap_range
			"mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 4095),
			},
			// default:TF Default "auto",API default "manual" → api_gap_default
			"mode_value": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "auto",
			},
			// type:TF list vs API string(标量↔集合硬冲突,容差表不含)→ api_gap_type
			"conflict_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// 抑制:client_token→ClientToken 在 suppress_params,虽 API required 也不报,入 suppressed[]
			"client_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// TF 枚举 ⊊ API 枚举(TF 更严)→ coverage note,NOT finding(方向安全)
			"safe_enum": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"a"}, false),
			},
			// 映射不上:renamed_field→RenamedField 在 API 无对应参数 → queue(unmapped),不报
			"renamed_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// 枚举不可解析:StringInSlice 传变量而非字面 slice → enum_status unknown → queue(enum_unparsed),不猜
			"opaque_enum": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice(someEnumList, false),
			},
		},
	}
}

func resourceAliCloudProbeMechCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateProbeMech"
	response, err = client.RpcPost("ProbeMech", "2024-01-01", action, query, request, true)
	_ = response
	_ = err
	return nil
}

func resourceAliCloudProbeMechUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	// LegacyAction 在 OpenAPI 标 deprecated → api_gap_deprecated_action
	action := "LegacyAction"
	response, err = client.RpcPost("ProbeMech", "2024-01-01", action, query, request, true)
	_ = response
	_ = err
	return nil
}

func resourceAliCloudProbeMechRead(d *schema.ResourceData, meta interface{}) error {
	action := "DescribeProbeMech"
	_ = action
	return nil
}
