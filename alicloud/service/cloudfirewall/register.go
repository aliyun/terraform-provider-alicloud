package cloudfirewall

import (
	conns "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cloud Firewall
// service. It registers the whole migrated product: the SDKv2 resources and
// data sources that previously lived in the legacy root-package schema maps.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "CloudFirewall",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_cloud_firewall_instance_v2", Factory: ResourceCloudFirewallInstanceV2},
			{TypeName: "alicloud_cloud_firewall_instance", Factory: ResourceCloudFirewallInstance},
			{TypeName: "alicloud_cloud_firewall_instance_member", Factory: ResourceCloudFirewallInstanceMember},
			{TypeName: "alicloud_cloud_firewall_address_book", Factory: ResourceCloudFirewallAddressBook},
			{TypeName: "alicloud_cloud_firewall_control_policy", Factory: ResourceCloudFirewallControlPolicy},
			{TypeName: "alicloud_cloud_firewall_control_policy_order", Factory: ResourceCloudFirewallControlPolicyOrder},
			{TypeName: "alicloud_cloud_firewall_fw_switch", Factory: ResourceCloudFirewallFwSwitch},
			{TypeName: "alicloud_cloud_firewall_ips_config", Factory: ResourceCloudFirewallIPSConfig},
			{TypeName: "alicloud_cloud_firewall_nat_firewall", Factory: ResourceCloudFirewallNatFirewall},
			{TypeName: "alicloud_cloud_firewall_nat_firewall_control_policy", Factory: ResourceCloudFirewallNatFirewallControlPolicy},
			{TypeName: "alicloud_cloud_firewall_nat_firewall_control_policy_order", Factory: ResourceCloudFirewallNatFirewallControlPolicyOrder},
			{TypeName: "alicloud_cloud_firewall_policy_advanced_config", Factory: ResourceCloudFirewallPolicyAdvancedConfig},
			{TypeName: "alicloud_cloud_firewall_private_dns", Factory: ResourceCloudFirewallPrivateDns},
			{TypeName: "alicloud_cloud_firewall_threat_intelligence_switch", Factory: ResourceCloudFirewallThreatIntelligenceSwitch},
			{TypeName: "alicloud_cloud_firewall_user_alarm_config", Factory: ResourceCloudFirewallUserAlarmConfig},
			{TypeName: "alicloud_cloud_firewall_vpc_cen_tr_firewall", Factory: ResourceCloudFirewallVpcCenTrFirewall},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall", Factory: ResourceCloudFirewallVpcFirewall},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_acl_engine_mode", Factory: ResourceCloudFirewallVpcFirewallAclEngineMode},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_cen", Factory: ResourceCloudFirewallVpcFirewallCen},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_control_policy", Factory: ResourceCloudFirewallVpcFirewallControlPolicy},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_control_policy_order", Factory: ResourceCloudFirewallVpcFirewallControlPolicyOrder},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_ips_config", Factory: ResourceCloudFirewallVpcFirewallIpsConfig},
			{TypeName: "alicloud_cloud_firewall_ai_traffic_analysis_status", Factory: ResourceCloudFirewallAiTrafficAnalysisStatus},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_cloud_firewall_instances", Factory: DataSourceCloudFirewallInstances},
			{TypeName: "alicloud_cloud_firewall_instance_members", Factory: DataSourceCloudFirewallInstanceMembers},
			{TypeName: "alicloud_cloud_firewall_address_books", Factory: DataSourceCloudFirewallAddressBooks},
			{TypeName: "alicloud_cloud_firewall_control_policies", Factory: DataSourceCloudFirewallControlPolicies},
			{TypeName: "alicloud_cloud_firewall_nat_firewalls", Factory: DataSourceCloudFirewallNatFirewalls},
			{TypeName: "alicloud_cloud_firewall_tls_inspect_ca_certificates", Factory: DataSourceCloudFirewallTlsInspectCaCertificates},
			{TypeName: "alicloud_cloud_firewall_vpc_cen_tr_firewalls", Factory: DataSourceCloudFirewallVpcCenTrFirewalls},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_cens", Factory: DataSourceCloudFirewallVpcFirewallCens},
			{TypeName: "alicloud_cloud_firewall_vpc_firewalls", Factory: DataSourceCloudFirewallVpcFirewalls},
			{TypeName: "alicloud_cloud_firewall_vpc_firewall_control_policies", Factory: DataSourceCloudFirewallVpcFirewallControlPolicies},
		},
	}
}
