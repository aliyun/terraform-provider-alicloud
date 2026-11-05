package wafopenapi

// convertClusterTypeResponse and its siblings mirror the root-package helpers
// used by the legacy waf data sources (data_source_waf_domains). They are
// root-only (defined in resource_alicloud_waf_domain.go, which stays in the
// root package for the out-of-scope alicloud_waf_domain resource) and cannot
// be imported from a service subpackage, so they are copied verbatim here to
// preserve the original business behavior of the migrated data source.

func convertClusterTypeResponse(source int) string {
	switch source {
	case 0:
		return "PhysicalCluster"
	case 1:
		return "VirtualCluster"
	}
	return ""
}

func convertHttpToUserIpResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}

func convertHttpsRedirectResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}

func convertIsAccessProductResponse(source int) string {
	switch source {
	case 0:
		return "Off"
	case 1:
		return "On"
	}
	return ""
}

func convertLoadBalancingResponse(source int) string {
	switch source {
	case 0:
		return "IpHash"
	case 1:
		return "RoundRobin"
	}
	return ""
}
