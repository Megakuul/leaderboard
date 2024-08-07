import { RequestTokens } from "./auth";

/**
 * @typedef {Object} FetchUserResponseUser
 * @property {string} username
 * @property {boolean} disabled
 * @property {string} region
 * @property {string} title
 * @property {string} iconurl
 * @property {number} elo
 */

/**
 * @typedef {Object} FetchUserResponse
 * @property {string} message
 * @property {string} newpagekey
 * @property {FetchUserResponseUser[]} users
 */

/**
 * Fetches data from the api based on the query params provided.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {string} region
 * @param {string} pagesize
 * @param {string} username
 * @param {string} elo
 * @param {string} lastpagekey
 * @returns {Promise<FetchUserResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const FetchUser = async (region, pagesize="", username="", elo="", lastpagekey="") => {
  const params = new URLSearchParams({
    username: username,
    region: region,
    pagesize: pagesize,
    elo: elo,
    lastpagekey: lastpagekey,
  })
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/user/fetch?${params.toString()}`, {
    method: "GET"
  })
  if (res.ok) {
    return await res.json();
  } else {
    throw new Error(await res.text())
  }
}

/**
 * @typedef {Object} UpdateUserRequestUser
 * @property {string} title
 * @property {string} iconurl
 * @property {boolean} disabled
 */

/**
 * @typedef {Object} UpdateUserRequest
 * @property {UpdateUserRequestUser} user_updates
 */

/**
 * @typedef {Object} UpdateUserResponseUser
 * @property {string} username
 * @property {boolean} disabled
 * @property {string} region
 * @property {string} title
 * @property {string} email
 * @property {string} iconurl
 * @property {number} elo
 */

/**
 * @typedef {Object} UpdateUserResponse
 * @property {string} message
 * @property {UpdateUserResponseUser} updated_user
 */

/**
 * Updates or registers the user based on the cognito user profile.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {UpdateUserRequest} userData
 * @returns {Promise<UpdateUserResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const UpdateUser = async (userData) => {
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/user/update`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${localStorage.getItem("id_token")}`
    },
    body: JSON.stringify(userData),
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
 * @typedef {Object} FetchGameResponseParticipant
 * @property {string} username
 * @property {boolean} underdog
 * @property {number} team
 * @property {number} placement
 * @property {number} points
 * @property {number} elo
 * @property {number} elo_update
 * @property {boolean} confirmed
 */

/**
 * @typedef {Object} FetchGameResponseGame
 * @property {string} gameid
 * @property {string} date
 * @property {boolean} readonly
 * @property {number} expires_in
 * @property {Object.<string, FetchGameResponseParticipant>} participants
 */

/**
 * @typedef {Object} FetchGameResponse
 * @property {string} message
 * @property {FetchGameResponseGame[]} games
 */

/**
 * Fetches game data from the api based on the query params provided.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {string} gameid
 * @param {string} date
 * @returns {Promise<FetchGameResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const FetchGame = async (gameid="", date="") => {
  const params = new URLSearchParams({
    gameid: gameid,
    date: date,
  })
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/game/fetch?${params.toString()}`, {
    method: "GET"
  })
  if (res.ok) {
    return await res.json();
  } else {
    throw new Error(await res.text())
  }
}

/**
 * @typedef {Object} AddGameRequestParticipant
 * @property {string} username
 * @property {number} team
 * @property {number} placement
 * @property {number} points
 */

/**
 * @typedef {Object} AddGameRequest
 * @property {number} placement_points
 * @property {AddGameRequestParticipant[]} participants
 */

/**
 * @typedef {Object} AddGameResponse
 * @property {string} message
 * @property {string} gameid
 */

/**
 * Adds a new game to the api.
 * https://github.com/Megakuul/leaderboard/blob/main/README.md#api
 * @param {AddGameRequest} gameData
 * @returns {Promise<AddGameResponse>} if api call succeeds.
 * @throws {Error} if api call failed.
 */
export const AddGame = async (gameData) => {
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const res = await fetch(`${devUrl?devUrl:""}/api/game/add`, {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${localStorage.getItem("id_token")}`
    },
    body: JSON.stringify(gameData)
  })
  if (res.ok) {
    return await res.json();
  } else if (res.status === 401) {
    RequestTokens()
  } else {
    throw new Error(await res.text())
  }
}