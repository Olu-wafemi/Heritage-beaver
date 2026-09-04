"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import type { FamilyMember, MythChapter, Story, WisdomExtract } from "@/lib/types";
import { useAuthStore } from "@/store/useAuthStore";

type Counts = {
  members: number | null;
  stories: number | null;
  wisdom: number | null;
  chapters: number | null;
};

export default function DashboardPage() {
  const user = useAuthStore((s) => s.user);
  const [counts, setCounts] = useState<Counts>({
    members: null,
    stories: null,
    wisdom: null,
    chapters: null,
  });

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;
    (async () => {
      try {
        const [tree, stories, wisdom, chapters] = await Promise.all([
          apiFetch<{ members: FamilyMember[] }>("/family/tree", { token }),
          apiFetch<Story[]>("/stories", { token }),
          apiFetch<WisdomExtract[]>("/wisdom-extracts", { token }),
          apiFetch<MythChapter[]>("/mythology/chapters", { token }),
        ]);
        setCounts({
          members: tree.members.length,
          stories: stories.length,
          wisdom: wisdom.length,
          chapters: chapters.length,
        });
      } catch {
        setCounts({ members: 0, stories: 0, wisdom: 0, chapters: 0 });
      }
    })();
  }, []);

  const firstName = user?.display_name?.split(" ")?.[0] ?? "Storykeeper";

  const stages = [
    {
      numeral: "I",
      name: "The People",
      count: counts.members,
      href: "/dashboard/family",
      cta: "Name someone",
      done: (counts.members ?? 0) > 0,
    },
    {
      numeral: "II",
      name: "The Stories",
      count: counts.stories,
      href: "/dashboard/stories",
      cta: "Tell a story",
      done: (counts.stories ?? 0) > 0,
    },
    {
      numeral: "III",
      name: "The Wisdom",
      count: counts.wisdom,
      href: "/dashboard/wisdom",
      cta: "Draw wisdom",
      done: (counts.wisdom ?? 0) > 0,
    },
    {
      numeral: "IV",
      name: "The Myth",
      count: counts.chapters,
      href: "/dashboard/mythology",
      cta: "Weave a chapter",
      done: (counts.chapters ?? 0) > 0,
    },
  ];

  const next = stages.find((s) => !s.done);

  return (
    <div>
      <p className="font-mono text-xs tracking-[0.25em] text-ember-600">THE HEARTH</p>
      <h1 className="mt-3 font-display text-4xl font-medium tracking-tight">
        Welcome back, {firstName}<span className="text-ember-600">.</span>
      </h1>
      <p className="mt-2 max-w-xl leading-7 text-ink-700">
        This is where your telling stands. Each stage feeds the next — people become
        stories, stories yield wisdom, wisdom becomes myth.
      </p>

      {next && (
        <Link
          href={next.href}
          className="mt-6 flex items-center justify-between gap-4 rounded-2xl bg-ink-950 p-6 text-parchment-50 transition hover:bg-ink-900"
        >
          <div>
            <p className="font-mono text-xs tracking-[0.2em] text-gold-200">
              NEXT IN THE TELLING
            </p>
            <p className="mt-2 font-display text-2xl font-medium">
              {next.cta} — continue stage {next.numeral}, {next.name}
            </p>
          </div>
          <span
            aria-hidden="true"
            className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-ember-600 text-xl"
          >
            →
          </span>
        </Link>
      )}

      <div className="mt-8 grid gap-px overflow-hidden rounded-2xl border border-parchment-300 bg-parchment-300 sm:grid-cols-2 lg:grid-cols-4">
        {stages.map((s) => (
          <Link key={s.numeral} href={s.href} className="group bg-parchment-50 p-6 transition hover:bg-parchment-100">
            <div className="flex items-center justify-between">
              <span className="font-mono text-xs tracking-widest text-gold-600">
                {s.numeral}
              </span>
              <span
                aria-hidden="true"
                className={`h-2.5 w-2.5 rotate-45 ${
                  s.done ? "bg-palm-600" : "border border-ink-400"
                }`}
                title={s.done ? "Underway" : "Not yet begun"}
              />
            </div>
            <p className="mt-3 font-display text-xl font-semibold">{s.name}</p>
            <p className="mt-1 font-display text-4xl font-medium text-ink-900">
              {s.count === null ? "…" : s.count}
            </p>
            <p className="mt-2 text-sm font-medium text-ember-700 opacity-0 transition group-hover:opacity-100">
              {s.cta} →
            </p>
          </Link>
        ))}
      </div>

      <div className="knot-divider mt-10" aria-hidden="true">
        <span />
      </div>
      <p className="mx-auto mt-6 max-w-xl text-center font-display text-lg italic text-ink-500">
        “The fire does not go out where the stories are still being told.”
      </p>
    </div>
  );
}
