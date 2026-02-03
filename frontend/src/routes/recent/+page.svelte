<script lang="ts">
  import { t } from 'svelte-i18n';
  import { toast } from 'svelte-sonner';
  import { api } from '$lib/services/api';
  import PromptTable from '$lib/components/PromptTable.svelte';
  import { onMount } from 'svelte';
  import { Input } from '$lib/components/ui/input';
  import { Search } from 'lucide-svelte';

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
  let searchQuery = $state('');

  async function loadRecent() {
    try {
      loading = true;
      const data = await api.getRecentPrompts(20);
      prompts = data as Prompt[];
    } catch (error) {
      console.error('Failed to load recent prompts:', error);
      toast.error($t('action.confirm') as string, {
        description: 'Failed to load recent prompts'
      });
    } finally {
      loading = false;
    }
  }

  function handleDeletePrompt(id: string) {
    handleDeletePromptAsync(id);
  }

  async function handleDeletePromptAsync(id: string) {
    if (!confirm($t('prompt.deleteConfirm') as string)) return;
    try {
      await api.deletePrompt(id);
      toast.success($t('prompt.delete') as string, {
        description: 'Prompt deleted successfully'
      });
      await loadRecent();
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
      await loadRecent();
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
    await loadRecent();
  }

  let filteredPrompts = $derived(
    prompts.filter(p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.path.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  onMount(() => {
    loadRecent();
  });
</script>

<div class="flex flex-1 flex-col gap-4 p-4">

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <p>{$t('prompt.loading')}</p>
    </div>
  {:else if filteredPrompts.length === 0}
    <div class="flex flex-col items-center justify-center py-12 text-center">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="mb-4 text-muted-foreground">
        <circle cx="12" cy="12" r="10"/>
        <polyline points="12 6 12 12 16 14"/>
      </svg>
      <h3 class="text-lg font-semibold">{$t('sidebar.recent')}</h3>
      <p class="text-muted-foreground">还没有最近访问的提示词</p>
    </div>
  {:else}
    <PromptTable
      prompts={filteredPrompts}
      showFavorite={true}
      ondelete={handleDeletePrompt}
      onTogglePublish={handlePublishToggle}
      onToggleFavorite={handleToggleFavorite}
    />
  {/if}
</div>
