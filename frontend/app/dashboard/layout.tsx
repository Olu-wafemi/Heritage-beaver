"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import { ProtectedRoute } from "@/components/ProtectedRoute";

import { logout as serverLogout } from "@/lib/auth";

const navItems = [
  { href: "/dashboard", label: "Hearth", hint: "Where the telling stands" },
  { href: "/dashboard/family", label: "The People", hint: "Name your family" },
  { href: "/dashboard/relationships", label: "The Bonds", hint: "Tie the tree together" },
  { href: "/dashboard/stories", label: "The Stories", hint: "Tell what happened" },
  { href: "/dashboard/wisdom", label: "The Wisdom", hint: "What the stories teach" },
  { href: "/dashboard/mythology", label: "The Myth", hint: "Chapters of your lineage" },
];

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.clearAuth);

  function logout() {
    const refreshToken = useAuthStore.getState().refreshToken;
    setAuth();
    router.push("/login");
    void serverLogout(refreshToken);
  }

  return (
    <ProtectedRoute>
      <div className="flex min-h-screen bg-parchment-50">
        <aside className="relative flex w-64 shrink-0 flex-col bg-ink-950 text-parchment-100">
          <div className="kente-strip absolute inset-y-0 left-0 w-1.5" aria-hidden="true" />
          <div className="px-6 pb-5 pt-7">
            <Link
              href="/dashboard"
              className="font-display text-xl font-semibold tracking-tight text-parchment-50"
            >
              Hearthside
            </Link>
            <p className="mt-1 font-mono text-[11px] tracking-[0.2em] text-parchment-200/60">
              YOUR LIVING FAMILY MYTH
            </p>
          </div>

          <nav className="flex-1 space-y-1 px-3">
            {navItems.map((item) => {
              const active =
                item.href === "/dashboard"
                  ? pathname === "/dashboard"
                  : pathname.startsWith(item.href);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`group block rounded-xl px-4 py-3 transition ${
                    active
                      ? "bg-ember-600 text-parchment-50"
                      : "text-parchment-200/80 hover:bg-white/5 hover:text-parchment-50"
                  }`}
                >
                  <span className="block text-sm font-semibold">{item.label}</span>
                  <span
                    className={`mt-0.5 block text-xs ${
                      active ? "text-parchment-100/80" : "text-parchment-200/50"
                    }`}
                  >
                    {item.hint}
                  </span>
                </Link>
              );
            })}
          </nav>

          <div className="p-4">
            <button
              onClick={logout}
              className="w-full rounded-xl border border-white/15 px-3 py-2.5 text-sm font-medium text-parchment-200 transition hover:border-white/30 hover:text-parchment-50"
            >
              Leave the fireside
            </button>
          </div>
        </aside>

        <main className="min-w-0 flex-1 overflow-y-auto">
          <div className="mx-auto max-w-5xl px-8 py-10">{children}</div>
        </main>
      </div>
    </ProtectedRoute>
  );
}
