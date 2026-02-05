<script lang="ts">
  import {onMount, onDestroy, createEventDispatcher} from 'svelte';
  import {EditorView, minimalSetup, basicSetup} from 'codemirror';
  import {EditorState, StateEffect} from '@codemirror/state';
  import {keymap} from '@codemirror/view';
  import {defaultKeymap, indentWithTab} from '@codemirror/commands';
  import {markdown} from '@codemirror/lang-markdown';
  import {oneDark} from '@codemirror/theme-one-dark';
  import {highlightSelectionMatches, searchKeymap} from '@codemirror/search';
  import {closeBrackets, autocompletion, closeBracketsKeymap, completionKeymap} from '@codemirror/autocomplete';
  import {history, historyKeymap} from '@codemirror/commands';
  import {indentOnInput, bracketMatching} from '@codemirror/language';
  import {cn} from '$lib/utils';
  import Badge from '$lib/components/ui/badge/badge.svelte';
  import {Compartment} from '@codemirror/state';

  interface Props {
    class?: string;
    key: string;
    language?: string;
    initialCode?: string;
    readonly?: boolean;
    theme?: 'light' | 'dark';
  }

  let {
    class: className,
    language = 'markdown',
    initialCode = '',
    readonly = false,
    theme = 'light'
  }: Props = $props();

  const dispatch = createEventDispatcher();

  let dom: HTMLDivElement | undefined = $state();
  let view: EditorView | null = $state(null);
  let _mounted = false;
  let code = $state(initialCode);
  let copied = $state(false);
  const readOnlyCompartment = new Compartment();
  const themeCompartment = new Compartment();

  // 创建编辑器
  function createEditor() {
    if (!dom) return;

    const extensions = [
      minimalSetup,
      basicSetup,
      keymap.of([indentWithTab, ...defaultKeymap, ...historyKeymap, ...closeBracketsKeymap, ...completionKeymap, ...searchKeymap]),
      highlightSelectionMatches(),
      closeBrackets(),
      autocompletion(),
      bracketMatching(),
      indentOnInput(),
      EditorView.lineWrapping,

      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          code = update.state.doc.toString();
          dispatch('change', {code});
        }
      }),

      readOnlyCompartment.of(EditorState.readOnly.of(readonly)),
      themeCompartment.of(theme === 'dark' ? oneDark : [])
    ];

    view = new EditorView({
      state: EditorState.create({
        doc: initialCode,
        extensions
      }),
      parent: dom
    });
    code = initialCode;
  }

  onMount(() => {
    _mounted = true;
    createEditor();
  });

  onDestroy(() => {
    if (view) {
      view.destroy();
      view = null;
    }
  });

  // 监听只读状态变化
  $effect(() => {
    if (view) {
      view.dispatch({
        effects: readOnlyCompartment.reconfigure(
            EditorState.readOnly.of(readonly)
        )
      });
    }
  });

  // 监听主题变化
  $effect(() => {
    if (view) {
      view.dispatch({
        effects: themeCompartment.reconfigure(
            theme === 'dark' ? oneDark : []
        )
      });
    }
  });

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(code);
      copied = true;
      setTimeout(() => {
        copied = false;
      }, 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  }

  // 语言标签映射
  const languageLabels: Record<string, string> = {
    markdown: 'Markdown',
    javascript: 'JavaScript',
    typescript: 'TypeScript',
    python: 'Python',
    json: 'JSON',
    text: 'Plain Text'
  };
</script>

<div class={cn(
  "mx-auto w-full overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm ring-1 ring-slate-200",
  className
)}>
  <!-- Editor Header -->
  <div class="flex items-center justify-between border-b border-slate-200 bg-slate-50/50 px-4 py-2">
    <div class="flex items-center gap-3">
      <span class="flex items-center gap-1.5 text-sm font-medium text-slate-700">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
             stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
             class="text-slate-500">
          <polyline points="16 18 22 12 16 6"/>
          <polyline points="8 6 2 12 8 18"/>
        </svg>
        {languageLabels[language] || language}
      </span>
      <div class="h-4 w-px bg-slate-300"></div>
      <span class="text-xs text-slate-500">
        {code.length} characters
      </span>
      <div class="h-4 w-px bg-slate-300"></div>
      <span class="text-xs text-slate-500">
        {code.split('\n').length} lines
      </span>
      {#if readonly}
        <div class="h-4 w-px bg-slate-300"></div>
        <Badge variant="default" class="text-[10px]">Read only</Badge>
      {/if}
    </div>

    <div class="flex items-center gap-2">
      <button
          onclick={copyToClipboard}
          class="flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-medium text-slate-600 transition-colors hover:bg-slate-200 hover:text-slate-900 active:bg-slate-300"
          title="Copy code"
      >
        {#if copied}
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
               class="text-green-600">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
          <span class="text-green-600">Copied!</span>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect width="14" height="14" x="8" y="8" rx="2" ry="2"/>
            <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>
          </svg>
          Copy
        {/if}
      </button>
    </div>
  </div>

  <!-- CodeMirror Editor -->
  <div
      class="codemirror-editor min-h-fit overflow-hidden"
      class:opacity-50={readonly}
      class:cursor-not-allowed={readonly}
  >
    <div bind:this={dom} class="h-full w-full"></div>
  </div>

</div>

<style>
    @import "tailwindcss";

    .codemirror-editor :global(.cm-editor) {
        height: 100%;
        font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace;
        font-size: 14px;
    }

    .codemirror-editor :global(.cm-editor.cm-focused) {
        outline: none;
    }

    .codemirror-editor :global(.cm-scroller) {
        overflow: auto;
        font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace;
    }

    /* 自定义滚动条 */
    .codemirror-editor :global(::-webkit-scrollbar) {
        width: 8px;
        height: 8px;
    }

    .codemirror-editor :global(::-webkit-scrollbar-track) {
        background: transparent;
    }

    .codemirror-editor :global(::-webkit-scrollbar-thumb) {
        background-color: #cbd5e1;
        border-radius: 4px;
    }

    .codemirror-editor :global(::-webkit-scrollbar-thumb:hover) {
        background-color: #94a3b8;
    }

    /* 只读状态 */
    .cursor-not-allowed {
        cursor: not-allowed;
    }

    .opacity-50 {
        opacity: 0.5;
    }
</style>
