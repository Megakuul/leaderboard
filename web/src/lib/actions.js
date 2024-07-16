import { RequestTokens } from "./auth";

/**
 * @typedef {Object} User
 * @property {string} username
 * @property {string} title
 * @property {string} iconurl
 * @property {number} elo
 */

/**
 * @typedef {Object} FetchResponse
 * @property {string} message
 * @property {string} newpagekey
 * @property {User[]} users
 */

/**
 * Fetches data from the api based on the query params provided.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {string} username evaluated first 
 * @param {string} elo only evaluated if no username is provided
 * @param {string} lastpagekey only evaluated if no elo is provided
 * @returns {Promise<FetchResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const Fetch = async (username="", elo="", lastpagekey="") => {
  const params = new URLSearchParams({
    username: username,
    elo: elo,
    lastpagekey: lastpagekey,
  })
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/fetch?${params.toString()}`, {
    method: "GET"
  })
  if (res.ok) {
    return await res.json();
  } else {
    throw new Error(await res.text())
  }
}


/**
 * @typedef {Object} UpdateResponse
 * @property {string} message
 */

/**
 * Updates or registers the user based on the cognito user profile.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {string} idToken
 * @returns {Promise<UpdateResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const Update = async (idToken) => {
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/update`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${idToken}`
    },
  })
  if (res.ok) {
    return await res.json();
  } else if (res.status === 401) {
    RequestTokens()
  } else {
    throw new Error(await res.text());
  }
}

/**
 * @typedef {Object} UserResult
 * @property {string} username
 * @property {number} placement
 */

/**
 * @typedef {Object} AddGameRequest
 * @property {UserResult[]} results
 */

/**
 * @typedef {Object} AddGameResponse
 * @property {string} message
 */

/**
 * Adds a game based on the provided request.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {string} accessToken 
 * @param {AddGameRequest} request 
 * @returns {Promise<AddGameResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const AddGame = async (accessToken, request) => {
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/addgame`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${accessToken}`
    },
    body: JSON.stringify(request),
  })
  if (res.ok) {
    return await res.json();
  } else if (res.status === 401) {
    RequestTokens()
  } else {
    throw new Error(await res.text());
  }
}