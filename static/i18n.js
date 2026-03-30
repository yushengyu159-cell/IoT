/**
 * 统一的多语言管理类
 */
class I18n {
    constructor() {
        this.currentLanguage = 'en';
        this.translations = {};
        this.storageKey = 'esg_language';
        this.loadPromise = null;
    }

    async init() {
        const savedLang = localStorage.getItem(this.storageKey);
        if (savedLang && this.isValidLanguage(savedLang)) {
            this.currentLanguage = savedLang;
        } else {
            const browserLang = this.detectBrowserLanguage();
            this.currentLanguage = browserLang;
        }
        await this.loadTranslations();
        this.applyTranslations();
    }

    detectBrowserLanguage() {
        const lang = navigator.language || navigator.userLanguage;
        const langMap = {
            'zh-CN': 'zh-cn',
            'zh-TW': 'zh-tw',
            'zh-HK': 'zh-tw',
            'en': 'en',
            'en-US': 'en',
            'en-GB': 'en'
        };
        return langMap[lang] || 'en';
    }

    isValidLanguage(lang) {
        return ['en', 'zh-cn', 'zh-tw'].includes(lang);
    }

    async loadTranslations() {
        if (this.loadPromise) {
            return this.loadPromise;
        }
        this.loadPromise = fetch('/static/i18n.json?v=20250323')
            .then(response => response.json())
            .then(data => {
                this.translations = data;
                return data;
            })
            .catch(error => {
                console.error('Failed to load translations:', error);
                this.translations = this.getDefaultTranslations();
            });
        return this.loadPromise;
    }

    getDefaultTranslations() {
        return {
            en: { t: { title: 'Welcome to ESG VISA', continue: 'Continue' } },
            'zh-cn': { t: { title: '欢迎来到 ESG VISA', continue: '继续' } },
            'zh-tw': { t: { title: '歡迎來到 ESG VISA', continue: '繼續' } }
        };
    }

    async setLanguage(lang) {
        if (!this.isValidLanguage(lang)) {
            console.warn('Invalid language: ' + lang);
            return false;
        }
        this.currentLanguage = lang;
        localStorage.setItem(this.storageKey, lang);
        await this.loadTranslations();
        this.applyTranslations();
        window.dispatchEvent(new CustomEvent('languageChanged', { 
            detail: { language: lang }
        }));
        return true;
    }

    getLanguage() {
        return this.currentLanguage;
    }

    t(key, fallback) {
        const langData = this.translations[this.currentLanguage];
        if (!langData) {
            return fallback || key;
        }
        const keys = key.split('.');
        let value = langData;
        for (const k of keys) {
            if (value && typeof value === 'object') {
                value = value[k];
            } else {
                return fallback || key;
            }
        }
        return value || fallback || key;
    }

    applyTranslations() {
        const langData = this.translations[this.currentLanguage];
        if (!langData) return;
        const title = langData.t?.title;
        if (title) {
            const h1 = document.querySelector('.header h1');
            if (h1) h1.textContent = title;
        }
        const subtitle = langData.t?.subtitle;
        if (subtitle) {
            const p = document.querySelector('.header p');
            if (p) p.textContent = subtitle;
        }
        const emailInput = document.getElementById('email');
        if (emailInput && langData.t?.emailPh) {
            emailInput.placeholder = langData.t.emailPh;
        }
    }
}

window.i18n = new I18n();
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => window.i18n.init());
} else {
    window.i18n.init();
}
