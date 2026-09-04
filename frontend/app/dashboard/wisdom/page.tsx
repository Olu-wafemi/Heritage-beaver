"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import type { WisdomExtract } from "@/lib/types";

export default function WisdomPage() {
  const [extracts, setExtracts] = useState<WisdomExtract[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const list = await apiFetch<WisdomExtract[]>("/wisdom-extracts", { token });
      setExtracts(list);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load wisdom");
      setExtracts([]);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  return (
    <div>
      <p className="font-mono text-xs tracking-[0.25em] text-ember-600">STAGE III — LISTEN</p>
      <h1 className="mt-2 font-display text-4xl font-medium tracking-tight">The Wisdom</h1>
      <p className="mt-2 max-w-lg leading-7 text-ink-700">
        Proverbs, warnings, and values Hearthside has heard inside your stories — kept
        here like kola nuts in a bowl.
      </p>

      {error && (
        <div
          role="alert"
          className="mt-6 rounded-xl border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
        >
          {error}
        </div>
      )}

      <div className="mt-8">
        {extracts === null ? (
          <div className="grid gap-4 sm:grid-cols-2">
            {[0, 1, 2, 3].map((i) => (
              <div key={i} className="h-36 animate-pulse rounded-2xl bg-parchment-200/70" />
            ))}
          </div>
        ) : extracts.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-parchment-300 bg-parchment-100 p-12 text-center">
            <p className="font-display text-2xl font-medium">The bowl is empty</p>
            <p className="mx-auto mt-2 max-w-sm text-sm leading-6 text-ink-700">
              Open one of your stories and press “Draw out the wisdom” to fill it.
            </p>
            <Link
              href="/dashboard/stories"
              className="mt-5 inline-block rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
            >
              Go to the stories
            </Link>
          </div>
        ) : (
          <ul className="grid gap-4 sm:grid-cols-2">
            {extracts.map((w) => (
              <li key={w.id} className="flex flex-col rounded-2xl bg-ink-950 p-6 text-parchment-100">
                <span className="w-fit rounded-full bg-ember-600 px-2.5 py-0.5 font-mono text-[11px] tracking-widest text-parchment-50">
                  {w.wisdom_type.toUpperCase()}
                </span>
                <p className="mt-4 font-display text-xl font-medium italic leading-snug text-parchment-50">
                  “{w.excerpt}”
                </p>
                <p className="mt-3 flex-1 text-sm leading-6 text-parchment-200/85">{w.meaning}</p>
                <Link
                  href={`/dashboard/stories/${w.story_id}`}
                  className="mt-4 text-xs font-semibold tracking-wide text-gold-200 hover:underline"
                >
                  HEAR IT IN ITS STORY →
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
