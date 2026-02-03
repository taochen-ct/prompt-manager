<script lang="ts">
  import { t } from 'svelte-i18n';
  import { toast } from 'svelte-sonner';
  import { api } from '$lib/services/api';
  import PromptCard from '$lib/components/PromptCard.svelte';
  import { onMount } from 'svelte';
  import type { PageData } from './$types';

  interface Prompt {
    id: string;
    name: string;
    path: string;
    latestVersion: string;
    isPublish: boolean;
    createBy: string;
    username: string;
    createAt: string;
    updateAt: string;
    isFavorite?: boolean;
  }

  let prompts = $state<Prompt[]>([]);
  let loading = $state(true);
  let { data }: { data: PageData } = $props();
  let categoryName = $derived(data.category)

  async function loadCategoryPrompts() {
    try {
      loading = true;
      const data = await api.getPromptsByCategory(categoryName);
      prompts = data as Prompt[];
    } catch (error) {
      console.error('Failed to load category prompts:', error);
      toast.error($t('action.confirm') as string, {
        description: 'Failed to load prompts'
      });
    } finally {
      loading = false;
    }
  }

  function handleDeletePrompt(id: string) {
    handleDeletePromptAsync(id);
  }

  async function handleDeletePromptAsync(id: string) {
    try {
      await api.deletePrompt(id);
      toast.success($t('prompt.delete') as string, {
        description: 'Prompt deleted successfully'
      });
      await loadCategoryPrompts();
    } catch (error) {
      toast.error($t('action.confirm') as string, {
        description: 'Failed to delete prompt'
      });
    }
  }

  function handlePublishToggle(prompt: Prompt) {
    handlePublishToggleAsync(prompt);
  }

  async function handlePublishToggleAsync(prompt: Prompt) {
    try {
      await api.updatePrompt({
        id: prompt.id,
        name: prompt.name,
        isPublish: !prompt.isPublish
      });
      await loadCategoryPrompts();
    } catch (error) {
      toast.error($t('action.confirm') as string, {
        description: 'Failed to update prompt'
      });
    }
  }

  function handleToggleFavorite(prompt: Prompt) {
    handleToggleFavoriteAsync(prompt);
  }

  async function handleToggleFavoriteAsync(prompt: Prompt) {
    await api.toggleFavorite(prompt.id);
    await loadCategoryPrompts();
  }

  function getCategoryTitle(name: string): string {
    const titles: Record<string, string> = {
      'copywriting': '文案生成',
      'coding': '代码助手',
      'translation': '翻译工具',
      'analysis': '数据分析'
    };
    return titles[name] || name;
  }

  onMount(() => {
    loadCategoryPrompts();
  });

  $effect(() => {
    if (categoryName) {
      loadCategoryPrompts();
    }
  });
</script>

<div class="flex flex-1 flex-col gap-4 p-4">
  <div class="flex items-center gap-2">
    <a href="/" class="text-sm text-muted-foreground hover:text-foreground">{$t('sidebar.all')}</a>
    <span class="text-muted-foreground">/</span>
    <span class="font-medium">{getCategoryTitle(categoryName)}</span>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <p>{$t('prompt.loading')}</p>
    </div>
  {:else if prompts.length === 0}
    <div class="flex flex-col items-center justify-center py-12 text-center">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="mb-4 text-muted-foreground">
        <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
      </svg>
      <h3 class="text-lg font-semibold">{getCategoryTitle(categoryName)}</h3>
      <p class="text-muted-foreground">该分类下还没有提示词</p>
    </div>
  {:else}
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {#each prompts as prompt (prompt.id)}
        <PromptCard
          {prompt}
          showFavorite={true}
          ondelete={handleDeletePrompt}
          onTogglePublish={handlePublishToggle}
          onToggleFavorite={handleToggleFavorite}
        />
      {/each}
    </div>
  {/if}
</div>
