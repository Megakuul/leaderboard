<script>
  import { FetchUser } from "$lib/api/actions";
  import { buttonVariants } from "$lib/components/ui/button";
  import LoaderCircle from "lucide-svelte/icons/loader-circle";
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import * as Select from "$lib/components/ui/select";
  import CircleAlert from "lucide-svelte/icons/circle-alert";
  import * as Alert from "$lib/components/ui/alert/index.js";
  import * as Avatar from "$lib/components/ui/avatar/index.js";
  import { toast } from "svelte-sonner";
  import { onMount } from "svelte";

  /** @type {string} */
  const REGIONS = import.meta.env.VITE_LEADERBOARD_REGIONS;

  /** @type {import("bits-ui/dist").Selected<any>} */
  let selectedRegion = { value: "" };

  /**
   * @readonly
   * @enum {string}
   */
  const QUERYTYPES = {
    DEFAULT: "default",
    ELO: "elo",
    USERNAME: "username"
  }

  /** @type {import("bits-ui/dist").Selected<any>} */
  let selectedQueryType = { value: "default" };

  /** @type {string} */
  let queryString;

  /** @type {boolean} */
  let queryButtonState;

  /** @type {string} */
  let queryEntryCount = "50";

  /** @type {string} */
  let queryResultError;

  /** @type {import("$lib/api/actions").FetchUserResponseUser[]} */
  let queryResults;

  /** @type {string} */
  let lastPageKey;

  onMount(async () => {
    try {
      const response = await FetchUser("", queryEntryCount, "", "", "")
      lastPageKey = response.newpagekey;
      queryResults = response.users;
      queryResultError = "";
    } catch (err) {
      queryResultError = err.message;
      toast.error("Failed to fetch entries", {
        description: err.message,
      })
    }
  })
</script>

<div class="bg-gray-950 bg-opacity-70 my-12 flex flex-col lg:flex-row items-start lg:items-center gap-2 justify-between p-3 rounded-lg lg:w-9/12">
  <Dialog.Root>
    <Dialog.Trigger class="w-full lg:w-40 {buttonVariants({ variant: "outline" })}">Add Game</Dialog.Trigger>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Add new Game</Dialog.Title>
        <Dialog.Description>
          Add a game to the leaderboard.
        </Dialog.Description>
      </Dialog.Header>
      <Dialog.Footer>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>

  <Input bind:value={queryEntryCount} type="number" placeholder="Entries" class="w-full lg:w-20" />

  <Input bind:value={queryString} type="text" placeholder="Query String" class="w-full lg:max-w-40 xl:max-w-xs" />

  <Select.Root portal={null} bind:selected={selectedQueryType}>
    <Select.Trigger class="w-full lg:w-[180px]">
      <Select.Value placeholder="Select a query type" />
    </Select.Trigger>
    <Select.Content>
      <Select.Group>
        <Select.Label>Query Type</Select.Label>
        {#each Object.entries(QUERYTYPES) as queryType}
          <Select.Item value={queryType[1]} label={queryType[1]}>{queryType[1]}</Select.Item>
        {/each}
      </Select.Group>
    </Select.Content>
  </Select.Root>

  <Select.Root portal={null} bind:selected={selectedRegion}>
    <Select.Trigger class="w-full lg:w-[180px]">
      <Select.Value placeholder="Select a region" />
    </Select.Trigger>
    <Select.Content>
      <Select.Group>
        <Select.Label>Region</Select.Label>
        {#each REGIONS.split(",") as region}
          <Select.Item value={region} label={region}>{region}</Select.Item>
        {/each}
      </Select.Group>
    </Select.Content>
  </Select.Root>

  <Button class="w-full lg:w-40" on:click={async () => {
    try {
      queryButtonState = true;
      const response = await FetchUser(
        selectedRegion.value||"",
        queryEntryCount,
        selectedQueryType.value===QUERYTYPES.USERNAME ? queryString : "",
        selectedQueryType.value===QUERYTYPES.ELO ? queryString : "",
        "",
      )
      lastPageKey = response.newpagekey;
      queryResults = response.users;
      queryResultError = "";
      toast.success("Entries fetched")
    } catch (err) {
      queryResultError = err.message;
      toast.error("Failed to fetch entries", {
        description: err.message,
      })
    }
    queryButtonState = false;
  }}>
    Query 
    {#if queryButtonState}
      <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
    {/if}
  </Button>
</div>

<div class="lg:w-9/12 my-10">
  {#if queryResults}
    {#each queryResults as result}
      <div class="flex flex-row gap-4 items-center w-full bg-slate-950 bg-opacity-50 rounded-xl p-5">
        <Avatar.Root>
          <Avatar.Image src="{result.iconurl}" alt="{result.username} icon" />
          <Avatar.Fallback>AN</Avatar.Fallback>
        </Avatar.Root>
        <p class="text-xl uppercase font-bold">{result.username}</p>
        <p class="text-slate-50 text-opacity-50"># {result.title}</p>
        <p class="text-slate-50 text-opacity-50 ml-auto mr-2">{result.region}</p>
        <p class="text-xl font-bold mr-2">{result.elo}</p>
      </div>
    {/each}
  {:else if queryResultError}
    <Alert.Root variant="destructive" class="w-full bg-slate-950">
      <CircleAlert class="h-4 w-4" />
      <Alert.Title>Error</Alert.Title>
      <Alert.Description>{queryResultError}</Alert.Description>
    </Alert.Root>
  {:else}
    <LoaderCircle class="h-8 w-8 animate-spin" />
  {/if}
</div>