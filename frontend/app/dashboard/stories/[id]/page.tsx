"use client";

import { use, useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import type { Story, WisdomExtract } from "@/lib/types";

export default function StoryDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [story, setStory] = useState<Story | null>(null);
  const [extracts, setExtracts] = useState<WisdomExtract[]>([]);
  const [wisdomMessage, setWisdomMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [extracting, setExtracting] = useState(false);

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const s = await apiFetch<Story>(`/stories/${id}`, { token });
      setStory(s);
      const w = await apiFetch<WisdomExtract[]>(`/wisdom-extracts?story_id=${id}`, { token });
      setExtracts(w);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load story");
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  async function extractWisdom() {
    setExtracting(true);
    setError(null);
    setWisdomMessage(null);
    try {
      const token = localStorage.getItem("token") ?? "";
      const res = await apiFetch<{ extracts: WisdomExtract[]; message?: string }>(
        `/stories/${id}/process-wisdom`,
        { method: "POST", token }
      );
      setExtracts(res.extracts ?? []);
      if (res.message) setWisdomMessage(res.message);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to extract wisdom");
    } finally {
      setExtracting(false);
    }
  }

  if (error && !story) {
    return (
      <div>
        <div
          role="alert"
          className="rounded-xl border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
        >
          {error}
        </div>
        <Link
          href="/dashboard/stories"
          className="mt-4 inline-block text-sm font-medium text-ember-700 hover:underline"
        >
          ← Back to the stories
        </Link>
      </div>
    );
  }

  if (!story) {
    return (
      <div className="mx-auto max-w-3xl space-y-4">
        <div className="h-8 w-2/3 animate-pulse rounded-lg bg-parchment-200/70" />
        <div className="h-64 animate-pulse rounded-2xl bg-parchment-200/70" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl">
      <Link
        href="/dashboard/stories"
        className="text-sm font-medium text-ink-500 transition hover:text-ember-700"
      >
        ← Back to the stories
      </Link>

      <article className="mt-6 overflow-hidden rounded-2xl border border-parchment-300 bg-parchment-100">
        <div className="kente-strip" aria-hidden="true" />
        <div className="p-8 sm:p-10">
          <p className="font-mono text-xs tracking-[0.25em] text-ember-600">A TELLING</p>
          <h1 className="mt-3 font-display text-4xl font-medium leading-tight tracking-tight">
            {story.title}
          </h1>
          <p className="mt-3 font-mono text-xs tracking-widest text-ink-400">
            {new Date(story.created_at).toLocaleDateString()}
            {story.source_language ? ` · TOLD IN ${story.source_language.toUpperCase()}` : ""}
          </p>

          {story.summary && (
            <p className="mt-5 border-l-2 border-gold-500 pl-4 font-display text-lg italic leading-8 text-ink-700">
              {story.summary}
            </p>
          )}

          <div className="manuscript mt-6 whitespace-pre-wrap">{story.content}</div>
        </div>
      </article>

      <section className="mt-10">
        <div className="knot-divider" aria-hidden="true">
          <span />
        </div>
        <div className="mt-6 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 className="font-display text-2xl font-medium">What this story teaches</h2>
            <p className="mt-1 text-sm text-ink-500">
              Hearthside listens for proverbs, warnings, and values — and keeps them here.
            </p>
          </div>
          <button
            onClick={extractWisdom}
            disabled={extracting}
            className="rounded-full bg-ink-900 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:opacity-60"
          >
            {extracting ? "Listening…" : extracts.length > 0 ? "Listen again" : "Draw out the wisdom"}
          </button>
        </div>

        {error && (
          <div
            role="alert"
            className="mt-4 rounded-xl border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
          >
            {error}
          </div>
        )}

        {extracts.length === 0 ? (
          <div className="mt-4 rounded-2xl border border-dashed border-parchment-300 bg-parchment-100 p-8 text-center">
            <p className="font-display text-lg font-medium">
              {wisdomMessage ?? "No wisdom drawn from this story yet."}
            </p>
            {!wisdomMessage && (
              <p className="mx-auto mt-2 max-w-md text-sm leading-6 text-ink-500">
                Stories that teach mention what happened, what someone said — proverbs,
                advice, warnings, values — and how it changed things.
              </p>
            )}
          </div>
        ) : (
          <ul className="mt-4 space-y-3">
            {extracts.map((w) => (
              <li
                key={w.id}
                className="rounded-2xl bg-ink-950 p-6 text-parchment-100"
              >
                <p className="font-display text-xl font-medium italic leading-snug text-parchment-50">
                  “{w.excerpt}”
                </p>
                <p className="mt-3 text-sm leading-6 text-parchment-200/85">{w.meaning}</p>
                <p className="mt-3 font-mono text-[11px] tracking-[0.2em] text-gold-200">
                  {w.wisdom_type.toUpperCase()} · KEPT WITH {(w.confidence * 100).toFixed(0)}% CERTAINTY
                </p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
