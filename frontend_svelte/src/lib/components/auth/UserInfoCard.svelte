<script lang="ts">
  import { onDestroy } from 'svelte';
  import { getUserProfile, getAuthSession } from '$lib/context/auth-context';

  const { displayName, email, roleName } = getUserProfile();
  const { isAuthenticated } = getAuthSession();

  let localSensitiveData: string | null = null;

  onDestroy(() => {
    localSensitiveData = null;
  });
</script>

{#if $isAuthenticated}
  <div class="bg-white border border-gray-200 rounded-xl p-4 shadow-sm">
    <div class="flex items-center gap-3">
      <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
        <span class="text-green-700 font-bold text-sm">
          {$displayName?.charAt(0)?.toUpperCase() || '?'}
        </span>
      </div>
      <div>
        <p class="text-sm font-bold text-gray-900">{$displayName || 'User'}</p>
        <p class="text-xs text-gray-500">{$email || '-'}</p>
        <span class="inline-block mt-0.5 text-xs font-semibold px-2 py-0.5 rounded-full bg-green-100 text-green-700">
          {$roleName || 'unknown'}
        </span>
      </div>
    </div>
  </div>
{/if}
