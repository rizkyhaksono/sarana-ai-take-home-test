export function setAuthCookie(token: string, expiresInDays: number = 7) {
  if (typeof window === 'undefined') return;
  const maxAge = expiresInDays * 24 * 60 * 60; // in seconds
  document.cookie = `sarana-note-app-auth=${token}; path=/; max-age=${maxAge}; SameSite=Lax`;
}

export function getAuthCookie(): { token: string } | null {
  if (typeof window === 'undefined') return null;
  const value = `; ${document.cookie}`;
  const parts = value.split(`; sarana-note-app-auth=`);
  if (parts.length === 2) {
    const token = parts.pop()?.split(';').shift();
    if (token) return { token };
  }
  return null;
}

