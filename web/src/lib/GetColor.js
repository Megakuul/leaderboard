  /**
   * Generate a hsl color (0-360) based on a index.
   * @param {number} index
  */
  export const GetColor = (index) => {
    if (Number.isNaN(index)) return undefined;

    // Starting with index 1 instead of 0
    index++;

    const COLORSPACE = 360;

    // The color space is split up into layers which always contain 
    // twice as much colors as the previous layer. To obtain the layer in constant time
    // log2 is used which returns the power on 2 required to get the input number it takes.
    // This number is usually going to have some fraction in it, however we need the layer not the exact power,
    // so we floor the fraction and get the power of 2 that leads us to the first index of the layer.
    const LAYER = Math.floor(Math.log2(index));
    // For later calculations we need two values:
    // 1. The amount of numbers in this layer.
    // 2. The starting index of this layer.
    // Because we used log2 and the layer represents a power of two, both values can be obtained by just using the layer as power of two.
    // This returns the starting index of this layer, and because every layer doubles in size, this index is also the length of this layer. 
    const BASE = Math.pow(2, LAYER);

    // As explained above we have the first index of the layer, now we also need the offset so that we can get the
    // exact index inside the layer. For this we just take the starting index and subtract it from the index.
    const LAYER_OFFSET = index - BASE;

    // Now that we have the base and the layer offset of the index, we can use this information and apply it to the hue 360 color space.
    // To do this, we calculate the multiplier that is used for the layer offset. We do now need to acquire numbers between the previous layer,
    // therefore we get a multiplier which is just half the size of our actual increment steps per index.
    const MULTIPLIER = COLORSPACE / (BASE * 2);
    
    // Now that we have a multiplier with the half size of the actual distance between index steps, we need to multiply the layer offset by two.
    // This would give us the values that were present on the previous layer, we do now add + 1 to the layer offset to shift the values.
    // This means that we do now essentially get the number between multiplier i and multiplier j from the previous layer.
    return MULTIPLIER * (LAYER_OFFSET * 2 + 1);
  }