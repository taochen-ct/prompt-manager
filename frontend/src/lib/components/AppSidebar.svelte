<script module lang="ts">
  import type {IconNode} from "@lucide/svelte";

  export interface IMenuItem {
    title: string
    icon?: IconNode
    url?: string
    subItems?: ISubMenuItem[]
  }

  export interface ISubMenuItem {
    title: string
    url?: string
  }

  export type IModeToggle = "published" | "draft" | "developing"
</script>

<script lang="ts">
  import {t} from 'svelte-i18n';
  import {ChevronUp, Plus} from 'lucide-svelte';
  import FolderOpen from "@lucide/svelte/icons/folder";
  import FileText from "@lucide/svelte/icons/file-text";
  import Star from "@lucide/svelte/icons/star";
  import Clock from "@lucide/svelte/icons/clock";
  import Sparkles from "@lucide/svelte/icons/sparkles";
  import {goto} from '$app/navigation';
  import {Button} from "$lib/components/ui/button";
  import {Input} from "$lib/components/ui/input";

  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
  import * as Dialog from '$lib/components/ui/dialog';
  import * as Collapsible from "$lib/components/ui/collapsible/index.js";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import {api, auth} from "$lib/services/api";
  import type {Category} from "$lib/services/type";
  import {toast} from "svelte-sonner";
  import {onMount} from "svelte";

  interface Props {
    className?: string;
  }

  let {className}: Props = $props();
  let newPrompt = $state({
    name: '',
    path: '',
    createdBy: auth.getUser()?.username as string,
    username: auth.getUser()?.nickname as string,
  });

  // Auth state
  let isLoggedIn = $state(auth.isLoggedIn());
  let currentUser = $state(auth.getUser());

  async function handleLogout() {
    try {
      await api.logout();
      toast.success('Logged out successfully');
    } catch (error) {
      // Ignore logout errors
    }
    isLoggedIn = auth.isLoggedIn();
    currentUser = auth.getUser();
    await goto('/login');
  }


  let promptCategories = $state<Category[]>([]);
  let createDialogOpen = $state<boolean>(false);

  async function handleCreatePrompt() {
    try {
      const res = await api.createPrompt(newPrompt);
      toast.success($t('action.save') as string, {
        description: 'Prompt created successfully'
      });
      createDialogOpen = false;
      newPrompt = {
        name: '',
        path: '',
        createdBy: auth.getUser()?.username as string,
        username: auth.getUser()?.nickname as string
      };
      await goto(`/editor/${res.id}`, {})
    } catch (error) {
      toast.error($t('action.confirm') as string, {
        description: 'Failed to create prompt'
      });
    }
  }

  function getIconComponent(iconName: string) {
    const iconMap: Record<string, any> = {
      'file': FileText,
      'folder': FolderOpen,
      'sparkles': Sparkles
    };
    return iconMap[iconName] || FileText;
  }

  onMount(async () => {
    promptCategories = await api.getCategoryList();
  });
</script>

<Sidebar.Root class={className}>
  <Sidebar.Content>
    <Sidebar.Group>
      <Dialog.Root bind:open={createDialogOpen}>
        <Dialog.Trigger>
          <Button class="w-full">
            <Plus class="h-4 w-4 mr-2"/>
            {$t('prompt.create')}
          </Button>
        </Dialog.Trigger>
        <Dialog.Content>
          <Dialog.Header>
            <Dialog.Title>{$t('prompt.create')}</Dialog.Title>
            <Dialog.Description>
              Fill in the information to create a new prompt.
            </Dialog.Description>
          </Dialog.Header>
          <div class="grid gap-4 py-4">
            <div class="grid gap-2">
              <label for="name">{$t('prompt.name')}</label>
              <Input
                  id="name"
                  bind:value={newPrompt.name}
                  placeholder="Enter prompt name"
              />
            </div>
            <div class="grid gap-2">
              <label for="path">{$t('prompt.path')}</label>
              <Input
                  id="path"
                  bind:value={newPrompt.path}
                  placeholder="Enter prompt path"
              />
            </div>
          </div>
          <Dialog.Footer>
            <Button variant="outline" onclick={() => createDialogOpen = false}>
              {$t('action.cancel')}
            </Button>
            <Button onclick={handleCreatePrompt}>
              {$t('action.save')}
            </Button>
          </Dialog.Footer>
        </Dialog.Content>
      </Dialog.Root>
    </Sidebar.Group>

    <Sidebar.Group>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          <Sidebar.MenuItem>
            <a href="/"
               class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors">
              <FolderOpen class="h-4 w-4"/>
              <span>{$t('sidebar.all')}</span>
            </a>
          </Sidebar.MenuItem>

          <Sidebar.MenuItem>
            <a href="/favorites"
               class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors">
              <Star class="h-4 w-4"/>
              <span>{$t('sidebar.favorites')}</span>
            </a>
          </Sidebar.MenuItem>

          <Sidebar.MenuItem>
            <a href="/recent"
               class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors">
              <Clock class="h-4 w-4"/>
              <span>{$t('sidebar.recent')}</span>
            </a>
          </Sidebar.MenuItem>


          <Sidebar.MenuItem>
            <Collapsible.Root open class="group/collapsible">
              <Collapsible.Trigger>
                {#snippet child({props})}
                  <a {...props}
                     class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors">
                    <FolderOpen class="h-4 w-4"/>
                    <span>{$t('prompt.tags')}</span>
                  </a>
                {/snippet}
              </Collapsible.Trigger>
              <Collapsible.Content>
                <Sidebar.MenuSub>
                  {#each promptCategories as category}
                    {@const Component = getIconComponent(category.icon)}
                    <Sidebar.MenuSubItem>
                      <a href={category.url}
                         class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors">
                        <Component class="h-4 w-4"/>
                        <span class="flex-1">{category.title}</span>
                        <span class="text-xs text-muted-foreground">{category.count}</span>
                      </a>
                    </Sidebar.MenuSubItem>
                  {/each}
                </Sidebar.MenuSub>
              </Collapsible.Content>
            </Collapsible.Root>
          </Sidebar.MenuItem>
        </Sidebar.Menu>
      </Sidebar.GroupContent>
    </Sidebar.Group>
  </Sidebar.Content>

  <Sidebar.Footer>
    {#if isLoggedIn && currentUser}
      <Sidebar.Menu>
        <Sidebar.MenuItem>
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              {#snippet child({props})}
                <Sidebar.MenuButton
                    {...props}
                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <div class="flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 text-primary">
                    <span
                        class="text-xs font-medium">{currentUser?.nickname?.charAt(0) || currentUser?.username?.charAt(0) || 'U'}</span>
                  </div>
                  <span class="font-medium">{currentUser?.nickname || currentUser?.username}</span>
                  <ChevronUp class="ms-auto h-4 w-4"/>
                </Sidebar.MenuButton>
              {/snippet}
            </DropdownMenu.Trigger>
            <DropdownMenu.Content side="top" class="w-(--bits-dropdown-menu-anchor-width)">
              <DropdownMenu.Item>
                <span>{$t('user.profile')}</span>
              </DropdownMenu.Item>
              <DropdownMenu.Item>
                <span>{$t('settings.title')}</span>
              </DropdownMenu.Item>
              <DropdownMenu.Separator/>
              <DropdownMenu.Item class="text-destructive" onclick={handleLogout}>
                <span>{$t('auth.logout')}</span>
              </DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Sidebar.MenuItem>
      </Sidebar.Menu>
    {:else}
      <Sidebar.Menu>
        <Sidebar.MenuItem>
          <a href="/login"
             class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent hover:text-accent-foreground transition-colors">
            <div class="flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 text-primary">
              <span class="text-xs font-medium">?</span>
            </div>
            <span class="font-medium">Login</span>
          </a>
        </Sidebar.MenuItem>
      </Sidebar.Menu>
    {/if}
  </Sidebar.Footer>
</Sidebar.Root>
