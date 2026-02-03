<script lang="ts">
  import { t } from 'svelte-i18n';
  import { toast } from 'svelte-sonner';
  import { api } from '$lib/services/api';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import * as Dialog from '$lib/components/ui/dialog';
  import PromptTable from '$lib/components/PromptTable.svelte';
  import { Plus, Search } from 'lucide-svelte';

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

  async function loadPrompts() {
    try {
      loading = true;
      const data = await api.getPromptList({ offset: 0, limit: 100 });
      prompts = data as Prompt[];
    } catch (error) {
      console.error('Failed to load prompts:', error);
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
    if (!confirm($t('prompt.deleteConfirm') as string)) return;
    try {
      await api.deletePrompt(id);
      toast.success($t('prompt.delete') as string, {
        description: 'Prompt deleted successfully'
      });
      await loadPrompts();
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
      await loadPrompts();
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
    await loadPrompts();
  }

  let filteredPrompts = $derived(
    prompts.filter(p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.path.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  $effect(() => {
    loadPrompts();
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
        <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/>
        <polyline points="14 2 14 8 20 8"/>
        <line x1="12" x2="12" y1="18" y2="18"/>
        <line x1="9" x2="15" y1="15" y2="15"/>
      </svg>
      <h3 class="text-lg font-semibold">{$t('prompt.noPrompts')}</h3>
      <p class="text-muted-foreground">{$t('prompt.createFirst')}</p>
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
