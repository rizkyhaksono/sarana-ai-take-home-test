import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { ThemeToggle } from '@/components/theme-toggle'

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="border-b">
        <div className="container mx-auto flex items-center justify-between py-4 px-4">
          <h1 className="text-2xl font-bold">Notes App</h1>
          <ThemeToggle />
        </div>
      </header>

      <main className="flex flex-1 flex-col items-center justify-center p-4 text-center">
        <div className="space-y-6 max-w-2xl">
          <h2 className="text-4xl font-bold tracking-tight sm:text-6xl">
            Your notes, <br />
            <span className="text-primary">beautifully organized</span>
          </h2>
          <p className="text-lg text-muted-foreground">
            A modern note-taking application with dark mode support,
            image uploads, and a clean interface.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
            <Button asChild size="lg">
              <Link href="/login">Get Started</Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link href="/register">Create Account</Link>
            </Button>
          </div>
        </div>
      </main>

      <footer className="border-t py-6 text-center text-sm text-muted-foreground">
        <p>Built with Next.js, Go, and PostgreSQL</p>
      </footer>
    </div>
  )
}
