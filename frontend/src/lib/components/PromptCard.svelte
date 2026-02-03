<script lang="ts">
  import { t } from 'svelte-i18n';
  import { api } from '$lib/services/api';
  import { Button } from '$lib/components/ui/button';
  import Card from '$lib/components/ui/card/card.svelte';
  import CardHeader from '$lib/components/ui/card/card-header.svelte';
  import CardTitle from '$lib/components/ui/card/card-title.svelte';
  import CardDescription from '$lib/components/ui/card/card-description.svelte';
  import CardContent from '$lib/components/ui/card/card-content.svelte';
  import CardFooter from '$lib/components/ui/card/card-footer.svelte';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import { Star } from 'lucide-svelte';

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

  interface Props {
    prompt: Prompt;
    showFavorite?: boolean;
    ondelete?: (id: string) => void;
    onTogglePublish?: (prompt: Prompt) => void;
    onToggleFavorite?: (prompt: Prompt) => void;
  }

  let { prompt, showFavorite = false, ondelete, onTogglePublish, onToggleFavorite }: Props = $props();

  function handleDelete() {
    if (!confirm($t('prompt.deleteConfirm') as string)) return;
    ondelete?.(prompt.id);
  }

  function handlePublish() {
    onTogglePublish?.(prompt);
  }

  function handleFavorite() {
    onToggleFavorite?.(prompt);
  }
</script>

<Card>
  <CardHeader>
    <div class="flex items-start justify-between">
      <div class="flex items-start gap-2">
        {#if showFavorite}
          <button
            class="mt-1 {prompt.isFavorite ? 'text-yellow-500' : 'text-muted-foreground hover:text-yellow-500'}"
            onclick={handleFavorite}
          >
            <Star class="h-4 w-4" fill={prompt.isFavorite ? 'currentColor' : 'none'} />
          </button>
        {/if}
        <div class="space-y-1">
          <CardTitle class="text-lg">{prompt.name}</CardTitle>
          <CardDescription>{prompt.path}</CardDescription>
        </div>
      </div>
      <Badge variant={prompt.isPublish ? 'default' : 'secondary'}>
        {prompt.isPublish ? $t('prompt.published') : $t('prompt.unpublished')}
      </Badge>
    </div>
  </CardHeader>
  <CardContent>
    <div class="space-y-2 text-sm">
      <div class="flex justify-between">
        <span class="text-muted-foreground">{$t('version.version')}</span>
        <span class="font-medium">{prompt.latestVersion || '-'}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-muted-foreground">{$t('prompt.createdBy')}</span>
        <span class="font-medium">{prompt.username}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-muted-foreground">{$t('prompt.updatedAt')}</span>
        <span class="font-medium">{prompt.updateAt}</span>
      </div>
    </div>
  </CardContent>
  <CardFooter class="flex justify-end gap-2">
    <Button variant="outline" size="sm" onclick={handlePublish}>
      {prompt.isPublish ? $t('version.unpublish') : $t('version.publish')}
    </Button>
    <Button variant="outline" size="sm" onclick={handleDelete}>
      {$t('prompt.delete')}
    </Button>
    <a href="/editor/{prompt.id}">
      <Button size="sm">{$t('prompt.edit')}</Button>
    </a>
  </CardFooter>
</Card>
