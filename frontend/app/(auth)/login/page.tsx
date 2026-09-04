"use client";

import { FormEvent, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { login, resendVerification } from "@/lib/auth";
import { useAuthStore } from "@/store/useAuthStore";

export default function LoginPage() {
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [resent, setResent] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await login(email, password);
      setAuth(res.token, res.user, res.refresh_token);
      router.push("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setLoading(false);
    }
  }

  const inputClass =
    "mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition placeholder:text-ink-400 focus:border-ember-600 focus:ring-2 focus:ring-ember-100";

  return (
    <main className="grid min-h-screen bg-parchment-50 lg:grid-cols-2">
      <div className="relative hidden flex-col justify-between overflow-hidden bg-ink-950 p-12 text-parchment-100 lg:flex">
        <div className="kente-strip absolute inset-x-0 top-0" aria-hidden="true" />
        <Link href="/" className="font-display text-xl font-semibold text-parchment-50">
          Hearthside
        </Link>
        <figure>
          <blockquote className="font-display text-3xl font-medium italic leading-snug text-parchment-50">
            “When an elder dies, a library burns to the ground.”
          </blockquote>
          <figcaption className="mt-4 font-mono text-xs tracking-[0.2em] text-gold-200">
            AFRICAN PROVERB — THE REASON THIS EXISTS
          </figcaption>
        </figure>
        <p className="max-w-sm text-sm leading-6 text-parchment-200/70">
          Every login returns you to the telling: your people, their stories, the wisdom
          already drawn from them.
        </p>
      </div>

      <div className="flex items-center justify-center px-6 py-12">
        <div className="w-full max-w-md">
          <h1 className="font-display text-4xl font-medium tracking-tight text-ink-900">
            Welcome back<span className="text-ember-600">.</span>
          </h1>
          <p className="mt-2 text-sm text-ink-700">
            Sign in to continue weaving your family story.
          </p>

          <form onSubmit={onSubmit} className="mt-8 space-y-5">
            {error && (
              <div
                role="alert"
                className="rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
              >
                {error}
                {/not verified/i.test(error) && (
                  <span>
                    {" "}
                    <button
                      type="button"
                      onClick={() =>
                        resendVerification(email)
                          .then(() => setResent(true))
                          .catch(() => setResent(true))
                      }
                      className="font-semibold underline"
                    >
                      Send the link again
                    </button>
                    {resent && " — check your inbox."}
                  </span>
                )}
              </div>
            )}

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-ink-900">
                Email
              </label>
              <input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                className={inputClass}
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-ink-900">
                Password
              </label>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className={inputClass}
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-full bg-ember-600 py-3 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? "Returning…" : "Return to the telling"}
            </button>
          </form>

          <p className="mt-6 text-center text-sm text-ink-700">
            New to Hearthside?{" "}
            <Link href="/register" className="font-semibold text-ember-700 hover:underline">
              Begin your family myth
            </Link>
          </p>
        </div>
      </div>
    </main>
  );
}
