import Link from "next/link";

const journey = [
  {
    step: "I — Gather",
    title: "Name the people",
    body: "Build the family tree: the living and the departed, each with a face, a place, a tongue.",
  },
  {
    step: "II — Remember",
    title: "Tell one true story",
    body: "A moment, an event, a memory — told the way you remember hearing it, in your own language.",
  },
  {
    step: "III — Listen",
    title: "Draw out the wisdom",
    body: "Proverbs, warnings, blessings and values surface from every story, automatically kept.",
  },
  {
    step: "IV — Become",
    title: "Weave the living myth",
    body: "Stories gather into chapters of your family's mythology — guidance for generations to come.",
  },
];

export default function Home() {
  return (
    <main className="min-h-screen bg-parchment-50 text-ink-900">
      <div className="kente-strip" aria-hidden="true" />

      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
        <span className="font-display text-xl font-semibold tracking-tight">
          Hearthside
        </span>
        <nav className="flex items-center gap-5 text-sm font-medium">
          <Link href="/login" className="text-ink-700 transition hover:text-ember-700">
            Sign in
          </Link>
          <Link
            href="/register"
            className="rounded-full bg-ink-900 px-5 py-2 text-parchment-50 transition hover:bg-ember-700"
          >
            Begin the telling
          </Link>
        </nav>
      </header>

      {/* Decide/Learn surface: one idea per section, editorial not centered-stack */}
      <section className="mx-auto grid max-w-6xl items-end gap-10 px-6 pb-16 pt-10 sm:pt-16 lg:grid-cols-[1.2fr_0.8fr]">
        <div>
          <p className="text-xs font-semibold uppercase tracking-[0.25em] text-ember-600">
            An ancestral wisdom platform
          </p>
          <h1 className="mt-5 font-display text-5xl font-medium leading-[1.05] tracking-tight sm:text-6xl">
            Every family is
            <br />
            a <em className="text-ember-600">mythology</em>
            <br />
            waiting to be told.
          </h1>
          <p className="mt-6 max-w-md text-lg leading-8 text-ink-700">
            Hearthside keeps the stories, proverbs, and voices of your people — and
            weaves them into guidance your grandchildren&apos;s grandchildren can hold.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              href="/register"
              className="rounded-full bg-ember-600 px-7 py-3 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
            >
              Start your family myth
            </Link>
            <Link
              href="/login"
              className="rounded-full border border-ink-900/20 px-7 py-3 text-sm font-semibold text-ink-900 transition hover:border-ink-900 hover:bg-parchment-100"
            >
              Continue the story
            </Link>
          </div>
        </div>

        <figure className="border border-parchment-300 bg-parchment-100 p-8 shadow-[8px_8px_0_0_var(--color-parchment-200)]">
          <blockquote className="font-display text-2xl font-medium italic leading-snug">
            “A child who is not embraced by the village will burn it down to feel its
            warmth.”
          </blockquote>
          <figcaption className="mt-4 text-xs font-semibold uppercase tracking-[0.2em] text-ink-500">
            — Drawn from a grandmother&apos;s story
          </figcaption>
          <p className="mt-3 text-sm leading-6 text-ink-700">
            This is what Hearthside does: it listens to stories like yours and keeps the
            wisdom inside them — automatically, in your language.
          </p>
        </figure>
      </section>

      <section className="border-y border-parchment-300 bg-parchment-100">
        <div className="mx-auto max-w-6xl px-6 py-14">
          <p className="text-xs font-semibold uppercase tracking-[0.25em] text-ember-600">
            How the weaving works
          </p>
          <ol className="mt-8 grid gap-px overflow-hidden rounded-2xl border border-parchment-300 bg-parchment-300 sm:grid-cols-2 lg:grid-cols-4">
            {journey.map((j) => (
              <li key={j.step} className="bg-parchment-50 p-7">
                <p className="font-mono text-xs tracking-widest text-gold-600">{j.step}</p>
                <h2 className="mt-3 font-display text-xl font-semibold">{j.title}</h2>
                <p className="mt-2 text-sm leading-6 text-ink-700">{j.body}</p>
              </li>
            ))}
          </ol>
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 py-16">
        <div className="flex flex-col items-start justify-between gap-6 sm:flex-row sm:items-center">
          <h2 className="max-w-md font-display text-3xl font-medium leading-tight">
            The ones who came before you are still speaking.{" "}
            <span className="italic text-ink-500">Are you listening?</span>
          </h2>
          <Link
            href="/register"
            className="shrink-0 rounded-full bg-ink-900 px-7 py-3 text-sm font-semibold text-parchment-50 transition hover:bg-ember-700"
          >
            Begin the telling
          </Link>
        </div>
      </section>

      <footer>
        <div className="kente-strip" aria-hidden="true" />
        <p className="bg-ink-950 py-6 text-center font-mono text-xs tracking-widest text-parchment-200">
          HEARTHSIDE — PERSONAL MYTHOLOGY &amp; ANCESTRAL WISDOM
        </p>
      </footer>
    </main>
  );
}
