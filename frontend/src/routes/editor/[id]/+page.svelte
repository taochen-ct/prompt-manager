<script lang="ts">
  import {t} from 'svelte-i18n';
  import {toast} from 'svelte-sonner';
  import {api, auth} from '$lib/services/api';
  import CodeEditor from '$lib/components/CodeEditor.svelte';
  import {cn} from '$lib/utils';
  import type {PageData} from './$types';
  import * as Breadcrumb from "$lib/components/ui/breadcrumb";
  import {Separator} from "$lib/components/ui/separator";
  import {Button} from "$lib/components/ui/button";
  import {Input} from "$lib/components/ui/input";
  import {Textarea} from "$lib/components/ui/textarea";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Select from "$lib/components/ui/select";
  import * as Tabs from "$lib/components/ui/tabs";
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import {Checkbox} from "$lib/components/ui/checkbox";
  import {Label} from "$lib/components/ui/label";

  let {data}: { data: PageData } = $props();

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
    variables: null,
    changeLog: null,
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
        createdBy: auth.getUser()?.username as string,
        username: auth.getUser()?.nickname as string,
        isPublish: newVersion.isPublish
      });

      toast.success($t('action.save') as string, {
        description: 'Version saved successfully'
      });

      createVersionDialogOpen = false;
      newVersion = {
        version: incrementVersion(newVersion.version),
        content: '',
        variables: null,
        changeLog: null,
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

<header class="flex h-16 shrink-0 items-center justify-between border-b bg-white px-6">
  <div class="flex items-center gap-4">
    <Breadcrumb.Root>
      <Breadcrumb.List>
        <Breadcrumb.Item>
          <Breadcrumb.Page>
            {prompt?.name || data.id}
          </Breadcrumb.Page>
        </Breadcrumb.Item>
        {#if selectedVersion}
          <Breadcrumb.Separator/>
          <Breadcrumb.Item>
            <Breadcrumb.Page>
              {selectedVersion.version}
            </Breadcrumb.Page>
          </Breadcrumb.Item>
        {/if}
      </Breadcrumb.List>
    </Breadcrumb.Root>
    {#if prompt}
      <Badge variant={prompt.isPublish ? 'default' : 'secondary'}>
        {prompt.isPublish ? $t('prompt.published') : $t('prompt.unpublished')}
      </Badge>
    {/if}
  </div>

  <div class="flex items-center gap-3">
    <!-- Version Selector -->
    <Dialog.Root bind:open={createVersionDialogOpen}>
      <Dialog.Trigger>
        <Button>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2">
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
            {$t('version.dialog.description')}
          </Dialog.Description>
        </Dialog.Header>
        <div class="grid gap-4 py-4">
          <div class="grid gap-2">
            <Label for="version">{$t('version.version')}</Label>
            <Input id="version" bind:value={newVersion.version} placeholder="1.0.0"/>
          </div>
          <div class="grid gap-2">
            <Label for="content">{$t('version.content')}</Label>
            <Textarea id="content" bind:value={newVersion.content} class="font-mono text-sm"/>
          </div>
          <div class="grid gap-2">
            <Label for="variables">{$t('version.variables')}</Label>
            <Input id="variables" bind:value={newVersion.variables} placeholder='["topic", "tone"]'/>
          </div>
          <div class="grid gap-2">
            <Label for="changeLog">{$t('version.changeLog')}</Label>
            <Textarea id="changeLog" bind:value={newVersion.changeLog} placeholder="What changed?"/>
          </div>
          <div class="flex justify-between content-center gap-2">
            <Label>{$t('version.isPublish')}</Label>
            <Checkbox id="isPublish" bind:checked={newVersion.isPublish} />
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

<!-- Main Content Area -->
<div class="flex flex-1 overflow-hidden">
  <!-- Code Editor Panel -->
  <div class="flex-1 overflow-auto border-r border-slate-200 bg-white p-4">
    {#if selectedVersion}
      <div class="mb-4 flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">{$t('version.editor.content')}</h2>
          <p class="text-sm text-slate-500">
            Last updated by {selectedVersion.username} · {selectedVersion.updatedAt}
          </p>
        </div>
        <div class="flex items-center gap-2">
          <Badge variant="outline">
            {selectedVersion.version}
          </Badge>
          {#if selectedVersion.isPublish}
            <Badge variant="default">Published</Badge>
          {:else}
            <Badge variant="secondary">Draft</Badge>
          {/if}
        </div>
      </div>
      {#key selectedVersion.id}
        <CodeEditor
            key={selectedVersion.id}
            language="markdown"
            initialCode={selectedVersion.content}
            on:change={e => {newVersion.content = e.detail.code}}
            class="min-h-[500px] h-full"
        />
      {/key}
    {:else}
      <div class="flex flex-col items-center justify-center py-20 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none"
             stroke="currentColor" stroke-width="1" class="mb-4 text-slate-300">
          <polyline points="16 18 22 12 16 6"/>
          <polyline points="8 6 2 12 8 18"/>
        </svg>
        <h3 class="mb-2 text-lg font-medium text-slate-900">No version selected</h3>
        <p class="mb-4 text-slate-500">Select a version from the sidebar or create a new one</p>
        <Button onclick={() => createVersionDialogOpen = true}>
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2">
            <path d="M5 12h14"/>
            <path d="M12 5v14"/>
          </svg>
          {$t('version.createFirst')}
        </Button>
      </div>
    {/if}
  </div>

  <!-- Version History Sidebar -->
  <div class="w-96 shrink-0 overflow-auto bg-slate-50">
    <div class="border-b border-slate-200 bg-white px-4 py-3">
      <h3 class="font-medium text-slate-900">{$t('version.editor.history')}</h3>
      <p class="text-xs text-slate-500">{versions.length} versions</p>
    </div>

    {#if loading}
      <div class="flex items-center justify-center py-12">
        <p class="text-sm text-slate-500">{$t('prompt.loading')}</p>
      </div>
    {:else if versions.length === 0}
      <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
        <p class="mb-4 text-sm text-slate-500">{$t('version.noVersions')}</p>
        <Button class="w-full" onclick={() => createVersionDialogOpen = true}>
          {$t('version.createFirst')}
        </Button>
      </div>
    {:else}
      <div class="divide-y divide-slate-200">
        {#each versions as version (version.id)}
          <button
              class="w-full cursor-pointer px-4 py-3 text-left transition-colors hover:bg-slate-100
              {selectedVersion?.id === version.id ? 'bg-indigo-50 hover:bg-indigo-100' : ''}"
              onclick={() => selectVersion(version)}
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class={cn(
                  "flex h-8 w-8 items-center justify-center rounded-full text-xs font-medium",
                  version.isPublish
                    ? "bg-green-100 text-green-700"
                    : "bg-slate-200 text-slate-600"
                )}>
                  v{version.version}
                </div>
                <div>
                  <p class="text-sm font-medium text-slate-900">
                    {version.changeLog || 'No description'}
                  </p>
                  <p class="text-xs text-slate-500">
                    {version.username} · {version.updatedAt}
                  </p>
                </div>
              </div>
              {#if version.isPublish}
                <Badge variant="secondary" class="text-[10px]"}>
                  Published
                </Badge>
              {/if}
              {#if !version.isPublish}
                <Badge variant="secondary" class="text-[10px]">
                  Draft
                </Badge>
              {/if}
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>
