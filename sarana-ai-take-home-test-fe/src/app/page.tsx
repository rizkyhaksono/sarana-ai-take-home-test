import Link from 'next/link'
import { Button } from '@/components/ui/button'

export default function Home() {
  return (
    <main className="flex flex-1 flex-col items-center justify-center p-4 text-center">
      <div className="space-y-4 max-w-2xl">
        <h2 className="text-base font-bold tracking-tight sm:text-3xl">
          Your notes, <br />
          <span className="text-primary">beautifully organized</span>
        </h2>
        <p className="text-sm text-muted-foreground">
          A modern note-taking application with dark mode support,
          image uploads, and a clean interface.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
          <Button asChild size="lg">
            <Link href="/login">Get Started</Link>
          </Button>
          <Button asChild variant="neutral" size="lg">
            <Link href="/register">Create Account</Link>
          </Button>
        </div>
      </div>
    </main>
  )
}
