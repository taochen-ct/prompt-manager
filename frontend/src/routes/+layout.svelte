<script lang="ts">
  import './i18n.ts';
  import './layout.css';
  import favicon from '$lib/assets/favicon.svg';
  import {waitLocale} from "svelte-i18n";
  import {ModeWatcher} from "mode-watcher";
  import {Toaster, toast} from 'svelte-sonner'
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";
  import {Separator} from "$lib/components/ui/separator/index.js";
  import AppSidebar from "$lib/components/AppSidebar.svelte";
  import Header from "$lib/components/Header.svelte";
  import type {LayoutData} from './$types';
  import type {Snippet} from "svelte";
  import {page} from '$app/state';

  let {data, children}: { data: LayoutData, children: Snippet } = $props();
  let isLoginPage = $derived(page.url.pathname === '/login');
</script>

<svelte:head>
  <link rel="icon" href={favicon}/>
</svelte:head>

{#await waitLocale()}
  <div class="loading-container">
    <p>Loading language...</p>
  </div>
{:then}
  {#if !isLoginPage}
    <Sidebar.Provider>
      <Toaster/>
      <ModeWatcher/>
      <AppSidebar/>
      <Sidebar.Inset>
        <Header/>
        {@render children()}
      </Sidebar.Inset>
    </Sidebar.Provider>
  {:else}
    {@render children()}
  {/if}
{/await}
