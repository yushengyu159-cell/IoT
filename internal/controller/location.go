package controller

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"fabric-sdk/internal/service"
)

// LocationController 位置信息控制器
type LocationController struct {
	locationService *service.LocationService
}

// NewLocationController 创建位置信息控制器
func NewLocationController() *LocationController {
	return &LocationController{
		locationService: service.NewLocationService(),
	}
}

// GetProvinces 获取省份列表
func (c *LocationController) GetProvinces(r *ghttp.Request) {
	ctx := r.Context()
	
	provinces, err := c.locationService.GetProvinces(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取省份列表失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取省份列表失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取省份列表成功",
		Data:    provinces,
	})
}

// GetCities 获取城市列表
func (c *LocationController) GetCities(r *ghttp.Request) {
	ctx := r.Context()
	
	provinceIDStr := r.Get("province_id").String()
	if provinceIDStr == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "省份ID不能为空",
			Data:    nil,
		})
		return
	}
	
	provinceID, err := strconv.Atoi(provinceIDStr)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "省份ID格式错误",
			Data:    nil,
		})
		return
	}
	
	cities, err := c.locationService.GetCitiesByProvince(ctx, provinceID)
	if err != nil {
		g.Log().Error(ctx, "获取城市列表失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取城市列表失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取城市列表成功",
		Data:    cities,
	})
}

// GetDistricts 获取区县列表
func (c *LocationController) GetDistricts(r *ghttp.Request) {
	ctx := r.Context()
	
	cityIDStr := r.Get("city_id").String()
	if cityIDStr == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "城市ID不能为空",
			Data:    nil,
		})
		return
	}
	
	cityID, err := strconv.Atoi(cityIDStr)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "城市ID格式错误",
			Data:    nil,
		})
		return
	}
	
	districts, err := c.locationService.GetDistrictsByCity(ctx, cityID)
	if err != nil {
		g.Log().Error(ctx, "获取区县列表失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取区县列表失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取区县列表成功",
		Data:    districts,
	})
}

// SearchPOI 搜索兴趣点
func (c *LocationController) SearchPOI(r *ghttp.Request) {
	ctx := r.Context()
	
	keyword := r.Get("keyword").String()
	cityIDStr := r.Get("city_id").String()
	poiType := r.Get("type").String()
	
	var cityID int
	if cityIDStr != "" {
		var err error
		cityID, err = strconv.Atoi(cityIDStr)
		if err != nil {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    400,
				Message: "城市ID格式错误",
				Data:    nil,
			})
			return
		}
	}
	
	pois, err := c.locationService.SearchPOI(ctx, keyword, cityID, poiType)
	if err != nil {
		g.Log().Error(ctx, "搜索兴趣点失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "搜索兴趣点失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "搜索兴趣点成功",
		Data:    pois,
	})
}

// SaveLocation 保存用户位置信息
func (c *LocationController) SaveLocation(r *ghttp.Request) {
	ctx := r.Context()
	
	var req struct {
		UserEmail        string  `json:"user_email" v:"required#用户邮箱不能为空"`
		Country          string  `json:"country"`
		CountryCode      string  `json:"country_code"`
		ProvinceID       int     `json:"province_id"`
		ProvinceName     string  `json:"province_name"`
		CityID           int     `json:"city_id"`
		CityName         string  `json:"city_name"`
		DistrictID       int     `json:"district_id"`
		DistrictName     string  `json:"district_name"`
		MetroStation     string  `json:"metro_station"`
		BusinessDistrict string  `json:"business_district"`
		Latitude         float64 `json:"latitude"`
		Longitude        float64 `json:"longitude"`
		AddressDetail    string  `json:"address_detail"`
		FormattedAddress string  `json:"formatted_address"`
		BuildingName     string  `json:"building_name"`
		BuildingType     string  `json:"building_type"`
	}
	
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	location := &service.UserLocation{
		UserEmail:        req.UserEmail,
		Country:          req.Country,
		CountryCode:      req.CountryCode,
		ProvinceID:       req.ProvinceID,
		ProvinceName:     req.ProvinceName,
		CityID:           req.CityID,
		CityName:         req.CityName,
		DistrictID:       req.DistrictID,
		DistrictName:     req.DistrictName,
		MetroStation:     req.MetroStation,
		BusinessDistrict: req.BusinessDistrict,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		AddressDetail:    req.AddressDetail,
		FormattedAddress: req.FormattedAddress,
		BuildingName:     req.BuildingName,
		BuildingType:     req.BuildingType,
		IsPrimary:        true,
	}
	
	err := c.locationService.SaveUserLocation(ctx, location)
	if err != nil {
		g.Log().Error(ctx, "保存用户位置信息失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "保存用户位置信息失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "保存用户位置信息成功",
		Data:    location,
	})
}

// GetUserLocation 获取用户位置信息
func (c *LocationController) GetUserLocation(r *ghttp.Request) {
	ctx := r.Context()
	
	userEmail := r.Get("user_email").String()
	if userEmail == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "用户邮箱不能为空",
			Data:    nil,
		})
		return
	}
	
	location, err := c.locationService.GetUserLocation(ctx, userEmail)
	if err != nil {
		g.Log().Error(ctx, "获取用户位置信息失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取用户位置信息失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取用户位置信息成功",
		Data:    location,
	})
}

// GetNearbyPOIs 获取附近兴趣点
func (c *LocationController) GetNearbyPOIs(r *ghttp.Request) {
	ctx := r.Context()
	
	latitudeStr := r.Get("latitude").String()
	longitudeStr := r.Get("longitude").String()
	radiusStr := r.Get("radius").String()
	poiType := r.Get("type").String()
	
	if latitudeStr == "" || longitudeStr == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "经纬度不能为空",
			Data:    nil,
		})
		return
	}
	
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "纬度格式错误",
			Data:    nil,
		})
		return
	}
	
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "经度格式错误",
			Data:    nil,
		})
		return
	}
	
	radius := 1000 // 默认1公里
	if radiusStr != "" {
		radius, err = strconv.Atoi(radiusStr)
		if err != nil {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    400,
				Message: "半径格式错误",
				Data:    nil,
			})
			return
		}
	}
	
	pois, err := c.locationService.GetNearbyPOIs(ctx, latitude, longitude, radius, poiType)
	if err != nil {
		g.Log().Error(ctx, "获取附近兴趣点失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取附近兴趣点失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取附近兴趣点成功",
		Data:    pois,
	})
}

// GeocodeAddress 地理编码
func (c *LocationController) GeocodeAddress(r *ghttp.Request) {
	ctx := r.Context()
	
	address := r.Get("address").String()
	if address == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "地址不能为空",
			Data:    nil,
		})
		return
	}
	
	latitude, longitude, err := c.locationService.GeocodeAddress(ctx, address)
	if err != nil {
		g.Log().Error(ctx, "地理编码失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "地理编码失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "地理编码成功",
		Data: map[string]interface{}{
			"latitude":  latitude,
			"longitude": longitude,
		},
	})
}

// ReverseGeocode 逆地理编码
func (c *LocationController) ReverseGeocode(r *ghttp.Request) {
	ctx := r.Context()
	
	latitudeStr := r.Get("latitude").String()
	longitudeStr := r.Get("longitude").String()
	
	if latitudeStr == "" || longitudeStr == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "经纬度不能为空",
			Data:    nil,
		})
		return
	}
	
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "纬度格式错误",
			Data:    nil,
		})
		return
	}
	
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "经度格式错误",
			Data:    nil,
		})
		return
	}
	
	address, err := c.locationService.ReverseGeocode(ctx, latitude, longitude)
	if err != nil {
		g.Log().Error(ctx, "逆地理编码失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "逆地理编码失败",
			Data:    nil,
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "逆地理编码成功",
		Data: map[string]interface{}{
			"address": address,
		},
	})
}
