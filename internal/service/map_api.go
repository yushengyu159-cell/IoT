package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// MapAPIService 地图API服务
type MapAPIService struct {
	amapKey    string
	baiduKey   string
	httpClient *http.Client
}

// NewMapAPIService 创建地图API服务实例
func NewMapAPIService() *MapAPIService {
	return &MapAPIService{
		amapKey:  g.Cfg().MustGet(context.Background(), "map.amap.key").String(),
		baiduKey: g.Cfg().MustGet(context.Background(), "map.baidu.key").String(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GeocodeResult 地理编码结果
type GeocodeResult struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	Formatted string  `json:"formatted"`
}

// ReverseGeocodeResult 逆地理编码结果
type ReverseGeocodeResult struct {
	Address   string `json:"address"`
	Formatted string `json:"formatted"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
}

// AmapGeocodeResponse 高德地图地理编码响应
type AmapGeocodeResponse struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Count    string `json:"count"`
	Geocodes []struct {
		FormattedAddress string `json:"formatted_address"`
		Province          string `json:"province"`
		City              string `json:"city"`
		District          string `json:"district"`
		Location          string `json:"location"`
		Level             string `json:"level"`
	} `json:"geocodes"`
}

// AmapReverseGeocodeResponse 高德地图逆地理编码响应
type AmapReverseGeocodeResponse struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Regeocode struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
}

// BaiduGeocodeResponse 百度地图地理编码响应
type BaiduGeocodeResponse struct {
	Status string `json:"status"`
	Result []struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
		} `json:"address_component"`
	} `json:"result"`
}

// BaiduReverseGeocodeResponse 百度地图逆地理编码响应
type BaiduReverseGeocodeResponse struct {
	Status string `json:"status"`
	Result struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
		} `json:"address_component"`
	} `json:"result"`
}

// GeocodeAddress 地理编码（地址转坐标）
func (s *MapAPIService) GeocodeAddress(ctx context.Context, address string) (*GeocodeResult, error) {
	// 优先使用 Nominatim（OpenStreetMap），免费且无需API密钥
	result, err := s.nominatimGeocode(ctx, address)
	if err == nil {
		return result, nil
	}
	g.Log().Warning(ctx, "Nominatim地理编码失败，尝试其他服务:", err)

	// 备用：高德地图API
	if s.amapKey != "" {
		result, err := s.amapGeocode(ctx, address)
		if err == nil {
			return result, nil
		}
		g.Log().Warning(ctx, "高德地图地理编码失败，尝试百度地图:", err)
	}

	// 备用：百度地图API
	if s.baiduKey != "" {
		result, err := s.baiduGeocode(ctx, address)
		if err == nil {
			return result, nil
		}
		g.Log().Warning(ctx, "百度地图地理编码失败:", err)
	}

	return nil, fmt.Errorf("所有地图API都不可用")
}

// ReverseGeocode 逆地理编码（坐标转地址）
func (s *MapAPIService) ReverseGeocode(ctx context.Context, latitude, longitude float64) (*ReverseGeocodeResult, error) {
	// 优先使用 Nominatim（OpenStreetMap），免费且无需API密钥
	result, err := s.nominatimReverseGeocode(ctx, latitude, longitude)
	if err == nil {
		return result, nil
	}
	g.Log().Warning(ctx, "Nominatim逆地理编码失败，尝试其他服务:", err)

	// 备用：高德地图API
	if s.amapKey != "" {
		result, err := s.amapReverseGeocode(ctx, latitude, longitude)
		if err == nil {
			return result, nil
		}
		g.Log().Warning(ctx, "高德地图逆地理编码失败，尝试百度地图:", err)
	}

	// 备用：百度地图API
	if s.baiduKey != "" {
		result, err := s.baiduReverseGeocode(ctx, latitude, longitude)
		if err == nil {
			return result, nil
		}
		g.Log().Warning(ctx, "百度地图逆地理编码失败:", err)
	}

	return nil, fmt.Errorf("所有地图API都不可用")
}

// amapGeocode 高德地图地理编码
func (s *MapAPIService) amapGeocode(ctx context.Context, address string) (*GeocodeResult, error) {
	baseURL := "https://restapi.amap.com/v3/geocode/geo"
	params := url.Values{}
	params.Set("key", s.amapKey)
	params.Set("address", address)
	params.Set("output", "json")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result AmapGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("高德地图API错误: %s", result.Info)
	}

	if len(result.Geocodes) == 0 {
		return nil, fmt.Errorf("未找到地址对应的坐标")
	}

	geocode := result.Geocodes[0]
	location := geocode.Location
	coords := make([]string, 2)
	for i, coord := range location {
		if coord == ',' {
			coords = []string{location[:i], location[i+1:]}
			break
		}
	}

	lat, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return nil, fmt.Errorf("解析纬度失败: %v", err)
	}

	lng, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return nil, fmt.Errorf("解析经度失败: %v", err)
	}

	return &GeocodeResult{
		Latitude:  lat,
		Longitude: lng,
		Address:   geocode.FormattedAddress,
		Formatted: geocode.FormattedAddress,
	}, nil
}

// amapReverseGeocode 高德地图逆地理编码
func (s *MapAPIService) amapReverseGeocode(ctx context.Context, latitude, longitude float64) (*ReverseGeocodeResult, error) {
	baseURL := "https://restapi.amap.com/v3/geocode/regeo"
	params := url.Values{}
	params.Set("key", s.amapKey)
	params.Set("location", fmt.Sprintf("%.6f,%.6f", longitude, latitude))
	params.Set("output", "json")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result AmapReverseGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("高德地图API错误: %s", result.Info)
	}

	return &ReverseGeocodeResult{
		Address:   result.Regeocode.FormattedAddress,
		Formatted: result.Regeocode.FormattedAddress,
		Province:  result.Regeocode.AddressComponent.Province,
		City:      result.Regeocode.AddressComponent.City,
		District:  result.Regeocode.AddressComponent.District,
	}, nil
}

// baiduGeocode 百度地图地理编码
func (s *MapAPIService) baiduGeocode(ctx context.Context, address string) (*GeocodeResult, error) {
	baseURL := "https://api.map.baidu.com/geocoding/v2/"
	params := url.Values{}
	params.Set("ak", s.baiduKey)
	params.Set("address", address)
	params.Set("output", "json")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result BaiduGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "0" {
		return nil, fmt.Errorf("百度地图API错误: %s", result.Status)
	}

	if len(result.Result) == 0 {
		return nil, fmt.Errorf("未找到地址对应的坐标")
	}

	location := result.Result[0].Location

	return &GeocodeResult{
		Latitude:  location.Lat,
		Longitude: location.Lng,
		Address:   result.Result[0].FormattedAddress,
		Formatted: result.Result[0].FormattedAddress,
	}, nil
}

// baiduReverseGeocode 百度地图逆地理编码
func (s *MapAPIService) baiduReverseGeocode(ctx context.Context, latitude, longitude float64) (*ReverseGeocodeResult, error) {
	baseURL := "https://api.map.baidu.com/reverse_geocoding/v2/"
	params := url.Values{}
	params.Set("ak", s.baiduKey)
	params.Set("location", fmt.Sprintf("%.6f,%.6f", latitude, longitude))
	params.Set("output", "json")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result BaiduReverseGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "0" {
		return nil, fmt.Errorf("百度地图API错误: %s", result.Status)
	}

	return &ReverseGeocodeResult{
		Address:   result.Result.FormattedAddress,
		Formatted: result.Result.FormattedAddress,
		Province:  result.Result.AddressComponent.Province,
		City:      result.Result.AddressComponent.City,
		District:  result.Result.AddressComponent.District,
	}, nil
}

// NominatimGeocodeResponse Nominatim 地理编码响应
type NominatimGeocodeResponse []struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Address     struct {
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		State       string `json:"state"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Road        string `json:"road"`
		Suburb      string `json:"suburb"`
		Postcode    string `json:"postcode"`
	} `json:"address"`
}

// NominatimReverseGeocodeResponse Nominatim 逆地理编码响应
type NominatimReverseGeocodeResponse struct {
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Address     struct {
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		State       string `json:"state"`
		Province    string `json:"province"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Road        string `json:"road"`
		Suburb      string `json:"suburb"`
		Postcode    string `json:"postcode"`
		HouseNumber string `json:"house_number"`
		Building    string `json:"building"`
	} `json:"address"`
}

// nominatimGeocode 使用 OpenStreetMap Nominatim 服务进行地理编码（地址转坐标）
func (s *MapAPIService) nominatimGeocode(ctx context.Context, address string) (*GeocodeResult, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", address)
	params.Add("format", "json")
	params.Add("limit", "1")
	params.Add("accept-language", "zh-CN,zh-TW,en") // 支持简体中文、繁体中文、英文

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// ⚠️ 必须设置 User-Agent，否则 Nominatim 会拒绝请求
	req.Header.Set("User-Agent", "ESG-VISA/1.0 (fabric-sdk)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Nominatim返回状态码 %d: %s", resp.StatusCode, string(body))
	}

	var results NominatimGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("未找到地址对应的坐标")
	}

	result := results[0]
	lat, err := strconv.ParseFloat(result.Lat, 64)
	if err != nil {
		return nil, fmt.Errorf("解析纬度失败: %v", err)
	}

	lon, err := strconv.ParseFloat(result.Lon, 64)
	if err != nil {
		return nil, fmt.Errorf("解析经度失败: %v", err)
	}

	return &GeocodeResult{
		Latitude:  lat,
		Longitude: lon,
		Address:   result.DisplayName,
		Formatted: result.DisplayName,
	}, nil
}

// nominatimReverseGeocode 使用 OpenStreetMap Nominatim 服务进行逆地理编码（坐标转地址）
func (s *MapAPIService) nominatimReverseGeocode(ctx context.Context, latitude, longitude float64) (*ReverseGeocodeResult, error) {
	baseURL := "https://nominatim.openstreetmap.org/reverse"
	params := url.Values{}
	params.Add("format", "json")
	params.Add("lat", fmt.Sprintf("%.6f", latitude))
	params.Add("lon", fmt.Sprintf("%.6f", longitude))
	params.Add("zoom", "18")                        // zoom=18 = building level precision
	params.Add("addressdetails", "1")
	params.Add("accept-language", "zh-CN,zh-TW,en") // 支持简体中文、繁体中文、英文

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// ⚠️ 必须设置 User-Agent，否则 Nominatim 会拒绝请求
	req.Header.Set("User-Agent", "ESG-VISA/1.0 (fabric-sdk)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Nominatim返回状态码 %d: %s", resp.StatusCode, string(body))
	}

	var result NominatimReverseGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 提取省市区信息
	province := result.Address.Province
	if province == "" {
		province = result.Address.State
	}

	city := result.Address.City
	if city == "" {
		city = result.Address.Town
	}
	if city == "" {
		city = result.Address.Village
	}

	district := result.Address.Suburb

	// Build concise address: building/road + house_number + district + city
	address := result.DisplayName // fallback to full display name
	if result.Name != "" && result.Name != result.Address.Road {
		// Nominatim returned a POI/building name
		parts := []string{result.Name}
		if result.Address.HouseNumber != "" {
			parts = append(parts, result.Address.HouseNumber)
		}
		if result.Address.Road != "" {
			parts = append(parts, result.Address.Road)
		}
		if district != "" {
			parts = append(parts, district)
		}
		if city != "" {
			parts = append(parts, city)
		}
		address = strings.Join(parts, ", ")
	}

	return &ReverseGeocodeResult{
		Address:   address,
		Formatted: address,
		Province:  province,
		City:      city,
		District:  district,
	}, nil
}