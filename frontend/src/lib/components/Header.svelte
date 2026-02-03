<script module lang="ts">
  export interface IHeaderProps {
    className?: string;
  }
</script>

<script lang="ts">
  import {t} from 'svelte-i18n';
  import {page} from '$app/state';
  import {goto} from '$app/navigation';
  import {ChevronUp, Bell, Search, Menu, Settings, LogOut, User} from 'lucide-svelte';
  import {Button} from "$lib/components/ui/button";
  import {Input} from "$lib/components/ui/input";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import * as Tooltip from "$lib/components/ui/tooltip/index.js";

  let {className}: { className?: string } = $props();

  let searchQuery = $state('');

  let user = $state({
    username: 'Admin',
    nickname: '管理员'
  });

  let notificationCount = $state(3);

  let currentMode = $derived(page.url.searchParams.get('mode') || 'published');

  function handleSearch(e: Event) {
    e.preventDefault();
    if (searchQuery.trim()) {
      const url = new URL(page.url);
      url.searchParams.set('q', searchQuery);
      goto(url.toString());
    }
  }

  function switchMode(mode: string) {
    const url = new URL(page.url);
    url.searchParams.set('mode', mode);
    goto(url.toString(), {keepFocus: true, noScroll: true});
  }
</script>

<header
    class={`
    sticky top-0 z-50 w-full flex h-12 items-center justify-between gap-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 px-4
    ${className || ''}
  `}
>
  <!-- 移动端菜单按钮 -->
  <Button variant="ghost" size="icon" class="md:hidden shrink-0">
    <Menu class="h-5 w-5"/>
  </Button>

  <!-- 右侧区域 -->
  <div class="flex items-center justify-between flex-1 gap-4">
    <!-- 搜索栏 -->
    <form onsubmit={handleSearch} class="flex-1 max-w-md">
      <div class="relative">
        <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground"/>
        <Input
            type="search"
            bind:value={searchQuery}
            placeholder={$t('action.search')}
            class="pl-9 pr-4 bg-muted/50 focus:bg-background transition-colors h-9"
        />
      </div>
    </form>

    <!-- 通知按钮 -->
    <Button variant="ghost" size="icon" class="relative shrink-0 h-9 w-9">
      <Bell class="h-4 w-4"/>
      {#if notificationCount > 0}
        <span
            class="absolute -top-1 -right-1 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-destructive text-[9px] font-medium text-destructive-foreground">
          {notificationCount}
        </span>
      {/if}
    </Button>
  </div>
</header>
