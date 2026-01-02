import Link from 'next/link';

export default function Header() {
  return (
    <header className="border-b border-border bg-white">
      <div className="container mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-8">
            <Link href="/" className="text-2xl font-bold text-primary-600">
              Marvin
            </Link>
            <nav className="hidden md:flex space-x-6">
              <Link 
                href="/" 
                className="text-foreground hover:text-primary-600 transition-colors"
              >
                Dashboard
              </Link>
              <Link 
                href="/checks" 
                className="text-foreground hover:text-primary-600 transition-colors"
              >
                All Checks
              </Link>
            </nav>
          </div>
          <div className="flex items-center space-x-4">
            <span className="text-sm text-muted">Documentation QA</span>
          </div>
        </div>
      </div>
    </header>
  );
}
