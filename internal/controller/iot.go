package controller

import (
	"encoding/json"
	"fabric-sdk/internal/model"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"gorm.io/gorm"
)

type IoTController struct{}

var IoT = new(IoTController)
var IoTDB *gorm.DB

// POST /api/esg-data
// Receives decoded IoT uplink data from LoRaWAN CS01-LB current sensor via ChirpStack
// Supports fPort=2 (realtime), fPort=3 (history), fPort=5 (device info), fPort=7 (simple log)
func (c *IoTController) ReceiveUplink(r *ghttp.Request) {
	body := r.GetBody()
	if body == nil || len(body) == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "empty body"})
		return
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "invalid JSON: " + err.Error()})
		return
	}

	// --- Extract deviceInfo (ChirpStack nested) ---
	deviceName, deviceEUI := "", ""
	if di := getMap(raw, "deviceInfo"); di != nil {
		deviceName = getString(di, "deviceName", "")
		deviceEUI = getString(di, "devEui", "")
	}
	if deviceName == "" {
		deviceName = getString(raw, "deviceName", "")
	}
	if deviceEUI == "" {
		deviceEUI = getString(raw, "devEUI", getString(raw, "devEui", ""))
	}

	fPort := getInt(raw, "fPort", 0)
	fCnt := getInt(raw, "fCnt", 0)

	// --- Extract signal info ---
	var rssiVal, snrVal float64
	frequency := 0
	if rxList, ok := raw["rxInfo"].([]interface{}); ok && len(rxList) > 0 {
		if rx, ok := rxList[0].(map[string]interface{}); ok {
			rssiVal = getFloat(rx, "rssi", 0)
			snrVal = getFloat(rx, "snr", 0)
		}
	}
	if tx := getMap(raw, "txInfo"); tx != nil {
		frequency = getInt(tx, "frequency", 0)
	}

	// --- Extract decoded object ---
	obj := getMap(raw, "object")
	decodedData := ""
	if obj != nil {
		if b, err := json.Marshal(obj); err == nil {
			decodedData = string(b)
		}
	}

	// --- Parse by fPort ---
	record := model.IoTData{
		DeviceName:  deviceName,
		DeviceEUI:   deviceEUI,
		FPort:       fPort,
		FCnt:        fCnt,
		RSSI:        rssiVal,
		SNR:         snrVal,
		Frequency:   frequency,
		RawPayload:  string(body),
		DecodedData: decodedData,
		ReceivedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	switch fPort {
	case 2:
		parseFPort2(obj, &record)
		record.UplinkType = "realtime"
	case 3:
		record.UplinkType = "history"
	case 5:
		parseFPort5(obj, &record)
		record.UplinkType = "device_info"
	case 7:
		record.UplinkType = "simple_log"
	default:
		record.UplinkType = "unknown"
	}

	// --- Save ---
	if db := getIoTDB(); db != nil {
		if err := db.Create(&record).Error; err != nil {
			g.Log().Errorf(r.Context(), "[IoT] DB save failed: %v", err)
		} else {
			g.Log().Infof(r.Context(), "[IoT] Saved id=%d device=%s eui=%s fPort=%d type=%s batV=%.1f i1=%.2f i2=%.2f i3=%.2f i4=%.2f",
				record.ID, deviceName, deviceEUI, fPort, record.UplinkType, record.BatV,
				record.Current1A, record.Current2A, record.Current3A, record.Current4A)
		}
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "OK",
		Data: map[string]interface{}{
			"id":         record.ID,
			"deviceName": deviceName,
			"deviceEUI":  deviceEUI,
			"fPort":      fPort,
			"uplinkType": record.UplinkType,
			"batV":       record.BatV,
			"current1A":  record.Current1A,
			"current2A":  record.Current2A,
			"current3A":  record.Current3A,
			"current4A":  record.Current4A,
			"receivedAt": record.ReceivedAt,
		},
	})
}

// parseFPort2 parses realtime current sensor data (fPort=2)
func parseFPort2(obj map[string]interface{}, record *model.IoTData) {
	if obj == nil {
		return
	}
	record.NodeType = getString(obj, "Node_type", "")
	record.BatV = getFloat(obj, "BatV", 0)
	record.Current1A = getFloat(obj, "Current1_A", 0)
	record.Current2A = getFloat(obj, "Current2_A", 0)
	record.Current3A = getFloat(obj, "Current3_A", 0)
	record.Current4A = getFloat(obj, "Current4_A", 0)
	record.EXTITrigger = getString(obj, "EXTI_Trigger", "")
	record.EXTILevel = getString(obj, "EXTI_Level", "")
	record.Cur1LStatus = getString(obj, "Cur1L_status", "")
	record.Cur1HStatus = getString(obj, "Cur1H_status", "")
	record.Cur2LStatus = getString(obj, "Cur2L_status", "")
	record.Cur2HStatus = getString(obj, "Cur2H_status", "")
	record.Cur3LStatus = getString(obj, "Cur3L_status", "")
	record.Cur3HStatus = getString(obj, "Cur3H_status", "")
	record.Cur4LStatus = getString(obj, "Cur4L_status", "")
	record.Cur4HStatus = getString(obj, "Cur4H_status", "")

	// Cumulative current (extended payload, bytes.length==28 or 32)
	record.CurTotalMod = getInt(obj, "curtotal_mod", 0)
	if record.CurTotalMod > 0 {
		record.CurTotal1MAmin = getUint(obj, "curtotal1_mA_min", 0)
		record.CurTotal2MAmin = getUint(obj, "curtotal2_mA_min", 0)
		record.CurTotal3MAmin = getUint(obj, "curtotal3_mA_min", 0)
		record.CurTotal4MAmin = getUint(obj, "curtotal4_mA_min", 0)
	}
}

// parseFPort5 parses device info (fPort=5)
func parseFPort5(obj map[string]interface{}, record *model.IoTData) {
	if obj == nil {
		return
	}
	record.SensorModel = getString(obj, "SENSOR_MODEL", "")
	record.FirmwareVersion = getString(obj, "FIRMWARE_VERSION", "")
	record.FrequencyBand = getString(obj, "FREQUENCY_BAND", "")
	record.SubBand = getString(obj, "SUB_BAND", "")
	record.BatV = getFloat(obj, "BAT", 0)
	record.NodeType = "CS01-LB"
}

// GET /api/esg-data/list
func (c *IoTController) ListData(r *ghttp.Request) {
	deviceEUI := r.Get("deviceEUI").String()
	deviceName := r.Get("deviceName").String()
	fPort := r.Get("fPort").Int()
	uplinkType := r.Get("uplinkType").String()
	limit := r.Get("limit").Int()
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	db := getIoTDB()
	if db == nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "database not initialized"})
		return
	}

	query := db.Order("id DESC").Limit(int(limit))
	if deviceEUI != "" {
		query = query.Where("device_eui = ?", deviceEUI)
	}
	if deviceName != "" {
		query = query.Where("device_name = ?", deviceName)
	}
	if fPort > 0 {
		query = query.Where("f_port = ?", fPort)
	}
	if uplinkType != "" {
		query = query.Where("uplink_type = ?", uplinkType)
	}

	var records []model.IoTData
	if err := query.Find(&records).Error; err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "success", Data: records})
}

// --- helper functions ---

func getString(m map[string]interface{}, key string, fallback string) string {
	if v, ok := m[key]; ok && v != nil {
		s := fmt.Sprintf("%v", v)
		if s != "" {
			return s
		}
	}
	return fallback
}

func getFloat(m map[string]interface{}, key string, fallback float64) float64 {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return fallback
}

func getInt(m map[string]interface{}, key string, fallback int) int {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return fallback
}

func getUint(m map[string]interface{}, key string, fallback uint) uint {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return uint(f)
		}
	}
	return fallback
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if dm, ok := v.(map[string]interface{}); ok {
			return dm
		}
	}
	return nil
}

func getIoTDB() *gorm.DB {
	if IoTDB != nil {
		return IoTDB
	}
	return DB
}
