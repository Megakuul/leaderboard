<script>
  import { Toaster } from "$lib/components/ui/sonner";
  import Rift from "$lib/components/custom/Rift.svelte";
  import "./app.css";
  import Toolbar from "./Toolbar.svelte";
  import { onMount } from "svelte";
  import { GetTokens } from "$lib/api/auth";

  onMount(() => {
    // Fetch access and id_token if provided by the cognito callback
    // and moves them into the local storage.
    GetTokens();
  })

  /** @type {HTMLHeadingElement} */
  let leaderboardTitle;
</script>

<Toaster />

<main class="m-0 p-0 bg-slate-900 min-h-screen w-screen overflow-hidden">
  <Rift height="1000" width="1000"></Rift>
  
  <section class="absolute max-w-full bg-gray-950 top-1/4 translate-y-1/2 right-1/2 xl:right-1/4 translate-x-1/2 transition-all p-10 sm:p-16 rounded-lg bg-opacity-100 xl:bg-opacity-80 shadow-xl select-none">
    <p class="text-sm sm:text-lg text-red-900 font-bold">Join the official</p>
    <h1 class="text-4xl sm:text-5xl xl:text-7xl 2xl:text-8xl text-center text-transparent font-bold text-shadow lg:mx-4 lg:my-2 bg-clip-text bg-gradient-to-br from-red-900 to-indigo-600 text-opacity-70">Leaderboard</h1>
    <p class="w-full text-sm sm:text-lg text-indigo-600 text-end font-bold">fighting for glory and honor</p>
  </section>
  <section class="flex flex-col items-center">
    <Toolbar leaderboardTitle={leaderboardTitle}></Toolbar>
    <h1 bind:this={leaderboardTitle} class="top-96 text-9xl">Leaderboard</h1>
  </section>

</main>

<style>
  main {
    @apply bg-repeat;
    background-image: url('./assets/grid.svg');
  }

  .text-shadow {
    text-shadow: 5px 5px 5px rgba(0,0,0,0.1);
  }
</style>