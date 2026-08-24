"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import type { FamilyMember } from "@/lib/types";

export default function DashboardPage() {
  const [memberCount, setMemberCount] = useState<number | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;
    apiFetch<{ members: FamilyMember[] }>("/family/tree", { token })
      .then((res) => setMemberCount(res.members.length))
      .catch(() => setMemberCount(0));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold tracking-tight text-stone-900">Dashboard</h1>
      <p className="mt-1 text-sm text-stone-600">
        The heart of your family&apos;s living mythology.
      </p>

      <div className="mt-8 grid gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-stone-200 bg-white p-6 shadow-sm">
          <p className="text-sm font-medium text-stone-500">Family members</p>
          <p className="mt-2 text-3xl font-bold text-stone-900">
            {memberCount === null ? "…" : memberCount}
          </p>
        </div>
        <div className="rounded-xl border border-stone-200 bg-white p-6 shadow-sm opacity-60">
          <p className="text-sm font-medium text-stone-500">Stories</p>
          <p className="mt-2 text-3xl font-bold text-stone-900">—</p>
          <p className="mt-1 text-xs text-stone-400">Coming in Phase 1 UI</p>
        </div>
        <div className="rounded-xl border border-stone-200 bg-white p-6 shadow-sm opacity-60">
          <p className="text-sm font-medium text-stone-500">Wisdom extracts</p>
          <p className="mt-2 text-3xl font-bold text-stone-900">—</p>
          <p className="mt-1 text-xs text-stone-400">Coming in Phase 1 UI</p>
        </div>
      </div>
    </div>
  );
}
