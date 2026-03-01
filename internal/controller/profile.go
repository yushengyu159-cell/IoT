package controller

import (
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

type ProfileController struct{}

var Profile = new(ProfileController)

// GetBuildingProfile returns the profile data for a smart building
// GET /api/profile/building
func (c *ProfileController) GetBuildingProfile(r *ghttp.Request) {
	// In a real application, this would query the database or blockchain
	// For now, we return mock data as requested

	// Mock data structure
	type BuildingProfile struct {
		BuildingName string `json:"buildingName"`
		Address      string `json:"address"`
		Type         string `json:"type"`
		YearBuilt    int    `json:"yearBuilt"`
		GFA          string `json:"gfa"` // Gross Floor Area
		Grade        string `json:"grade"`
		AccreditedBy string `json:"accreditedBy"`
		Verified     bool   `json:"verified"`
		DigitalID    string `json:"digitalID"`
		CarbonID     string `json:"carbonID"`
		Coordinates  struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"coordinates"`
	}

	profile := BuildingProfile{
		BuildingName: "International Commerce Centre",
		Address:      "1 Austin Road West, West Kowloon, Hong Kong",
		Type:         "Commercial Complex",
		YearBuilt:    2010,
		GFA:          "312,917 m²",
		Grade:        "A",
		AccreditedBy: "BEAM Plus",
		Verified:     true,
		DigitalID:    "12325",
		CarbonID:     generateRandomHash(), // Simulate dynamic hash
		Coordinates: struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}{
			Lat: 22.3034,
			Lng: 114.1602,
		},
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    profile,
	})
}

// GetBlockchainRecord returns the blockchain record details for a given hash
// GET /api/profile/blockchain-record
func (c *ProfileController) GetBlockchainRecord(r *ghttp.Request) {
	hashID := r.Get("hashId").String()
	if hashID == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "hashId is required"})
		return
	}

	// Mock blockchain record
	record := map[string]interface{}{
		"txId":      "0x" + generateRandomHex(64),
		"blockNum":  rand.Intn(10000000),
		"timestamp": time.Now().Format(time.RFC3339),
		"signer":    "CN=Admin,OU=Fabric,O=Hyperledger,C=US",
		"status":    "VALID",
		"dataHash":  hashID,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    record,
	})
}

// Helper functions
func generateRandomHash() string {
	return "0x8A2" + generateRandomHex(3) + "..." + generateRandomHex(3) + "fo1"
}

func generateRandomHex(n int) string {
	const letters = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetCertificate returns the BEAM Plus certification details
// GET /api/profile/certificate
func (c *ProfileController) GetCertificate(r *ghttp.Request) {
	// Mock certificate data
	certificate := map[string]interface{}{
		"buildingName":     "International Commerce Centre",
		"address":          "1 Austin Road West, West Kowloon, Hong Kong",
		"certificationType": "NB", // NB = New Building, EB = Existing Building
		"grade":            "A",
		"validityPeriod":   "2020-01-01 to 2025-12-31",
		"txId":             "0x" + generateRandomHex(64),
		"certificateImage": "", // Can be a URL to certificate image
		"issuedDate":       "2020-01-01",
		"expiryDate":       "2025-12-31",
		"issuer":           "BEAM Plus",
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    certificate,
	})
}
