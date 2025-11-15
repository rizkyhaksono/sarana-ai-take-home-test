import { ThemeToggle } from "./theme-toggle";

export default function Navbar() {
  return (
    <header className="border-b">
      <div className="container mx-auto flex items-center justify-between py-4 px-4">
        <h1 className="text-base font-bold">Sarana Notes App</h1>
        <ThemeToggle />
      </div>
    </header>
  )
}