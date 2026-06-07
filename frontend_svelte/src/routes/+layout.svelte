<script lang="ts">
  import '../app.css';
  import { Toaster } from '$lib/components';
  import logo from '$lib/assets/logo.png';

  let { children } = $props();
</script>

<svelte:head>
  <title>Al-Anwar Payment | Sistem Administrasi Keuangan</title>
  <link rel="icon" type="image/png" href={logo} />
</svelte:head>

{@render children()}

<Toaster />
