package paillmtrace

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the PAI LLM Trace
// (PaiLLMTrace) service. New types are Plugin Framework native and attach here.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiLLMTrace",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_pai_llm_trace_evaluation_templates",
				Factory:  NewEvaluationTemplatesDataSource,
			},
		},
	}
}
