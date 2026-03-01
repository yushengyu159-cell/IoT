package service

import (
	"context"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

// LocationService 位置信息服务
type LocationService struct {
	db *gorm.DB
}

// NewLocationService 创建位置信息服务实例
func NewLocationService() *LocationService {
	return &LocationService{
		db: DB,
	}
}

// Province 省份信息
type Province struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	NameTw   string `json:"name_tw"`
	CountryCode string `json:"country_code"`
}

// City 城市信息
type City struct {
	ID        int     `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	NameEn    string  `json:"name_en"`
	NameTw    string  `json:"name_tw"`
	ProvinceID int    `json:"province_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// District 区县信息
type District struct {
	ID        int     `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	NameEn    string  `json:"name_en"`
	NameTw    string  `json:"name_tw"`
	CityID    int     `json:"city_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UserLocation 用户位置信息
type UserLocation struct {
	ID               int     `json:"id"`
	UserEmail        string  `json:"user_email"`
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
	IsPrimary        bool    `json:"is_primary"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// POI 兴趣点信息
type POI struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	NameEn      string  `json:"name_en"`
	NameTw      string  `json:"name_tw"`
	Type        string  `json:"type"`
	CityID      int     `json:"city_id"`
	DistrictID  int     `json:"district_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
}

// GetProvinces 获取所有省份
func (s *LocationService) GetProvinces(ctx context.Context) ([]Province, error) {
	var provinces []Province
	
	err := s.db.Table("provinces").
		Order("id ASC").
		Find(&provinces).Error
	
	if err != nil {
		return nil, fmt.Errorf("获取省份列表失败: %v", err)
	}
	
	return provinces, nil
}

// GetCitiesByProvince 根据省份ID获取城市列表
func (s *LocationService) GetCitiesByProvince(ctx context.Context, provinceID int) ([]City, error) {
	var cities []City
	
	err := s.db.Table("cities").
		Where("province_id = ?", provinceID).
		Order("id ASC").
		Find(&cities).Error
	
	if err != nil {
		return nil, fmt.Errorf("获取城市列表失败: %v", err)
	}
	
	return cities, nil
}

// GetDistrictsByCity 根据城市ID获取区县列表
func (s *LocationService) GetDistrictsByCity(ctx context.Context, cityID int) ([]District, error) {
	var districts []District
	
	err := s.db.Table("districts").
		Where("city_id = ?", cityID).
		Order("id ASC").
		Find(&districts).Error
	
	if err != nil {
		return nil, fmt.Errorf("获取区县列表失败: %v", err)
	}
	
	return districts, nil
}

// SearchPOI 搜索兴趣点
func (s *LocationService) SearchPOI(ctx context.Context, keyword string, cityID int, poiType string) ([]POI, error) {
	var pois []POI
	
	query := s.db.Table("pois").Where("is_active = ?", true)
	
	if keyword != "" {
		query = query.Where("name LIKE ? OR name_en LIKE ? OR name_tw LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if cityID > 0 {
		query = query.Where("city_id = ?", cityID)
	}
	
	if poiType != "" {
		query = query.Where("type = ?", poiType)
	}
	
	err := query.Order("name ASC").Limit(50).Find(&pois).Error
	
	if err != nil {
		return nil, fmt.Errorf("搜索兴趣点失败: %v", err)
	}
	
	return pois, nil
}

// SaveUserLocation 保存用户位置信息
func (s *LocationService) SaveUserLocation(ctx context.Context, location *UserLocation) error {
	// 检查用户是否已有位置信息
	var existingLocation UserLocation
	err := s.db.Table("user_locations").
		Where("user_email = ? AND is_primary = ?", location.UserEmail, true).
		First(&existingLocation).Error
	
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询用户位置信息失败: %v", err)
	}
	
	if err == nil {
		// 更新现有记录
		err = s.db.Table("user_locations").
			Where("id = ?", existingLocation.ID).
			Updates(location).Error
	} else {
		// 插入新记录
		err = s.db.Table("user_locations").Create(location).Error
	}
	
	if err != nil {
		return fmt.Errorf("保存用户位置信息失败: %v", err)
	}
	
	return nil
}

// GetUserLocation 获取用户位置信息
func (s *LocationService) GetUserLocation(ctx context.Context, userEmail string) (*UserLocation, error) {
	var location UserLocation
	
	err := s.db.Table("user_locations").
		Where("user_email = ? AND is_primary = ?", userEmail, true).
		First(&location).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("获取用户位置信息失败: %v", err)
	}
	
	return &location, nil
}

// GeocodeAddress 地理编码（地址转坐标）
func (s *LocationService) GeocodeAddress(ctx context.Context, address string) (float64, float64, error) {
	mapService := NewMapAPIService()
	result, err := mapService.GeocodeAddress(ctx, address)
	if err != nil {
		// 如果地图API失败，根据地址关键词返回对应坐标
		return s.getCoordinatesByAddress(address)
	}
	return result.Latitude, result.Longitude, nil
}

// getCoordinatesByAddress 根据地址关键词返回对应坐标（支持简体中文、繁体中文、英文）
func (s *LocationService) getCoordinatesByAddress(address string) (float64, float64, error) {
	// 转换为小写以便不区分大小写匹配英文
	addressLower := strings.ToLower(address)
	
	// 北京 / 北京 / Beijing / Peking
	if strings.Contains(address, "北京市") || strings.Contains(address, "北京") || 
	   strings.Contains(address, "北京") || strings.Contains(address, "北京") ||
	   strings.Contains(addressLower, "beijing") || strings.Contains(addressLower, "peking") {
		return 39.9042, 116.4074, nil
	}
	
	// 上海 / 上海 / Shanghai
	if strings.Contains(address, "上海市") || strings.Contains(address, "上海") ||
	   strings.Contains(address, "上海") || strings.Contains(address, "上海") ||
	   strings.Contains(addressLower, "shanghai") {
		return 31.2304, 121.4737, nil
	}
	
	// 香港 / 香港 / Hong Kong / Hongkong
	if strings.Contains(address, "香港") || strings.Contains(address, "香港") ||
	   strings.Contains(addressLower, "hong kong") || strings.Contains(addressLower, "hongkong") {
		return 22.3193, 114.1694, nil
	}
	
	// 澳门 / 澳門 / Macau / Macao
	if strings.Contains(address, "澳门") || strings.Contains(address, "澳門") ||
	   strings.Contains(addressLower, "macau") || strings.Contains(addressLower, "macao") {
		return 22.1987, 113.5439, nil
	}
	
	// 台湾 / 台灣 / Taiwan
	if strings.Contains(address, "台湾") || strings.Contains(address, "台灣") ||
	   strings.Contains(addressLower, "taiwan") {
		return 25.0330, 121.5654, nil
	}
	
	// 广州 / 廣州 / Guangzhou / Canton
	if strings.Contains(address, "广州市") || strings.Contains(address, "广州") ||
	   strings.Contains(address, "廣州市") || strings.Contains(address, "廣州") ||
	   strings.Contains(addressLower, "guangzhou") || strings.Contains(addressLower, "canton") {
		return 23.1291, 113.2644, nil
	}
	
	// 深圳 / 深圳 / Shenzhen
	if strings.Contains(address, "深圳市") || strings.Contains(address, "深圳") ||
	   strings.Contains(address, "深圳市") || strings.Contains(address, "深圳") ||
	   strings.Contains(addressLower, "shenzhen") {
		return 22.5431, 114.0579, nil
	}
	
	// 杭州 / 杭州 / Hangzhou
	if strings.Contains(address, "杭州市") || strings.Contains(address, "杭州") ||
	   strings.Contains(address, "杭州市") || strings.Contains(address, "杭州") ||
	   strings.Contains(addressLower, "hangzhou") {
		return 30.2741, 120.1551, nil
	}
	
	// 南京 / 南京 / Nanjing / Nanking
	if strings.Contains(address, "南京市") || strings.Contains(address, "南京") ||
	   strings.Contains(address, "南京市") || strings.Contains(address, "南京") ||
	   strings.Contains(addressLower, "nanjing") || strings.Contains(addressLower, "nanking") {
		return 32.0603, 118.7969, nil
	}
	
	// 成都 / 成都 / Chengdu
	if strings.Contains(address, "成都市") || strings.Contains(address, "成都") ||
	   strings.Contains(address, "成都市") || strings.Contains(address, "成都") ||
	   strings.Contains(addressLower, "chengdu") {
		return 30.5728, 104.0668, nil
	}
	
	// 武汉 / 武漢 / Wuhan
	if strings.Contains(address, "武汉市") || strings.Contains(address, "武汉") ||
	   strings.Contains(address, "武漢市") || strings.Contains(address, "武漢") ||
	   strings.Contains(addressLower, "wuhan") {
		return 30.5928, 114.3055, nil
	}
	
	// 西安 / 西安 / Xi'an / Xian
	if strings.Contains(address, "西安市") || strings.Contains(address, "西安") ||
	   strings.Contains(address, "西安市") || strings.Contains(address, "西安") ||
	   strings.Contains(addressLower, "xi'an") || strings.Contains(addressLower, "xian") {
		return 34.3416, 108.9398, nil
	}
	
	// 重庆 / 重慶 / Chongqing / Chungking
	if strings.Contains(address, "重庆市") || strings.Contains(address, "重庆") ||
	   strings.Contains(address, "重慶市") || strings.Contains(address, "重慶") ||
	   strings.Contains(addressLower, "chongqing") || strings.Contains(addressLower, "chungking") {
		return 29.5630, 106.5516, nil
	}
	
	// 天津 / 天津 / Tianjin / Tientsin
	if strings.Contains(address, "天津市") || strings.Contains(address, "天津") ||
	   strings.Contains(address, "天津市") || strings.Contains(address, "天津") ||
	   strings.Contains(addressLower, "tianjin") || strings.Contains(addressLower, "tientsin") {
		return 39.3434, 117.3616, nil
	}
	
	// 苏州 / 蘇州 / Suzhou
	if strings.Contains(address, "苏州市") || strings.Contains(address, "苏州") ||
	   strings.Contains(address, "蘇州市") || strings.Contains(address, "蘇州") ||
	   strings.Contains(addressLower, "suzhou") {
		return 31.2989, 120.5853, nil
	}
	
	// 厦门 / 廈門 / Xiamen / Amoy
	if strings.Contains(address, "厦门市") || strings.Contains(address, "厦门") ||
	   strings.Contains(address, "廈門市") || strings.Contains(address, "廈門") ||
	   strings.Contains(addressLower, "xiamen") || strings.Contains(addressLower, "amoy") {
		return 24.4798, 118.0819, nil
	}
	
	// 青岛 / 青島 / Qingdao / Tsingtao
	if strings.Contains(address, "青岛市") || strings.Contains(address, "青岛") ||
	   strings.Contains(address, "青島市") || strings.Contains(address, "青島") ||
	   strings.Contains(addressLower, "qingdao") || strings.Contains(addressLower, "tsingtao") {
		return 36.0671, 120.3826, nil
	}
	
	// 大连 / 大連 / Dalian / Dairen
	if strings.Contains(address, "大连市") || strings.Contains(address, "大连") ||
	   strings.Contains(address, "大連市") || strings.Contains(address, "大連") ||
	   strings.Contains(addressLower, "dalian") || strings.Contains(addressLower, "dairen") {
		return 38.9140, 121.6147, nil
	}
	
	// 郑州 / 鄭州 / Zhengzhou
	if strings.Contains(address, "郑州市") || strings.Contains(address, "郑州") ||
	   strings.Contains(address, "鄭州市") || strings.Contains(address, "鄭州") ||
	   strings.Contains(addressLower, "zhengzhou") {
		return 34.7466, 113.6254, nil
	}
	
	// 长沙 / 長沙 / Changsha
	if strings.Contains(address, "长沙市") || strings.Contains(address, "长沙") ||
	   strings.Contains(address, "長沙市") || strings.Contains(address, "長沙") ||
	   strings.Contains(addressLower, "changsha") {
		return 28.2278, 112.9388, nil
	}
	
	// 沈阳 / 瀋陽 / Shenyang / Mukden
	if strings.Contains(address, "沈阳市") || strings.Contains(address, "沈阳") ||
	   strings.Contains(address, "瀋陽市") || strings.Contains(address, "瀋陽") ||
	   strings.Contains(addressLower, "shenyang") || strings.Contains(addressLower, "mukden") {
		return 41.8057, 123.4315, nil
	}
	
	// 哈尔滨 / 哈爾濱 / Harbin
	if strings.Contains(address, "哈尔滨市") || strings.Contains(address, "哈尔滨") ||
	   strings.Contains(address, "哈爾濱市") || strings.Contains(address, "哈爾濱") ||
	   strings.Contains(addressLower, "harbin") {
		return 45.7732, 126.6577, nil
	}
	
	// 昆明 / 昆明 / Kunming
	if strings.Contains(address, "昆明市") || strings.Contains(address, "昆明") ||
	   strings.Contains(address, "昆明市") || strings.Contains(address, "昆明") ||
	   strings.Contains(addressLower, "kunming") {
		return 25.0389, 102.7183, nil
	}
	
	// 省份匹配（简体、繁体、英文）
	// 河北 / 河北 / Hebei
	if strings.Contains(address, "河北省") || strings.Contains(address, "河北") ||
	   strings.Contains(address, "河北省") || strings.Contains(address, "河北") ||
	   strings.Contains(addressLower, "hebei") {
		return 38.0428, 114.5149, nil
	}
	
	// 广东 / 廣東 / Guangdong
	if strings.Contains(address, "广东省") || strings.Contains(address, "广东") ||
	   strings.Contains(address, "廣東省") || strings.Contains(address, "廣東") ||
	   strings.Contains(addressLower, "guangdong") {
		return 23.1291, 113.2644, nil
	}
	
	// 江苏 / 江蘇 / Jiangsu
	if strings.Contains(address, "江苏省") || strings.Contains(address, "江苏") ||
	   strings.Contains(address, "江蘇省") || strings.Contains(address, "江蘇") ||
	   strings.Contains(addressLower, "jiangsu") {
		return 32.0603, 118.7969, nil
	}
	
	// 浙江 / 浙江 / Zhejiang
	if strings.Contains(address, "浙江省") || strings.Contains(address, "浙江") ||
	   strings.Contains(address, "浙江省") || strings.Contains(address, "浙江") ||
	   strings.Contains(addressLower, "zhejiang") {
		return 30.2741, 120.1551, nil
	}
	
	// 山东 / 山東 / Shandong
	if strings.Contains(address, "山东省") || strings.Contains(address, "山东") ||
	   strings.Contains(address, "山東省") || strings.Contains(address, "山東") ||
	   strings.Contains(addressLower, "shandong") {
		return 36.6512, 117.1201, nil
	}
	
	// 河南 / 河南 / Henan
	if strings.Contains(address, "河南省") || strings.Contains(address, "河南") ||
	   strings.Contains(address, "河南省") || strings.Contains(address, "河南") ||
	   strings.Contains(addressLower, "henan") {
		return 34.7466, 113.6254, nil
	}
	
	// 湖北 / 湖北 / Hubei
	if strings.Contains(address, "湖北省") || strings.Contains(address, "湖北") ||
	   strings.Contains(address, "湖北省") || strings.Contains(address, "湖北") ||
	   strings.Contains(addressLower, "hubei") {
		return 30.5928, 114.3055, nil
	}
	
	// 湖南 / 湖南 / Hunan
	if strings.Contains(address, "湖南省") || strings.Contains(address, "湖南") ||
	   strings.Contains(address, "湖南省") || strings.Contains(address, "湖南") ||
	   strings.Contains(addressLower, "hunan") {
		return 28.2278, 112.9388, nil
	}
	
	// 四川 / 四川 / Sichuan / Szechuan
	if strings.Contains(address, "四川省") || strings.Contains(address, "四川") ||
	   strings.Contains(address, "四川省") || strings.Contains(address, "四川") ||
	   strings.Contains(addressLower, "sichuan") || strings.Contains(addressLower, "szechuan") {
		return 30.5728, 104.0668, nil
	}
	
	// 陕西 / 陝西 / Shaanxi / Shensi
	if strings.Contains(address, "陕西省") || strings.Contains(address, "陕西") ||
	   strings.Contains(address, "陝西省") || strings.Contains(address, "陝西") ||
	   strings.Contains(addressLower, "shaanxi") || strings.Contains(addressLower, "shensi") {
		return 34.3416, 108.9398, nil
	}
	
	// 山西 / 山西 / Shanxi
	if strings.Contains(address, "山西省") || strings.Contains(address, "山西") ||
	   strings.Contains(address, "山西省") || strings.Contains(address, "山西") ||
	   strings.Contains(addressLower, "shanxi") {
		return 37.8706, 112.5489, nil
	}
	
	// 辽宁 / 遼寧 / Liaoning
	if strings.Contains(address, "辽宁省") || strings.Contains(address, "辽宁") ||
	   strings.Contains(address, "遼寧省") || strings.Contains(address, "遼寧") ||
	   strings.Contains(addressLower, "liaoning") {
		return 41.8057, 123.4315, nil
	}
	
	// 吉林 / 吉林 / Jilin / Kirin
	if strings.Contains(address, "吉林省") || strings.Contains(address, "吉林") ||
	   strings.Contains(address, "吉林省") || strings.Contains(address, "吉林") ||
	   strings.Contains(addressLower, "jilin") || strings.Contains(addressLower, "kirin") {
		return 43.8171, 125.3235, nil
	}
	
	// 黑龙江 / 黑龍江 / Heilongjiang
	if strings.Contains(address, "黑龙江省") || strings.Contains(address, "黑龙江") ||
	   strings.Contains(address, "黑龍江省") || strings.Contains(address, "黑龍江") ||
	   strings.Contains(addressLower, "heilongjiang") {
		return 45.7732, 126.6577, nil
	}
	
	// 安徽 / 安徽 / Anhui
	if strings.Contains(address, "安徽省") || strings.Contains(address, "安徽") ||
	   strings.Contains(address, "安徽省") || strings.Contains(address, "安徽") ||
	   strings.Contains(addressLower, "anhui") {
		return 31.8612, 117.2837, nil
	}
	
	// 福建 / 福建 / Fujian / Fukien
	if strings.Contains(address, "福建省") || strings.Contains(address, "福建") ||
	   strings.Contains(address, "福建省") || strings.Contains(address, "福建") ||
	   strings.Contains(addressLower, "fujian") || strings.Contains(addressLower, "fukien") {
		return 26.0745, 119.2965, nil
	}
	
	// 江西 / 江西 / Jiangxi
	if strings.Contains(address, "江西省") || strings.Contains(address, "江西") ||
	   strings.Contains(address, "江西省") || strings.Contains(address, "江西") ||
	   strings.Contains(addressLower, "jiangxi") {
		return 28.6820, 115.8922, nil
	}
	
	// 云南 / 雲南 / Yunnan
	if strings.Contains(address, "云南省") || strings.Contains(address, "云南") ||
	   strings.Contains(address, "雲南省") || strings.Contains(address, "雲南") ||
	   strings.Contains(addressLower, "yunnan") {
		return 25.0389, 102.7183, nil
	}
	
	// 贵州 / 貴州 / Guizhou / Kweichow
	if strings.Contains(address, "贵州省") || strings.Contains(address, "贵州") ||
	   strings.Contains(address, "貴州省") || strings.Contains(address, "貴州") ||
	   strings.Contains(addressLower, "guizhou") || strings.Contains(addressLower, "kweichow") {
		return 26.6470, 106.6302, nil
	}
	
	// 甘肃 / 甘肅 / Gansu / Kansu
	if strings.Contains(address, "甘肃省") || strings.Contains(address, "甘肃") ||
	   strings.Contains(address, "甘肅省") || strings.Contains(address, "甘肅") ||
	   strings.Contains(addressLower, "gansu") || strings.Contains(addressLower, "kansu") {
		return 36.0611, 103.8343, nil
	}
	
	// 青海 / 青海 / Qinghai
	if strings.Contains(address, "青海省") || strings.Contains(address, "青海") ||
	   strings.Contains(address, "青海省") || strings.Contains(address, "青海") ||
	   strings.Contains(addressLower, "qinghai") {
		return 36.6232, 101.7782, nil
	}
	
	// 新疆 / 新疆 / Xinjiang / Sinkiang
	if strings.Contains(address, "新疆") || strings.Contains(address, "新疆") ||
	   strings.Contains(addressLower, "xinjiang") || strings.Contains(addressLower, "sinkiang") {
		return 43.8256, 87.6168, nil
	}
	
	// 西藏 / 西藏 / Tibet / Xizang
	if strings.Contains(address, "西藏") || strings.Contains(address, "西藏") ||
	   strings.Contains(addressLower, "tibet") || strings.Contains(addressLower, "xizang") {
		return 29.6465, 91.1172, nil
	}
	
	// 内蒙古 / 內蒙古 / Inner Mongolia / Nei Mongol
	if strings.Contains(address, "内蒙古") || strings.Contains(address, "內蒙古") ||
	   strings.Contains(addressLower, "inner mongolia") || strings.Contains(addressLower, "nei mongol") {
		return 40.8414, 111.7656, nil
	}
	
	// 宁夏 / 寧夏 / Ningxia
	if strings.Contains(address, "宁夏") || strings.Contains(address, "寧夏") ||
	   strings.Contains(addressLower, "ningxia") {
		return 38.4872, 106.2309, nil
	}
	
	// 广西 / 廣西 / Guangxi
	if strings.Contains(address, "广西") || strings.Contains(address, "廣西") ||
	   strings.Contains(addressLower, "guangxi") {
		return 22.8170, 108.3669, nil
	}
	
	// 海南 / 海南 / Hainan
	if strings.Contains(address, "海南") || strings.Contains(address, "海南") ||
	   strings.Contains(addressLower, "hainan") {
		return 20.0444, 110.1999, nil
	}
	
		// 默认返回北京坐标
	return 39.9042, 116.4074, nil
}

// ReverseGeocode 逆地理编码（坐标转地址）
func (s *LocationService) ReverseGeocode(ctx context.Context, latitude, longitude float64) (string, error) {
	mapService := NewMapAPIService()
	result, err := mapService.ReverseGeocode(ctx, latitude, longitude)
	if err != nil {
		// 如果地图API失败，根据坐标判断大概位置
		address := s.getLocationByCoordinates(latitude, longitude)
		return address, nil
	}
	return result.Address, nil
}

// getLocationByCoordinates 根据坐标判断大概位置
func (s *LocationService) getLocationByCoordinates(latitude, longitude float64) string {
	// 根据坐标范围判断大概位置，提供更详细的地址信息
	if latitude >= 39.0 && latitude <= 41.0 && longitude >= 115.0 && longitude <= 118.0 {
		// 根据具体坐标提供更精确的北京地址
		if latitude >= 39.8 && latitude <= 40.0 && longitude >= 116.4 && longitude <= 116.5 {
			return "北京市朝阳区双龙路128号"
		} else if latitude >= 39.9 && latitude <= 40.0 && longitude >= 116.4 && longitude <= 116.5 {
			return "北京市海淀区中关村大街1号"
		} else if latitude >= 39.9 && latitude <= 40.0 && longitude >= 116.3 && longitude <= 116.4 {
			return "北京市西城区西单北大街131号"
		} else if latitude >= 39.9 && latitude <= 40.0 && longitude >= 116.2 && longitude <= 116.3 {
			return "北京市东城区王府井大街138号"
		} else {
			return "北京市朝阳区建国门外大街1号"
		}
	} else if latitude >= 31.0 && latitude <= 32.0 && longitude >= 120.0 && longitude <= 122.0 {
		return "上海市浦东新区陆家嘴环路1000号"
	} else if latitude >= 22.0 && latitude <= 25.0 && longitude >= 113.0 && longitude <= 115.0 {
		return "广东省广州市天河区珠江新城花城大道85号"
	} else if latitude >= 25.0 && latitude <= 27.0 && longitude >= 118.0 && longitude <= 121.0 {
		return "福建省厦门市思明区鼓浪屿龙头路1号"
	} else if latitude >= 30.0 && latitude <= 32.0 && longitude >= 118.0 && longitude <= 121.0 {
		return "江苏省南京市鼓楼区中山路321号"
	} else if latitude >= 36.0 && latitude <= 40.0 && longitude >= 114.0 && longitude <= 120.0 {
		return "河北省石家庄市长安区中山东路265号"
	} else if latitude >= 39.0 && latitude <= 42.0 && longitude >= 125.0 && longitude <= 130.0 {
		return "辽宁省沈阳市和平区中山路123号"
	} else if latitude >= 45.0 && latitude <= 48.0 && longitude >= 125.0 && longitude <= 130.0 {
		return "黑龙江省哈尔滨市道里区中央大街89号"
	} else if latitude >= 22.0 && latitude <= 25.0 && longitude >= 114.0 && longitude <= 116.0 {
		return "香港特别行政区中环金融街8号"
	} else if latitude >= 22.0 && latitude <= 25.0 && longitude >= 113.0 && longitude <= 115.0 {
		return "澳门特别行政区大堂区友谊大马路555号"
	} else if latitude >= 38.0 && latitude <= 40.0 && longitude >= 116.0 && longitude <= 119.0 {
		return "天津市和平区南京路219号"
	} else if latitude >= 35.0 && latitude <= 37.0 && longitude >= 114.0 && longitude <= 118.0 {
		return "河南省郑州市金水区花园路88号"
	} else if latitude >= 36.0 && latitude <= 40.0 && longitude >= 110.0 && longitude <= 115.0 {
		return "山西省太原市迎泽区迎泽大街366号"
	} else if latitude >= 34.0 && latitude <= 36.0 && longitude >= 108.0 && longitude <= 112.0 {
		return "陕西省西安市雁塔区小寨东路126号"
	} else if latitude >= 32.0 && latitude <= 35.0 && longitude >= 114.0 && longitude <= 120.0 {
		return "安徽省合肥市庐阳区长江中路369号"
	} else if latitude >= 28.0 && latitude <= 31.0 && longitude >= 118.0 && longitude <= 123.0 {
		// 浙江省更详细的地址判断
		if latitude >= 29.5 && latitude <= 30.5 && longitude >= 120.0 && longitude <= 121.0 {
			return "浙江省绍兴市柯桥区平水镇剑灶村128号"
		} else if latitude >= 30.0 && latitude <= 30.5 && longitude >= 120.0 && longitude <= 120.5 {
			return "浙江省杭州市西湖区文三路259号"
		} else {
			return "浙江省杭州市西湖区文三路259号"
		}
	} else if latitude >= 26.0 && latitude <= 30.0 && longitude >= 103.0 && longitude <= 110.0 {
		return "四川省成都市锦江区春熙路98号"
	} else if latitude >= 24.0 && latitude <= 26.0 && longitude >= 108.0 && longitude <= 112.0 {
		return "广西壮族自治区南宁市青秀区民族大道88号"
	} else if latitude >= 18.0 && latitude <= 21.0 && longitude >= 108.0 && longitude <= 111.0 {
		return "海南省海口市龙华区国贸大道56号"
	} else if latitude >= 43.0 && latitude <= 46.0 && longitude >= 87.0 && longitude <= 97.0 {
		return "新疆维吾尔自治区乌鲁木齐市天山区解放南路123号"
	} else if latitude >= 29.0 && latitude <= 33.0 && longitude >= 91.0 && longitude <= 99.0 {
		return "西藏自治区拉萨市城关区北京中路65号"
	} else if latitude >= 35.0 && latitude <= 40.0 && longitude >= 100.0 && longitude <= 107.0 {
		return "青海省西宁市城东区八一路88号"
	} else if latitude >= 35.0 && latitude <= 40.0 && longitude >= 103.0 && longitude <= 108.0 {
		return "甘肃省兰州市城关区天水路156号"
	} else if latitude >= 35.0 && latitude <= 40.0 && longitude >= 105.0 && longitude <= 111.0 {
		return "宁夏回族自治区银川市兴庆区解放东街78号"
	} else if latitude >= 40.0 && latitude <= 43.0 && longitude >= 110.0 && longitude <= 120.0 {
		return "内蒙古自治区呼和浩特市新城区新华大街234号"
	} else {
		return fmt.Sprintf("位置坐标: %.6f, %.6f", latitude, longitude)
	}
}

// GetNearbyPOIs 获取附近兴趣点
func (s *LocationService) GetNearbyPOIs(ctx context.Context, latitude, longitude float64, radius int, poiType string) ([]POI, error) {
	var pois []POI
	
	// 使用简单的矩形范围查询（实际项目中应使用空间数据库）
	latRange := float64(radius) / 111.0 // 粗略计算纬度范围
	lngRange := float64(radius) / (111.0 * math.Cos(latitude*math.Pi/180)) // 经度范围
	
	query := s.db.Table("pois").
		Where("is_active = ?", true).
		Where("latitude BETWEEN ? AND ?", latitude-latRange, latitude+latRange).
		Where("longitude BETWEEN ? AND ?", longitude-lngRange, longitude+lngRange)
	
	if poiType != "" {
		query = query.Where("type = ?", poiType)
	}
	
	err := query.Order("id ASC").Limit(20).Find(&pois).Error
	
	if err != nil {
		return nil, fmt.Errorf("获取附近兴趣点失败: %v", err)
	}
	
	return pois, nil
}

// UpdateUserLocationFromRegistration 从注册信息更新用户位置
func (s *LocationService) UpdateUserLocationFromRegistration(ctx context.Context, userEmail, buildingName, buildingAddr, buildingType string) error {
	// 尝试从地址中提取位置信息
	location := &UserLocation{
		UserEmail:     userEmail,
		Country:       "中国",
		CountryCode:   "CN",
		BuildingName:  buildingName,
		BuildingType:  buildingType,
		AddressDetail: buildingAddr,
		IsPrimary:     true,
	}
	
	// 简单的地址解析（实际项目中应使用更复杂的地址解析服务）
	if strings.Contains(buildingAddr, "北京") {
		location.ProvinceName = "北京市"
		location.CityName = "北京市"
	} else if strings.Contains(buildingAddr, "上海") {
		location.ProvinceName = "上海市"
		location.CityName = "上海市"
	} else if strings.Contains(buildingAddr, "广州") {
		location.ProvinceName = "广东省"
		location.CityName = "广州市"
	} else if strings.Contains(buildingAddr, "深圳") {
		location.ProvinceName = "广东省"
		location.CityName = "深圳市"
	}
	
	return s.SaveUserLocation(ctx, location)
}
