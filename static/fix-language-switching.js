// 批量修复页面语言转换功能的脚本
// 使用方法：在浏览器控制台运行此脚本

(function() {
    console.log('开始修复页面语言转换功能...');
    
    // 需要修复的页面列表
    const pagesToFix = [
        'register1.html',
        'register12.html', 
        'register-3-2.html',
        'fix-registration.html',
        'debug-data.html',
        'file-management.html'
    ];
    
    // 语言转换器HTML模板
    const languageSwitcherHTML = `
        <!-- 语言转换器 -->
        <div class="language-switcher" aria-label="Language switcher">
            <button class="lang-btn" id="langBtn" type="button">
                <i class="fas fa-globe"></i>
                <span id="currentLang">EN</span>
                <i class="fas fa-chevron-down"></i>
            </button>
            <div class="lang-menu" id="langMenu" role="menu">
                <div class="lang-option active" data-lang="en" role="menuitem">
                    <span class="lang-flag">🇺🇸</span>
                    <span class="lang-name">English</span>
                </div>
                <div class="lang-option" data-lang="zh-cn" role="menuitem">
                    <span class="lang-flag">🇨🇳</span>
                    <span class="lang-name">简体</span>
                </div>
                <div class="lang-option" data-lang="zh-tw" role="menuitem">
                    <span class="lang-flag">🇭🇰</span>
                    <span class="lang-name">繁體</span>
                </div>
            </div>
        </div>`;
    
    // 语言转换器CSS样式
    const languageSwitcherCSS = `
        <style>
        .language-switcher {
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 1000;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
        
        .lang-btn {
            background: white;
            border: 1px solid #e1e5e9;
            border-radius: 8px;
            padding: 8px 12px;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 14px;
            font-weight: 500;
            color: #333;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            transition: all 0.2s ease;
        }
        
        .lang-btn:hover {
            background: #f8f9fa;
            border-color: #d1d5db;
        }
        
        .lang-btn i {
            font-size: 12px;
        }
        
        #currentLang {
            min-width: 30px;
            text-align: center;
        }
        
        .lang-menu {
            position: absolute;
            top: 100%;
            right: 0;
            background: white;
            border: 1px solid #e1e5e9;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            min-width: 140px;
            margin-top: 4px;
            display: none;
            overflow: hidden;
        }
        
        .lang-option {
            padding: 10px 12px;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 8px;
            transition: background-color 0.2s ease;
            border-bottom: 1px solid #f1f3f4;
        }
        
        .lang-option:last-child {
            border-bottom: none;
        }
        
        .lang-option:hover {
            background: #f8f9fa;
        }
        
        .lang-option.active {
            background: #e3f2fd;
            color: #1976d2;
        }
        
        .lang-flag {
            font-size: 16px;
        }
        
        .lang-name {
            font-size: 14px;
            font-weight: 500;
        }
        </style>`;
    
    // 语言转换器JavaScript
    const languageSwitcherJS = `
        <script>
        // 语言转换器交互逻辑
        document.addEventListener('DOMContentLoaded', function() {
            const langBtn = document.getElementById('langBtn');
            const langMenu = document.getElementById('langMenu');
            
            if (!langBtn || !langMenu) return;
            
            // 切换菜单显示
            function toggleLanguageMenu() {
                langMenu.style.display = (langMenu.style.display === 'block') ? 'none' : 'block';
            }
            
            // 点击按钮切换菜单
            langBtn.addEventListener('click', function(e) {
                e.stopPropagation();
                toggleLanguageMenu();
            });
            
            // 点击选项切换语言
            langMenu.querySelectorAll('.lang-option').forEach(option => {
                option.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const lang = this.getAttribute('data-lang');
                    
                    // 调用全局语言切换函数
                    if (window.switchLanguage) {
                        window.switchLanguage(lang);
                    }
                    
                    // 关闭菜单
                    langMenu.style.display = 'none';
                });
            });
            
            // 点击外部关闭菜单
            document.addEventListener('click', function(e) {
                if (!langBtn.contains(e.target) && !langMenu.contains(e.target)) {
                    langMenu.style.display = 'none';
                }
            });
            
            // 初始化语言选择器状态
            if (window.LanguageConfig) {
                window.LanguageConfig.updateLanguageSwitcher();
            }
        });
        </script>`;
    
    // 检查当前页面是否需要修复
    const currentPage = window.location.pathname.split('/').pop();
    if (pagesToFix.includes(currentPage)) {
        console.log(`正在修复当前页面: ${currentPage}`);
        
        // 检查是否已经有语言转换器
        const existingSwitcher = document.querySelector('.language-switcher');
        if (existingSwitcher) {
            console.log('页面已有语言转换器，跳过修复');
            return;
        }
        
        // 添加语言转换器到页面
        const container = document.querySelector('.container');
        if (container) {
            // 在grid-pattern后添加语言转换器
            const gridPattern = container.querySelector('.grid-pattern');
            if (gridPattern) {
                gridPattern.insertAdjacentHTML('afterend', languageSwitcherHTML);
            } else {
                container.insertAdjacentHTML('afterbegin', languageSwitcherHTML);
            }
            
            // 添加CSS样式
            document.head.insertAdjacentHTML('beforeend', languageSwitcherCSS);
            
            // 添加JavaScript
            document.body.insertAdjacentHTML('beforeend', languageSwitcherJS);
            
            // 引入语言配置文件
            if (!document.querySelector('script[src="/static/language-config.js"]')) {
                const script = document.createElement('script');
                script.src = '/static/language-config.js';
                document.head.appendChild(script);
            }
            
            console.log('语言转换器已添加到当前页面');
        } else {
            console.log('未找到容器元素，无法添加语言转换器');
        }
    } else {
        console.log(`当前页面 ${currentPage} 不需要修复或已在修复列表中`);
    }
    
    console.log('页面语言转换功能修复完成');
})();
