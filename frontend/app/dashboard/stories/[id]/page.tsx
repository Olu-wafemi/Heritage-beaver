"use client";

import { use, useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import type { Story, WisdomExtract } from "@/lib/types";

export default function StoryDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [story, setStory] = useState<Story | null>(null);
  const [extracts, setExtracts] = useState<WisdomExtract[]>([]);
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
    try {
      const token = localStorage.getItem("token") ?? "";
      await apiFetch(`/stories/${id}/process-wisdom`, { method: "POST", token });
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to extract wisdom");
    } finally {
      setExtracting(false);
    }
  }

  if (error && !story) {
    return (
      <div>
        <div role="alert" className="rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
          {error}
        </div>
        <Link href="/dashboard/stories" className="mt-4 inline-block text-sm text-amber-700 hover:underline">
          ← Back to stories
        </Link>
      </div>
    );
  }

  if (!story) {
    return <div className="h-64 animate-pulse rounded-xl bg-stone-200/70" />;
  }

  return (
    <div className="mx-auto max-w-3xl">
      <Link href="/dashboard/stories" className="text-sm text-stone-500 hover:text-stone-800">
        ← Back to stories
      </Link>

      <article className="mt-6 rounded-xl border border-stone-200 bg-white p-8 shadow-sm">
        <h1 className="text-3xl font-bold tracking-tight text-stone-900">{story.title}</h1>
        <p className="mt-2 text-xs uppercase tracking-wide text-stone-400">
          {new Date(story.created_at).toLocaleDateString()}
          {story.source_language ? ` · ${story.source_language}` : ""}
        </p>

        {story.summary && (
          <p className="mt-4 border-l-2 border-amber-300 pl-4 text-sm italic text-stone-600">
            {story.summary}
          </p>
        )}

        <div className="mt-6 whitespace-pre-wrap leading-7 text-stone-800">{story.content}</div>
      </article>

      <section className="mt-8">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-stone-900">Wisdom extracted</h2>
          <button
            onClick={extractWisdom}
            disabled={extracting}
            className="rounded-lg border border-amber-700 px-3 py-1.5 text-sm font-medium text-amber-800 transition hover:bg-amber-50 disabled:opacity-60"
          >
            {extracting ? "Listening…" : "Extract wisdom"}
          </button>
        </div>

        {error && (
          <div role="alert" className="mt-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        {extracts.length === 0 ? (
          <p className="mt-4 rounded-xl border border-dashed border-stone-300 bg-white p-8 text-center text-sm text-stone-600">
            No wisdom drawn from this story yet.
          </p>
        ) : (
          <ul className="mt-4 space-y-3">
            {extracts.map((w) => (
              <li key={w.id} className="rounded-xl border border-stone-200 bg-white p-5 shadow-sm">
                <p className="font-medium text-stone-900">“{w.excerpt}”</p>
                <p className="mt-2 text-sm text-stone-600">{w.meaning}</p>
                <p className="mt-2 text-xs uppercase tracking-wide text-stone-400">
                  {w.wisdom_type} · confidence {(w.confidence * 100).toFixed(0)}%
                </p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
