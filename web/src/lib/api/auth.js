const COGNITO_DOMAIN = import.meta.env.VITE_COGNITO_DOMAIN
const COGNITO_CLIENT_ID = import.meta.env.VITE_COGNITO_CLIENT_ID
const COGNITO_SCOPES = "openid profile email"
const COGNITO_RESPONSE_TYPE = "token";


/**
 * Fetches the access_token & id_token from the url and inserts them to localstore.
 */
export const GetTokens = () => {
  // Params are provided as url fragment so that they are not sent to the server.
  const url = new URL(window.location.href);
  const params = new URLSearchParams(url.hash.slice(1));

  const idToken = params.get("id_token");
  const accessToken = params.get("access_token");
  // Removing search param from history.
  window.history.replaceState({}, "", `${window.location.origin}${window.location.pathname}`)
  
  if (idToken) {
    localStorage.setItem("id_token", idToken);
  }
  if (accessToken) {
    localStorage.setItem("access_token", accessToken);
  }
}

/**
 * Redirects the user to the cognito endpoint to fetch the access_token & id_token.
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