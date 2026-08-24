"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import { ProtectedRoute } from "@/components/ProtectedRoute";

const navItems = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/dashboard/family", label: "Family" },
  { href: "/dashboard/relationships", label: "Relationships" },
  { href: "/dashboard/stories", label: "Stories" },
  { href: "/dashboard/wisdom", label: "Wisdom" },
];

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.clearAuth);

  function logout() {
    setAuth();
    router.push("/login");
  }

  return (
    <ProtectedRoute>
      <div className="flex min-h-screen bg-stone-50">
        <aside className="flex w-60 flex-col border-r border-stone-200 bg-white">
          <div className="border-b border-stone-200 px-6 py-5">
            <Link href="/dashboard" className="text-lg font-bold tracking-tight text-stone-900">
              Heritage Weaver
            </Link>
            <p className="mt-0.5 text-xs text-stone-500">Your living family myth</p>
          </div>

          <nav className="flex-1 space-y-1 px-3 py-4">
            {navItems.map((item) => {
              const active =
                item.href === "/dashboard"
                  ? pathname === "/dashboard"
                  : pathname.startsWith(item.href);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`block rounded-lg px-3 py-2 text-sm font-medium transition ${
                    active
                      ? "bg-amber-50 text-amber-800"
                      : "text-stone-600 hover:bg-stone-100 hover:text-stone-900"
                  }`}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>

          <div className="border-t border-stone-200 p-4">
            <button
              onClick={logout}
              className="w-full rounded-lg border border-stone-300 px-3 py-2 text-sm font-medium text-stone-700 transition hover:bg-stone-100"
            >
              Sign out
            </button>
          </div>
        </aside>

        <main className="flex-1 overflow-y-auto">
          <div className="mx-auto max-w-5xl px-8 py-10">{children}</div>
        </main>
      </div>
    </ProtectedRoute>
  );
}
