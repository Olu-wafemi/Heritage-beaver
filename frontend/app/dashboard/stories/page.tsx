"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import { useAuthStore } from "@/store/useAuthStore";
import type { FamilyMember, Story } from "@/lib/types";

type StoryForm = {
  title: string;
  content: string;
  source_type: string;
  source_language: string;
  summary: string;
  family_member_id: string;
};

const emptyForm: StoryForm = {
  title: "",
  content: "",
  source_type: "text",
  source_language: "",
  summary: "",
  family_member_id: "",
};

const inputClass =
  "mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition placeholder:text-ink-400 focus:border-ember-600 focus:ring-2 focus:ring-ember-100";

export default function StoriesPage() {
  const user = useAuthStore((s) => s.user);
  const [stories, setStories] = useState<Story[] | null>(null);
  const [members, setMembers] = useState<FamilyMember[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<StoryForm>(emptyForm);
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const [storyList, tree] = await Promise.all([
        apiFetch<Story[]>("/stories", { token }),
        apiFetch<{ members: FamilyMember[] }>("/family/tree", { token }),
      ]);
      setStories(storyList);
      setMembers(tree.members);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load stories");
      setStories([]);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError(null);
    try {
      const token = localStorage.getItem("token") ?? "";
      const body = JSON.stringify({
        title: form.title,
        content: form.content,
        source_type: form.source_type || "text",
        source_language: form.source_language,
        summary: form.summary,
        ...(form.family_member_id ? { family_member_id: form.family_member_id } : {}),
      });
      await apiFetch("/stories", { method: "POST", body, token });
      setForm(emptyForm);
      setShowForm(false);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save story");
    } finally {
      setSaving(false);
    }
  }

  void user;

  return (
    <div>
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="font-mono text-xs tracking-[0.25em] text-ember-600">STAGE II — REMEMBER</p>
          <h1 className="mt-2 font-display text-4xl font-medium tracking-tight">The Stories</h1>
          <p className="mt-2 max-w-lg leading-7 text-ink-700">
            Every family myth begins with a true story — told the way you remember hearing it.
          </p>
        </div>
        <button
          onClick={() => setShowForm((v) => !v)}
          className="rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
        >
          Tell a story
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
          onSubmit={onSubmit}
          className="mt-6 space-y-4 rounded-2xl border border-parchment-300 bg-parchment-100 p-6"
        >
          <div className="rounded-xl bg-ink-950 p-5 text-sm leading-6 text-parchment-100">
            <p className="font-display text-base font-semibold text-gold-200">
              What makes a story worth keeping?
            </p>
            <ul className="mt-2 list-disc space-y-1 pl-5 text-parchment-200/90">
              <li><strong>What happened</strong> — a moment, event, or memory, not just a title.</li>
              <li><strong>What someone said or taught</strong> — proverbs, advice, warnings, blessings.</li>
              <li><strong>How it changed things</strong> — what it taught you or the family.</li>
            </ul>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-ink-900">
                Title<span className="text-ember-600"> *</span>
              </label>
              <input
                type="text"
                required
                value={form.title}
                onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
                placeholder="How grandma survived the flood"
                className={inputClass}
              />
            </div>

            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-ink-900">
                Story<span className="text-ember-600"> *</span>
              </label>
              <textarea
                required
                rows={6}
                value={form.content}
                onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))}
                placeholder="Tell it the way you remember hearing it…"
                className={inputClass}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-ink-900">Tongue</label>
              <input
                type="text"
                value={form.source_language}
                onChange={(e) => setForm((f) => ({ ...f, source_language: e.target.value }))}
                placeholder="Yoruba, Igbo, English…"
                className={inputClass}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-ink-900">Whose story</label>
              <select
                value={form.family_member_id}
                onChange={(e) => setForm((f) => ({ ...f, family_member_id: e.target.value }))}
                className={inputClass}
              >
                <option value="">— A general family story —</option>
                {members.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.display_name}
                  </option>
                ))}
              </select>
            </div>

            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-ink-900">Summary</label>
              <input
                type="text"
                value={form.summary}
                onChange={(e) => setForm((f) => ({ ...f, summary: e.target.value }))}
                placeholder="One line, for the record"
                className={inputClass}
              />
            </div>
          </div>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-full bg-ink-900 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:opacity-60"
            >
              {saving ? "Keeping…" : "Keep this story"}
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
        {stories === null ? (
          <div className="space-y-3">
            {[0, 1, 2].map((i) => (
              <div key={i} className="h-24 animate-pulse rounded-2xl bg-parchment-200/70" />
            ))}
          </div>
        ) : stories.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-parchment-300 bg-parchment-100 p-12 text-center">
            <p className="font-display text-2xl font-medium">The page is blank</p>
            <p className="mx-auto mt-2 max-w-sm text-sm leading-6 text-ink-700">
              Write down a memory — a proverb, a journey, a lesson. It becomes part of your myth.
            </p>
            <button
              onClick={() => setShowForm(true)}
              className="mt-5 rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
            >
              Write the first story
            </button>
          </div>
        ) : (
          <ul className="space-y-3">
            {stories.map((s) => (
              <li key={s.id}>
                <Link
                  href={`/dashboard/stories/${s.id}`}
                  className="group flex gap-4 rounded-2xl border border-parchment-300 bg-parchment-100 p-5 transition hover:border-gold-500"
                >
                  <span className="kente-strip-vertical rounded-full" aria-hidden="true" />
                  <span className="min-w-0 flex-1">
                    <span className="flex items-baseline justify-between gap-4">
                      <span className="truncate font-display text-xl font-semibold group-hover:text-ember-700">
                        {s.title}
                      </span>
                      <span className="shrink-0 font-mono text-xs text-ink-400">
                        {new Date(s.created_at).toLocaleDateString()}
                      </span>
                    </span>
                    {s.summary ? (
                      <span className="mt-1 line-clamp-2 block text-sm leading-6 text-ink-700">
                        {s.summary}
                      </span>
                    ) : (
                      <span className="mt-1 line-clamp-1 block text-sm italic text-ink-500">
                        {s.content}
                      </span>
                    )}
                    {s.source_language && (
                      <span className="mt-2 inline-block rounded-full bg-parchment-200 px-2.5 py-0.5 text-xs font-medium text-ink-700">
                        {s.source_language}
                      </span>
                    )}
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
