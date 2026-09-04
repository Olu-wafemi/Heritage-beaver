"use client";

import { FormEvent, use, useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import type { MythChapter } from "@/lib/types";

export default function ChapterDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const [chapter, setChapter] = useState<MythChapter | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [title, setTitle] = useState("");
  const [theme, setTheme] = useState("");
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const c = await apiFetch<MythChapter>(`/mythology/chapters/${id}`, { token });
      setChapter(c);
      setTitle(c.title);
      setTheme(c.theme);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load chapter");
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  async function onSave(e: FormEvent) {
    e.preventDefault();
    if (!chapter) return;
    setSaving(true);
    setError(null);
    try {
      const token = localStorage.getItem("token") ?? "";
      const updated = await apiFetch<MythChapter>(`/mythology/chapters/${id}`, {
        method: "PATCH",
        token,
        body: JSON.stringify({
          title,
          theme,
          chapter_type: chapter.chapter_type,
          narrative: chapter.narrative,
        }),
      });
      setChapter(updated);
      setEditing(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save chapter");
    } finally {
      setSaving(false);
    }
  }

  if (error && !chapter) {
    return (
      <div>
        <div
          role="alert"
          className="rounded-xl border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
        >
          {error}
        </div>
        <Link
          href="/dashboard/mythology"
          className="mt-4 inline-block text-sm font-medium text-ember-700 hover:underline"
        >
          ← Back to the myth
        </Link>
      </div>
    );
  }

  if (!chapter) {
    return (
      <div className="mx-auto max-w-3xl space-y-4">
        <div className="h-8 w-2/3 animate-pulse rounded-lg bg-parchment-200/70" />
        <div className="h-96 animate-pulse rounded-2xl bg-parchment-200/70" />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl">
      <Link
        href="/dashboard/mythology"
        className="text-sm font-medium text-ink-500 transition hover:text-ember-700"
      >
        ← Back to the myth
      </Link>

      <article className="mt-6 overflow-hidden rounded-2xl bg-ink-950 text-parchment-100">
        <div className="kente-strip" aria-hidden="true" />
        <div className="p-8 sm:p-12">
          <p className="font-mono text-xs tracking-[0.25em] text-gold-200">
            {chapter.theme.toUpperCase()} · {chapter.chapter_type.toUpperCase()}
          </p>
          {editing ? (
            <form onSubmit={onSave} className="mt-4 space-y-4">
              <input
                type="text"
                required
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full rounded-lg border border-white/20 bg-white/5 px-3 py-2.5 font-display text-2xl font-medium text-parchment-50 outline-none focus:border-gold-500"
              />
              <input
                type="text"
                value={theme}
                onChange={(e) => setTheme(e.target.value)}
                placeholder="theme"
                className="w-full rounded-lg border border-white/20 bg-white/5 px-3 py-2.5 text-sm text-parchment-100 outline-none placeholder:text-parchment-200/40 focus:border-gold-500"
              />
              <div className="flex gap-3">
                <button
                  type="submit"
                  disabled={saving}
                  className="rounded-full bg-ember-600 px-5 py-2 text-sm font-semibold text-parchment-50 transition hover:bg-ember-500 disabled:opacity-60"
                >
                  {saving ? "Keeping…" : "Keep changes"}
                </button>
                <button
                  type="button"
                  onClick={() => setEditing(false)}
                  className="rounded-full border border-white/20 px-5 py-2 text-sm text-parchment-200 transition hover:border-white/40"
                >
                  Cancel
                </button>
              </div>
            </form>
          ) : (
            <>
              <h1 className="mt-3 font-display text-4xl font-medium leading-tight tracking-tight text-parchment-50 sm:text-5xl">
                {chapter.title}
              </h1>
              <div className="manuscript mt-8 whitespace-pre-wrap !text-parchment-100/90">
                {chapter.narrative}
              </div>
              <div className="mt-10 flex items-center justify-between border-t border-white/10 pt-5">
                <p className="font-mono text-xs tracking-widest text-parchment-200/50">
                  WOVEN {new Date(chapter.created_at).toLocaleDateString().toUpperCase()}
                </p>
                <button
                  onClick={() => setEditing(true)}
                  className="rounded-full border border-white/20 px-4 py-1.5 text-xs font-medium text-parchment-200 transition hover:border-gold-500 hover:text-gold-200"
                >
                  Retell the title
                </button>
              </div>
            </>
          )}
        </div>
      </article>

      {error && (
        <div
          role="alert"
          className="mt-4 rounded-xl border border-red-300 bg-red-50 px-4 py-3 text-sm text-red-800"
        >
          {error}
        </div>
      )}
    </div>
  );
}
