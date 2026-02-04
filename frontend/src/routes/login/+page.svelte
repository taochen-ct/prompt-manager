<script lang="ts">
  import {goto} from '$app/navigation';
  import {api, auth} from '$lib/services/api';
  import {Button} from '$lib/components/ui/button';
  import {Input} from '$lib/components/ui/input';
  import {toast} from 'svelte-sonner';

  let username = $state('');
  let password = $state('');
  let loading = $state(false);

  async function handleLogin() {
    if (!username || !password) {
      toast.error('Please enter username and password');
      return;
    }

    loading = true;
    try {
      await api.login({username, password});
      toast.success('Login successful');
      await goto('/');
    } catch (error) {
      toast.error('Login failed', {
        description: error instanceof Error ? error.message : 'Invalid credentials'
      });
    } finally {
      loading = false;
    }
  }

  // Redirect if already logged in
  $effect(() => {
    if (auth.isLoggedIn()) {
      goto('/');
    }
  });
</script>

<div class="flex items-center justify-center bg-background">
  <div class="w-full max-w-sm p-6 space-y-6">
    <div class="space-y-2 text-center">
      <h1 class="text-2xl font-bold">Login</h1>
      <p class="text-muted-foreground">Enter your credentials to access</p>
    </div>

    <div class="space-y-4">
      <div class="space-y-2">
        <label for="username" class="text-sm font-medium">Username</label>
        <Input
            id="username"
            bind:value={username}
            placeholder="Enter username"
            disabled={loading}
        />
      </div>

      <div class="space-y-2">
        <label for="password" class="text-sm font-medium">Password</label>
        <Input
            id="password"
            type="password"
            bind:value={password}
            placeholder="Enter password"
            disabled={loading}
        />
      </div>

      <Button class="w-full" onclick={handleLogin} disabled={loading}>
        {loading ? 'Logging in...' : 'Login'}
      </Button>
    </div>
  </div>
</div>
