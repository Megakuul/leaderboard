<script>
  import { AddGame, FetchGame, FetchUser } from "$lib/api/actions";
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
  import { ScrollArea } from "$lib/components/ui/scroll-area";
  import { fade } from "svelte/transition";
  import { Badge } from "$lib/components/ui/badge";

  /** @type {string} */
  const REGIONS = import.meta.env.VITE_LEADERBOARD_REGIONS;

  /** @type {boolean} */
  let addButtonState;

  /** @type {number} */
  let addPlacementPoints = 100;

  /** @type {import("$lib/api/actions").AddGameRequestParticipant[]} */
  let addParticipants = [];

  /** @type {import("$lib/api/actions").FetchGameResponseGame} */
  let addGameResult = undefined;

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
  let queryResults = [
    {
      disabled: true,
      elo: 15,
      iconurl: "",
      region: "eu-central-1",
      title: "superkuul",
      username: "Salami"
    },
    {
      disabled: false,
      elo: 15,
      iconurl: "",
      region: "eu-central-1",
      title: "superkuul",
      username: "Salami"
    }
  ];

  /** @type {boolean} */
  let loadPageButtonState;

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
      {#if !addGameResult}
        <Dialog.Header>
          <Dialog.Title>Add new Game</Dialog.Title>
          <Dialog.Description>
            Add a game to the leaderboard.
          </Dialog.Description>
        </Dialog.Header>

        <div class="flex flex-row gap-2">
          <Input bind:value={addPlacementPoints} on:input={(_) => addPlacementPoints = +addPlacementPoints} type="number" placeholder="Placement Points" />

          <Button variant="outline" class="w-full" on:click={() => {
            addParticipants = addParticipants.concat({
              username: "",
              placement: NaN,
              points: NaN,
              team: NaN,
            })
          }}>Add Participant</Button>
        </div>

        <ScrollArea class="max-h-96 w-full p-2">
          {#each addParticipants as participant}
            <div transition:fade={{ delay: 250, duration: 300 }} class="flex flex-col gap-4 m-1 p-4 my-4 bg-black bg-opacity-60 rounded-lg">
              <Input bind:value={participant.username} type="text" placeholder="Username" class="w-full" />
              <div class="flex flex-row gap-2">
                <Select.Root portal={null} bind:selected={participant.team}>
                  <Select.Trigger class="w-full">
                    <Select.Value placeholder="Select a query type" />
                  </Select.Trigger>
                  <Select.Content>
                    <Select.Group>
                      <Select.Label>Query Type</Select.Label>
                      {#each addParticipants as _, i}
                        <Select.Item value={i} label={i}>{"Team " + i}</Select.Item>
                      {/each}
                    </Select.Group>
                  </Select.Content>
                </Select.Root>
                <Input bind:value={participant.team} on:input={(_) => participant.team = +participant.team} type="number" placeholder="Team" class="w-full" />
                <Input bind:value={participant.placement} on:input={(_) => participant.placement = +participant.placement} type="number" placeholder="Placement" class="w-full" />
                <Input bind:value={participant.points} on:input={(_) => participant.points = +participant.points} type="number" placeholder="Points" class="w-full" />
              </div>
            </div>
          {/each}
        </ScrollArea>

        <Dialog.Footer>
          <Button type="submit" on:click={async () => {
            try {
              addButtonState = true;
              const addResponse = await AddGame({
                placement_points: addPlacementPoints,
                participants: addParticipants,
              })
              toast.success("Game added")
              const fetchResponse = await FetchGame(addResponse.gameid)
              addGameResult = fetchResponse.games[0];
              toast.success("Preview loaded")
            } catch (err) {
              toast.error("Failed to add game", {
                description: err.message,
              })
            }
            addButtonState = false;
          }}>
            {#if addButtonState}
              <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
            {/if}
            Add Game
          </Button>
        </Dialog.Footer>
      {:else}
        <Dialog.Header>
          <Dialog.Title>Review Game {addGameResult.gameid}</Dialog.Title>
          <Dialog.Description>
            Participants received a confirmation email. 
            All must confirm by {new Date(addGameResult.expires_in * 1000).toLocaleString()} or the game will be canceled.
          </Dialog.Description>
        </Dialog.Header>

        <ScrollArea class="max-h-96 w-full p-2">
          {#each Object.entries(addGameResult.participants) as [_, participant]}
            <div class="relative flex flex-col gap-4 m-1 p-4 my-4 bg-black bg-opacity-60 rounded-lg">
              <div class="absolute z-40 top-0 right-0 flex flex-row gap-4">
                {#if !participant.confirmed}
                  <Badge class="bg-orange-500">Not Confirmed</Badge>
                {/if}
                {#if participant.underdog}
                  <Badge class="bg-indigo-700">Underdog</Badge>
                {/if}
                {#if participant.elo_update >= 0}
                  <Badge class="bg-green-500">+{participant.elo_update}</Badge>
                {:else}
                  <Badge class="bg-red-600">{participant.elo_update}</Badge>
                {/if}
              </div>
              <Input disabled value={participant.username} type="text" placeholder="Username" class="w-full" />
              <div class="flex flex-row gap-2">
                <Input disabled value={participant.team} type="number" placeholder="Team" class="w-full" />
                <Input disabled value={participant.placement} type="number" placeholder="Placement" class="w-full" />
                <Input disabled value={participant.points} type="number" placeholder="Points" class="w-full" />
              </div>
            </div>
          {/each}
        </ScrollArea>

        <Dialog.Footer>
          <Button on:click={() => {
            addGameResult = undefined;
            addParticipants = [];
          }}>
            Create new Game
          </Button>
        </Dialog.Footer>
      {/if}
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
      <div class="flex flex-row gap-4 items-center w-full bg-slate-950 bg-opacity-50 rounded-xl my-4 p-5" 
        class:opacity-50={result.disabled} title="{result.disabled?"disabled":""}">
        <Avatar.Root>
          <Avatar.Image src="{result.iconurl}" alt="{result.username} icon" />
          <Avatar.Fallback>AN</Avatar.Fallback>
        </Avatar.Root>
        <p class="text-xl uppercase font-bold">{result.username}</p>
        <p class="hidden sm:block text-slate-50 text-opacity-50 overflow-hidden"># {result.title}</p>
        <p class="hidden sm:block text-slate-50 text-opacity-50 overflow-hidden ml-auto mr-2">{result.region}</p>
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
    <center>
      <LoaderCircle class="h-8 w-8 animate-spin" />
    </center>
  {/if}
  {#if lastPageKey}
    <Button variant="ghost" class="w-full my-6" on:click={async () => {
      try {
        loadPageButtonState = true;
        const response = await FetchUser(
          selectedRegion.value||"",
          queryEntryCount,
          "",
          "",
          lastPageKey,
        )
        lastPageKey = response.newpagekey;
        queryResults = queryResults.concat(response.users);
        queryResultError = "";
        toast.success("Entries fetched")
      } catch (err) {
        queryResultError = err.message;
        toast.error("Failed to fetch entries", {
          description: err.message,
        })
      }
      loadPageButtonState = false;
    }}>
      Load more...
      {#if loadPageButtonState}
        <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
      {/if}
    </Button>
    {/if}
</div>