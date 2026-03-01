package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// LiveSenseAPIService LiveSense GraphQL API服务
type LiveSenseAPIService struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
	token      string
	tokenExp   time.Time
}

// NewLiveSenseAPIService 创建LiveSense API服务实例
func NewLiveSenseAPIService() *LiveSenseAPIService {
	return &LiveSenseAPIService{
		baseURL:  "https://graphql-api.livesense.com.au",
		username: "belucksapi",
		password: "i6a6LKS783W7ffsr#",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token        string      `json:"token"`
	ResponseBody interface{} `json:"responseBody"`
	ContextID    int         `json:"contextId"` // API返回的是数字类型
	Permissions  interface{} `json:"permissions"`
}

// Authenticate 获取访问令牌（公开方法）
func (s *LiveSenseAPIService) Authenticate(ctx context.Context) error {
	return s.authenticate(ctx)
}

// authenticate 获取访问令牌（内部方法）
func (s *LiveSenseAPIService) authenticate(ctx context.Context) error {
	// 如果token未过期，直接返回
	if s.token != "" && time.Now().Before(s.tokenExp) {
		return nil
	}

	url := fmt.Sprintf("%s/authenticate", s.baseURL)
	payload := map[string]string{
		"username": s.username,
		"password": s.password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化认证请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建认证请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("认证请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("认证失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("解析认证响应失败: %w", err)
	}

	if authResp.Token == "" {
		return fmt.Errorf("认证响应中未找到token")
	}

	s.token = authResp.Token
	// 设置token过期时间为1小时后（JWT通常有效期较长，这里保守设置为1小时）
	s.tokenExp = time.Now().Add(1 * time.Hour)

	g.Log().Info(ctx, "LiveSense API认证成功")
	return nil
}

// GraphQLRequest GraphQL请求结构
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse GraphQL响应结构
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"` // 使用RawMessage以便延迟解析
	Errors []interface{}   `json:"errors,omitempty"`
}

// executeGraphQL 执行GraphQL查询
func (s *LiveSenseAPIService) executeGraphQL(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	// 确保已认证
	if err := s.authenticate(ctx); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/graphql", s.baseURL)
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化GraphQL请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建GraphQL请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	// 添加调试日志：记录请求内容
	g.Log().Debug(ctx, "发送GraphQL请求，URL:", url, "Variables:", variables)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GraphQL请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 添加调试日志：记录响应内容（前500字符）
	bodyStr := string(body)
	if len(bodyStr) > 500 {
		g.Log().Debug(ctx, "GraphQL响应（前500字符）:", bodyStr[:500])
	} else {
		g.Log().Debug(ctx, "GraphQL完整响应:", bodyStr)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(body, &graphqlResp); err != nil {
		return nil, fmt.Errorf("解析GraphQL响应失败: %w", err)
	}

	if len(graphqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL返回错误: %v", graphqlResp.Errors)
	}

	return &graphqlResp, nil
}

// SensorValueItem 单个传感器数值项
type SensorValueItem struct {
	Sensor     string `json:"sensor"`
	Timezone   string `json:"timezone"`
	DataPoints []struct {
		Value        interface{} `json:"value"` // 可能是字符串或数字
		TimestampUtc string      `json:"timestampUtc"`
	} `json:"dataPoints"`
}

// ValuesResponse 传感器数值响应
// 注意：API可能返回values为对象或数组，或null
type ValuesResponse struct {
	Values interface{} `json:"values"` // 可能是null、对象或数组
}

// GetValuesList 获取传感器数值列表（处理values为对象、数组或null的情况）
func (v *ValuesResponse) GetValuesList() []SensorValueItem {
	if v.Values == nil {
		return nil
	}

	var items []SensorValueItem

	// 尝试解析为数组
	if valuesArray, ok := v.Values.([]interface{}); ok {
		for _, item := range valuesArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				parsedItem := parseSensorValueItem(itemMap)
				if parsedItem != nil {
					items = append(items, *parsedItem)
				}
			}
		}
		return items
	}

	// 尝试解析为单个对象（这是API实际返回的格式）
	if valueObj, ok := v.Values.(map[string]interface{}); ok {
		parsedItem := parseSensorValueItem(valueObj)
		if parsedItem != nil {
			items = append(items, *parsedItem)
		}
		return items
	}

	// 如果都不匹配，记录类型信息用于调试
	g.Log().Warning(context.Background(), "GetValuesList: values类型不匹配，类型:", fmt.Sprintf("%T", v.Values))
	return nil
}

// parseSensorValueItem 解析传感器数值项
func parseSensorValueItem(data map[string]interface{}) *SensorValueItem {
	item := &SensorValueItem{}

	if sensor, ok := data["sensor"].(string); ok {
		item.Sensor = sensor
	}
	if timezone, ok := data["timezone"].(string); ok {
		item.Timezone = timezone
	}

	if dataPoints, ok := data["dataPoints"].([]interface{}); ok {
		for _, dp := range dataPoints {
			if dpMap, ok := dp.(map[string]interface{}); ok {
				var dataPoint struct {
					Value        interface{} `json:"value"`
					TimestampUtc string      `json:"timestampUtc"`
				}
				if value, exists := dpMap["value"]; exists {
					dataPoint.Value = value
				}
				if ts, ok := dpMap["timestampUtc"].(string); ok {
					dataPoint.TimestampUtc = ts
				}
				item.DataPoints = append(item.DataPoints, dataPoint)
			}
		}
	}

	return item
}

// ParseValue 解析value字段（支持字符串和数字）
func ParseValue(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		var f float64
		if _, err := fmt.Sscanf(v, "%f", &f); err == nil {
			return f, true
		}
	}
	return 0, false
}

// GetSensorValues 获取传感器数值
func (s *LiveSenseAPIService) GetSensorValues(ctx context.Context, sensorID, from, to, aggregation, group, limit string) (*ValuesResponse, error) {
	query := `query values ($id: String, $from: String, $to: String, $aggregation: String, $group: String, $limit: String) {
		values (id: $id, from: $from, to: $to, aggregation: $aggregation, group: $group, limit: $limit) {
			sensor
			timezone
			dataPoints {
				value
				timestampUtc
			}
		}
	}`

	// 根据Postman文件，所有参数都应该传递，包括空字符串
	// 测试发现：传递空字符串aggregation和limit时能获取数据，不传这些参数时返回null
	variables := map[string]interface{}{
		"id": sensorID,
	}
	// from和to：如果为空字符串，不传；如果有值，传递
	if from != "" {
		variables["from"] = from
	}
	if to != "" {
		variables["to"] = to
	}
	// aggregation：如果有时间范围，必须传递（即使是空字符串）
	if from != "" && to != "" {
		variables["aggregation"] = aggregation // 即使是空字符串也传递
	} else if aggregation != "" {
		variables["aggregation"] = aggregation
	}
	// group：必须传递（Postman文件中是"allValues"）
	if group != "" {
		variables["group"] = group
	}
	// limit：如果有时间范围，必须传递（即使是空字符串）
	if from != "" && to != "" {
		variables["limit"] = limit // 即使是空字符串也传递
	} else if limit != "" {
		variables["limit"] = limit
	}

	resp, err := s.executeGraphQL(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	// 将响应数据转换为ValuesResponse
	// resp.Data是json.RawMessage类型，直接解析为map
	var dataMap map[string]interface{}
	if err := json.Unmarshal(resp.Data, &dataMap); err != nil {
		return nil, fmt.Errorf("解析响应数据失败: %w", err)
	}

	// 提取values字段
	valuesField, exists := dataMap["values"]
	if !exists {
		g.Log().Warning(ctx, "API响应中未找到values字段")
		return &ValuesResponse{Values: nil}, nil
	}

	// 构建ValuesResponse
	valuesResp := ValuesResponse{
		Values: valuesField, // 直接使用values字段，可能是对象、数组或null
	}

	// 添加调试日志
	if valuesResp.Values == nil {
		g.Log().Warning(ctx, "API返回values为null")
	} else {
		switch v := valuesResp.Values.(type) {
		case map[string]interface{}:
			g.Log().Info(ctx, "API返回values为对象，包含字段数:", len(v))
		case []interface{}:
			g.Log().Info(ctx, "API返回values为数组，长度:", len(v))
		default:
			g.Log().Warning(ctx, "API返回values为未知类型:", fmt.Sprintf("%T", v))
		}
	}

	return &valuesResp, nil
}

// SensorsUnderContextResponse 上下文下传感器响应
type SensorsUnderContextResponse struct {
	SensorsUnderContext []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Sensors []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"sensors"`
	} `json:"sensorsUnderContext"`
}

// GetSensorsUnderContext 获取上下文下的传感器列表
func (s *LiveSenseAPIService) GetSensorsUnderContext(ctx context.Context, contextID int, contextTypeID, groupByTypeID string) (*SensorsUnderContextResponse, error) {
	query := `query sensorsUnderContext($id: Int, $contextTypeId: String, $groupByTypeId: String) {
		sensorsUnderContext(id: $id, contextTypeId: $contextTypeId, groupByTypeId: $groupByTypeId) {
			id
			name
			sensors {
				id
				name
			}
		}
	}`

	variables := map[string]interface{}{
		"id": contextID,
	}
	if contextTypeID != "" {
		variables["contextTypeId"] = contextTypeID
	}
	if groupByTypeID != "" {
		variables["groupByTypeId"] = groupByTypeID
	}

	resp, err := s.executeGraphQL(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	// 将响应数据转换为SensorsUnderContextResponse
	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化响应数据失败: %w", err)
	}

	var sensorsResp SensorsUnderContextResponse
	if err := json.Unmarshal(dataBytes, &sensorsResp); err != nil {
		return nil, fmt.Errorf("解析传感器列表响应失败: %w", err)
	}

	return &sensorsResp, nil
}

// 全局服务实例
var LiveSenseAPI = NewLiveSenseAPIService()

