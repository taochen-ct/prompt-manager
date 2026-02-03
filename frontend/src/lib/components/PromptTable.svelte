<script lang="ts">
  import { t } from 'svelte-i18n';
  import * as Table from '$lib/components/ui/table/index.js';
  import { Button } from '$lib/components/ui/button';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import { Star, MoreHorizontal, Edit, Trash2, Eye, EyeOff } from 'lucide-svelte';
  import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

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
    prompts: Prompt[];
    showFavorite?: boolean;
    ondelete?: (id: string) => void;
    onTogglePublish?: (prompt: Prompt) => void;
    onToggleFavorite?: (prompt: Prompt) => void;
  }

  let { prompts, showFavorite = false, ondelete, onTogglePublish, onToggleFavorite }: Props = $props();

  function handleDelete(id: string) {
    ondelete?.(id);
  }

  function handlePublish(prompt: Prompt) {
    onTogglePublish?.(prompt);
  }

  function handleFavorite(prompt: Prompt) {
    onToggleFavorite?.(prompt);
  }
</script>

<div class="rounded-md border">
  <Table.Root>
    <Table.Header>
      <Table.Row>
        {#if showFavorite}
          <Table.Head class="w-12"></Table.Head>
        {/if}
        <Table.Head>{$t('prompt.name')}</Table.Head>
        <Table.Head>{$t('prompt.path')}</Table.Head>
        <Table.Head>{$t('version.version')}</Table.Head>
        <Table.Head>{$t('prompt.status')}</Table.Head>
        <Table.Head>{$t('prompt.createdBy')}</Table.Head>
        <Table.Head>{$t('prompt.updatedAt')}</Table.Head>
        <Table.Head class="w-24">{$t('action.actions')}</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#each prompts as prompt (prompt.id)}
        <Table.Row>
          {#if showFavorite}
            <Table.Cell>
              <button
                class="flex items-center justify-center {prompt.isFavorite ? 'text-yellow-500' : 'text-muted-foreground hover:text-yellow-500'}"
                onclick={() => handleFavorite(prompt)}
              >
                <Star class="h-4 w-4" fill={prompt.isFavorite ? 'currentColor' : 'none'} />
              </button>
            </Table.Cell>
          {/if}
          <Table.Cell class="font-medium">{prompt.name}</Table.Cell>
          <Table.Cell class="text-muted-foreground">{prompt.path}</Table.Cell>
          <Table.Cell>{prompt.latestVersion || '-'}</Table.Cell>
          <Table.Cell>
            <Badge variant={prompt.isPublish ? 'default' : 'secondary'}>
              {prompt.isPublish ? $t('prompt.published') : $t('prompt.unpublished')}
            </Badge>
          </Table.Cell>
          <Table.Cell>{prompt.username}</Table.Cell>
          <Table.Cell class="text-muted-foreground">{prompt.updateAt}</Table.Cell>
          <Table.Cell>
            <DropdownMenu.Root>
              <DropdownMenu.Trigger>
                <Button variant="ghost" size="icon" class="h-8 w-8">
                  <MoreHorizontal class="h-4 w-4" />
                </Button>
              </DropdownMenu.Trigger>
              <DropdownMenu.Content align="end">
                <DropdownMenu.Item onclick={() => handlePublish(prompt)}>
                  {#if prompt.isPublish}
                    <EyeOff class="mr-2 h-4 w-4" />
                    {$t('version.unpublish')}
                  {:else}
                    <Eye class="mr-2 h-4 w-4" />
                    {$t('version.publish')}
                  {/if}
                </DropdownMenu.Item>
                <DropdownMenu.Item>
                  <a href="/editor/{prompt.id}" class="flex items-center">
                    <Edit class="mr-2 h-4 w-4" />
                    {$t('prompt.edit')}
                  </a>
                </DropdownMenu.Item>
                <DropdownMenu.Separator />
                <DropdownMenu.Item class="text-destructive" onclick={() => handleDelete(prompt.id)}>
                  <Trash2 class="mr-2 h-4 w-4" />
                  {$t('prompt.delete')}
                </DropdownMenu.Item>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </Table.Cell>
        </Table.Row>
      {/each}
    </Table.Body>
  </Table.Root>
</div>
