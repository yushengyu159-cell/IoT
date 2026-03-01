// 自动翻译脚本 - 用于没有语言转换按键的页面
// 这些页面会自动跟随index页面的语言设置

(function() {
    'use strict';
    
    // 等待语言配置加载
    function waitForLanguageConfig() {
        return new Promise((resolve) => {
            if (window.LanguageConfig) {
                resolve();
            } else {
                const checkInterval = setInterval(() => {
                    if (window.LanguageConfig) {
                        clearInterval(checkInterval);
                        resolve();
                    }
                }, 100);
            }
        });
    }
    
    // 自动翻译函数
    function autoTranslate() {
        if (!window.LanguageConfig) return;
        
        const currentLang = window.LanguageConfig.currentLanguage;
        const langData = window.LanguageConfig.data[currentLang];
        
        if (!langData) return;
        
        // 翻译所有带有data-translate属性的元素
        document.querySelectorAll('[data-translate]').forEach(element => {
            const key = element.getAttribute('data-translate');
            const translation = langData.translations[key];
            
            if (translation) {
                if (element.tagName === 'INPUT' && element.type !== 'submit' && element.type !== 'button') {
                    // 输入框更新placeholder
                    element.placeholder = translation;
                } else if (element.tagName === 'BUTTON' || element.type === 'submit') {
                    // 按钮更新文本内容
                    element.textContent = translation;
                } else {
                    // 其他元素更新文本内容
                    element.textContent = translation;
                }
            }
        });
    }
    
    // 监听语言变化事件
    function setupLanguageListener() {
        // 监听localStorage变化（当index页面切换语言时）
        window.addEventListener('storage', function(e) {
            if (e.key === 'esg_language' && e.newValue) {
                if (window.LanguageConfig) {
                    window.LanguageConfig.currentLanguage = e.newValue;
                    autoTranslate();
                }
            }
        });
        
        // 监听自定义语言变化事件
        window.addEventListener('languageChanged', function(e) {
            autoTranslate();
        });
    }
    
    // 页面加载完成后执行
    document.addEventListener('DOMContentLoaded', function() {
        waitForLanguageConfig().then(() => {
            // 立即执行一次翻译
            autoTranslate();
            
            // 设置监听器
            setupLanguageListener();
        });
    });
    
    // 如果页面已经加载完成，立即执行
    if (document.readyState === 'loading') {
        // 页面还在加载，等待DOMContentLoaded事件
    } else {
        // 页面已经加载完成，立即执行
        waitForLanguageConfig().then(() => {
            autoTranslate();
            setupLanguageListener();
        });
    }
})();
