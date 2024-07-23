<script>
  /** @type {string} */
   export let height;
  /** @type {string} */
  export let width;
</script>

<svg width="{width}" height="{height}" viewBox="10 10 512 581" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path class="path" filter="url(#innerShadow)" d="M6.5 6C40.0195 44.8375 69.4878 83.0518 94.8377 118.5M199.5 285.5L288.25 299M199.5 285.5L150.75 350M199.5 285.5L273.5 363.5L242.5 475L354.5 505.5L332.625 575.5M199.5 285.5C188.85 265.753 170.127 232.222 143.504 190.5M102 414.5L154 468.5M102 414.5L30.5 455.5M102 414.5L150.75 350M456.5 251L377 312.5L332.625 305.75M456.5 251L506.5 267M456.5 251L471.5 155M288.25 299L319 241L291.5 186.5M288.25 299L332.625 305.75M332.625 305.75L354.5 359.5L397 375.5M150.75 350L72.5 355.5L37.5 333.5M94.8377 118.5L58.5 225L98.5 244M94.8377 118.5C113.238 144.231 129.469 168.504 143.504 190.5M143.504 190.5L172.5 125L224.5 109.5"/>
  <defs>
    <filter id="innerShadow">
      <feFlood flood-color="black"/>
      <feComposite operator="out" in2="SourceGraphic"/>
      <feGaussianBlur stdDeviation="4"/>
      <feComposite operator="atop" in2="SourceGraphic"/>
    </filter>
  </defs>
</svg>

<style>
  .path {
    @apply stroke-indigo-500;
    filter: url("#innerShadow") brightness(80%);
    animation: 
      draw 8s ease-in forwards,
      pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 1;
    }
    50% {
      opacity: 0.7;
    }
    100% {
      opacity: 1;
    }
  }

  @keyframes draw {
    0% {
      stroke-width: 2;
      stroke-dasharray: 1000;
      stroke-dashoffset: 1000;
    }
    10% {
      stroke-dashoffset: 1000;
    }
    100% {
      stroke-width: 10;
      stroke-dashoffset: 0;
    }
  }
</style>