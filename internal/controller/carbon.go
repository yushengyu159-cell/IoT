package controller

import (
	"fabric-sdk/internal/service"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type CarbonController struct{}

var Carbon = new(CarbonController)

// GetCarbonOverview returns the overall carbon emission data
// GET /api/carbon/overview
// 对于belucksapi用户，尝试从LiveSense API获取真实数据
func (c *CarbonController) GetCarbonOverview(r *ghttp.Request) {
	type Overview struct {
		TotalEmission   float64 `json:"totalEmission"` // Total emissions in tons
		Unit            string  `json:"unit"`          // Unit (e.g., "tCO2e")
		ChangeRate      float64 `json:"changeRate"`    // Change rate vs last period (percentage)
		Trend           string  `json:"trend"`         // "up" or "down"
		LastUpdated     string  `json:"lastUpdated"`
		ComplianceScore int     `json:"complianceScore"` // 0-100
	}

	// 从请求中获取用户邮箱（支持多种方式）
	email := r.Get("email").String()
	if email == "" {
		email = r.GetHeader("X-User-Email")
	}
	if email == "" {
		// 尝试从Cookie获取
		cookieEmail := r.Cookie.Get("user_email")
		if cookieEmail != nil {
			email = cookieEmail.String()
		}
	}
	if email == "" {
		// 尝试从请求上下文获取（如果之前通过中间件设置）
		if ctxEmail := r.GetCtxVar("user_email"); ctxEmail != nil {
			email = ctxEmail.String()
		}
	}
	
	// 如果是belucksapi用户，尝试从LiveSense API获取真实数据
	// 支持多个用户，只要用户在数据库中的name字段为"belucksapi"
	if email != "" && isBelucksAPIUser(email) {
		// 尝试从多个传感器获取数据并聚合
		sensorIDs := []string{"27732", "27728", "27729"} // CO2, Temperature, Humidity
		
		var totalEmission float64
		var dataPointCount int
		var hasData bool
		
		// 根据Postman文件，不预设时间范围，先尝试不传时间参数获取数据
		// 如果返回数据，则从数据中的timestampUtc提取时间范围
		// 如果返回null，再尝试使用Postman文件中验证成功的时间范围
		
		// 尝试多个传感器
		for _, sensorID := range sensorIDs {
			// 策略1: 先不传时间参数，让API返回默认数据
			values, err := service.LiveSenseAPI.GetSensorValues(r.Context(), sensorID, "", "", "", "allValues", "")
			if err == nil && values != nil {
				valuesList := values.GetValuesList()
				if len(valuesList) > 0 && len(valuesList[0].DataPoints) > 0 {
					// 从返回的数据中提取时间范围
					dataPoints := valuesList[0].DataPoints
					for _, dp := range dataPoints {
						if val, ok := service.ParseValue(dp.Value); ok {
							totalEmission += val
							dataPointCount++
							hasData = true
						}
					}
					if hasData {
						// 提取时间范围信息
						firstTs := dataPoints[0].TimestampUtc
						lastTs := dataPoints[len(dataPoints)-1].TimestampUtc
						g.Log().Info(r.Context(), "✓ 成功从传感器", sensorID, "获取真实数据（无时间参数），数据点数量:", dataPointCount, "时间范围:", firstTs, "至", lastTs)
						break
					}
				}
			}
			
			// 策略2: 如果策略1失败，使用Postman文件中验证成功的时间范围
			if !hasData {
				g.Log().Info(r.Context(), "策略1失败，尝试策略2：使用Postman验证的时间范围，传感器:", sensorID)
				// Postman文件中验证成功的时间范围：from=1769327847000, to=1769414247000
				values, err := service.LiveSenseAPI.GetSensorValues(r.Context(), sensorID, "1769327847000", "1769414247000", "", "allValues", "")
				if err != nil {
					g.Log().Error(r.Context(), "策略2调用API失败，传感器:", sensorID, "错误:", err)
				} else if values != nil {
					// 检查values.Values是否为nil
					if values.Values == nil {
						g.Log().Warning(r.Context(), "策略2: values.Values为nil，传感器:", sensorID)
					} else {
						// 检查values.Values的类型
						switch v := values.Values.(type) {
						case map[string]interface{}:
							g.Log().Info(r.Context(), "策略2: values.Values是对象类型，传感器:", sensorID, "包含字段数:", len(v))
							if dataPoints, ok := v["dataPoints"].([]interface{}); ok {
								g.Log().Info(r.Context(), "策略2: 找到dataPoints数组，数量:", len(dataPoints))
							}
						case []interface{}:
							g.Log().Info(r.Context(), "策略2: values.Values是数组类型，传感器:", sensorID, "长度:", len(v))
						default:
							g.Log().Warning(r.Context(), "策略2: values.Values是未知类型:", fmt.Sprintf("%T", v), "传感器:", sensorID)
						}
					}
					
					valuesList := values.GetValuesList()
					g.Log().Info(r.Context(), "策略2: GetValuesList返回，传感器:", sensorID, "数据项数量:", len(valuesList))
					if len(valuesList) > 0 {
						for _, v := range valuesList {
							g.Log().Info(r.Context(), "处理数据项，传感器:", v.Sensor, "数据点数量:", len(v.DataPoints))
							for _, dp := range v.DataPoints {
								if val, ok := service.ParseValue(dp.Value); ok {
									totalEmission += val
									dataPointCount++
									hasData = true
								}
							}
						}
						if hasData {
							g.Log().Info(r.Context(), "✓ 成功从传感器", sensorID, "获取真实数据（使用Postman验证的时间范围），数据点数量:", dataPointCount)
							break
						} else {
							g.Log().Warning(r.Context(), "策略2获取到数据但无法解析数值，传感器:", sensorID)
						}
					} else {
						g.Log().Warning(r.Context(), "策略2返回空数据列表，传感器:", sensorID)
					}
				} else {
					g.Log().Warning(r.Context(), "策略2返回nil响应，传感器:", sensorID)
				}
			}
			
			if hasData {
				break
			}
		}
		
		// 如果获取到数据，转换为碳排放
		if hasData && dataPointCount > 0 {
			// 转换为碳排放单位
			avgValue := totalEmission / float64(dataPointCount)
			// 简化的转换公式（实际需要根据传感器类型和单位进行转换）
			carbonEmission := (avgValue * 0.001) * 30.0 // 假设转换系数
			
			data := Overview{
				TotalEmission:   carbonEmission,
				Unit:            "tCO2e",
				ChangeRate:      -5.2, // 可以从历史数据计算
				Trend:           "down",
				LastUpdated:     time.Now().Format("2006-01-02"),
				ComplianceScore: 92,
			}
			
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    200,
				Message: "Success (LiveSense Data)",
				Data:    data,
			})
			return
		}
		
		// 如果所有传感器都没有数据，记录详细日志但继续使用模拟数据
		g.Log().Warning(r.Context(), "LiveSense API所有传感器都返回空数据，已尝试", len(sensorIDs), "个传感器，使用模拟数据")
	}

	// 默认返回模拟数据
	data := Overview{
		TotalEmission:   1250.5,
		Unit:            "tCO2e",
		ChangeRate:      -5.2,
		Trend:           "down",
		LastUpdated:     time.Now().Format("2006-01-02"),
		ComplianceScore: 92,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    data,
	})
}

// isBelucksAPIUser 检查是否为belucksapi用户
func isBelucksAPIUser(email string) bool {
	db := service.GetDB()
	if db == nil {
		return false
	}
	
	var count int64
	result := db.Table("dids").Where("email = ? AND name = ?", email, "belucksapi").Count(&count)
	return result.Error == nil && count > 0
}

// GetCarbonTrackingTimeline returns detailed tracking data (timeline)
// GET /api/carbon/tracking/timeline
// 对于belucksapi用户，尝试从LiveSense API获取真实数据
func (c *CarbonController) GetCarbonTrackingTimeline(r *ghttp.Request) {
	// 从请求中获取用户邮箱（支持多种方式）
	email := r.Get("email").String()
	if email == "" {
		email = r.GetHeader("X-User-Email")
	}
	if email == "" {
		// 尝试从Cookie获取
		cookieEmail := r.Cookie.Get("user_email")
		if cookieEmail != nil {
			email = cookieEmail.String()
		}
	}
	if email == "" {
		// 尝试从请求上下文获取（如果之前通过中间件设置）
		if ctxEmail := r.GetCtxVar("user_email"); ctxEmail != nil {
			email = ctxEmail.String()
		}
	}
	
	// 如果是belucksapi用户，尝试从LiveSense API获取真实数据
	// 支持多个用户，只要用户在数据库中的name字段为"belucksapi"
	if email != "" && isBelucksAPIUser(email) {
		sensorID := r.Get("sensorId").String()
		if sensorID == "" {
			sensorID = "27732" // 默认传感器ID
		}
		
		// 根据Postman文件，不预设时间范围，先尝试不传时间参数获取数据
		// 如果返回数据，则从数据中的timestampUtc提取时间范围
		// 如果返回null，再尝试使用Postman文件中验证成功的时间范围
		
		var valuesList []service.SensorValueItem
		var lastErr error
		
		// 策略1: 先不传时间参数，让API返回默认数据
		values, err := service.LiveSenseAPI.GetSensorValues(r.Context(), sensorID, "", "", "", "allValues", "")
		if err == nil && values != nil {
			valuesList = values.GetValuesList()
			if len(valuesList) > 0 && len(valuesList[0].DataPoints) > 0 {
				// 从返回的数据中提取时间范围信息
				dataPoints := valuesList[0].DataPoints
				firstTs := dataPoints[0].TimestampUtc
				lastTs := dataPoints[len(dataPoints)-1].TimestampUtc
				g.Log().Info(r.Context(), "✓ 成功获取时间线真实数据（无时间参数），传感器:", sensorID, "数据点数量:", len(dataPoints), "时间范围:", firstTs, "至", lastTs)
			}
		} else if err != nil {
			lastErr = err
		}
		
		// 策略2: 如果策略1失败，使用Postman文件中验证成功的时间范围
		if len(valuesList) == 0 {
			values, err := service.LiveSenseAPI.GetSensorValues(r.Context(), sensorID, "1769327847000", "1769414247000", "", "allValues", "")
			if err == nil && values != nil {
				valuesList = values.GetValuesList()
				if len(valuesList) > 0 {
					dataPointCount := 0
					for _, v := range valuesList {
						dataPointCount += len(v.DataPoints)
					}
					g.Log().Info(r.Context(), "✓ 成功获取时间线真实数据（使用Postman验证的时间范围），传感器:", sensorID, "数据点数量:", dataPointCount)
				}
			} else if err != nil {
				lastErr = err
			}
		}
		
		if len(valuesList) > 0 {
			// 将传感器数据转换为时间线格式
			var labels []string
			var emissionValues []float64
			var baseline []float64
			var timeRangeStr string
			
			// 收集所有数据点并按时间排序
			type DataPointWithTime struct {
				Value     float64
				Timestamp time.Time
			}
			
			var allDataPoints []DataPointWithTime
			for _, v := range valuesList {
				for _, dp := range v.DataPoints {
					if val, ok := service.ParseValue(dp.Value); ok {
						if timestamp, err := strconv.ParseInt(dp.TimestampUtc, 10, 64); err == nil {
							t := time.Unix(timestamp/1000, 0)
							allDataPoints = append(allDataPoints, DataPointWithTime{
								Value:     val,
								Timestamp: t,
							})
						}
					}
				}
			}
			
			// 如果没有数据点，返回空数据
			if len(allDataPoints) == 0 {
				g.Log().Warning(r.Context(), "传感器数据点为空")
				timeRangeStr = "无数据"
			} else {
				// 按时间排序
				sort.Slice(allDataPoints, func(i, j int) bool {
					return allDataPoints[i].Timestamp.Before(allDataPoints[j].Timestamp)
				})
				
				// 确定时间范围（从第一个数据点到最后一个数据点）
				firstTime := allDataPoints[0].Timestamp
				lastTime := allDataPoints[len(allDataPoints)-1].Timestamp
				timeRangeStr = fmt.Sprintf("%s 至 %s", firstTime.Format("2006-01-02 15:04"), lastTime.Format("2006-01-02 15:04"))
				
				// 计算时间跨度（天数）
				timeSpan := lastTime.Sub(firstTime)
				daysDiff := int(timeSpan.Hours() / 24)
				
				// 根据时间跨度选择合适的分组粒度
				var timeFormat string
				var timeKeyFunc func(time.Time) string
				var timeIncrement func(time.Time) time.Time
				
				if daysDiff < 1 {
					// 小于1天，按小时分组（精确到小时）
					timeFormat = "01-02 15:04"
					timeKeyFunc = func(t time.Time) string {
						// 按小时分组，分钟设为0，确保与数据点匹配
						return t.Format("2006-01-02 15:00")
					}
					timeIncrement = func(t time.Time) time.Time {
						return t.Add(1 * time.Hour)
					}
				} else if daysDiff < 7 {
					// 小于7天，按小时分组（每3小时一组）
					timeFormat = "01-02 15:04"
					timeKeyFunc = func(t time.Time) string {
						// 按3小时分组
						hour := t.Hour() / 3 * 3
						return t.Format("2006-01-02") + fmt.Sprintf(" %02d:00", hour)
					}
					timeIncrement = func(t time.Time) time.Time {
						return t.Add(3 * time.Hour)
					}
				} else if daysDiff < 30 {
					// 小于30天，按天分组
					timeFormat = "01-02"
					timeKeyFunc = func(t time.Time) string {
						return t.Format("2006-01-02")
					}
					timeIncrement = func(t time.Time) time.Time {
						return t.AddDate(0, 0, 1)
					}
				} else {
					// 大于30天，按月分组
					timeFormat = "2006-01"
					timeKeyFunc = func(t time.Time) string {
						return t.Format("2006-01")
					}
					timeIncrement = func(t time.Time) time.Time {
						return t.AddDate(0, 1, 0)
					}
				}
				
				// 按选择的时间粒度分组数据
				dataPointsByTime := make(map[string][]float64)
				for _, dp := range allDataPoints {
					// 使用相同的时间键函数，确保分组键一致
					// 对于小时级别，需要将分钟设为0
					alignedTime := dp.Timestamp
					if daysDiff < 1 {
						// 按小时对齐
						alignedTime = time.Date(dp.Timestamp.Year(), dp.Timestamp.Month(), dp.Timestamp.Day(), dp.Timestamp.Hour(), 0, 0, 0, dp.Timestamp.Location())
					}
					timeKey := timeKeyFunc(alignedTime)
					dataPointsByTime[timeKey] = append(dataPointsByTime[timeKey], dp.Value)
				}
				
				// 添加调试日志
				g.Log().Info(r.Context(), "数据分组统计: 总数据点=", len(allDataPoints), ", 分组数=", len(dataPointsByTime))
				if len(dataPointsByTime) > 0 {
					// 显示前几个分组键和每个分组的数据点数量
					count := 0
					for k, v := range dataPointsByTime {
						if count < 5 {
							g.Log().Info(r.Context(), "分组键:", k, ", 数据点数量:", len(v))
							count++
						}
					}
	}

				// 计算基准值（所有数据的平均值）
				totalSum := 0.0
				totalCount := 0
				for _, points := range dataPointsByTime {
					for _, p := range points {
						totalSum += p
						totalCount++
					}
				}
				avgBaseline := 0.0
				if totalCount > 0 {
					avgBaseline = (totalSum / float64(totalCount) * 0.001) * 30.0 // 转换为碳排放
				} else {
					avgBaseline = 110.0 // 默认基准
				}
				
				// 生成时间序列（从第一个时间点到最后一个时间点）
				currentTime := firstTime
				// 对齐到时间粒度的起始点
				if daysDiff < 1 {
					// 按小时对齐
					currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), 0, 0, 0, currentTime.Location())
				} else if daysDiff < 7 {
					// 按小时对齐
					currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), 0, 0, 0, currentTime.Location())
				} else if daysDiff < 30 {
					// 按天对齐
					currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
				} else {
					// 按月对齐
					currentTime = time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, currentTime.Location())
				}
				
				endTime := lastTime
				if daysDiff < 30 {
					endTime = endTime.Add(24 * time.Hour) // 确保包含最后一天
				}
				
				// 限制最大数据点数量（避免图表过于密集）
				maxDataPoints := 50
				pointCount := 0
				
				// 生成时间标签和数据
				for (currentTime.Before(endTime) || currentTime.Equal(endTime)) && pointCount < maxDataPoints {
					timeKey := timeKeyFunc(currentTime)
					labels = append(labels, currentTime.Format(timeFormat))
					
					if points, ok := dataPointsByTime[timeKey]; ok && len(points) > 0 {
						// 计算该时间段的统计值
						sum := 0.0
						max := points[0]
						min := points[0]
						for _, p := range points {
							sum += p
							if p > max {
								max = p
							}
							if p < min {
								min = p
							}
						}
						avg := sum / float64(len(points))
						
						// CO2转换为碳排放（简化转换）
						carbonValue := (avg * 0.001) * 30.0
						emissionValues = append(emissionValues, carbonValue)
					} else {
						emissionValues = append(emissionValues, 0)
					}
					
					// 基准线使用计算出的平均值
					baseline = append(baseline, avgBaseline)
					
					// 移动到下一个时间点
					currentTime = timeIncrement(currentTime)
					pointCount++
				}
			}
			
			// 如果标签为空，生成默认的最近12个月
			if len(labels) == 0 {
				timeRangeStr = "无数据"
				for i := 11; i >= 0; i-- {
					date := time.Now().AddDate(0, -i, 0).Format("2006-01")
					labels = append(labels, date)
					emissionValues = append(emissionValues, 0)
					baseline = append(baseline, 110.0)
				}
			}
			
			data := g.Map{
				"labels":      labels,
				"values":      emissionValues,
				"baseline":    baseline,
				"dataSource":  "LiveSense API",
				"dataPoints":  len(allDataPoints),
				"timeRange":   timeRangeStr,
			}
			
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    200,
				Message: "Success (LiveSense Data)",
				Data:    data,
			})
			return
	}

		// 如果获取LiveSense数据失败，记录日志但继续使用模拟数据
		if len(valuesList) == 0 {
			if lastErr != nil {
				g.Log().Warning(r.Context(), "无法从LiveSense API获取时间线数据，使用模拟数据:", lastErr)
			} else {
				g.Log().Warning(r.Context(), "LiveSense API返回空时间线数据，使用模拟数据")
			}
		}
	}

	// 默认返回模拟数据
	var labels []string
	var values []float64
	var baseline []float64

	baseEmission := 100.0

	// Handle different periods if needed, for now default to 1Y/Monthly
	for i := 11; i >= 0; i-- {
		date := time.Now().AddDate(0, -i, 0).Format("2006-01")
		labels = append(labels, date)
		values = append(values, baseEmission+float64(rand.Intn(20)-10))
		baseline = append(baseline, 110.0)
	}

	data := g.Map{
		"labels":   labels,
		"values":   values,
		"baseline": baseline,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    data,
	})
}

// GetCarbonSourceDecomposition returns decomposition of carbon sources
// GET /api/carbon/tracking/source
func (c *CarbonController) GetCarbonSourceDecomposition(r *ghttp.Request) {
	type SourceData struct {
		Name       string  `json:"name"`
		Value      float64 `json:"value"`      // Emission value
		Percentage float64 `json:"percentage"` // Percentage of total
		Benchmark  float64 `json:"benchmark"`  // Industry average percentage
	}

	sources := []SourceData{
		{Name: "Energy", Value: 600.2, Percentage: 48.0, Benchmark: 45.0},
		{Name: "Materials", Value: 312.6, Percentage: 25.0, Benchmark: 30.0},
		{Name: "Transport", Value: 187.5, Percentage: 15.0, Benchmark: 15.0},
		{Name: "Waste", Value: 100.0, Percentage: 8.0, Benchmark: 5.0},
		{Name: "Water", Value: 50.2, Percentage: 4.0, Benchmark: 5.0},
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    sources,
	})
}

// GetCarbonVerification returns data verification status
// GET /api/carbon/verification
func (c *CarbonController) GetCarbonVerification(r *ghttp.Request) {
	type VerificationItem struct {
		Type       string `json:"type"`     // e.g., "Energy Data", "Material Data"
		Status     string `json:"status"`   // "verified", "pending", "failed"
		Method     string `json:"method"`   // "Automatic", "Third-party"
		Verifier   string `json:"verifier"` // Institution name or "System"
		Date       string `json:"date"`
		ValidUntil string `json:"validUntil"`
		ReportLink string `json:"reportLink"`
	}

	statusList := []VerificationItem{
		{
			Type:       "Energy Consumption 2024",
			Status:     "verified",
			Method:     "Automatic",
			Verifier:   "System IoT Monitor",
			Date:       "2024-11-30",
			ValidUntil: "2025-11-30",
			ReportLink: "#",
		},
		{
			Type:       "Building Materials Audit",
			Status:     "verified",
			Method:     "Third-party",
			Verifier:   "SGS Certification",
			Date:       "2024-10-15",
			ValidUntil: "2025-10-15",
			ReportLink: "#",
		},
		{
			Type:       "Waste Disposal Q4",
			Status:     "pending",
			Method:     "Third-party",
			Verifier:   "-",
			Date:       "-",
			ValidUntil: "-",
			ReportLink: "",
		},
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    statusList,
	})
}

// GetCarbonAnalysis returns analysis data (comparison and prediction)
// GET /api/carbon/analysis
func (c *CarbonController) GetCarbonAnalysis(r *ghttp.Request) {
	// Mock prediction data
	labels := []string{}
	predicted := []float64{}
	upper := []float64{}
	lower := []float64{}

	baseVal := 100.0
	for i := 1; i <= 6; i++ {
		date := time.Now().AddDate(0, i, 0).Format("2006-01")
		labels = append(labels, date)
		val := baseVal - float64(i)*2 // Declining trend
		predicted = append(predicted, val)
		upper = append(upper, val+5.0)
		lower = append(lower, val-5.0)
	}

	// Mock comparison data
	comparison := g.Map{
		"myBuilding":      1250.5,
		"industryAverage": 1400.0,
		"bestInClass":     1100.0,
	}

	data := g.Map{
		"labels":     labels,
		"predicted":  predicted,
		"upper":      upper,
		"lower":      lower,
		"comparison": comparison,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    data,
	})
}
