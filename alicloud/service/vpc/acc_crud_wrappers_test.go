package vpc_test

// CRUD call-through wrappers for the external test package vpc_test. The
// migrated unit tests call the production CRUD functions by their unexported
// names (e.g. resourceAliCloudVpcVswitchCreate(d, rawClient)); those are
// invisible from vpc_test. Each wrapper re-uses the exported factory\'s
// schema.Resource CRUD field, preserving the exact call signature, so the
// test call sites stay byte-identical and the factory\'s live schema wiring
// (timeouts, schema, etc.) is exercised rather than a stale copy.

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudVpcFlowLogCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcFlowLog().Create(d, m)
}

func resourceAliCloudVpcFlowLogDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcFlowLog().Delete(d, m)
}

func resourceAliCloudVpcFlowLogRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcFlowLog().Read(d, m)
}

func resourceAliCloudVpcFlowLogUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcFlowLog().Update(d, m)
}

func resourceAliCloudVpcIpv4CidrBlockCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4CidrBlock().Create(d, m)
}

func resourceAliCloudVpcIpv4CidrBlockDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4CidrBlock().Delete(d, m)
}

func resourceAliCloudVpcIpv4CidrBlockRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4CidrBlock().Read(d, m)
}

func resourceAliCloudVpcIpv4GatewayCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4Gateway().Create(d, m)
}

func resourceAliCloudVpcIpv4GatewayDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4Gateway().Delete(d, m)
}

func resourceAliCloudVpcIpv4GatewayRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4Gateway().Read(d, m)
}

func resourceAliCloudVpcIpv4GatewayUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv4Gateway().Update(d, m)
}

func resourceAliCloudVpcIpv6EgressRuleCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6EgressRule().Create(d, m)
}

func resourceAliCloudVpcIpv6EgressRuleDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6EgressRule().Delete(d, m)
}

func resourceAliCloudVpcIpv6EgressRuleRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6EgressRule().Read(d, m)
}

func resourceAliCloudVpcIpv6InternetBandwidthCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6InternetBandwidth().Create(d, m)
}

func resourceAliCloudVpcIpv6InternetBandwidthDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6InternetBandwidth().Delete(d, m)
}

func resourceAliCloudVpcIpv6InternetBandwidthRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6InternetBandwidth().Read(d, m)
}

func resourceAliCloudVpcIpv6InternetBandwidthUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6InternetBandwidth().Update(d, m)
}

func resourceAliCloudVpcNetworkAclAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcNetworkAclAttachment().Create(d, m)
}

func resourceAliCloudVpcNetworkAclAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcNetworkAclAttachment().Delete(d, m)
}

func resourceAliCloudVpcNetworkAclAttachmentRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcNetworkAclAttachment().Read(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPoolCidrBlock().Create(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPoolCidrBlock().Delete(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPoolCidrBlock().Read(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPool().Create(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPool().Delete(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPool().Read(d, m)
}

func resourceAliCloudVpcPublicIpAddressPoolUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPublicIpAddressPool().Update(d, m)
}

func resourceAliCloudVpcRouteTableCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcRouteTable().Create(d, m)
}

func resourceAliCloudVpcRouteTableDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcRouteTable().Delete(d, m)
}

func resourceAliCloudVpcRouteTableRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcRouteTable().Read(d, m)
}

func resourceAliCloudVpcRouteTableUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcRouteTable().Update(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilter().Create(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilter().Delete(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterEgressRuleCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterEgressRule().Create(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterEgressRuleDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterEgressRule().Delete(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterEgressRuleRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterEgressRule().Read(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterEgressRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterEgressRule().Update(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterIngressRule().Create(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterIngressRule().Delete(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterIngressRule().Read(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterIngressRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilterIngressRule().Update(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilter().Read(d, m)
}

func resourceAliCloudVpcTrafficMirrorFilterUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcTrafficMirrorFilter().Update(d, m)
}

func resourceAliCloudVpcVpcCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVpc().Create(d, m)
}

func resourceAliCloudVpcVpcDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVpc().Delete(d, m)
}

func resourceAliCloudVpcVpcRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVpc().Read(d, m)
}

func resourceAliCloudVpcVpcUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVpc().Update(d, m)
}

func resourceAliCloudVpcVswitchCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVswitch().Create(d, m)
}

func resourceAliCloudVpcVswitchDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVswitch().Delete(d, m)
}

func resourceAliCloudVpcVswitchRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVswitch().Read(d, m)
}

func resourceAliCloudVpcVswitchUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcVswitch().Update(d, m)
}

func resourceAlicloudHavipCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudHavip().Create(d, m)
}

func resourceAlicloudHavipDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudHavip().Delete(d, m)
}

func resourceAlicloudHavipRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudHavip().Read(d, m)
}

func resourceAlicloudHavipUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudHavip().Update(d, m)
}

func resourceAlicloudVpcDhcpOptionsAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudVpcDhcpOptionsSetAttachement().Create(d, m)
}

func resourceAlicloudVpcDhcpOptionsSetAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudVpcDhcpOptionsSetAttachement().Delete(d, m)
}

func resourceAlicloudVpcDhcpOptionsSetAttachmentRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudVpcDhcpOptionsSetAttachement().Read(d, m)
}

func resourceAlicloudVpcDhcpOptionsSetAttachmentUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAlicloudVpcDhcpOptionsSetAttachement().Update(d, m)
}

func resourceAlicloudVpcIpv6GatewayCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6Gateway().Create(d, m)
}

func resourceAlicloudVpcIpv6GatewayDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6Gateway().Delete(d, m)
}

func resourceAlicloudVpcIpv6GatewayRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6Gateway().Read(d, m)
}

func resourceAlicloudVpcIpv6GatewayUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcIpv6Gateway().Update(d, m)
}

func resourceAlicloudVpcPrefixListCreate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPrefixList().Create(d, m)
}

func resourceAlicloudVpcPrefixListDelete(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPrefixList().Delete(d, m)
}

func resourceAlicloudVpcPrefixListRead(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPrefixList().Read(d, m)
}

func resourceAlicloudVpcPrefixListUpdate(d *schema.ResourceData, m interface{}) error {
	return vpc.ResourceAliCloudVpcPrefixList().Update(d, m)
}
