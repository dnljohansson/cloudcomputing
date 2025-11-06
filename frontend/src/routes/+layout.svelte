<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

	import './styles.css';
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<nav class="navbar">
    	<a href="/">HOME</a>
    	<a href="/about">ABOUT</a>
	</nav>
</svelte:head>

{@render children()}
