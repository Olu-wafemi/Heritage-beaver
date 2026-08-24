import Link from "next/link";

const features = [
  {
    title: "Preserve their voices",
    body: "Record the stories, proverbs, and advice of the people who shaped you — before they fade.",
  },
  {
    title: "Weave a living myth",
    body: "Your family history becomes an evolving mythology, chapter by chapter, in your own words.",
  },
  {
    title: "Guidance from the ancestors",
    body: "Ask questions of AI-guided ancestors rooted in your family's wisdom and culture.",
  },
];

export default function Home() {
  return (
    <main className="min-h-screen bg-stone-50">
      <header className="mx-auto flex max-w-5xl items-center justify-between px-6 py-6">
        <span className="text-lg font-bold tracking-tight text-stone-900">Heritage Weaver</span>
        <nav className="flex items-center gap-4 text-sm font-medium">
          <Link href="/login" className="text-stone-600 transition hover:text-stone-900">
            Sign in
          </Link>
          <Link
            href="/register"
            className="rounded-lg bg-amber-700 px-4 py-2 text-white transition hover:bg-amber-800"
          >
            Get started
          </Link>
        </nav>
      </header>

      <section className="mx-auto max-w-3xl px-6 pb-20 pt-16 text-center sm:pt-24">
        <h1 className="text-4xl font-bold tracking-tight text-stone-900 sm:text-5xl">
          Your family&apos;s story,
          <br />
          woven into a living myth.
        </h1>
        <p className="mx-auto mt-6 max-w-xl text-lg leading-8 text-stone-600">
          Heritage Weaver preserves the stories, wisdom, and voices of your family — and turns them
          into culturally grounded guidance for the generations to come.
        </p>
        <div className="mt-10 flex flex-col items-center justify-center gap-3 sm:flex-row">
          <Link
            href="/register"
            className="w-full rounded-lg bg-amber-700 px-6 py-3 text-sm font-semibold text-white transition hover:bg-amber-800 sm:w-auto"
          >
            Start your family myth
          </Link>
          <Link
            href="/login"
            className="w-full rounded-lg border border-stone-300 bg-white px-6 py-3 text-sm font-semibold text-stone-800 transition hover:bg-stone-100 sm:w-auto"
          >
            I already have an account
          </Link>
        </div>
      </section>

      <section className="border-t border-stone-200 bg-white">
        <div className="mx-auto grid max-w-5xl gap-8 px-6 py-16 sm:grid-cols-3">
          {features.map((f) => (
            <div key={f.title}>
              <h2 className="font-semibold text-stone-900">{f.title}</h2>
              <p className="mt-2 text-sm leading-6 text-stone-600">{f.body}</p>
            </div>
          ))}
        </div>
      </section>

      <footer className="border-t border-stone-200 py-8 text-center text-xs text-stone-500">
        Heritage Weaver AI — personal mythology &amp; ancestral wisdom
      </footer>
    </main>
  );
}
