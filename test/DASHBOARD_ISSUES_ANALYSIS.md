# Dashboard问题分析和解决方案

## 分析日期
2026-03-26

## 问题清单

### 问题1：用户信息Modal没有实现三种语言转换 ❌

**现状：**
- Dashboard添加了独立语言切换器（不符合要求）
- 未使用index.html的applyTranslations模式
- 未监听localStorage语言变化

**LOGIN界面实现：**
- 使用applyTranslations(data)函数
- 在switchLanguage()时调用applyTranslations(languageData[lang])

### 问题2：DID信息缺少注册地址等内容 ⚠️

**后端API返回的完整数据：**
- Phone, Age, DID, BuildingName, BuildingAddr, BuildingType等

**localStorage存储的简化数据：**
- 只有Email, Name, Role, Phone, Age, DID, RegisterTime

**缺失字段：**
- BuildingName, BuildingAddr, BuildingType, PropertyName, Occupation, Institution

### 问题3：Dashboard不应该有语言切换按键 ✅

**用户要求：**
- Dashboard不需要语言切换按键
- 统一使用index界面的语言切换

## 解决方案

### 1. 移除Dashboard语言切换器UI
删除header中的language-switcher div

### 2. 更新register1.html
存储完整的API响应数据（包括地址、职业等）

### 3. Dashboard添加语言支持
- 添加applyTranslations函数
- 监听esg_language变化
- 显示所有用户字段

## 实施计划

1. 移除language-switcher UI
2. 更新register1.html存储完整数据
3. Dashboard添加语言监听
4. 更新showUserDetail显示所有字段
5. 重新部署测试

---
**状态：** 待实施
