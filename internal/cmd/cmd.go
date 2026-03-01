package cmd

import (
	"context"
	"fabric-sdk/internal/controller"
	"fabric-sdk/internal/service"

	// "fabric-sdk/pkg/fabric" // 已移除

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

var (
	Main = gcmd.Command{
		Name:  "fabric-sdk",
		Usage: "fabric-sdk",
		Brief: "Fabric SDK with GoFrame Framework",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// 静态资源与首页跳转
			s.AddStaticPath("/static", "static")
			s.BindHandler("/", func(r *ghttp.Request) {
				r.Response.RedirectTo("/static/")
			})

			// 自动初始化数据库
			if err := service.InitMySQL(); err != nil {
				g.Log().Fatal(ctx, "数据库初始化失败: ", err)
				return err
			}

			// 自动初始化链码服务
			if err := service.Chaincode.InitChaincodeService(ctx); err != nil {
				g.Log().Fatal(ctx, "链码服务初始化失败: ", err)
				return err
			}

			// 只注册链码、ESG等API路由
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Bind(
					controller.Hello,
				)
			})

			s.Group("/api/esg", func(group *ghttp.RouterGroup) {
				group.POST("/upload", controller.ESG.UploadFile)
				group.GET("/query", controller.ESG.QueryFile)
				group.GET("/list", controller.ESG.ListFiles)
				group.POST("/batch-list", controller.ESG.BatchListFiles)
				group.POST("/upload-encrypted", controller.ESG.UploadEncrypted)
				group.POST("/download-encrypted", controller.ESG.DownloadEncrypted)
			})

			s.Group("/api/chaincode", func(group *ghttp.RouterGroup) {
				group.ALL("/init", controller.Chaincode.Init)
				group.ALL("/write-record", controller.Chaincode.WriteRecord)
				group.ALL("/read-record", controller.Chaincode.ReadRecord)
				group.ALL("/create-asset", controller.Chaincode.CreateAsset)
				group.ALL("/read-asset", controller.Chaincode.ReadAsset)
				group.ALL("/transfer-asset", controller.Chaincode.TransferAsset)
				group.ALL("/asset-history", controller.Chaincode.GetAssetHistory)
				group.ALL("/get-all-assets", controller.Chaincode.GetAllAssets)
				group.ALL("/test", controller.Chaincode.TestFunctions)
				group.ALL("/close", controller.Chaincode.Close)
			})

			s.Group("/api/ipfs", func(group *ghttp.RouterGroup) {
				group.POST("/upload", controller.Ipfs.Upload)
				group.GET("/download", controller.Ipfs.Download)
				group.POST("/upload-encrypted", controller.Ipfs.UploadEncrypted)
				group.POST("/download-decrypted", controller.Ipfs.DownloadDecrypted)
			})

			s.Group("/api/did", func(group *ghttp.RouterGroup) {
				group.POST("/register", controller.DID.RegisterDID)
				group.POST("/verify", controller.DID.VerifyDID)
			})

			// 邮箱验证相关路由
			s.Group("/api/email", func(group *ghttp.RouterGroup) {
				group.POST("/send-code", controller.Email.SendVerificationCode)
				group.POST("/verify-code", controller.Email.VerifyCode)
				group.GET("/status", controller.Email.GetVerificationStatus)
				group.POST("/resend-code", controller.Email.ResendCode)
				group.POST("/send-password-reset", controller.Email.SendPasswordResetCode)
				group.POST("/check-email-exists", controller.Email.CheckEmailExists)
				group.POST("/test-connection", controller.Email.TestConnection)
			})

			// Smart Building Profile Routes
			s.Group("/api/profile", func(group *ghttp.RouterGroup) {
				group.GET("/building", controller.Profile.GetBuildingProfile)
				group.GET("/blockchain-record", controller.Profile.GetBlockchainRecord)
				group.GET("/certificate", controller.Profile.GetCertificate)
			})

			// Carbon Transparency Routes
			s.Group("/api/carbon", func(group *ghttp.RouterGroup) {
				group.GET("/overview", controller.Carbon.GetCarbonOverview)
				group.GET("/tracking/timeline", controller.Carbon.GetCarbonTrackingTimeline)
				group.GET("/tracking/source", controller.Carbon.GetCarbonSourceDecomposition)
				group.GET("/verification", controller.Carbon.GetCarbonVerification)
				group.GET("/analysis", controller.Carbon.GetCarbonAnalysis)
			})

			// ESG Analytics Routes
			s.Group("/api/esg", func(group *ghttp.RouterGroup) {
				group.POST("/report/generate", controller.ESG.GenerateReport)
				group.GET("/analysis", controller.ESG.GetAIAnalysis)
				group.GET("/reports", controller.ESG.GetReportHistory)
				group.GET("/report/:id", controller.ESG.GetReportDetail)
				group.GET("/report/:id/export", controller.ESG.ExportReport)
			})

			// 管理端审核路由
			s.Group("/api/admin", func(group *ghttp.RouterGroup) {
				group.POST("/review-approve", controller.Admin.ReviewApprove)
				group.POST("/review-reject", controller.Admin.ReviewReject)
			})

			// 微信登录相关路由
			s.Group("/api/wechat", func(group *ghttp.RouterGroup) {
				group.POST("/check-and-generate-qrcode", controller.Wechat.CheckAndGenerateQRCode)
				group.GET("/check-login-status", controller.Wechat.CheckLoginStatus)
				group.POST("/confirm-login", controller.Wechat.ConfirmWeChatLogin)
				group.POST("/scan-qrcode", controller.Wechat.ScanQRCode)
			})

			// 失败注册记录管理路由
			failedRegController := &controller.FailedRegistrationController{}
			s.Group("/api/failed-registration", func(group *ghttp.RouterGroup) {
				group.POST("/add", failedRegController.AddFailedRegistration)
				group.GET("/get", failedRegController.GetFailedRegistration)
				group.GET("/list", failedRegController.GetAllFailedRegistrations)
				group.POST("/update-status", failedRegController.UpdateFailedRegistrationStatus)
				group.DELETE("/delete", failedRegController.DeleteFailedRegistration)
				group.GET("/stats", failedRegController.GetFailedRegistrationStats)
			})

			s.Group("/api/register", func(group *ghttp.RouterGroup) {
				group.POST("/step1", controller.Register.Step1)
				group.POST("/step2", controller.Register.Step2)
				group.POST("/step3", controller.Register.Step3)
				group.POST("/complete", controller.Register.Complete)
				group.GET("/status", controller.Register.GetStatus)
				group.POST("/update", controller.Register.Update)
				group.POST("/verify", controller.Register.Verify)
				group.POST("/check-email", controller.Register.CheckEmail)
				group.POST("/reset-password", controller.Register.ResetPassword)
				group.POST("/clear-cache", controller.Register.ClearCache)
				// 新增：获取所有注册成功的邮箱账号
				group.GET("/emails", controller.Register.GetAllEmails)
				// 新增：获取所有用户（管理员用，只显示已审核通过的用户）
				group.GET("/admin/users", controller.Register.GetAllUsersForAdmin)
				// 新增：获取待审核用户（管理员审核用）
				group.GET("/admin/pending-users", controller.Register.GetPendingUsersForAdmin)
				// 新增：清理未完成注册的账号
				group.POST("/admin/cleanup", controller.Register.CleanupIncompleteRegistrations)
				// 新增：根据邮箱获取用户详细信息
				group.GET("/user-detail", controller.Register.GetUserDetail)
			})

			// 位置信息相关路由
			locationController := controller.NewLocationController()
			s.Group("/api/location", func(group *ghttp.RouterGroup) {
				group.GET("/provinces", locationController.GetProvinces)
				group.GET("/cities", locationController.GetCities)
				group.GET("/districts", locationController.GetDistricts)
				group.GET("/pois", locationController.SearchPOI)
				group.POST("/save", locationController.SaveLocation)
				group.GET("/user", locationController.GetUserLocation)
				group.GET("/nearby", locationController.GetNearbyPOIs)
				group.GET("/geocode", locationController.GeocodeAddress)
				group.GET("/reverse-geocode", locationController.ReverseGeocode)
			})

			// LiveSense API路由（仅belucksapi用户可用）
			s.Group("/api/livesense", func(group *ghttp.RouterGroup) {
				// 添加用户验证中间件
				group.Middleware(controller.LiveSense.CheckBelucksAPIUser)
				group.GET("/values", controller.LiveSense.GetSensorValues)
				group.GET("/sensors", controller.LiveSense.GetSensorsUnderContext)
				group.POST("/authenticate", controller.LiveSense.Authenticate)
			})

			// 配置CORS中间件
			s.Use(func(r *ghttp.Request) {
				r.Response.CORSDefault()
				r.Middleware.Next()
			})

			// 启动HTTP服务
			s.Run()
			return nil
		},
	}
)
