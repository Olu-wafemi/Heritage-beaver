import { config } from "./config";

type FetchOptions = RequestInit & { token?: string; _retried?: boolean };

let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      const refreshToken = localStorage.getItem("refresh_token");
      if (!refreshToken) throw new Error("No session to refresh");

      const res = await fetch(`${config.apiUrl}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
      if (!res.ok) throw new Error("Session expired");

      const data = (await res.json()) as { token: string; refresh_token: string; user: unknown };
      localStorage.setItem("token", data.token);
      localStorage.setItem("refresh_token", data.refresh_token);
      localStorage.setItem("user", JSON.stringify(data.user));
      return data.token;
    })();

    try {
      return await refreshPromise;
    } finally {
      refreshPromise = null;
    }
  }
  return refreshPromise;
}

export async function apiFetch<T>(path: string, opts: FetchOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(opts.headers as Record<string, string>),
  };

  const token = opts.token ?? localStorage.getItem("token") ?? "";
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${config.apiUrl}${path}`, {
    ...opts,
    headers,
  });

  // One transparent retry: short-lived access token may have lapsed.
  if (res.status === 401 && !opts._retried && !path.startsWith("/auth/")) {
    try {
      const fresh = await refreshAccessToken();
      return apiFetch<T>(path, { ...opts, token: fresh, _retried: true });
    } catch {
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("user");
    }
  }

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new Error(body?.error ?? `Request failed: ${res.status}`);
  }

  if (res.status === 204) return null as T;
  return res.json() as Promise<T>;
}

export async function checkHealth(): Promise<{ status: string }> {
  return apiFetch("/health");
}
