"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
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

const inputClass =
  "mt-1.5 w-full rounded-lg border border-parchment-300 bg-parchment-50 px-3 py-2.5 text-sm text-ink-900 outline-none transition focus:border-ember-600 focus:ring-2 focus:ring-ember-100";

export default function RelationshipsPage() {
  const [relationships, setRelationships] = useState<Relationship[] | null>(null);
  const [members, setMembers] = useState<FamilyMember[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [saving, setSaving] = useState(false);
  const [sourceId, setSourceId] = useState("");
  const [targetId, setTargetId] = useState("");
  const [relType, setRelType] = useState("");

  const load = useCallback(async () => {
    try {
      const token = localStorage.getItem("token") ?? "";
      const [rels, tree] = await Promise.all([
        apiFetch<Relationship[]>("/relationships", { token }),
        apiFetch<{ members: FamilyMember[] }>("/family/tree", { token }),
      ]);
      setRelationships(rels);
      setMembers(tree.members);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load relationships");
      setRelationships([]);
    }
  }, []);

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
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="font-mono text-xs tracking-[0.25em] text-ember-600">
            STAGE I — GATHER, CONTINUED
          </p>
          <h1 className="mt-2 font-display text-4xl font-medium tracking-tight">The Bonds</h1>
          <p className="mt-2 max-w-lg leading-7 text-ink-700">
            Tie your people together — parent to child, sibling to sibling. The tree takes
            shape from these threads.
          </p>
        </div>
        <button
          onClick={() => setShowForm((v) => !v)}
          disabled={members.length < 2}
          title={members.length < 2 ? "Name at least two family members first" : undefined}
          className="rounded-full bg-ember-600 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Tie a thread
        </button>
      </div>

      {members.length < 2 && (
        <p className="mt-6 rounded-xl border border-gold-500/40 bg-gold-200/20 px-4 py-3 text-sm text-ink-700">
          Name at least two family members before tying threads between them.
        </p>
      )}

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
          <div className="grid items-end gap-4 sm:grid-cols-[1fr_auto_1fr] sm:gap-2">
            <div>
              <label className="block text-sm font-medium text-ink-900">
                From<span className="text-ember-600"> *</span>
              </label>
              <select
                required
                value={sourceId}
                onChange={(e) => setSourceId(e.target.value)}
                className={inputClass}
              >
                <option value="">Choose a person</option>
                {members.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.display_name}
                  </option>
                ))}
              </select>
            </div>

            <div className="pb-1 text-center">
              <span className="inline-block rounded-full bg-ember-100 px-3 py-1 font-mono text-xs text-ember-700">
                is the
              </span>
            </div>

            <div>
              <label className="block text-sm font-medium text-ink-900">
                Of<span className="text-ember-600"> *</span>
              </label>
              <select
                required
                value={targetId}
                onChange={(e) => setTargetId(e.target.value)}
                className={inputClass}
              >
                <option value="">Choose a person</option>
                {members.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.display_name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-ink-900">
              Bond<span className="text-ember-600"> *</span>
            </label>
            <select
              required
              value={relType}
              onChange={(e) => setRelType(e.target.value)}
              className={inputClass}
            >
              <option value="">What is the bond?</option>
              {RELATIONSHIP_TYPES.map((t) => (
                <option key={t} value={t}>
                  {t.replace("_", " ")}
                </option>
              ))}
            </select>
          </div>

          <p className="font-display text-lg italic text-ink-700">
            {memberName(sourceId) || "…"} is the {relType ? relType.replace("_", " ") : "…"} of{" "}
            {memberName(targetId) || "…"}.
          </p>

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={saving}
              className="rounded-full bg-ink-900 px-5 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700 disabled:opacity-60"
            >
              {saving ? "Tying…" : "Tie the thread"}
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
        {relationships === null ? (
          <div className="space-y-3">
            {[0, 1].map((i) => (
              <div key={i} className="h-16 animate-pulse rounded-2xl bg-parchment-200/70" />
            ))}
          </div>
        ) : relationships.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-parchment-300 bg-parchment-100 p-12 text-center">
            <p className="font-display text-2xl font-medium">No threads tied yet</p>
            <p className="mx-auto mt-2 max-w-sm text-sm leading-6 text-ink-700">
              Connect parents to children, siblings to each other — the tree takes shape
              from these threads.
            </p>
          </div>
        ) : (
          <ul className="space-y-3">
            {relationships.map((r) => (
              <li
                key={r.id}
                className="flex items-center justify-between gap-4 rounded-2xl border border-parchment-300 bg-parchment-100 p-4"
              >
                <p className="text-sm text-ink-900">
                  <span className="font-display text-base font-semibold">
                    {memberName(r.source_member_id)}
                  </span>
                  <span className="mx-2 inline-block rounded-full bg-ember-600 px-2.5 py-0.5 align-middle font-mono text-[11px] text-parchment-50">
                    {r.relationship_type.replace("_", " ")}
                  </span>
                  <span className="font-display text-base font-semibold">
                    {memberName(r.target_member_id)}
                  </span>
                </p>
                <button
                  onClick={() => onDelete(r.id)}
                  className="shrink-0 rounded-full border border-red-300 px-3 py-1.5 text-xs font-medium text-red-700 transition hover:bg-red-50"
                >
                  Untie
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
