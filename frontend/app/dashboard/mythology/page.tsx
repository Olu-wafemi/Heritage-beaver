"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import type { MythChapter, Story } from "@/lib/types";

export default function MythologyPage() {
  const [chapters, setChapters] = useState<MythChapter[] | null>(null);
  const [stories, setStories] = useState<Story[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [theme, setTheme] = useState("");
  const [chapterType, setChapterType] = useState("origin");
  const [selected, setSelected] = useState<string[]>([]);
  const [weaving, setWeaving] = useState(false);

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const [chapterList, storyList] = await Promise.all([
        apiFetch<MythChapter[]>("/mythology/chapters", { token }),
        apiFetch<Story[]>("/stories", { token }),
      ]);
      setChapters(chapterList);
      setStories(storyList);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load mythology");
      setChapters([]);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  function toggleStory(id: string) {
    setSelected((prev) => (prev.includes(id) ? prev.filter((s) => s !== id) : [...prev, id]));
  }

  async function onWeave(e: FormEvent) {
    e.preventDefault();
    setWeaving(true);
    setError(null);
    try {
      const token = localStorage.getItem("token") ?? "";
      await apiFetch("/mythology/chapters", {
        method: "POST",
        token,
        body: JSON.stringify({
          theme,
          chapter_type: chapterType || "origin",
          story_ids: selected,
        }),
      });
      setTheme("");
      setSelected([]);
      setShowForm(false);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to weave chapter");
    } finally {
      setWeaving(false);
    }
  }

  return (
    <div>
      <p className="font-mono text-xs tracking-[0.25em] text-ember-600">STAGE IV — BECOME</p>
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="mt-2 font-display text-4xl font-medium tracking-tight">The Myth</h1>
          <p className="mt-2 max-w-lg leading-7 text-ink-700">
            Your stories, gathered into chapters of a living mythology — the book of your
            lineage, still being written.
          </p>
        </div>
        <button
          onClick={() => setShowForm((v) => !v)}
          disabled={stories.length === 0}
          title={stories.length === 0 ? "Tell at least one story first" : undefined}
          className="rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Weave a chapter
        </button>
      </div>

      {error && (
        <div
          role="alert"
          className="mt-6 rounded-xl border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
        >
          {error}
        </div>
      )}

      {showForm && (
        <form
          onSubmit={onWeave}
          className="mt-6 space-y-4 rounded-2xl border border-parchment-300 bg-parchment-100 p-6"
        >
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label className="block text-sm font-medium text-ink-900">Theme</label>
              <input
                type="text"
                value={theme}
                onChange={(e) => setTheme(e.target.value)}
                placeholder="resilience, homecoming, sacrifice…"
                className="mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition placeholder:text-ink-400 focus:border-ember-600 focus:ring-2 focus:ring-ember-100"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-ink-900">Chapter kind</label>
              <select
                value={chapterType}
                onChange={(e) => setChapterType(e.target.value)}
                className="mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition focus:border-ember-600 focus:ring-2 focus:ring-ember-100"
              >
                <option value="origin">Origin — how things began</option>
                <option value="trial">Trial — what was endured</option>
                <option value="homecoming">Homecoming — what was found</option>
                <option value="blessing">Blessing — what is passed on</option>
              </select>
            </div>
          </div>

          <div>
            <p className="text-sm font-medium text-ink-900">
              Weave from <span className="text-ink-500">(leave empty to use all stories)</span>
            </p>
            <ul className="mt-2 max-h-44 space-y-1.5 overflow-y-auto rounded-xl border border-parchment-300 bg-parchment-50 p-3">
              {stories.map((s) => (
                <li key={s.id}>
                  <label className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2 py-1.5 text-sm transition hover:bg-parchment-100">
                    <input
                      type="checkbox"
                      checked={selected.includes(s.id)}
                      onChange={() => toggleStory(s.id)}
                      className="h-4 w-4 accent-[#b4430f]"
                    />
                    <span className="truncate">{s.title}</span>
                  </label>
                </li>
              ))}
            </ul>
          </div>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={weaving}
              className="rounded-full bg-ink-900 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:opacity-60"
            >
              {weaving ? "Weaving…" : "Weave this chapter"}
            </button>
            <button
              type="button"
              onClick={() => setShowForm(false)}
              className="rounded-full border border-ink-900/20 px-5 py-2.5 text-sm font-medium text-ink-700 transition hover:bg-parchment-50"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      <div className="mt-8">
        {chapters === null ? (
          <div className="space-y-3">
            {[0, 1].map((i) => (
              <div key={i} className="h-32 animate-pulse rounded-2xl bg-parchment-200/70" />
            ))}
          </div>
        ) : chapters.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-parchment-300 bg-parchment-100 p-12 text-center">
            <p className="font-display text-2xl font-medium">The book awaits its first chapter</p>
            <p className="mx-auto mt-2 max-w-sm text-sm leading-6 text-ink-700">
              {stories.length === 0
                ? "Tell a story first — chapters are woven from the stories you keep."
                : "Gather your stories into the first chapter of your family's mythology."}
            </p>
            {stories.length === 0 ? (
              <Link
                href="/dashboard/stories"
                className="mt-5 inline-block rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
              >
                Tell a story first
              </Link>
            ) : (
              <button
                onClick={() => setShowForm(true)}
                className="mt-5 rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
              >
                Weave the first chapter
              </button>
            )}
          </div>
        ) : (
          <ol className="relative space-y-6 border-l-2 border-parchment-300 pl-0">
            {chapters.map((c, i) => (
              <li key={c.id} className="relative pl-10">
                <span
                  aria-hidden="true"
                  className="absolute -left-[9px] top-7 h-4 w-4 rotate-45 border-2 border-parchment-50 bg-ember-600"
                />
                <Link
                  href={`/dashboard/mythology/${c.id}`}
                  className="group block overflow-hidden rounded-2xl border border-parchment-300 bg-parchment-100 transition hover:border-gold-500"
                >
                  <div className="p-6 sm:p-7">
                    <p className="font-mono text-xs tracking-[0.25em] text-ember-600">
                      CHAPTER {chapters.length - i} · {c.theme.toUpperCase()} ·{" "}
                      {c.chapter_type.toUpperCase()}
                    </p>
                    <h2 className="mt-2 font-display text-2xl font-medium tracking-tight group-hover:text-ember-700 sm:text-3xl">
                      {c.title}
                    </h2>
                    <p className="mt-3 line-clamp-3 font-display text-base italic leading-7 text-ink-700">
                      {c.narrative}
                    </p>
                    <p className="mt-3 font-mono text-xs text-ink-400">
                      {new Date(c.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </Link>
              </li>
            ))}
          </ol>
        )}
      </div>
    </div>
  );
}
