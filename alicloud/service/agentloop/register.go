package agentloop

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "AgentLoop",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_agentloop_agent_space", Factory: ResourceAgentSpace},
			{TypeName: "alicloud_agentloop_context_store", Factory: ResourceContextStore},
			{TypeName: "alicloud_agentloop_context_store_api_key", Factory: ResourceContextStoreApiKey},
			{TypeName: "alicloud_agentloop_dataset", Factory: ResourceDataset},
			{TypeName: "alicloud_agentloop_endpoint_connector", Factory: ResourceEndpointConnector},
			{TypeName: "alicloud_agentloop_evaluation_task", Factory: ResourceEvaluationTask},
			{TypeName: "alicloud_agentloop_evaluator", Factory: ResourceEvaluator},
			{TypeName: "alicloud_agentloop_evaluator_skill", Factory: ResourceEvaluatorSkill},
			{TypeName: "alicloud_agentloop_optimize_task", Factory: ResourceOptimizeTask},
			{TypeName: "alicloud_agentloop_pipeline", Factory: ResourcePipeline},
		},
	}
}
