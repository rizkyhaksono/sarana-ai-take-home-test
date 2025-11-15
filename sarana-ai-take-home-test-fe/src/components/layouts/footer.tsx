import Link from "next/link"

export default function Footer() {
  return (
    <footer className="border-t py-6 text-center text-sm text-muted-foreground">
      <Link
        href="https://github.com/rizkyhaksono"
        className="hover:underline hover:underline-offset-4"
        target="_blank"
      >
        Built by Rizky Haksono
      </Link>
    </footer>
  )
}