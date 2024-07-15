import { get } from 'svelte/store';
import { page } from '$app/stores';

const COGNITO_DOMAIN = import.meta.env.VITE_COGNITO_DOMAIN
const COGNITO_CLIENT_ID = import.meta.env.VITE_COGNITO_CLIENT_ID
const COGNITO_SCOPES = "openid profile"
const COGNITO_RESPONSE_TYPE = "token";

/**
 * Fetches the access_token from the url and returns it.
 * @returns {string}
 */
export const GetAccessToken = () => {
  const currentPage = get(page);
  const token = currentPage.url.searchParams.get("access_token");
  // Removing access_token search param from history.
  window.history.replaceState({}, "", currentPage.url.searchParams.delete("access_token"))
  return token;
}

/**
 * Redirects the user to the cognito endpoint to fetch the access_token.
 * The redirect_uri is set to the current window location.
 */
export const RequestAccessToken = () => {
  const params = new URLSearchParams({
    client_id: COGNITO_CLIENT_ID,
    response_type: COGNITO_RESPONSE_TYPE,
    scope: COGNITO_SCOPES,
    redirect_uri: window.location
  })
  window.location.replace(
    `https://${COGNITO_DOMAIN}/login?${params.toString()}`
  )
}