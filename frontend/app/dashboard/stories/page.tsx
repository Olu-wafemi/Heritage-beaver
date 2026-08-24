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

  return (
    <div>
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-stone-900">Stories</h1>
          <p className="mt-1 text-sm text-stone-600">
            Every family myth begins with a true story.
          </p>
        </div>
        <button
          onClick={() => setShowForm((v) => !v)}
          className="rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-amber-800"
        >
          New story
        </button>
      </div>

      {error && (
        <div role="alert" className="mt-6 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">
          {error}
        </div>
      )}

      {showForm && (
        <form
          onSubmit={onSubmit}
          className="mt-6 space-y-4 rounded-xl border border-stone-200 bg-white p-6 shadow-sm"
        >
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-stone-800">
                Title<span className="text-red-500"> *</span>
              </label>
              <input
                type="text"
                required
                value={form.title}
                onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
                placeholder="How grandma survived the flood"
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              />
            </div>

            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-stone-800">
                Story<span className="text-red-500"> *</span>
              </label>
              <textarea
                required
                rows={6}
                value={form.content}
                onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))}
                placeholder="Tell it the way you remember hearing it…"
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-stone-800">Language</label>
              <input
                type="text"
                value={form.source_language}
                onChange={(e) => setForm((f) => ({ ...f, source_language: e.target.value }))}
                placeholder="Yoruba, Igbo, English…"
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-stone-800">About whom</label>
              <select
                value={form.family_member_id}
                onChange={(e) => setForm((f) => ({ ...f, family_member_id: e.target.value }))}
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              >
                <option value="">— General family story —</option>
                {members.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.display_name}
                  </option>
                ))}
              </select>
            </div>

            <div className="sm:col-span-2">
              <label className="block text-sm font-medium text-stone-800">Summary</label>
              <input
                type="text"
                value={form.summary}
                onChange={(e) => setForm((f) => ({ ...f, summary: e.target.value }))}
                placeholder="One line, for the record"
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              />
            </div>
          </div>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-amber-800 disabled:opacity-60"
            >
              {saving ? "Saving…" : "Save story"}
            </button>
            <button
              type="button"
              onClick={() => setShowForm(false)}
              className="rounded-lg border border-stone-300 px-4 py-2 text-sm font-medium text-stone-700 hover:bg-stone-100"
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
              <div key={i} className="h-20 animate-pulse rounded-xl bg-stone-200/70" />
            ))}
          </div>
        ) : stories.length === 0 ? (
          <div className="rounded-xl border border-dashed border-stone-300 bg-white p-12 text-center">
            <p className="text-lg font-medium text-stone-900">No stories yet</p>
            <p className="mx-auto mt-1 max-w-sm text-sm text-stone-600">
              Write down a memory — a proverb, a journey, a lesson. It becomes part of your myth.
            </p>
            <button
              onClick={() => setShowForm(true)}
              className="mt-5 rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-800"
            >
              Write your first story
            </button>
          </div>
        ) : (
          <ul className="space-y-3">
            {stories.map((s) => (
              <li key={s.id}>
                <Link
                  href={`/dashboard/stories/${s.id}`}
                  className="block rounded-xl border border-stone-200 bg-white p-5 shadow-sm transition hover:border-amber-300 hover:bg-amber-50/40"
                >
                  <div className="flex items-center justify-between gap-4">
                    <h2 className="font-medium text-stone-900">{s.title}</h2>
                    <span className="shrink-0 text-xs text-stone-400">
                      {new Date(s.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  {s.summary && (
                    <p className="mt-1 line-clamp-2 text-sm text-stone-600">{s.summary}</p>
                  )}
                  <p className="mt-2 line-clamp-1 text-xs italic text-stone-400">{s.content}</p>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
