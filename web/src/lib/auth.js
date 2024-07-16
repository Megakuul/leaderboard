const COGNITO_DOMAIN = import.meta.env.VITE_COGNITO_DOMAIN
const COGNITO_CLIENT_ID = import.meta.env.VITE_COGNITO_CLIENT_ID
const COGNITO_SCOPES = "openid profile"
const COGNITO_RESPONSE_TYPE = "id_token token";

/**
 * @typedef {Object} AuthTokens
 * @property {string} id_token
 * @property {string} access_token
 */

/**
 * Fetches the access_token & id_token from the url and returns it.
 * @returns {AuthTokens}
 */
export const GetTokens = () => {
  const params = new URLSearchParams(window.location.search);
  const idToken = params.get("id_token");
  const accessToken = params.get("access_token");
  // Removing search param from history.
  window.history.replaceState({}, "", `${window.location.origin}${window.location.pathname}`)
  return {
    id_token: idToken,
    access_token: accessToken
  };
}

/**
 * Redirects the user to the cognito endpoint to fetch the access_token & id_token.
 * The redirect_uri is set to the current window location.
 */
export const RequestTokens = () => {
  const params = new URLSearchParams({
    client_id: COGNITO_CLIENT_ID,
    // prompt=none for silent auth
    response_type: COGNITO_RESPONSE_TYPE,
    scope: COGNITO_SCOPES,
    redirect_uri: window.location.href
  })
  window.location.replace(
    `https://${COGNITO_DOMAIN}/login?${params.toString()}`
  )
}