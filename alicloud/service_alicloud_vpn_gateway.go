package alicloud

import (
	"time"

	"strings"

	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func (client *AliyunClient) DescribeVpnGateway(vpnId string) (v vpc.DescribeVpnGatewayResponse, err error) {
	request := vpc.CreateDescribeVpnGatewayRequest()
	request.VpnGatewayId = vpnId

	resp, err := client.vpcconn.DescribeVpnGateway(request)
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, VpnNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", vpnId))
		}
		return
	}
	if resp == nil || resp.VpnGatewayId != vpnId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", vpnId))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeCustomerGateway(cgwId string) (v vpc.DescribeCustomerGatewayResponse, err error) {
	request := vpc.CreateDescribeCustomerGatewayRequest()
	request.CustomerGatewayId = cgwId

	resp, err := client.vpcconn.DescribeCustomerGateway(request)
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, CgwNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN customer gateway", cgwId))
		}
		return
	}
	if resp == nil || resp.CustomerGatewayId != cgwId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN customer gateway", cgwId))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeVpnConnection(id string) (v vpc.DescribeVpnConnectionResponse, err error) {
	request := vpc.CreateDescribeVpnConnectionRequest()
	request.VpnConnectionId = id

	resp, err := client.vpcconn.DescribeVpnConnection(request)
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, VpnConnNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN connection", id))
		}
		return
	}
	if resp == nil || resp.VpnConnectionId != id {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN connection", id))
	}
	return *resp, nil
}

func (client *AliyunClient) DescribeSslVpnServers(vpnId string, sslId string) (v vpc.DescribeSslVpnServersResponse, err error) {
	request := vpc.CreateDescribeSslVpnServersRequest()
	if sslId != "" {
		request.SslVpnServerId = sslId
	}

	if vpnId != "" {
		request.VpnGatewayId = vpnId
	}
	resp, err := client.vpcconn.DescribeSslVpnServers(request)
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, SslVpnServerNotFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", sslId))
		}
		return
	}

	if resp == nil || 0 == len(resp.SslVpnServers.SslVpnServer) {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", sslId))
	}

	if sslId != "" && sslId != resp.SslVpnServers.SslVpnServer[0].SslVpnServerId {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", sslId))
	}

	return *resp, nil
}

func (client *AliyunClient) DescribeSslVpnClientCert(id string) (v vpc.DescribeSslVpnClientCertResponse, err error) {
	request := vpc.CreateDescribeSslVpnClientCertRequest()
	request.SslVpnClientCertId = id

	resp, err := client.vpcconn.DescribeSslVpnClientCert(request)
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, SslVpnClientCertNofFound}) {
			return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
		}
		return
	}
	if resp == nil || resp.SslVpnClientCertId != id {
		return v, GetNotFoundErrorFromString(GetNotFoundMessage("VPN", id))
	}
	return *resp, nil
}

func (client *AliyunClient) WaitForVpn(vpnId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		//wait the order effective
		vpn, err := client.DescribeVpnGateway(vpnId)
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

func (client *AliyunClient) WaitForCustomerGateway(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		_, err := client.DescribeCustomerGateway(id)
		if err != nil {
			return err
		} else {
			break
		}
	}
	return nil
}

func (client *AliyunClient) WaitForSslVpnClientCert(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("VPN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		resp, err := client.DescribeSslVpnClientCert(id)
		if err != nil {
			return err
		}

		if strings.ToLower(resp.Status) == strings.ToLower(string(status)) {
			break
		}
	}
	return nil
}

func ParseIkeConfig(ike vpc.IkeConfig) (ikeConfigs []map[string]interface{}) {
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

func ParseIpsecConfig(ipsec vpc.IpsecConfig) (ipsecConfigs []map[string]interface{}) {
	item := map[string]interface{}{
		"ipsec_auth_alg": ipsec.IpsecAuthAlg,
		"ipsec_enc_alg":  ipsec.IpsecEncAlg,
		"ipsec_lifetime": ipsec.IpsecLifetime,
		"ipsec_pfs":      ipsec.IpsecPfs,
	}

	ipsecConfigs = append(ipsecConfigs, item)
	return
}

func AssembleIkeConfig(ikeCfgParam []interface{}) (string, error) {
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

func AssembleIpsecConfig(ipsecCfgParam []interface{}) (string, error) {
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

func AssembleNetworkSubnetToString(list []interface{}) string {
	if len(list) < 1 {
		return ""
	}
	var items []string
	for _, id := range list {
		items = append(items, fmt.Sprintf("%s", id))
	}
	return fmt.Sprintf("%s", strings.Join(items, COMMA_SEPARATED))
}
