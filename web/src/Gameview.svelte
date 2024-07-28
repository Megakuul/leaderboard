<script>
  import { Input } from "$lib/components/ui/input";
  import { ScrollArea } from "$lib/components/ui/scroll-area";
  import { Badge } from "$lib/components/ui/badge";
  import { GetColor } from "$lib/GetColor";

  /** @type {import("$lib/api/actions").FetchGameResponseGame} */
  export let Game;
</script>

<ScrollArea class="max-h-96 w-full p-2">
  {#each Object.entries(Game.participants) as [_, participant]}
    <div 
      class="relative p-[1px]">
      <div 
        style="background-color: rgba(0,0,0,0.2); background-color: hsl({GetColor(participant.team)}, 40%, 20%);"
        class="flex flex-col gap-4 m-1 p-4 my-4 bg-opacity-60 rounded-lg">
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
          <Input disabled value={participant.placement} type="text" class="w-full font-bold" />
          <Input disabled value={participant.points} type="text" class="w-full font-bold" />
        </div>
      </div>
    </div>
  {/each}
</ScrollArea>
