package bailian

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// Service provides typed query methods for Bailian resources.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService returns a Service wrapping the given client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// GetPromptTemplate queries GetPromptTemplate and returns the response object,
// or an error when the prompt template does not exist.
func (s *Service) GetPromptTemplate(workspaceId, promptTemplateId string) (map[string]interface{}, error) {
	return getPromptTemplate(s.client, workspaceId, promptTemplateId)
}

// GetMemory queries GetMemory and returns the response object,
// or an error when the memory does not exist.
func (s *Service) GetMemory(workspaceId, memoryId string) (map[string]interface{}, error) {
	return getMemory(s.client, workspaceId, memoryId)
}

// GetMemoryNode queries GetMemoryNode and returns the response object,
// or an error when the memory node does not exist.
func (s *Service) GetMemoryNode(workspaceId, memoryId, memoryNodeId string) (map[string]interface{}, error) {
	return getMemoryNode(s.client, workspaceId, memoryId, memoryNodeId)
}

// GetConnector queries GetConnector and returns the Data object,
// or an error when the connector does not exist.
func (s *Service) GetConnector(workspaceId, connectorId string) (map[string]interface{}, error) {
	return getConnector(s.client, workspaceId, connectorId)
}
