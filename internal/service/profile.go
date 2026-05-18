package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fabric-sdk/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

// GetBuildingProfile reads building profile from MySQL by email
func GetBuildingProfile(email string) (*model.BuildingAsset, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var asset model.BuildingAsset
	result := DB.Where("email = ?", email).First(&asset)
	if result.Error != nil {
		return nil, result.Error
	}
	return &asset, nil
}

// UpsertBuildingProfile creates or updates a building profile and writes reference to chain
func UpsertBuildingProfile(ctx context.Context, req *model.BuildingProfileRequest) (*model.BuildingAsset, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	tenantsJSON, _ := json.Marshal(req.MajorTenants)
	certsJSON, _ := json.Marshal(req.Certifications)

	// Compute data hash for integrity verification
	hashPayload := map[string]interface{}{
		"email": req.Email, "buildingName": req.BuildingName,
		"buildingAddr": req.BuildingAddr, "gfa": req.GFA, "nfa": req.NFA,
		"propertyType": req.PropertyType, "yearBuilt": req.YearBuilt,
	}
	hashBytes, _ := json.Marshal(hashPayload)
	hash := sha256.Sum256(hashBytes)
	dataHash := hex.EncodeToString(hash[:])

	// Chain asset key
	chainKey := fmt.Sprintf("BLDG-%s", strings.ReplaceAll(req.Email, "@", "_at_"))

	// Write to blockchain
	var txID string
	chainResult, chainErr := Chaincode.CreateAssetWithMetadata(ctx,
		chainKey, "BLDG_HASH", 0, req.Email, 0,
	)
	if chainErr != nil {
		g.Log().Warning(ctx, "chaincode write failed (continuing with DB save):", chainErr)
	} else if chainResult != nil {
		if v, ok := chainResult["txID"].(string); ok {
			txID = v
		} else if v, ok := chainResult["txid"].(string); ok {
			txID = v
		}
	}

	// Auto-generate IDs if empty
	assetID := req.AssetID
	if assetID == "" {
		assetID = fmt.Sprintf("AST-%s", strings.ToUpper(uuid.New().String()[:6]))
	}
	digitalAssetId := req.DigitalAssetId
	if digitalAssetId == "" {
		digitalAssetId = fmt.Sprintf("0xDA%s", uuid.New().String()[:12])
	}
	carbonAssetId := req.CarbonAssetId
	if carbonAssetId == "" {
		carbonAssetId = fmt.Sprintf("0xCA%s", uuid.New().String()[:12])
	}

	// Check existing
	var existing model.BuildingAsset
	result := DB.Where("email = ?", req.Email).First(&existing)

	if result.Error == nil {
		// Update
		updates := map[string]interface{}{}
		if req.BuildingName != "" { updates["building_name"] = req.BuildingName }
		if req.BuildingAddr != "" { updates["building_addr"] = req.BuildingAddr }
		if req.BuildingType != "" { updates["building_type"] = req.BuildingType }
		if req.Grade != "" { updates["grade"] = req.Grade }
		if req.AccreditedBy != "" { updates["accredited_by"] = req.AccreditedBy }
		if req.GFA != "" { updates["gfa"] = req.GFA }
		if req.NFA != "" { updates["nfa"] = req.NFA }
		if req.PropertyType != "" { updates["property_type"] = req.PropertyType }
		if req.YearBuilt != "" { updates["year_built"] = req.YearBuilt }
		if req.HVACSystem != "" { updates["hvac_system"] = req.HVACSystem }
		if req.FireSafety != "" { updates["fire_safety"] = req.FireSafety }
		if req.Developer != "" { updates["developer"] = req.Developer }
		if req.PropertyMgt != "" { updates["property_mgt"] = req.PropertyMgt }
		if string(tenantsJSON) != "null" { updates["major_tenants"] = string(tenantsJSON) }
		if string(certsJSON) != "null" { updates["certifications"] = string(certsJSON) }
		if req.Latitude != 0 { updates["latitude"] = req.Latitude }
		if req.Longitude != 0 { updates["longitude"] = req.Longitude }
		if req.OccupancyRate != 0 { updates["occupancy_rate"] = req.OccupancyRate }
		if req.AvgLeaseDuration != 0 { updates["avg_lease_duration"] = req.AvgLeaseDuration }
		updates["digital_asset_id"] = digitalAssetId
		updates["carbon_asset_id"] = carbonAssetId
		updates["data_hash"] = dataHash
		updates["updated_at"] = time.Now()
		if txID != "" {
			updates["chain_tx_id"] = txID
			updates["chain_asset_key"] = chainKey
			updates["verified"] = true
		}
		if err := DB.Model(&model.BuildingAsset{}).Where("email = ?", req.Email).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update building profile: %v", err)
		}
		DB.Where("email = ?", req.Email).First(&existing)
		return &existing, nil
	}

	// Create new
	asset := model.BuildingAsset{
		Email:            req.Email,
		BuildingName:     req.BuildingName,
		BuildingAddr:     req.BuildingAddr,
		BuildingType:     req.BuildingType,
		AssetID:          assetID,
		Grade:            req.Grade,
		AccreditedBy:     req.AccreditedBy,
		GFA:              req.GFA,
		NFA:              req.NFA,
		PropertyType:     req.PropertyType,
		YearBuilt:        req.YearBuilt,
		HVACSystem:       req.HVACSystem,
		FireSafety:       req.FireSafety,
		Developer:        req.Developer,
		PropertyMgt:      req.PropertyMgt,
		DigitalAssetId:   digitalAssetId,
		CarbonAssetId:    carbonAssetId,
		OccupancyRate:    req.OccupancyRate,
		AvgLeaseDuration: req.AvgLeaseDuration,
		MajorTenants:     string(tenantsJSON),
		Certifications:   string(certsJSON),
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		ChainTxID:        txID,
		ChainAssetKey:    chainKey,
		DataHash:         dataHash,
		Verified:         txID != "",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := DB.Create(&asset).Error; err != nil {
		return nil, fmt.Errorf("failed to create building profile: %v", err)
	}
	return &asset, nil
}

// GetBlockchainRecord reads on-chain asset to verify data integrity
func GetBlockchainRecord(ctx context.Context, email string) (map[string]interface{}, error) {
	chainKey := fmt.Sprintf("BLDG-%s", strings.ReplaceAll(email, "@", "_at_"))
	result, err := Chaincode.ReadAsset(ctx, chainKey)
	if err != nil {
		return nil, fmt.Errorf("blockchain read failed: %v", err)
	}
	return result, nil
}

// VerifyDataIntegrity checks on-chain asset matches MySQL record
func VerifyDataIntegrity(ctx context.Context, email string) (bool, error) {
	asset, err := GetBuildingProfile(email)
	if err != nil {
		return false, err
	}
	chainData, err := GetBlockchainRecord(ctx, email)
	if err != nil {
		return false, err
	}
	if chainAsset, ok := chainData["asset"].(map[string]interface{}); ok {
		return chainAsset["ID"] == asset.ChainAssetKey && chainAsset["Owner"] == email, nil
	}
	return false, nil
}

// AutoInitFromRegistration creates initial building profile from dids table data
func AutoInitFromRegistration(ctx context.Context, email string) (*model.BuildingAsset, error) {
	existing, err := GetBuildingProfile(email)
	if err == nil && existing != nil {
		return existing, nil
	}

	var did model.DID
	if err := DB.Where("email = ?", email).First(&did).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	req := &model.BuildingProfileRequest{
		Email:         email,
		BuildingName:  did.BuildingName,
		BuildingAddr:  did.BuildingAddr,
		BuildingType:  did.BuildingType,
	}
	return UpsertBuildingProfile(ctx, req)
}
