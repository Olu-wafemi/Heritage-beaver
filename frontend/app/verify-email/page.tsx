"use client";

import { Suspense, useEffect, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { verifyEmail } from "@/lib/auth";
import { useAuthStore } from "@/store/useAuthStore";

function VerifyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const setAuth = useAuthStore((s) => s.setAuth);
  const [state, setState] = useState<"working" | "done" | "failed">("working");
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    const token = searchParams.get("token");
    if (!token) {
      setState("failed");
      setMessage("This confirmation link is missing its token.");
      return;
    }
    verifyEmail(token)
      .then((res) => {
        setAuth(res.token, res.user, res.refresh_token);
        setState("done");
        setTimeout(() => router.push("/dashboard"), 1500);
      })
      .catch((err: unknown) => {
        setState("failed");
        setMessage(err instanceof Error ? err.message : "Confirmation failed");
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <main className="flex min-h-screen items-center justify-center bg-parchment-50 px-6">
      <div className="w-full max-w-md rounded-2xl border border-parchment-300 bg-parchment-100 p-10 text-center">
        <div className="kente-strip mx-auto w-24 rounded-full" aria-hidden="true" />
        {state === "working" && (
          <>
            <h1 className="mt-6 font-display text-3xl font-medium">Confirming…</h1>
            <p className="mt-2 text-sm text-ink-700">Lighting your place at the fireside.</p>
          </>
        )}
        {state === "done" && (
          <>
            <h1 className="mt-6 font-display text-3xl font-medium">
              Welcome home<span className="text-ember-600">.</span>
            </h1>
            <p className="mt-2 text-sm text-ink-700">
              Your email is confirmed — taking you to the hearth…
            </p>
          </>
        )}
        {state === "failed" && (
          <>
            <h1 className="mt-6 font-display text-3xl font-medium">Link spent<span className="text-ember-600">.</span></h1>
            <p className="mt-2 text-sm leading-6 text-ink-700">
              {message ?? "This confirmation link is invalid or expired."}
            </p>
            <Link
              href="/login"
              className="mt-6 inline-block rounded-full bg-ember-600 px-6 py-2.5 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
            >
              Back to sign in
            </Link>
          </>
        )}
      </div>
    </main>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense>
      <VerifyContent />
    </Suspense>
  );
}
