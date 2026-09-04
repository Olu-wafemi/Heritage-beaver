"use client";

import { FormEvent, useState } from "react";
import Link from "next/link";
import { register, resendVerification } from "@/lib/auth";

export default function RegisterPage() {
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [registeredEmail, setRegisteredEmail] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await register(email, password, displayName);
      setRegisteredEmail(res.user.email);
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
        <div>
          <p className="font-mono text-xs tracking-[0.25em] text-gold-200">
            HOW THE WEAVING WORKS
          </p>
          <ol className="mt-6 space-y-5">
            {[
              ["I — Gather", "Name the people of your family."],
              ["II — Remember", "Tell one true story."],
              ["III — Listen", "Keep the wisdom inside it."],
              ["IV — Become", "Weave your living myth."],
            ].map(([step, body]) => (
              <li key={step} className="flex gap-4">
                <span className="mt-2 h-2 w-2 shrink-0 rotate-45 bg-ember-500" aria-hidden="true" />
                <div>
                  <p className="font-mono text-xs tracking-widest text-parchment-200">{step}</p>
                  <p className="mt-0.5 font-display text-lg text-parchment-50">{body}</p>
                </div>
              </li>
            ))}
          </ol>
        </div>
        <p className="max-w-sm text-sm leading-6 text-parchment-200/70">
          Free to begin. Your stories stay yours — sharing is always your choice.
        </p>
      </div>

      <div className="flex items-center justify-center px-6 py-12">
        <div className="w-full max-w-md">
          <h1 className="font-display text-4xl font-medium tracking-tight text-ink-900">
            Begin the telling<span className="text-ember-600">.</span>
          </h1>
          <p className="mt-2 text-sm text-ink-700">
            One account holds your whole family&apos;s mythology.
          </p>

          {registeredEmail ? (
            <div className="mt-8 rounded-2xl border border-parchment-300 bg-parchment-100 p-8 text-center">
              <div className="kente-strip mx-auto w-24 rounded-full" aria-hidden="true" />
              <h2 className="mt-5 font-display text-2xl font-medium">Check your inbox<span className="text-ember-600">.</span></h2>
              <p className="mx-auto mt-2 max-w-xs text-sm leading-6 text-ink-700">
                A confirmation link is on its way to <strong>{registeredEmail}</strong>.
                Open it to take your place at the fireside.
              </p>
              <button
                onClick={() => registeredEmail && resendVerification(registeredEmail).catch(() => null)}
                className="mt-5 text-sm font-semibold text-ember-700 hover:underline"
              >
                Didn&apos;t arrive? Send it again
              </button>
              <p className="mt-4 text-sm text-ink-700">
                <Link href="/login" className="font-semibold text-ember-700 hover:underline">
                  Continue to sign in
                </Link>
              </p>
            </div>
          ) : (
          <form onSubmit={onSubmit} className="mt-8 space-y-5">
            {error && (
              <div
                role="alert"
                className="rounded-lg border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
              >
                {error}
              </div>
            )}

            <div>
              <label htmlFor="displayName" className="block text-sm font-medium text-ink-900">
                Your name
              </label>
              <input
                id="displayName"
                type="text"
                required
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                placeholder="What should the ancestors call you?"
                className={inputClass}
              />
            </div>

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
                minLength={8}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="At least 8 characters"
                className={inputClass}
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-full bg-ember-600 py-3 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? "Lighting the fire…" : "Light the fire"}
            </button>
          </form>
          )}

          <p className="mt-6 text-center text-sm text-ink-700">
            Already have an account?{" "}
            <Link href="/login" className="font-semibold text-ember-700 hover:underline">
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </main>
  );
}
