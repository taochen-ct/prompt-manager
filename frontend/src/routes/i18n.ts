import { init, register, locale } from 'svelte-i18n';

// 注册默认语言
register('en', () => import('../locales/en.json'));
register('zh', () => import('../locales/zh.json'));

// 初始化并设置默认语言
init({
  fallbackLocale: 'zh',
  initialLocale: 'zh',
});