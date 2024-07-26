<script>
  import { Toaster } from "$lib/components/ui/sonner";
  import Rift from "$lib/components/custom/Rift.svelte";
  import "./app.css";
  import Toolbar from "./Toolbar.svelte";
  import { onMount } from "svelte";
  import { GetTokens } from "$lib/api/auth";
  import { Button } from "$lib/components/ui/button";
  import Sparta from "./assets/Sparta.svelte";
  import Algorithm from "./assets/Algorithm.svelte";
  import Leaderboard from "./Leaderboard.svelte";

  onMount(() => {
    // Fetch id_token if provided by the cognito callback
    // and moves them into the local storage.
    GetTokens();

    tokenExpirationTime = new Date(localStorage.getItem("id_token_exp"))
    tokenUsername = localStorage.getItem("id_token_username");
  })

  /** @type {HTMLHeadingElement} */
  let leaderboardTitle;

  /** @type {Date} */
  let tokenExpirationTime;

  /** @type {string} */
  let tokenUsername;
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
    <Toolbar leaderboardTitle={leaderboardTitle} tokenExpirationTime={tokenExpirationTime} tokenUsername={tokenUsername}></Toolbar>

    <div class="w-full my-32 flex flex-col justify-around">

      <div class="my-16 flex flex-col lg:flex-row items-center justify-around">
        <h1 class="text-3xl sm:text-5xl font-bold text-shadow text-slate-200 animate-pulse">Track your Games</h1>
        <div class="ml-12 my-12 p-8 w-96 h-96 rounded-lg bg-gray-950 border-black border-4 card right">
          <p class="text-2xl text-transparent bg-clip-text bg-gradient-to-br from-red-900 to-indigo-600 text-opacity-70">
            Log your Game Results and watch them shape the Leaderboard
          </p>
          <br>
          <Button variant="outline" on:click={() => leaderboardTitle.scrollIntoView({ behavior: "smooth" })}>Getting started</Button>
        </div>
      </div>

      <div class="my-16 flex flex-col lg:flex-row-reverse items-center justify-around">
        <h1 class="text-3xl sm:text-5xl font-bold text-shadow text-slate-200 animate-pulse">Compete against others</h1>
        <div class="ml-12 my-12 p-8 w-96 h-96 rounded-lg bg-gray-950 border-black border-4 card left">
          <p class="text-2xl text-transparent bg-clip-text bg-gradient-to-br from-red-900 to-indigo-600 text-opacity-70">
            Challenge Opponents and let a sophisticated Algorithm refine your Rating based on their anticipated Strength
          </p>
          <br><br>
          <Algorithm />
        </div>
      </div>

      <div class="my-16 flex flex-col lg:flex-row items-center justify-around">
        <h1 class="text-3xl sm:text-5xl font-bold text-shadow text-slate-200 animate-pulse">Dominate the Leaderboard</h1>
        <div class="ml-12 my-12 p-8 w-96 h-96 rounded-lg bg-gray-950 border-black border-4 card right">
          <p class="text-2xl text-transparent bg-clip-text bg-gradient-to-br from-red-900 to-indigo-600 text-opacity-70">
            Conquer Matches, earn your Points, and storm the Leaderboard
          </p>
          <br><br>
          <Sparta />
        </div>
      </div>
    </div>

    <h1 bind:this={leaderboardTitle} class="top-96 text-5xl sm:text-7xl my-12">Leaderboard</h1>

    <Leaderboard></Leaderboard>
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

  .card {
    transition: all ease-in 1s;
    cursor: pointer;
  }

  .card.right {
    transform: perspective(500px) rotateX(5deg) rotateY(-10deg);
    box-shadow: 10px 10px 20px rgba(255, 255, 255, 0.1);
  }

  .card.right:hover {
    transform: translateX(-5%) perspective(none) rotateX(0deg) rotateY(0deg);
  }

  .card.left {
    transform: perspective(500px) rotateX(5deg) rotateY(10deg);
    box-shadow: -10px 10px 20px rgba(255, 255, 255, 0.1);
  }

  .card.left:hover {
    transform: translateX(5%) perspective(none) rotateX(0deg) rotateY(0deg);
  }

  .card:hover {
    transition: all ease-out 1s;
  }
</style>