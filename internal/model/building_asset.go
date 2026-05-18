package model

import (
	"time"

	"gorm.io/gorm"
)

// BuildingAsset stores full building profile data in MySQL
type BuildingAsset struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Email           string         `gorm:"uniqueIndex;size:128;not null" json:"email"`

	// Core Building Info
	BuildingName    string         `gorm:"size:128" json:"buildingName"`
	BuildingAddr    string         `gorm:"size:256" json:"buildingAddr"`
	BuildingType    string         `gorm:"size:64" json:"buildingType"`
	AssetID         string         `gorm:"size:64;uniqueIndex" json:"assetId"`
	Grade           string         `gorm:"size:16" json:"grade"`
	AccreditedBy    string         `gorm:"size:64" json:"accreditedBy"`

	// Physical Specifications
	GFA             string         `gorm:"size:64" json:"gfa"`
	NFA             string         `gorm:"size:64" json:"nfa"`
	PropertyType    string         `gorm:"size:64" json:"propertyType"`
	YearBuilt       string         `gorm:"size:32" json:"yearBuilt"`
	HVACSystem      string         `gorm:"size:128" json:"hvacSystem"`
	FireSafety      string         `gorm:"size:128" json:"fireSafety"`
	Developer       string         `gorm:"size:128" json:"developer"`
	PropertyMgt     string         `gorm:"size:128" json:"propertyMgt"`

	// Digital Identity
	DigitalAssetId  string         `gorm:"size:128" json:"digitalAssetId"`
	CarbonAssetId   string         `gorm:"size:128" json:"carbonAssetId"`

	// Occupancy Mix
	OccupancyRate    float64       `gorm:"default:0" json:"occupancyRate"`
	AvgLeaseDuration float64       `gorm:"default:0" json:"avgLeaseDuration"`
	MajorTenants     string        `gorm:"size:1000" json:"majorTenants"`

	// Certifications
	Certifications   string        `gorm:"size:1000" json:"certifications"`

	// Location
	Latitude         float64       `json:"latitude"`
	Longitude        float64       `json:"longitude"`

	// Blockchain Reference
	ChainTxID        string        `gorm:"size:128" json:"chainTxId"`
	ChainAssetKey    string        `gorm:"size:128" json:"chainAssetKey"`
	DataHash         string        `gorm:"size:128" json:"dataHash"`

	// Metadata
	Verified         bool          `gorm:"default:false" json:"verified"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BuildingAsset) TableName() string {
	return "building_assets"
}

// BuildingProfileRequest for POST /api/profile/building
type BuildingProfileRequest struct {
	Email            string   `json:"email"`
	BuildingName     string   `json:"buildingName,omitempty"`
	BuildingAddr     string   `json:"buildingAddr,omitempty"`
	BuildingType     string   `json:"buildingType,omitempty"`
	AssetID          string   `json:"assetId,omitempty"`
	Grade            string   `json:"grade,omitempty"`
	AccreditedBy     string   `json:"accreditedBy,omitempty"`
	GFA              string   `json:"gfa,omitempty"`
	NFA              string   `json:"nfa,omitempty"`
	PropertyType     string   `json:"propertyType,omitempty"`
	YearBuilt        string   `json:"yearBuilt,omitempty"`
	HVACSystem       string   `json:"hvacSystem,omitempty"`
	FireSafety       string   `json:"fireSafety,omitempty"`
	Developer        string   `json:"developer,omitempty"`
	PropertyMgt      string   `json:"propertyMgt,omitempty"`
	DigitalAssetId   string   `json:"digitalAssetId,omitempty"`
	CarbonAssetId    string   `json:"carbonAssetId,omitempty"`
	OccupancyRate    float64  `json:"occupancyRate,omitempty"`
	AvgLeaseDuration float64  `json:"avgLeaseDuration,omitempty"`
	MajorTenants     []string `json:"majorTenants,omitempty"`
	Certifications   []Certification `json:"certifications,omitempty"`
	Latitude         float64  `json:"latitude,omitempty"`
	Longitude        float64  `json:"longitude,omitempty"`
}

type Certification struct {
	Type  string `json:"type"`
	Level string `json:"level"`
}
