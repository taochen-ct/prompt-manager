<script lang="ts">
  import { t } from 'svelte-i18n';
  import { toast } from 'svelte-sonner';
  import { api } from '$lib/services/api';
  import RichTextEditor from '$lib/components/RichTextEditor.svelte';
  import type { PageData } from './$types';
  import * as Breadcrumb from "$lib/components/ui/breadcrumb";
  import { Separator } from "$lib/components/ui/separator";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Select from "$lib/components/ui/select";
  import * as Tabs from "$lib/components/ui/tabs";
  import Badge from '$lib/components/ui/badge/badge.svelte';

  let { data }: { data: PageData } = $props();

  interface Version {
    id: string;
    version: string;
    content: string;
    variables: string;
    changeLog: string;
    isPublish: boolean;
    username: string;
    createdAt: string;
    updatedAt: string;
  }

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
  }

  let prompt = $state<Prompt | null>(null);
  let versions = $state<Version[]>([]);
  let selectedVersion = $state<Version | null>(null);
  let loading = $state(true);
  let saving = $state(false);

  let createVersionDialogOpen = $state(false);
  let newVersion = $state({
    version: '',
    content: '',
    variables: '',
    changeLog: '',
    isPublish: false
  });

  async function loadPromptData() {
    try {
      loading = true;
      const [promptData, versionsData] = await Promise.all([
        api.getPrompt(data.id),
        api.getPromptVersions(data.id)
      ]);
      prompt = promptData as Prompt;
      versions = versionsData as Version[];

      if (versions.length > 0) {
        selectedVersion = versions[0];
        newVersion.content = selectedVersion?.content || '';
        newVersion.version = incrementVersion(selectedVersion?.version || '1.0.0');
      }
    } catch (error) {
      console.error('Failed to load prompt:', error);
      toast.error($t('action.confirm') as string, {
        description: 'Failed to load prompt data'
      });
    } finally {
      loading = false;
    }
  }

  function incrementVersion(version: string): string {
    const parts = version.split('.');
    const last = parseInt(parts[parts.length - 1]) + 1;
    parts[parts.length - 1] = last.toString();
    return parts.join('.');
  }

  async function handleSaveVersion() {
    if (!prompt) return;

    try {
      saving = true;
      await api.createVersion({
        promptId: prompt.id,
        version: newVersion.version,
        content: newVersion.content,
        variables: newVersion.variables,
        changeLog: newVersion.changeLog,
        createdBy: 'user-1',
        username: 'Admin',
        isPublish: newVersion.isPublish
      });

      toast.success($t('action.save') as string, {
        description: 'Version saved successfully'
      });

      createVersionDialogOpen = false;
      newVersion = {
        version: incrementVersion(newVersion.version),
        content: '',
        variables: '',
        changeLog: '',
        isPublish: false
      };

      await loadPromptData();
    } catch (error) {
      toast.error($t('action.confirm') as string, {
        description: 'Failed to save version'
      });
    } finally {
      saving = false;
    }
  }

  async function handlePublish(version: Version) {
    try {
      await api.updateVersion({
        id: version.id,
        version: version.version,
        content: version.content,
        variables: version.variables,
        changeLog: version.changeLog,
        isPublish: true
      });

      toast.success($t('version.publish') as string, {
        description: 'Version published successfully'
      });

      await loadPromptData();
    } catch (error) {
      toast.error($t('action.confirm') as string, {
        description: 'Failed to publish version'
      });
    }
  }

  function selectVersion(version: Version) {
    selectedVersion = version;
    newVersion.content = version.content;
  }

  $effect(() => {
    loadPromptData();
  });
</script>

<header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Page>
          {prompt?.name || data.id}
        </Breadcrumb.Page>
      </Breadcrumb.Item>
      {#if selectedVersion}
        <Breadcrumb.Separator />
        <Breadcrumb.Item>
          <Breadcrumb.Page>
            v{selectedVersion.version}
          </Breadcrumb.Page>
        </Breadcrumb.Item>
      {/if}
    </Breadcrumb.List>
  </Breadcrumb.Root>
  <div class="ml-auto flex items-center gap-2">
    {#if prompt}
      <Badge variant={prompt.isPublish ? 'default' : 'secondary'}>
        {prompt.isPublish ? $t('prompt.published') : $t('prompt.unpublished')}
      </Badge>
    {/if}
    <Dialog.Root bind:open={createVersionDialogOpen}>
      <Dialog.Trigger>
        <Button>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2">
            <path d="M5 12h14"/>
            <path d="M12 5v14"/>
          </svg>
          {$t('version.create')}
        </Button>
      </Dialog.Trigger>
      <Dialog.Content class="max-w-2xl">
        <Dialog.Header>
          <Dialog.Title>{$t('version.create')}</Dialog.Title>
          <Dialog.Description>
            Create a new version for this prompt.
          </Dialog.Description>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <label for="version">{$t('version.version')}</label>
              <Input id="version" bind:value={newVersion.version} placeholder="1.0.0" />
            </div>
            <div class="grid gap-2">
              <label for="publish">{$t('version.publish')}</label>
              <Select.Root type="single" value={newVersion.isPublish.toString()}>
                <Select.Trigger>
                  {newVersion.isPublish ? $t('version.publish') : $t('version.unpublish')}
                </Select.Trigger>
                <Select.Content>
                  <Select.Item value="true">{$t('version.publish')}</Select.Item>
                  <Select.Item value="false">{$t('version.unpublish')}</Select.Item>
                </Select.Content>
              </Select.Root>
            </div>
          </div>
          <div class="grid gap-2">
            <label for="content">{$t('version.content')}</label>
            <Textarea id="content" bind:value={newVersion.content} rows="6" />
          </div>
          <div class="grid gap-2">
            <label for="variables">{$t('version.variables')}</label>
            <Input id="variables" bind:value={newVersion.variables} placeholder='["topic", "tone"]' />
          </div>
          <div class="grid gap-2">
            <label for="changeLog">{$t('version.changeLog')}</label>
            <Input id="changeLog" bind:value={newVersion.changeLog} placeholder="What changed?" />
          </div>
        </div>
        <Dialog.Footer>
          <Button variant="outline" onclick={() => createVersionDialogOpen = false}>
            {$t('action.cancel')}
          </Button>
          <Button onclick={handleSaveVersion} disabled={saving}>
            {saving ? $t('prompt.loading') : $t('action.save')}
          </Button>
        </Dialog.Footer>
      </Dialog.Content>
    </Dialog.Root>
  </div>
</header>

<Tabs.Root value="editor" class="flex flex-1 flex-col">
  <Tabs.List class="px-4">
    <Tabs.Trigger value="editor">{$t('version.content')}</Tabs.Trigger>
    <Tabs.Trigger value="versions">
      {$t('version.title')}
      <Badge variant="secondary" class="ml-2">{versions.length}</Badge>
    </Tabs.Trigger>
  </Tabs.List>

  <Tabs.Content value="editor" class="flex-1 p-4">
    <RichTextEditor initialHtml={newVersion.content} />
  </Tabs.Content>

  <Tabs.Content value="versions" class="flex-1 p-4">
    {#if loading}
      <div class="flex items-center justify-center py-12">
        <p>{$t('prompt.loading')}</p>
      </div>
    {:else if versions.length === 0}
      <div class="flex flex-col items-center justify-center py-12 text-center">
        <p class="text-muted-foreground">{$t('version.noVersions')}</p>
        <Button class="mt-4" onclick={() => createVersionDialogOpen = true}>
          {$t('version.createFirst')}
        </Button>
      </div>
    {:else}
      <div class="space-y-2">
        {#each versions as version (version.id)}
          <div
            class="flex items-center justify-between rounded-lg border p-4 transition-colors hover:bg-muted/50
              {selectedVersion?.id === version.id ? 'border-primary bg-primary/5' : ''}"
          >
            <div class="flex items-center gap-4">
              <button
                class="flex h-10 w-10 items-center justify-center rounded-full text-sm font-medium
                  {version.isPublish ? 'bg-green-100 text-green-700' : 'bg-muted text-muted-foreground'}"
                onclick={() => selectVersion(version)}
              >
                v{version.version}
              </button>
              <div>
                <p class="font-medium">{version.changeLog || '-'}</p>
                <p class="text-sm text-muted-foreground">
                  {version.username} · {version.updatedAt}
                </p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <Badge variant={version.isPublish ? 'default' : 'secondary'}>
                {version.isPublish ? $t('prompt.published') : $t('prompt.unpublished')}
              </Badge>
              {#if !version.isPublish}
                <Button variant="outline" size="sm" onclick={() => handlePublish(version)}>
                  {$t('version.publish')}
                </Button>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </Tabs.Content>
</Tabs.Root>
