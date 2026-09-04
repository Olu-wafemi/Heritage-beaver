"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import type { FamilyMember } from "@/lib/types";

type FormState = {
  first_name: string;
  last_name: string;
  display_name: string;
  gender: string;
  birth_place: string;
  biography: string;
  is_living: boolean;
};

const emptyForm: FormState = {
  first_name: "",
  last_name: "",
  display_name: "",
  gender: "",
  birth_place: "",
  biography: "",
  is_living: true,
};

const inputClass =
  "mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition placeholder:text-ink-400 focus:border-ember-600 focus:ring-2 focus:ring-ember-100";
const labelClass = "block text-sm font-medium text-ink-900";

export default function FamilyPage() {
  const [members, setMembers] = useState<FamilyMember[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const tree = await apiFetch<{ members: FamilyMember[] }>("/family/tree", { token });
      setMembers(tree.members);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load family members");
      setMembers([]);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  function openCreate() {
    setEditingId(null);
    setForm(emptyForm);
    setShowForm(true);
  }

  function openEdit(member: FamilyMember) {
    setEditingId(member.id);
    setForm({
      first_name: member.first_name,
      last_name: member.last_name,
      display_name: member.display_name,
      gender: member.gender,
      birth_place: member.birth_place,
      biography: member.biography,
      is_living: member.is_living,
    });
    setShowForm(true);
  }

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError(null);
    try {
      const token = localStorage.getItem("token") ?? "";
      const body = JSON.stringify({ ...form });
      if (editingId) {
        await apiFetch(`/family-members/${editingId}`, { method: "PATCH", body, token });
      } else {
        await apiFetch("/family-members", { method: "POST", body, token });
      }
      setShowForm(false);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save family member");
    } finally {
      setSaving(false);
    }
  }

  async function onDelete(id: string) {
    if (!confirm("Remove this family member?")) return;
    try {
      const token = localStorage.getItem("token") ?? "";
      await apiFetch(`/family-members/${id}`, { method: "DELETE", token });
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete family member");
    }
  }

  return (
    <div>
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="font-mono text-xs tracking-[0.25em] text-ember-600">STAGE I — GATHER</p>
          <h1 className="mt-2 font-display text-4xl font-medium tracking-tight">The People</h1>
          <p className="mt-2 max-w-lg leading-7 text-ink-700">
            The people whose stories weave your mythology — the living and the remembered.
          </p>
        </div>
        <button
          onClick={openCreate}
          className="rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
        >
          Name someone
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
          <div className="grid gap-4 sm:grid-cols-2">
            <Field label="First name" required value={form.first_name}
              onChange={(v) => setForm((f) => ({ ...f, first_name: v }))} />
            <Field label="Last name" value={form.last_name}
              onChange={(v) => setForm((f) => ({ ...f, last_name: v }))} />
            <Field label="Display name" required value={form.display_name}
              onChange={(v) => setForm((f) => ({ ...f, display_name: v }))} />
            <Field label="Gender" value={form.gender}
              onChange={(v) => setForm((f) => ({ ...f, gender: v }))} />
            <Field label="Birth place" value={form.birth_place}
              onChange={(v) => setForm((f) => ({ ...f, birth_place: v }))} />
            <div>
              <label className={labelClass}>Standing</label>
              <select
                value={form.is_living ? "yes" : "no"}
                onChange={(e) => setForm((f) => ({ ...f, is_living: e.target.value === "yes" }))}
                className={inputClass}
              >
                <option value="yes">Living</option>
                <option value="no">In memory (deceased)</option>
              </select>
            </div>
          </div>

          <div>
            <label className={labelClass}>Biography</label>
            <textarea
              rows={3}
              value={form.biography}
              onChange={(e) => setForm((f) => ({ ...f, biography: e.target.value }))}
              placeholder="A few words about their life…"
              className={inputClass}
            />
          </div>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-full bg-ink-900 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:opacity-60"
            >
              {saving ? "Keeping…" : editingId ? "Save changes" : "Add to the tree"}
            </button>
            <button
              type="button"
              onClick={() => setShowForm(false)}
              className="rounded-full border border-ink-900/20 px-5 py-2.5 text-sm font-medium text-ink-700 transition hover:bg-parchment-100"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      <div className="mt-8">
        {members === null ? (
          <div className="grid gap-4 sm:grid-cols-2">
            {[0, 1, 2, 3].map((i) => (
              <div key={i} className="h-28 animate-pulse rounded-2xl bg-parchment-200/70" />
            ))}
          </div>
        ) : members.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-parchment-300 bg-parchment-100 p-12 text-center">
            <p className="font-display text-2xl font-medium">No faces on the wall yet</p>
            <p className="mx-auto mt-2 max-w-sm text-sm leading-6 text-ink-700">
              Every myth begins with its people. Name your first family member to start
              weaving.
            </p>
            <button
              onClick={openCreate}
              className="mt-5 rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
            >
              Name the first person
            </button>
          </div>
        ) : (
          <ul className="grid gap-4 sm:grid-cols-2">
            {members.map((m) => (
              <li
                key={m.id}
                className="flex items-start gap-4 rounded-2xl border border-parchment-300 bg-parchment-100 p-5 transition hover:border-gold-500"
              >
                <span
                  aria-hidden="true"
                  className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-full font-display text-xl font-semibold ${
                    m.is_living ? "bg-palm-600 text-parchment-50" : "bg-ink-900 text-gold-200"
                  }`}
                >
                  {(m.display_name || "?").trim().charAt(0).toUpperCase()}
                </span>
                <div className="min-w-0 flex-1">
                  <p className="font-display text-lg font-semibold leading-snug">
                    {m.display_name}
                  </p>
                  <p className="mt-0.5 text-xs text-ink-500">
                    {[m.gender, m.birth_place, m.is_living ? "living" : "in memory"]
                      .filter(Boolean)
                      .join(" · ")}
                  </p>
                  {m.biography && (
                    <p className="mt-2 line-clamp-2 text-sm leading-6 text-ink-700">
                      {m.biography}
                    </p>
                  )}
                  <div className="mt-3 flex gap-2">
                    <button
                      onClick={() => openEdit(m)}
                      className="rounded-full border border-ink-900/20 px-3 py-1 text-xs font-medium text-ink-700 transition hover:bg-parchment-50"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => onDelete(m.id)}
                      className="rounded-full border border-red-300 px-3 py-1 text-xs font-medium text-red-700 transition hover:bg-red-50"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

function Field({
  label,
  value,
  onChange,
  required,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  required?: boolean;
}) {
  return (
    <div>
      <label className="block text-sm font-medium text-ink-900">
        {label}
        {required && <span className="text-ember-600"> *</span>}
      </label>
      <input
        type="text"
        required={required}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition placeholder:text-ink-400 focus:border-ember-600 focus:ring-2 focus:ring-ember-100"
      />
    </div>
  );
}
