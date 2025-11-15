import Link from "next/link"
import { Github, Mail } from "lucide-react"
import { Card } from "../ui/card"

export default function Footer() {
  const currentYear = new Date().getFullYear()

  return (
    <Card>
      <footer className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
          {/* About */}
          <div>
            <h3 className="font-heading text-lg mb-3">Sarana Notes</h3>
            <p className="text-sm text-muted-foreground">
              A modern note-taking application built with Next.js and Go Fiber.
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="font-heading text-lg mb-3">Quick Links</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link href="/dashboard" className="text-muted-foreground hover:text-foreground transition-colors">
                  Dashboard
                </Link>
              </li>
              <li>
                <Link href="/logs" className="text-muted-foreground hover:text-foreground transition-colors">
                  Logs
                </Link>
              </li>
            </ul>
          </div>

          {/* Connect */}
          <div>
            <h3 className="font-heading text-lg mb-3">Connect</h3>
            <div className="flex items-center gap-4">
              <Link
                href="https://github.com/rizkyhaksono"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <Github className="h-5 w-5" />
                <span className="sr-only">GitHub</span>
              </Link>
              <Link
                href="mailto:rizkyhaksono@example.com"
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <Mail className="h-5 w-5" />
                <span className="sr-only">Email</span>
              </Link>
            </div>
          </div>
        </div>

        {/* Copyright */}
        <div className="pt-6 border-t border-border text-center text-sm text-muted-foreground">
          <p>© {currentYear} Built by Rizky Haksono. All rights reserved.</p>
        </div>
      </footer>
    </Card>
  )
}