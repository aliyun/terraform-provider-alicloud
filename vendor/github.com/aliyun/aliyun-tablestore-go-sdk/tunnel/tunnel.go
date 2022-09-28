package tunnel

type TunnelClient interface {
	TunnelMetaApi
	NewTunnelWorker(tunnelId string, workerConfig *TunnelWorkerConfig) (TunnelWorker, error)
}

type DefaultTunnelClient struct {
	api *TunnelApi
}

func NewTunnelClient(endpoint, instanceName, accessId, accessKey string, options ...ClientOption) TunnelClient {
	return NewTunnelClientWithConfig(endpoint, instanceName, accessId, accessKey, nil, options...)
}

func NewTunnelClientWithConfig(endpoint, instanceName, accessId, accessKey string, conf *TunnelConfig, options ...ClientOption) TunnelClient {
	return &DefaultTunnelClient{
		api: NewTunnelApi(endpoint, instanceName, accessId, accessKey, conf, options...),
	}
}

func NewTunnelClientWithToken(endpoint, instanceName, accessId, accessKey, token string, conf *TunnelConfig, options ...ClientOption) TunnelClient {
	return &DefaultTunnelClient{
		api: NewTunnelApiWithToken(endpoint, instanceName, accessId, accessKey, token, conf, options...),
	}
}

func NewTunnelClientWithExternalHeader(endpoint, instanceName, accessId, accessKey, token string, header map[string]string) TunnelClient {
	return NewTunnelClientWithConfigAndExternalHeader(endpoint, instanceName, accessId, accessKey, token, nil, header)
}

func NewTunnelClientWithConfigAndExternalHeader(endpoint, instanceName, accessId, accessKey, token string, conf *TunnelConfig, header map[string]string) TunnelClient {
	return &DefaultTunnelClient{
		api: NewTunnelApiWithExternalHeader(endpoint, instanceName, accessId, accessKey, token, conf, header),
	}
}

func (c *DefaultTunnelClient) CreateTunnel(req *CreateTunnelRequest) (*CreateTunnelResponse, error) {
	return c.api.CreateTunnel(req)
}

func (c *DefaultTunnelClient) DeleteTunnel(req *DeleteTunnelRequest) (*DeleteTunnelResponse, error) {
	return c.api.DeleteTunnel(req)
}

func (c *DefaultTunnelClient) ListTunnel(req *ListTunnelRequest) (*ListTunnelResponse, error) {
	return c.api.ListTunnel(req)
}

func (c *DefaultTunnelClient) DescribeTunnel(req *DescribeTunnelRequest) (*DescribeTunnelResponse, error) {
	return c.api.DescribeTunnel(req)
}

func (c *DefaultTunnelClient) GetRpo(req *GetRpoRequest) (*GetRpoResponse, error) {
	return c.api.GetRpo(req)
}

func (c *DefaultTunnelClient) GetRpoByOffset(req *GetRpoRequest) (*GetRpoResponse, error) {
	return c.api.GetRpoByOffset(req)
}

func (c *DefaultTunnelClient) Schedule(req *ScheduleRequest) (*ScheduleResponse, error) {
	return c.api.Schedule(req)
}

func (c *DefaultTunnelClient) NewTunnelWorker(tunnelId string, workerConfig *TunnelWorkerConfig) (TunnelWorker, error) {
	if workerConfig == nil {
		return nil, &TunnelError{Code: ErrCodeClientError, Message: "TunnelWorkerConfig can not be nil"}
	}
	conf := *workerConfig
	return newTunnelWorker(tunnelId, c.api, &conf)
}
