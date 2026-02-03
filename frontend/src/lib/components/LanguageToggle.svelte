<script>
  import {locale} from 'svelte-i18n';
  import {Button} from "$lib/components/ui/button/index.js";
  const toggleLanguage = () => {
    locale.set($locale === 'zh' ? 'en' : 'zh');
  }
</script>

<Button
    onclick={toggleLanguage}
    variant="outline"
    size="icon"
>
  <span
      class="!transition-all duration-500 h-[1.2rem] w-[1.2rem]
    {$locale === 'zh' ? 'scale-100 rotate-0' : 'scale-0 -rotate-90'}"
  >
    ZH
  </span>

  <span
      class="absolute !transition-all duration-500 h-[1.2rem] w-[1.2rem]
    {$locale === 'en' ? 'scale-100 rotate-0' : 'scale-0 rotate-90'}"
  >
    EN
  </span>

  <span class="sr-only">Toggle Language</span>
</Button>