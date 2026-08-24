"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import { useAuthStore } from "@/store/useAuthStore";
import type { FamilyMember, Relationship } from "@/lib/types";

const RELATIONSHIP_TYPES = [
  "parent",
  "child",
  "sibling",
  "spouse",
  "grandparent",
  "grandchild",
  "aunt_uncle",
  "niece_nephew",
  "cousin",
  "in_law",
];

export default function RelationshipsPage() {
  const user = useAuthStore((s) => s.user);
  const [relationships, setRelationships] = useState<Relationship[] | null>(null);
  const [members, setMembers] = useState<FamilyMember[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [saving, setSaving] = useState(false);
  const [sourceId, setSourceId] = useState("");
  const [targetId, setTargetId] = useState("");
  const [relType, setRelType] = useState("");

  const load = useCallback(async () => {
    if (!user?.id) return;
    try {
      const token = localStorage.getItem("token") ?? "";
      const [rels, tree] = await Promise.all([
        apiFetch<Relationship[]>(`/relationships?user_id=${user.id}`, { token }),
        apiFetch<{ members: FamilyMember[] }>("/family/tree", { token }),
      ]);
      setRelationships(rels);
      setMembers(tree.members);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load relationships");
      setRelationships([]);
    }
  }, [user?.id]);

  useEffect(() => {
    load();
  }, [load]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (sourceId === targetId) {
      setError("A person can't be related to themselves");
      return;
    }

    setSaving(true);
    try {
      const token = localStorage.getItem("token") ?? "";
      await apiFetch("/relationships", {
        method: "POST",
        token,
        body: JSON.stringify({
          user_id: user?.id,
          source_member_id: sourceId,
          target_member_id: targetId,
          relationship_type: relType,
        }),
      });
      setSourceId("");
      setTargetId("");
      setRelType("");
      setShowForm(false);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create relationship");
    } finally {
      setSaving(false);
    }
  }

  async function onDelete(id: string) {
    if (!confirm("Remove this relationship?")) return;
    try {
      const token = localStorage.getItem("token") ?? "";
      await apiFetch(`/relationships/${id}`, { method: "DELETE", token });
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete relationship");
    }
  }

  const memberName = (id: string) => members.find((m) => m.id === id)?.display_name ?? id;

  return (
    <div>
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-stone-900">Relationships</h1>
          <p className="mt-1 text-sm text-stone-600">
            Link your family members together — this becomes your family tree.
          </p>
        </div>
        <button
          onClick={() => setShowForm((v) => !v)}
          disabled={members.length < 2}
          title={members.length < 2 ? "Add at least two family members first" : undefined}
          className="rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-amber-800 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Link members
        </button>
      </div>

      {members.length < 2 && (
        <p className="mt-4 rounded-lg bg-amber-50 px-4 py-3 text-sm text-amber-800">
          Add at least two family members before linking them.
        </p>
      )}

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
          <div className="grid gap-4 sm:grid-cols-3">
            <div>
              <label className="block text-sm font-medium text-stone-800">
                From<span className="text-red-500"> *</span>
              </label>
              <select
                required
                value={sourceId}
                onChange={(e) => setSourceId(e.target.value)}
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              >
                <option value="">Select member</option>
                {members.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.display_name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-stone-800">
                Relationship<span className="text-red-500"> *</span>
              </label>
              <select
                required
                value={relType}
                onChange={(e) => setRelType(e.target.value)}
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              >
                <option value="">Select type</option>
                {RELATIONSHIP_TYPES.map((t) => (
                  <option key={t} value={t}>
                    {t.replace("_", " ")}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-stone-800">
                To<span className="text-red-500"> *</span>
              </label>
              <select
                required
                value={targetId}
                onChange={(e) => setTargetId(e.target.value)}
                className="mt-1.5 w-full rounded-lg border border-stone-300 px-3 py-2.5 text-sm outline-none focus:border-amber-600"
              >
                <option value="">Select member</option>
                {members.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.display_name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <p className="text-xs text-stone-500">
            Reads as: <span className="font-medium">[From]</span> is the{" "}
            <span className="font-medium">[relationship]</span> of{" "}
            <span className="font-medium">[To]</span>.
          </p>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-lg bg-amber-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-amber-800 disabled:opacity-60"
            >
              {saving ? "Linking…" : "Create link"}
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
        {relationships === null ? (
          <div className="space-y-3">
            {[0, 1].map((i) => (
              <div key={i} className="h-16 animate-pulse rounded-xl bg-stone-200/70" />
            ))}
          </div>
        ) : relationships.length === 0 ? (
          <div className="rounded-xl border border-dashed border-stone-300 bg-white p-12 text-center">
            <p className="text-lg font-medium text-stone-900">No links yet</p>
            <p className="mx-auto mt-1 max-w-sm text-sm text-stone-600">
              Connect parents to children, siblings to each other — the tree takes shape from
              these links.
            </p>
          </div>
        ) : (
          <ul className="space-y-3">
            {relationships.map((r) => (
              <li
                key={r.id}
                className="flex items-center justify-between rounded-xl border border-stone-200 bg-white p-4 shadow-sm"
              >
                <p className="text-sm text-stone-800">
                  <span className="font-medium">{memberName(r.source_member_id)}</span>
                  <span className="mx-2 text-stone-400">—</span>
                  <span className="rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
                    {r.relationship_type.replace("_", " ")}
                  </span>
                  <span className="mx-2 text-stone-400">→</span>
                  <span className="font-medium">{memberName(r.target_member_id)}</span>
                </p>
                <button
                  onClick={() => onDelete(r.id)}
                  className="rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50"
                >
                  Remove
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
