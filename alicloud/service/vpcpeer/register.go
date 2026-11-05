package vpcpeer

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the VpcPeer service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "VpcPeer",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_vpc_peer_connection", Factory: ResourcePeerConnection},
			{TypeName: "alicloud_vpc_peer_connection_accepter", Factory: ResourcePeerConnectionAccepter},
		},
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_vpc_peer_connections", Factory: DataSourcePeerConnections},
		},
	}
}
