package controller

import (
	"fabric-sdk/internal/service"
	"strconv"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ChaincodeController 链码控制器
type ChaincodeController struct{}

var Chaincode = new(ChaincodeController)

// Init 初始化链码服务
func (c *ChaincodeController) Init(r *ghttp.Request) {
	ctx := r.Context()

	// 初始化链码服务
	err := service.Chaincode.InitChaincodeService(ctx)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "初始化链码服务失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "链码服务初始化成功",
	})
}

// WriteRecord 写入记录
func (c *ChaincodeController) WriteRecord(r *ghttp.Request) {
	ctx := r.Context()

	// 获取参数
	key := r.Get("key").String()
	value := r.Get("value").String()

	if key == "" || value == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "key和value参数不能为空",
		})
		return
	}

	// 调用服务
	result, err := service.Chaincode.WriteRecord(ctx, key, value)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "写入记录失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "写入记录成功",
		Data:    result,
	})
}

// ReadRecord 读取记录
func (c *ChaincodeController) ReadRecord(r *ghttp.Request) {
	ctx := r.Context()

	// 获取参数
	key := r.Get("key").String()

	if key == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "key参数不能为空",
		})
		return
	}

	// 调用服务
	result, err := service.Chaincode.ReadRecord(ctx, key)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "读取记录失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "读取记录成功",
		Data:    result,
	})
}

// CreateAsset 创建资产
func (c *ChaincodeController) CreateAsset(r *ghttp.Request) {
	ctx := r.Context()

	// 获取参数
	assetID := r.Get("assetID").String()
	color := r.Get("color").String()
	sizeStr := r.Get("size").String()
	owner := r.Get("owner").String()
	appraisedValueStr := r.Get("appraisedValue").String()

	if assetID == "" || color == "" || sizeStr == "" || owner == "" || appraisedValueStr == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "所有参数都不能为空",
		})
		return
	}

	// 转换参数
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "size参数必须是数字",
		})
		return
	}

	appraisedValue, err := strconv.Atoi(appraisedValueStr)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "appraisedValue参数必须是数字",
		})
		return
	}

    // 调用服务（返回包含txID/timestamp的元数据）
    result, err := service.Chaincode.CreateAsset(ctx, assetID, color, size, owner, appraisedValue)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "创建资产失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "创建资产成功",
		Data:    result,
	})
}

// GetAssetHistory 获取资产历史
func (c *ChaincodeController) GetAssetHistory(r *ghttp.Request) {
    ctx := r.Context()

    assetID := r.Get("assetID").String()
    if assetID == "" {
        r.Response.WriteJson(ghttp.DefaultHandlerResponse{
            Code:    400,
            Message: "assetID参数不能为空",
        })
        return
    }

    result, err := service.Chaincode.GetAssetHistory(ctx, assetID)
    if err != nil {
        r.Response.WriteJson(ghttp.DefaultHandlerResponse{
            Code:    500,
            Message: "获取资产历史失败: " + err.Error(),
        })
        return
    }

    r.Response.WriteJson(ghttp.DefaultHandlerResponse{
        Code:    200,
        Message: "获取资产历史成功",
        Data:    result,
    })
}

// ReadAsset 读取资产
func (c *ChaincodeController) ReadAsset(r *ghttp.Request) {
	ctx := r.Context()

	// 获取参数
	assetID := r.Get("assetID").String()

	if assetID == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "assetID参数不能为空",
		})
		return
	}

	// 调用服务
	result, err := service.Chaincode.ReadAsset(ctx, assetID)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "读取资产失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "读取资产成功",
		Data:    result,
	})
}

// TransferAsset 转移资产
func (c *ChaincodeController) TransferAsset(r *ghttp.Request) {
	ctx := r.Context()

	// 获取参数
	assetID := r.Get("assetID").String()
	newOwner := r.Get("newOwner").String()

	if assetID == "" || newOwner == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "assetID和newOwner参数不能为空",
		})
		return
	}

	// 调用服务
	result, err := service.Chaincode.TransferAsset(ctx, assetID, newOwner)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "转移资产失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "转移资产成功",
		Data:    result,
	})
}

// GetAllAssets 获取所有资产
func (c *ChaincodeController) GetAllAssets(r *ghttp.Request) {
	ctx := r.Context()

	// 调用服务
	result, err := service.Chaincode.GetAllAssets(ctx)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取所有资产失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取所有资产成功",
		Data:    result,
	})
}

// TestFunctions 测试所有功能
func (c *ChaincodeController) TestFunctions(r *ghttp.Request) {
	ctx := r.Context()

	// 调用服务
	result, err := service.Chaincode.TestChaincodeFunctions(ctx)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "测试功能失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "测试功能完成",
		Data:    result,
	})
}

// Close 关闭链码服务
func (c *ChaincodeController) Close(r *ghttp.Request) {
	ctx := r.Context()

	// 关闭链码服务
	err := service.Chaincode.CloseChaincodeService(ctx)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "关闭链码服务失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "链码服务已关闭",
	})
}
