package model

type IoTData struct {
	ID          uint   `gorm:"primaryKey"`
	DeviceName  string `gorm:"size:128;index"`
	DeviceEUI   string `gorm:"size:32;index"`
	FPort       int    `gorm:"default:0;index"`
	FCnt        int    `gorm:"default:0"`
	UplinkType  string `gorm:"size:32;index"` // realtime / history / device_info / simple_log

	// Common fields (fPort=2 realtime)
	NodeType     string  `gorm:"size:64"`
	BatV         float64 `gorm:"default:0"`
	Current1A    float64 `gorm:"default:0;column:current1_a"`
	Current2A    float64 `gorm:"default:0;column:current2_a"`
	Current3A    float64 `gorm:"default:0;column:current3_a"`
	Current4A    float64 `gorm:"default:0;column:current4_a"`
	EXTITrigger  string  `gorm:"size:16;column:exti_trigger"`
	EXTILevel    string  `gorm:"size:16;column:exti_level"`
	Cur1LStatus  string  `gorm:"size:16;column:cur1l_status"`
	Cur1HStatus  string  `gorm:"size:16;column:cur1h_status"`
	Cur2LStatus  string  `gorm:"size:16;column:cur2l_status"`
	Cur2HStatus  string  `gorm:"size:16;column:cur2h_status"`
	Cur3LStatus  string  `gorm:"size:16;column:cur3l_status"`
	Cur3HStatus  string  `gorm:"size:16;column:cur3h_status"`
	Cur4LStatus  string  `gorm:"size:16;column:cur4l_status"`
	Cur4HStatus  string  `gorm:"size:16;column:cur4h_status"`

	// Cumulative current (fPort=2 extended, bytes.length==28 or 32)
	CurTotalMod    int  `gorm:"default:0;column:curtotal_mod"`
	CurTotal1MAmin uint `gorm:"default:0;column:curtotal1_ma_min"`
	CurTotal2MAmin uint `gorm:"default:0;column:curtotal2_ma_min"`
	CurTotal3MAmin uint `gorm:"default:0;column:curtotal3_ma_min"`
	CurTotal4MAmin uint `gorm:"default:0;column:curtotal4_ma_min"`

	// Device info (fPort=5)
	SensorModel     string `gorm:"size:64;column:sensor_model"`
	FirmwareVersion string `gorm:"size:32;column:firmware_version"`
	FrequencyBand   string `gorm:"size:32;column:frequency_band"`
	SubBand         string `gorm:"size:16;column:sub_band"`

	// Signal
	RSSI      float64 `gorm:"default:0"`
	SNR       float64 `gorm:"default:0"`
	Frequency int     `gorm:"default:0"`

	RawPayload  string `gorm:"type:text"`
	DecodedData string `gorm:"type:text"`
	ReceivedAt  string `gorm:"size:32;index"`
}
