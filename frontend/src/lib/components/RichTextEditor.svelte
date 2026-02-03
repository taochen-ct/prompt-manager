<script lang="ts">
  import {cn} from "$lib/utils";

  interface Props {
    class?: string;
    initialHtml?: string;
  }

  let {class: className, initialHtml = '<p>开始编辑...</p>'}: Props = $props();
  let html = $state(initialHtml);
  let editorRef: HTMLDivElement | undefined = $state();
  let fileInputRef: HTMLInputElement | undefined = $state();

  const formats = [
    {name: '正文', value: 'formatBlock', arg: '<p>'},
    {name: '标题 1', value: 'formatBlock', arg: '<h1>'},
    {name: '标题 2', value: 'formatBlock', arg: '<h2>'},
    {name: '无序列表', value: 'insertUnorderedList', arg: null},
    {name: '有序列表', value: 'insertOrderedList', arg: null}
  ];

  function execute(command: string, arg: string | null = null) {
    document.execCommand(command, false, arg);
    if (editorRef) html = editorRef.innerHTML;
  }

  function handleFileChange(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      const base64 = event.target?.result as string;
      execute('insertImage', base64);
      target.value = '';
    };
    reader.readAsDataURL(file);
  }
</script>


{#snippet editorHeader()}
  <div class="flex flex-wrap items-center gap-1 border-b border-slate-200 bg-slate-50/50 p-2">

    <select
        onchange={(e) => {
        const [cmd, arg] = e.currentTarget.value.split(':');
        execute(cmd, arg);
      }}
        class="h-9 rounded-md border-slate-200 bg-white px-2 text-sm font-medium text-slate-700 shadow-sm outline-none hover:bg-slate-50 focus:ring-2 focus:ring-indigo-500/20"
    >
      {#each formats as item}
        <option value={`${item.value}:${item.arg || ''}`}>{item.name}</option>
      {/each}
    </select>

    <div class="mx-1 h-6 w-px bg-slate-300"></div>

    <div class="flex gap-0.5">
      {#each [['bold', 'B', '加粗'], ['italic', 'I', '斜体'], ['underline', 'U', '下划线']] as [cmd, label, title]}
        <button
            onclick={() => execute(cmd)}
            title={title}
            class="flex h-9 w-9 items-center justify-center rounded-md text-slate-600 transition-colors hover:bg-slate-200 hover:text-slate-900 active:bg-slate-300"
        >
          <span class={cmd === 'bold' ? 'font-bold' : cmd === 'italic' ? 'italic' : 'underline'}>
            {label}
          </span>
        </button>
      {/each}
    </div>

    <div class="mx-1 h-6 w-px bg-slate-300"></div>

    <button
        onclick={() => fileInputRef?.click()}
        class="flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium text-indigo-600 transition-colors hover:bg-indigo-50 active:bg-indigo-100"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
           stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect width="18" height="18" x="3" y="3" rx="2" ry="2"/>
        <circle cx="9" cy="9" r="2"/>
        <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/>
      </svg>
      图片
    </button>

    <input
        type="file"
        accept="image/*"
        bind:this={fileInputRef}
        onchange={handleFileChange}
        class="hidden"
    />
  </div>
{/snippet}

<div
    class={cn(
        "mx-auto w-full overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm ring-1 ring-slate-200",
        className
    )}>
  {@render editorHeader()}
  <div
      contenteditable="true"
      bind:this={editorRef}
      bind:innerHTML={html}
      class="prose prose-slate max-w-none min-h-[400px] p-6 outline-none focus:ring-0"
      oninput={() => html = editorRef?.innerHTML || ''}
  ></div>
</div>

<div class="mt-2 flex justify-between text-[10px] uppercase tracking-wider text-slate-400">
  <!--  <span>Svelte 5 Editor</span>-->
  <span>Characters: {html.replace(/<[^>]*>/g, '').length}</span>
</div>

<style>
    @import "tailwindcss";

    .prose :global(h1) {
        @apply text-3xl font-extrabold mb-4 mt-2;
    }

    .prose :global(h2) {
        @apply text-2xl font-bold mb-3 mt-2;
    }

    .prose :global(ul) {
        @apply list-disc pl-5 mb-4;
    }

    .prose :global(ol) {
        @apply list-decimal pl-5 mb-4;
    }

    .prose :global(img) {
        @apply max-w-full rounded-lg shadow-md my-4;
    }

    .prose :global(p) {
        @apply mb-2 leading-relaxed;
    }
</style>