function App() {
  return (
    <main className="flex min-h-screen items-center justify-center px-6 py-16 text-slate-100">
      <section className="w-full max-w-3xl space-y-6 rounded-3xl border border-slate-800/80 bg-slate-900/70 p-10 shadow-2xl shadow-slate-950/40 backdrop-blur">
        <span className="inline-flex items-center gap-2 rounded-full border border-slate-700/80 bg-slate-800/60 px-4 py-1 text-xs font-medium uppercase tracking-[0.25em] text-slate-400">
          0xHub
        </span>
        <header className="space-y-4">
          <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">
            Hub Frontend Workspace
          </h1>
          <p className="text-base leading-relaxed text-slate-300">
            This Vite + React + TypeScript shell is ready for modular Hub
            features. Tailwind CSS, ESLint, and Prettier are preconfigured so we
            can focus on delivering a cohesive UI experience.
          </p>
        </header>
        <div className="grid gap-4 sm:grid-cols-2">
          <a
            className="rounded-xl border border-slate-700/60 bg-slate-800/70 px-5 py-4 transition hover:border-emerald-400/60 hover:bg-emerald-500/10 hover:text-emerald-200"
            href="https://vite.dev/guide/"
            target="_blank"
            rel="noreferrer"
          >
            <h2 className="text-lg font-medium">Vite Dev Server</h2>
            <p className="mt-1 text-sm text-slate-400">
              Fast HMR and first-class TypeScript support.
            </p>
          </a>
          <a
            className="rounded-xl border border-slate-700/60 bg-slate-800/70 px-5 py-4 transition hover:border-sky-400/60 hover:bg-sky-500/10 hover:text-sky-200"
            href="https://tailwindcss.com/docs/installation"
            target="_blank"
            rel="noreferrer"
          >
            <h2 className="text-lg font-medium">Tailwind Utility System</h2>
            <p className="mt-1 text-sm text-slate-400">
              Rapidly compose adaptive layouts with design tokens.
            </p>
          </a>
        </div>
        <footer className="text-sm text-slate-500">
          Update this shell as future modules land in the Hub UI.
        </footer>
      </section>
    </main>
  )
}

export default App
