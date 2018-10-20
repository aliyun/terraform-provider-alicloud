package alicloud

import (
	"time"

	"strings"

	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type VpnGatewayService struct {
	client *connectivity.AliyunClient
}

func (s *VpnGatewayService) DescribeVpnGateway(vpnId string) (v vpc.DescribeVpnGatewayResponse, err error) {
	request := vpc.CreateDescribeVpnGatewayRequest()
	request.VpnGatewayId = vpnId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpnGateway(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, VpnNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", vpnId))
		}
		return
	}
	resp, _ := raw.(*vpc.DescribeVpnGatewayResponse)
	if resp == nil || resp.VpnGatewayId != vpnId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", vpnId))
	}
	return *resp, nil
}

func (s *VpnGatewayService) DescribeCustomerGateway(cgwId string) (v vpc.DescribeCustomerGatewayResponse, err error) {
	request := vpc.CreateDescribeCustomerGatewayRequest()
	request.CustomerGatewayId = cgwId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeCustomerGateway(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, CgwNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN customer gateway", cgwId))
		}
		return
	}
	resp, _ := raw.(*vpc.DescribeCustomerGatewayResponse)
	if resp == nil || resp.CustomerGatewayId != cgwId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", cgwId))
	}
	return *resp, nil
}

func (s *VpnGatewayService) DescribeVpnConnection(id string) (v vpc.DescribeVpnConnectionResponse, err error) {
	request := vpc.CreateDescribeVpnConnectionRequest()
	request.VpnConnectionId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpnConnection(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, VpnConnNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN connection", id))
		}
		return
	}
	resp, _ := raw.(*vpc.DescribeVpnConnectionResponse)
	if resp == nil || resp.VpnConnectionId != id {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN connection", id))
	}
	return *resp, nil
}

func (s *VpnGatewayService) DescribeSslVpnServer(sslId string) (v vpc.SslVpnServer, err error) {
	request := vpc.CreateDescribeSslVpnServersRequest()
	request.SslVpnServerId = sslId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeSslVpnServers(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, SslVpnServerNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("SSL VPN server", sslId))
		}
		return
	}
	resp, _ := raw.(*vpc.DescribeSslVpnServersResponse)
	if resp == nil || 0 == len(resp.SslVpnServers.SslVpnServer) {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("SSL VPN server", sslId))
	}

	if sslId != resp.SslVpnServers.SslVpnServer[0].SslVpnServerId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("SSL VPN server", sslId))
	}

	return resp.SslVpnServers.SslVpnServer[0], nil
}

func (s *VpnGatewayService) DescribeSslVpnClientCert(id string) (v vpc.DescribeSslVpnClientCertResponse, err error) {
	request := vpc.CreateDescribeSslVpnClientCertRequest()
	request.SslVpnClientCertId = id

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeSslVpnClientCert(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, SslVpnClientCertNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
		}
		return
	}
	resp, _ := raw.(*vpc.DescribeSslVpnClientCertResponse)
	if resp == nil || resp.SslVpnClientCertId != id {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
	}
	return *resp, nil
}

func (s *VpnGatewayService) WaitForVpn(vpnId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		//wait the order effective
		vpn, err := s.DescribeVpnGateway(vpnId)
		if err != nil {
			return err
		}
		if strings.ToLower(vpn.Status) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *VpnGatewayService) WaitForCustomerGateway(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		_, err := s.DescribeCustomerGateway(id)
		if err != nil {
			return err
		} else {
			break
		}
	}
	return nil
}

func (s *VpnGatewayService) WaitForSslVpnClientCert(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("SSL VPN client cert", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		resp, err := s.DescribeSslVpnClientCert(id)
		if err != nil {
			return err
		}

		if strings.ToLower(resp.Status) == strings.ToLower(string(status)) {
			break
		}
	}
	return nil
}

func (s *VpnGatewayService) ParseIkeConfig(ike vpc.IkeConfig) (ikeConfigs []map[string]interface{}) {
	item := map[string]interface{}{
		"ike_auth_alg":  ike.IkeAuthAlg,
		"ike_enc_alg":   ike.IkeEncAlg,
		"ike_lifetime":  ike.IkeLifetime,
		"ike_local_id":  ike.LocalId,
		"ike_mode":      ike.IkeMode,
		"ike_pfs":       ike.IkePfs,
		"ike_remote_id": ike.RemoteId,
		"ike_version":   ike.IkeVersion,
		"psk":           ike.Psk,
	}

	ikeConfigs = append(ikeConfigs, item)
	return
}

func (s *VpnGatewayService) ParseIpsecConfig(ipsec vpc.IpsecConfig) (ipsecConfigs []map[string]interface{}) {
	item := map[string]interface{}{
		"ipsec_auth_alg": ipsec.IpsecAuthAlg,
		"ipsec_enc_alg":  ipsec.IpsecEncAlg,
		"ipsec_lifetime": ipsec.IpsecLifetime,
		"ipsec_pfs":      ipsec.IpsecPfs,
	}

	ipsecConfigs = append(ipsecConfigs, item)
	return
}

func (s *VpnGatewayService) AssembleIkeConfig(ikeCfgParam []interface{}) (string, error) {
	var ikeCfg IkeConfig
	v := ikeCfgParam[0]
	item := v.(map[string]interface{})
	ikeCfg = IkeConfig{
		IkeAuthAlg:  item["ike_auth_alg"].(string),
		IkeEncAlg:   item["ike_enc_alg"].(string),
		IkeLifetime: item["ike_lifetime"].(int),
		LocalId:     item["ike_local_id"].(string),
		IkeMode:     item["ike_mode"].(string),
		IkePfs:      item["ike_pfs"].(string),
		RemoteId:    item["ike_remote_id"].(string),
		IkeVersion:  item["ike_version"].(string),
		Psk:         item["psk"].(string),
	}

	data, err := json.Marshal(ikeCfg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *VpnGatewayService) AssembleIpsecConfig(ipsecCfgParam []interface{}) (string, error) {
	var ipsecCfg IpsecConfig
	v := ipsecCfgParam[0]
	item := v.(map[string]interface{})
	ipsecCfg = IpsecConfig{
		IpsecAuthAlg:  item["ipsec_auth_alg"].(string),
		IpsecEncAlg:   item["ipsec_enc_alg"].(string),
		IpsecLifetime: item["ipsec_lifetime"].(int),
		IpsecPfs:      item["ipsec_pfs"].(string),
	}

	data, err := json.Marshal(ipsecCfg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *VpnGatewayService) AssembleNetworkSubnetToString(list []interface{}) string {
	if len(list) < 1 {
		return ""
	}
	var items []string
	for _, id := range list {
		items = append(items, fmt.Sprintf("%s", id))
	}
	return fmt.Sprintf("%s", strings.Join(items, COMMA_SEPARATED))
}
