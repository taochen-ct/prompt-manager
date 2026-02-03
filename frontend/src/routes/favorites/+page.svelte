<script lang="ts">
  import { t } from 'svelte-i18n';
  import { toast } from 'svelte-sonner';
  import { api } from '$lib/services/api';
  import PromptTable from '$lib/components/PromptTable.svelte';
  import { Search } from 'lucide-svelte';
  import { Input } from '$lib/components/ui/input';

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

  async function loadFavorites() {
    try {
      loading = true;
      const data = await api.getFavorites();
      prompts = data as Prompt[];
    } catch (error) {
      console.error('Failed to load favorites:', error);
      toast.error($t('action.confirm') as string, {
        description: 'Failed to load favorites'
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
      await loadFavorites();
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
      await loadFavorites();
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
    await loadFavorites();
  }

  let filteredPrompts = $derived(
    prompts.filter(p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.path.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  $effect(() => {
    loadFavorites();
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
        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
      </svg>
      <h3 class="text-lg font-semibold">{$t('sidebar.favorites')}</h3>
      <p class="text-muted-foreground">还没有收藏任何提示词</p>
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
