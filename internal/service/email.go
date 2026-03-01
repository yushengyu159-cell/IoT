package service

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"

	"fabric-sdk/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"gopkg.in/yaml.v3"
)

// EmailConfig 邮箱配置
type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	Enabled  bool   `yaml:"enabled"`
	UseSSL   bool   `yaml:"use_ssl"` // 是否使用SSL
	UseTLS   bool   `yaml:"use_tls"` // 是否使用TLS
}

// EmailVerification 邮箱验证记录
type EmailVerification struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}

// 内存存储验证码（生产环境建议使用Redis）
var verificationCodes = make(map[string]*EmailVerification)

// 邮箱配置映射
var emailConfigs = make(map[string]*EmailConfig)

// 邮件模板多语言配置
type EmailTemplate struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailTemplates struct {
	Verification EmailTemplate `json:"verification"`
	Approval     EmailTemplate `json:"approval"`
	Rejection    EmailTemplate `json:"rejection"`
	PasswordReset EmailTemplate `json:"password_reset"`
}

var emailTemplates = map[string]EmailTemplates{
	"en": {
		Verification: EmailTemplate{
			Subject: "ESG VISA - Email Verification Code",
			Body: `
				<html>
				<body>
					<h2>ESG VISA Email Verification</h2>
					<p>Your verification code is: <strong style="font-size: 24px; color: #00BCD4;">%s</strong></p>
					<p>The verification code is valid for 10 minutes. Please use it promptly.</p>
					<p>If this was not your operation, please ignore this email.</p>
					<hr>
					<p style="color: #666; font-size: 12px;">
						This email is automatically sent by the ESG VISA system. Please do not reply.<br>
						Sent time: %s<br>
						Server IP: 47.238.159.234
					</p>
				</body>
				</html>
			`,
		},
		Approval: EmailTemplate{
			Subject: "ESG VISA - Registration Approved",
			Body: `
				<html>
				<body>
					<h2>Registration Approved</h2>
					<p>Email: %s</p>
					<p>Congratulations! Your registration has been approved. You can now log in to the system using your email and password.</p>
					<p style="color:#666;font-size:12px;">This email is automatically sent by the system. Please do not reply.</p>
				</body>
				</html>
			`,
		},
		Rejection: EmailTemplate{
			Subject: "ESG VISA - Registration Rejected",
			Body: `
				<html>
				<body>
					<h2>Registration Rejected</h2>
					<p>Email: %s</p>
					<p>Reason: %s</p>
					<p>Your registration has been rejected. You can register again with complete information.</p>
					<p style="color:#666;font-size:12px;">This email is automatically sent by the system. Please do not reply.</p>
				</body>
				</html>
			`,
		},
		PasswordReset: EmailTemplate{
			Subject: "ESG VISA - Password Reset Verification Code",
			Body: `
				<html>
				<body>
					<h2>ESG VISA Password Reset</h2>
					<p>Your password reset verification code is: <strong style="font-size: 24px; color: #00BCD4;">%s</strong></p>
					<p>The verification code is valid for 10 minutes. Please use it promptly.</p>
					<p>If this was not your operation, please ignore this email.</p>
					<hr>
					<p style="color: #666; font-size: 12px;">
						This email is automatically sent by the ESG VISA system. Please do not reply.<br>
						Sent time: %s<br>
						Server IP: 47.238.159.234
					</p>
				</body>
				</html>
			`,
		},
	},
	"zh": {
		Verification: EmailTemplate{
			Subject: "ESG VISA - 邮箱验证码",
			Body: `
				<html>
				<body>
					<h2>ESG VISA 邮箱验证</h2>
					<p>您的验证码是：<strong style="font-size: 24px; color: #00BCD4;">%s</strong></p>
					<p>验证码有效期为10分钟，请及时使用。</p>
					<p>如果这不是您的操作，请忽略此邮件。</p>
					<hr>
					<p style="color: #666; font-size: 12px;">
						此邮件由 ESG VISA 系统自动发送，请勿回复。<br>
						发送时间：%s<br>
						服务器IP：47.238.159.234
					</p>
				</body>
				</html>
			`,
		},
		Approval: EmailTemplate{
			Subject: "ESG VISA - 注册审核通过",
			Body: `
				<html>
				<body>
					<h2>注册审核通过</h2>
					<p>邮箱：%s</p>
					<p>恭喜！您的注册已审核通过。您现在可以使用邮箱和密码登录系统。</p>
					<p style="color:#666;font-size:12px;">此邮件由系统自动发送，请勿回复。</p>
				</body>
				</html>
			`,
		},
		Rejection: EmailTemplate{
			Subject: "ESG VISA - 注册审核拒绝",
			Body: `
				<html>
				<body>
					<h2>注册审核拒绝</h2>
					<p>邮箱：%s</p>
					<p>原因：%s</p>
					<p>您的注册已被拒绝。您可以重新注册并提供完整信息。</p>
					<p style="color:#666;font-size:12px;">此邮件由系统自动发送，请勿回复。</p>
				</body>
				</html>
			`,
		},
		PasswordReset: EmailTemplate{
			Subject: "ESG VISA - 密码重置验证码",
			Body: `
				<html>
				<body>
					<h2>ESG VISA 密码重置</h2>
					<p>您的密码重置验证码是：<strong style="font-size: 24px; color: #00BCD4;">%s</strong></p>
					<p>验证码有效期为10分钟，请及时使用。</p>
					<p>如果这不是您的操作，请忽略此邮件。</p>
					<hr>
					<p style="color: #666; font-size: 12px;">
						此邮件由 ESG VISA 系统自动发送，请勿回复。<br>
						发送时间：%s<br>
						服务器IP：47.238.159.234
					</p>
				</body>
				</html>
			`,
		},
	},
	"zh-tw": {
		Verification: EmailTemplate{
			Subject: "ESG VISA - 郵箱驗證碼",
			Body: `
				<html>
				<body>
					<h2>ESG VISA 郵箱驗證</h2>
					<p>您的驗證碼是：<strong style="font-size: 24px; color: #00BCD4;">%s</strong></p>
					<p>驗證碼有效期為10分鐘，請及時使用。</p>
					<p>如果這不是您的操作，請忽略此郵件。</p>
					<hr>
					<p style="color: #666; font-size: 12px;">
						此郵件由 ESG VISA 系統自動發送，請勿回覆。<br>
						發送時間：%s<br>
						伺服器IP：47.238.159.234
					</p>
				</body>
				</html>
			`,
		},
		Approval: EmailTemplate{
			Subject: "ESG VISA - 註冊審核通過",
			Body: `
				<html>
				<body>
					<h2>註冊審核通過</h2>
					<p>郵箱：%s</p>
					<p>恭喜！您的註冊已審核通過。您現在可以使用郵箱和密碼登入系統。</p>
					<p style="color:#666;font-size:12px;">此郵件由系統自動發送，請勿回覆。</p>
				</body>
				</html>
			`,
		},
		Rejection: EmailTemplate{
			Subject: "ESG VISA - 註冊審核拒絕",
			Body: `
				<html>
				<body>
					<h2>註冊審核拒絕</h2>
					<p>郵箱：%s</p>
					<p>原因：%s</p>
					<p>您的註冊已被拒絕。您可以重新註冊並提供完整資訊。</p>
					<p style="color:#666;font-size:12px;">此郵件由系統自動發送，請勿回覆。</p>
				</body>
				</html>
			`,
		},
		PasswordReset: EmailTemplate{
			Subject: "ESG VISA - 密碼重置驗證碼",
			Body: `
				<html>
				<body>
					<h2>ESG VISA 密碼重置</h2>
					<p>您的密碼重置驗證碼是：<strong style="font-size: 24px; color: #00BCD4;">%s</strong></p>
					<p>驗證碼有效期為10分鐘，請及時使用。</p>
					<p>如果這不是您的操作，請忽略此郵件。</p>
					<hr>
					<p style="color: #666; font-size: 12px;">
						此郵件由 ESG VISA 系統自動發送，請勿回覆。<br>
						發送時間：%s<br>
						伺服器IP：47.238.159.234
					</p>
				</body>
				</html>
			`,
		},
	},
}

// GetUserLanguage 获取用户语言偏好
func GetUserLanguage(email string) string {
	// 从数据库获取用户语言偏好
	if DB != nil {
		var did model.DID
		if err := DB.Where("email = ?", email).First(&did).Error; err == nil {
			// 如果用户设置了语言偏好，使用用户设置
			if did.Language != "" {
				return did.Language
			}
		}
	}
	
	// 默认返回中文
	return "zh"
}

// GetEmailTemplate 根据语言和类型获取邮件模板
func GetEmailTemplate(language, templateType string) EmailTemplate {
	templates, exists := emailTemplates[language]
	if !exists {
		// 如果语言不存在，使用中文作为默认
		templates = emailTemplates["zh"]
	}
	
	switch templateType {
	case "verification":
		return templates.Verification
	case "approval":
		return templates.Approval
	case "rejection":
		return templates.Rejection
	case "password_reset":
		return templates.PasswordReset
	default:
		return templates.Verification
	}
}

// 初始化邮箱配置
func initEmailConfigs() {
	// 尝试读取配置文件
	if err := loadEmailConfigFromFile(); err != nil {
		g.Log().Warning(nil, "无法读取邮箱配置文件，使用默认配置:", err)
		// 使用默认配置（仅用于测试）
		setDefaultEmailConfigs()
	}
}

// 从配置文件加载邮箱配置
func loadEmailConfigFromFile() error {
	configPath := "configs/email.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config struct {
		Email struct {
			QQ       EmailConfig `yaml:"qq"`
			Email163 EmailConfig `yaml:"163"`
			Gmail    EmailConfig `yaml:"gmail"`
		} `yaml:"email"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置配置
	if config.Email.QQ.Enabled {
		emailConfigs["qq.com"] = &config.Email.QQ
	}
	if config.Email.Email163.Enabled {
		emailConfigs["163.com"] = &config.Email.Email163
	}
	if config.Email.Gmail.Enabled {
		emailConfigs["gmail.com"] = &config.Email.Gmail
	}

	return nil
}

// 设置默认邮箱配置（仅用于测试）
func setDefaultEmailConfigs() {
	// 注意：这些是测试配置，生产环境必须使用真实的邮箱配置
	emailConfigs["qq.com"] = &EmailConfig{
		Host:     "smtp.qq.com",
		Port:     "587",
		Username: "test@qq.com",   // 测试用，需要替换为真实邮箱
		Password: "test_password", // 测试用，需要替换为真实SMTP密码
		From:     "test@qq.com",
		Enabled:  true,
		UseSSL:   false,
		UseTLS:   true,
	}

	emailConfigs["163.com"] = &EmailConfig{
		Host:     "smtp.163.com",
		Port:     "465",           // 改用465端口（SSL）
		Username: "test@163.com",  // 测试用，需要替换为真实邮箱
		Password: "test_password", // 测试用，需要替换为真实SMTP密码
		From:     "test@163.com",
		Enabled:  true,
		UseSSL:   true, // 使用SSL
		UseTLS:   false,
	}

	emailConfigs["gmail.com"] = &EmailConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: "test@gmail.com", // 测试用，需要替换为真实邮箱
		Password: "test_password",  // 测试用，需要替换为真实应用密码
		From:     "test@gmail.com",
		Enabled:  true,
		UseSSL:   false,
		UseTLS:   true,
	}

	g.Log().Warning(nil, "⚠️ 使用测试邮箱配置，邮件发送可能失败。生产环境请配置真实的邮箱信息。")
}

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code
}

// GetEmailConfig 根据邮箱地址获取配置
func GetEmailConfig(email string) (*EmailConfig, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的邮箱地址")
	}

	domain := parts[1]
	config, exists := emailConfigs[domain]
	if !exists {
		return nil, fmt.Errorf("不支持的邮箱类型: %s", domain)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("邮箱类型 %s 未启用", domain)
	}

	return config, nil
}

// GetUnifiedGmailConfig 获取统一的Gmail配置（用于发送所有邮件）
func GetUnifiedGmailConfig() (*EmailConfig, error) {
	config, exists := emailConfigs["gmail.com"]
	if !exists {
		return nil, fmt.Errorf("Gmail配置不存在")
	}

	if !config.Enabled {
		return nil, fmt.Errorf("Gmail配置未启用")
	}

	return config, nil
}

// SendVerificationEmail 发送验证码邮件
func SendVerificationEmail(email string) (string, error) {
	return SendVerificationEmailWithLanguage(email, "")
}

// SendVerificationEmailWithLanguage 发送验证码邮件（支持指定语言）
func SendVerificationEmailWithLanguage(email, language string) (string, error) {
	// 初始化邮箱配置
	initEmailConfigs()
	
	// 生成6位验证码
	code := GenerateVerificationCode()

	// 获取用户语言偏好
	if language == "" {
		language = GetUserLanguage(email)
	}

	g.Log().Info(nil, "📧 准备发送验证码到邮箱:", email, "验证码:", code, "语言:", language)

	// 使用统一的Gmail配置发送所有邮件
	config, err := GetUnifiedGmailConfig()
	if err != nil {
		g.Log().Error(nil, "❌ 获取Gmail配置失败:", err)
		return "", fmt.Errorf("获取Gmail配置失败: %v", err)
	}

	g.Log().Info(nil, "✅ 获取邮箱配置成功:", config.Host, "端口:", config.Port, "SSL:", config.UseSSL, "TLS:", config.UseTLS)

	// 检查是否是测试配置
	if isTestConfig(config) {
		g.Log().Warning(nil, "⚠️ 检测到测试配置，跳过邮件发送，直接返回验证码")
		// 存储验证码（10分钟有效期）
		verificationCodes[email] = &EmailVerification{
			Email:     email,
			Code:      code,
			ExpiresAt: time.Now().Add(10 * time.Minute),
			Used:      false,
		}
		return code, nil
	}

	// 获取邮件模板
	template := GetEmailTemplate(language, "verification")
	subject := template.Subject
	body := fmt.Sprintf(template.Body, code, time.Now().Format("2006-01-02 15:04:05"))
	
	g.Log().Info(nil, "📧 邮件模板信息 - 语言:", language, "主题:", subject)
	g.Log().Info(nil, "📧 邮件内容预览:", body[:100]+"...")

	// 构建邮件
	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", email, config.From, subject, body)

	g.Log().Info(nil, "📤 开始发送邮件...")

	// 根据配置选择发送方式
	if config.UseSSL {
		// 使用SSL连接
		err = sendMailSSL(config, email, message)
	} else if config.UseTLS {
		// 使用TLS连接
		err = sendMailTLS(config, email, message)
	} else {
		// 使用普通连接
		err = sendMailPlain(config, email, message)
	}

	if err != nil {
		g.Log().Error(nil, "❌ 发送邮件失败:", err)
		return "", fmt.Errorf("发送邮件失败: %v", err)
	}

	g.Log().Info(nil, "✅ 邮件发送成功:", email)

	// 存储验证码（10分钟有效期）
	verificationCodes[email] = &EmailVerification{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Used:      false,
	}

	g.Log().Info(nil, "💾 验证码已存储，有效期至:", verificationCodes[email].ExpiresAt.Format("2006-01-02 15:04:05"))
	return code, nil
}

// SendAuditResultEmail 发送审核结果通知
func SendAuditResultEmail(email string, approved bool, reason string) error {
	return SendAuditResultEmailWithLanguage(email, approved, reason, "")
}

// SendAuditResultEmailWithLanguage 发送审核结果通知（支持指定语言）
func SendAuditResultEmailWithLanguage(email string, approved bool, reason, language string) error {
	initEmailConfigs()
	
	// 获取用户语言偏好
	if language == "" {
		language = GetUserLanguage(email)
	}
	
	// 使用统一的Gmail配置发送所有邮件
	config, err := GetUnifiedGmailConfig()
	if err != nil {
		return err
	}
	
	// 获取邮件模板
	var template EmailTemplate
	if approved {
		template = GetEmailTemplate(language, "approval")
	} else {
		template = GetEmailTemplate(language, "rejection")
	}
	
	subject := template.Subject
	var body string
	if approved {
		body = fmt.Sprintf(template.Body, email)
	} else {
		body = fmt.Sprintf(template.Body, email, reason)
	}
	
	g.Log().Info(nil, "📧 审核邮件模板信息 - 语言:", language, "审核结果:", approved, "主题:", subject)
	g.Log().Info(nil, "📧 审核邮件内容预览:", body[:100]+"...")
	
	// 构建邮件
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", 
		email, config.From, subject, body)
	
	if isTestConfig(config) {
		statusText := "审核通过"
		if !approved {
			statusText = "审核拒绝"
		}
		g.Log().Info(nil, "[TEST] 审核结果邮件:", statusText, email, "语言:", language)
		return nil
	}
	
	if config.UseSSL {
		return sendMailSSL(config, email, msg)
	}
	if config.UseTLS {
		return sendMailTLS(config, email, msg)
	}
	return sendMailPlain(config, email, msg)
}

// SendPasswordResetEmail 发送密码重置验证码邮件
func SendPasswordResetEmail(email string) (string, error) {
	return SendPasswordResetEmailWithLanguage(email, "")
}

// SendPasswordResetEmailWithLanguage 发送密码重置验证码邮件（支持指定语言）
func SendPasswordResetEmailWithLanguage(email, language string) (string, error) {
	// 初始化邮箱配置
	initEmailConfigs()
	
	// 生成6位验证码
	code := GenerateVerificationCode()

	// 获取用户语言偏好
	if language == "" {
		language = GetUserLanguage(email)
	}

	g.Log().Info(nil, "📧 准备发送密码重置验证码到邮箱:", email, "验证码:", code, "语言:", language)

	// 使用统一的Gmail配置发送所有邮件
	config, err := GetUnifiedGmailConfig()
	if err != nil {
		g.Log().Error(nil, "❌ 获取Gmail配置失败:", err)
		return "", fmt.Errorf("获取Gmail配置失败: %v", err)
	}

	g.Log().Info(nil, "✅ 获取邮箱配置成功:", config.Host, "端口:", config.Port, "SSL:", config.UseSSL, "TLS:", config.UseTLS)

	// 检查是否是测试配置
	if isTestConfig(config) {
		g.Log().Warning(nil, "⚠️ 检测到测试配置，跳过邮件发送，直接返回验证码")
		// 存储验证码（10分钟有效期）
		verificationCodes[email] = &EmailVerification{
			Email:     email,
			Code:      code,
			ExpiresAt: time.Now().Add(10 * time.Minute),
			Used:      false,
		}
		return code, nil
	}

	// 获取邮件模板
	template := GetEmailTemplate(language, "password_reset")
	subject := template.Subject
	body := fmt.Sprintf(template.Body, code, time.Now().Format("2006-01-02 15:04:05"))

	// 构建邮件
	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", email, config.From, subject, body)

	g.Log().Info(nil, "📤 开始发送密码重置邮件...")

	// 根据配置选择发送方式
	if config.UseSSL {
		// 使用SSL连接
		err = sendMailSSL(config, email, message)
	} else if config.UseTLS {
		// 使用TLS连接
		err = sendMailTLS(config, email, message)
	} else {
		// 使用普通连接
		err = sendMailPlain(config, email, message)
	}

	if err != nil {
		g.Log().Error(nil, "❌ 发送密码重置邮件失败:", err)
		return "", fmt.Errorf("发送密码重置邮件失败: %v", err)
	}

	g.Log().Info(nil, "✅ 密码重置邮件发送成功:", email)

	// 存储验证码（10分钟有效期）
	verificationCodes[email] = &EmailVerification{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Used:      false,
	}

	g.Log().Info(nil, "💾 密码重置验证码已存储，有效期至:", verificationCodes[email].ExpiresAt.Format("2006-01-02 15:04:05"))
	return code, nil
}

// isTestConfig 检查是否是测试配置
func isTestConfig(config *EmailConfig) bool {
	// 检查密码是否是测试密码
	testPasswords := []string{
		"test_password",
		"your-qq-smtp-password",
		"your-163-smtp-password",
		"your-gmail-app-password",
	}

	for _, testPass := range testPasswords {
		if config.Password == testPass {
			return true
		}
	}

	return false
}

// sendMailSSL 使用SSL连接发送邮件
func sendMailSSL(config *EmailConfig, to, message string) error {
	// 建立SSL连接
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// 创建TLS配置
	tlsConfig := &tls.Config{
		ServerName:         config.Host,
		InsecureSkipVerify: false,
	}

	// 建立SSL连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("SSL连接失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return fmt.Errorf("SMTP客户端创建失败: %v", err)
	}
	defer client.Close()

	// 认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	// 设置发件人
	if err := client.Mail(config.From); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送邮件内容失败: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭邮件写入失败: %v", err)
	}

	return nil
}

// sendMailTLS 使用TLS连接发送邮件
func sendMailTLS(config *EmailConfig, to, message string) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 使用TLS连接
	return smtp.SendMail(addr, auth, config.From, []string{to}, []byte(message))
}

// sendMailPlain 使用普通连接发送邮件
func sendMailPlain(config *EmailConfig, to, message string) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 使用普通连接
	return smtp.SendMail(addr, auth, config.From, []string{to}, []byte(message))
}

// VerifyEmailCode 验证邮箱验证码
func VerifyEmailCode(email, code string) (bool, error) {
	verification, exists := verificationCodes[email]
	if !exists {
		return false, fmt.Errorf("验证码不存在或已过期")
	}

	// 检查是否已使用
	if verification.Used {
		return false, fmt.Errorf("验证码已被使用")
	}

	// 检查是否过期
	if time.Now().After(verification.ExpiresAt) {
		delete(verificationCodes, email)
		return false, fmt.Errorf("验证码已过期")
	}

	// 检查验证码是否匹配
	if verification.Code != code {
		return false, fmt.Errorf("验证码不正确")
	}

	// 标记为已使用
	verification.Used = true
	verificationCodes[email] = verification

	g.Log().Info(nil, "✅ 邮箱验证成功:", email)
	return true, nil
}

// CleanExpiredCodes 清理过期的验证码
func CleanExpiredCodes() {
	for email, verification := range verificationCodes {
		if time.Now().After(verification.ExpiresAt) {
			delete(verificationCodes, email)
			g.Log().Info(nil, "🗑️ 清理过期验证码:", email)
		}
	}
}

// GetVerificationStatus 获取验证码状态
func GetVerificationStatus(email string) (*EmailVerification, bool) {
	verification, exists := verificationCodes[email]
	if !exists {
		return nil, false
	}
	return verification, true
}

// CheckEmailExists 检查邮箱是否已完全注册成功
func CheckEmailExists(email string) bool {
	// 检查DID表中是否存在该邮箱且状态为completed
	var did model.DID
	if err := DB.Where("email = ? AND status = ?", email, "completed").First(&did).Error; err != nil {
		// 如果DID表中没有找到，检查是否在失败注册列表中
		// 如果在失败列表中，也认为可以重新注册（返回false）
		if IsEmailInFailedList(email) {
			g.Log().Info(nil, "📧 邮箱在失败注册列表中，可以重新注册:", email)
			return false
		}
		return false
	}
	return true
}

// TestSMTPConnection 测试SMTP连接
func TestSMTPConnection(email string) error {
	config, err := GetEmailConfig(email)
	if err != nil {
		return fmt.Errorf("获取邮箱配置失败: %v", err)
	}

	g.Log().Info(nil, "🧪 测试SMTP连接:", config.Host, "端口:", config.Port)

	// 测试网络连接
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// 尝试建立TCP连接
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		g.Log().Error(nil, "❌ TCP连接失败:", err)
		return fmt.Errorf("TCP连接失败: %v", err)
	}
	defer conn.Close()

	g.Log().Info(nil, "✅ TCP连接成功")

	// 尝试SMTP握手
	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		g.Log().Error(nil, "❌ SMTP客户端创建失败:", err)
		return fmt.Errorf("SMTP客户端创建失败: %v", err)
	}
	defer client.Close()

	g.Log().Info(nil, "✅ SMTP握手成功")

	// 尝试认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err := client.Auth(auth); err != nil {
		g.Log().Error(nil, "❌ SMTP认证失败:", err)
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	g.Log().Info(nil, "✅ SMTP认证成功")

	return nil
}

// 启动定时清理任务
func init() {
	initEmailConfigs() // 初始化邮箱配置
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			CleanExpiredCodes()
		}
	}()
}
