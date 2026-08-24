"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import { useAuthStore } from "@/store/useAuthStore";
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
      const userId = useAuthStore.getState().user?.id ?? "";
      const body = JSON.stringify({ ...form, user_id: userId });
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
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-stone-900">Family</h1>
          <p className="mt-1 text-sm text-stone-600">
            The people whose stories weave your mythology.
          </p>
        </div>
        <button
          onClick={openCreate}
          className="rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-amber-800"
        >
          Add member
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
              <label className="block text-sm font-medium text-stone-800">Living</label>
              <select
                value={form.is_living ? "yes" : "no"}
                onChange={(e) => setForm((f) => ({ ...f, is_living: e.target.value === "yes" }))}
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              >
                <option value="yes">Living</option>
                <option value="no">Deceased</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-stone-800">Biography</label>
            <textarea
              rows={3}
              value={form.biography}
              onChange={(e) => setForm((f) => ({ ...f, biography: e.target.value }))}
              placeholder="A few words about their life…"
              className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
            />
          </div>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-amber-800 disabled:opacity-60"
            >
              {saving ? "Saving…" : editingId ? "Save changes" : "Add member"}
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
        {members === null ? (
          <div className="space-y-3">
            {[0, 1, 2].map((i) => (
              <div key={i} className="h-16 animate-pulse rounded-xl bg-stone-200/70" />
            ))}
          </div>
        ) : members.length === 0 ? (
          <div className="rounded-xl border border-dashed border-stone-300 bg-white p-12 text-center">
            <p className="text-lg font-medium text-stone-900">No family members yet</p>
            <p className="mx-auto mt-1 max-w-sm text-sm text-stone-600">
              Every myth begins with its people. Add your first family member to start weaving.
            </p>
            <button
              onClick={openCreate}
              className="mt-5 rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white hover:bg-amber-800"
            >
              Add your first member
            </button>
          </div>
        ) : (
          <ul className="space-y-3">
            {members.map((m) => (
              <li
                key={m.id}
                className="flex items-center justify-between rounded-xl border border-stone-200 bg-white p-4 shadow-sm transition hover:border-stone-300"
              >
                <div>
                  <p className="font-medium text-stone-900">{m.display_name}</p>
                  <p className="text-sm text-stone-500">
                    {[m.gender, m.birth_place, m.is_living ? "living" : "in memory"]
                      .filter(Boolean)
                      .join(" · ")}
                  </p>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => openEdit(m)}
                    className="rounded-lg border border-stone-300 px-3 py-1.5 text-xs font-medium text-stone-700 hover:bg-stone-100"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(m.id)}
                    className="rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50"
                  >
                    Delete
                  </button>
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
      <label className="block text-sm font-medium text-stone-800">
        {label}
        {required && <span className="text-red-500"> *</span>}
      </label>
      <input
        type="text"
        required={required}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
      />
    </div>
  );
}
