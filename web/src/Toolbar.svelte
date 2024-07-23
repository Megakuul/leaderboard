<script>
  import { UpdateUser } from "$lib/api/actions";
  import { RequestTokens } from "$lib/api/auth";
  import { buttonVariants } from "$lib/components/ui/button";
  import LoaderCircle from "lucide-svelte/icons/loader-circle";
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Dialog from "$lib/components/ui/dialog";
  import { Input } from "$lib/components/ui/input";
  import { toast } from "svelte-sonner";

  /** @type {HTMLHeadingElement} */
  export let leaderboardTitle;

  /** @type {string} */
  let syncTitleInput;
  /** @type {string} */
  let syncIconInput;
  /** @type {boolean} */
  let syncButtonState = false;
</script>


<div class="bg-gray-950 bg-opacity-70 flex flex-col sm:flex-row items-center justify-between p-3 rounded-lg sm:w-9/12">
  <Button on:click={RequestTokens} variant="outline" class="w-60">Sign In</Button>
  <Button on:click={() => leaderboardTitle.scrollIntoView({ behavior: "smooth" })} variant="outline" class="w-60">
    Go to Leaderboard
  </Button>
  <Dialog.Root>
    <Dialog.Trigger class="w-60 {buttonVariants({ variant: "outline" })}">Sync Account</Dialog.Trigger>
    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Synchronize User?</Dialog.Title>
        <Dialog.Description>
          Update your account information. 
          <br>If never synchronized before, this will initialize your account.
        </Dialog.Description>
      </Dialog.Header>
      <Input bind:value={syncTitleInput} type="text" placeholder="Civ Jesus Nutshell" class="max-w-xs" />
      <Input bind:value={syncIconInput} type="url" placeholder="https://gravatar.com/avatar/xyz?size=256" class="max-w-xs" />
      <Dialog.Footer>
        <Button type="submit" on:click={async () => {
          try {
            syncButtonState = true;
            await UpdateUser({ user_update: {
              title: syncTitleInput,
              iconurl: syncIconInput,
            }})
            toast.success("Synchronized user")
          } catch (err) {
            toast.error("Failed to synchronize user", {
              description: err.message,
            })
          }
          syncButtonState = false;
        }}>
          Synchronize
          {#if syncButtonState}
            <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
          {/if}
        </Button>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
</div>