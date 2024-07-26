const COGNITO_DOMAIN = import.meta.env.VITE_COGNITO_DOMAIN
const COGNITO_CLIENT_ID = import.meta.env.VITE_COGNITO_CLIENT_ID
const COGNITO_SCOPES = "openid profile email"
const COGNITO_RESPONSE_TYPE = "token";

/**
 * Converts base64url to string
 * @param {string} str 
 * @returns {string}
 */
const decodeBase64URL = (str) => {
  const m = str.length % 4;
  return new TextDecoder().decode(Uint8Array.from(atob(
      str.replace(/-/g, '+')
          .replace(/_/g, '/')
          .padEnd(str.length + (m === 0 ? 0 : 4 - m), '=')
  ), c => c.charCodeAt(0)).buffer)
}

/**
 * Acquire expiration time as unix millisecond. Returns NaN on failure.
 * @param {string} idToken 
 * @returns {number}
 */
const getIDTokenExpiration = (idToken) => {
  try {
    const idTokenPayload = JSON.parse(
      decodeBase64URL(idToken.split(".")[1])
    )
    return idTokenPayload.exp * 1000;
  } catch {
    return NaN;
  }
}

/**
 * Acquire preferred_username. Returns "" on failure.
 * @param {string} idToken 
 * @returns {string}
 */
const getPreferredUsername = (idToken) => {
  try {
    const idTokenPayload = JSON.parse(
      decodeBase64URL(idToken.split(".")[1])
    )
    return idTokenPayload.preferred_username;
  } catch {
    return "";
  }
}

/**
 * Fetches the id_token from the url and inserts them to localstore.
 */
export const GetTokens = () => {
  // Params are provided as url fragment so that they are not sent to the server.
  const url = new URL(window.location.href);
  const params = new URLSearchParams(url.hash.slice(1));

  const idToken = params.get("id_token");
  const idTokenExp = getIDTokenExpiration(idToken);
  const username = getPreferredUsername(idToken);

  // Removing search param from history.
  window.history.replaceState({}, "", `${window.location.origin}${window.location.pathname}`)
  
  if (idToken) {
    localStorage.setItem("id_token", idToken);
  }
  if (idTokenExp) {
    localStorage.setItem("id_token_exp", new Date(idTokenExp).toISOString());
  }
  if (username) {
    localStorage.setItem("id_token_username", username);
  }
}

/**
 * Redirects the user to the cognito endpoint to fetch the id_token.
 * The redirect_uri is set to the current window location.
 */
export const RequestTokens = () => {
  const devUrl = import.meta.env.VITE_DEV_API_URL;
  const params = new URLSearchParams({
    client_id: COGNITO_CLIENT_ID,
    // prompt=none for silent auth
    response_type: COGNITO_RESPONSE_TYPE,
    scope: COGNITO_SCOPES,
    redirect_uri: devUrl ? devUrl : window.location.href
  })
  window.location.href = `${COGNITO_DOMAIN}/login?${params.toString()}`
}