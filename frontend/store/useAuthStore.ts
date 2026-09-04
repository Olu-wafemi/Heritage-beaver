import { create } from "zustand";
import type { User } from "@/lib/types";

type AuthState = {
  token: string | null;
  refreshToken: string | null;
  user: User | null;
  setAuth: (token: string, user: User, refreshToken: string) => void;
  clearAuth: () => void;
};

function read(key: string): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(key);
}

export const useAuthStore = create<AuthState>((set) => ({
  token: read("token"),
  refreshToken: read("refresh_token"),
  user:
    typeof window !== "undefined"
      ? (() => {
          try {
            const raw = localStorage.getItem("user");
            return raw ? (JSON.parse(raw) as User) : null;
          } catch {
            return null;
          }
        })()
      : null,
  setAuth: (token, user, refreshToken) => {
    localStorage.setItem("token", token);
    localStorage.setItem("refresh_token", refreshToken);
    localStorage.setItem("user", JSON.stringify(user));
    set({ token, user, refreshToken });
  },
  clearAuth: () => {
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user");
    set({ token: null, refreshToken: null, user: null });
  },
}));
