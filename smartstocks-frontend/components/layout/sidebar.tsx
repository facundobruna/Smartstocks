'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  Home,
  BookOpen,
  TrendingUp,
  Swords,
  Trophy,
  Gamepad2,
  User,
  MessageSquare,
  Settings,
  ChevronLeft,
  Menu,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet'
import { cn } from '@/lib/utils'
import { useState } from 'react'

const sidebarLinks = [
  { href: '/home', label: 'Inicio', icon: Home },
  { href: '/education', label: 'Educacion', icon: BookOpen },
  { href: '/simulator', label: 'Simulador', icon: TrendingUp },
  { href: '/pvp', label: 'PvP', icon: Swords },
  { href: '/rankings', label: 'Rankings', icon: Trophy },
  { href: '/tournaments', label: 'Torneos', icon: Gamepad2 },
  { href: '/profile', label: 'Perfil', icon: User },
  { href: '/forum', label: 'Foro', icon: MessageSquare },
]

interface SidebarProps {
  collapsed?: boolean
  onToggle?: () => void
}

function SidebarContent({ collapsed = false }: { collapsed?: boolean }) {
  const pathname = usePathname()

  return (
    <div className="flex flex-col h-full">
      {/* Logo */}
      <div className="h-16 flex items-center px-4 border-b">
        <Link href="/home" className="flex items-center gap-2">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center shrink-0">
            <span className="text-white font-bold text-lg">S</span>
          </div>
          {!collapsed && (
            <span className="font-bold text-lg text-foreground">SMARTSTOCKS</span>
          )}
        </Link>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        {sidebarLinks.map((link) => {
          const Icon = link.icon
          const isActive = pathname === link.href || pathname.startsWith(`${link.href}/`)

          return (
            <Link
              key={link.href}
              href={link.href}
              className={cn(
                'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                isActive
                  ? 'bg-primary/10 text-primary'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
              )}
            >
              <Icon className={cn('w-5 h-5 shrink-0', isActive && 'text-primary')} />
              {!collapsed && <span>{link.label}</span>}
            </Link>
          )
        })}
      </nav>

      {/* Settings at bottom */}
      <div className="px-3 py-4 border-t">
        <Link
          href="/settings"
          className={cn(
            'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
            pathname === '/settings'
              ? 'bg-primary/10 text-primary'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          )}
        >
          <Settings className="w-5 h-5 shrink-0" />
          {!collapsed && <span>Configuracion</span>}
        </Link>
      </div>
    </div>
  )
}

export function Sidebar({ collapsed = false, onToggle }: SidebarProps) {
  return (
    <>
      {/* Desktop Sidebar */}
      <aside
        className={cn(
          'hidden lg:flex flex-col h-screen bg-white border-r transition-all duration-300',
          collapsed ? 'w-[70px]' : 'w-[250px]'
        )}
      >
        <SidebarContent collapsed={collapsed} />

        {/* Collapse Toggle */}
        {onToggle && (
          <div className="absolute bottom-20 -right-3">
            <Button
              variant="outline"
              size="icon"
              className="h-6 w-6 rounded-full bg-white shadow-md"
              onClick={onToggle}
            >
              <ChevronLeft
                className={cn(
                  'h-4 w-4 transition-transform',
                  collapsed && 'rotate-180'
                )}
              />
            </Button>
          </div>
        )}
      </aside>
    </>
  )
}

export function MobileSidebar() {
  const [open, setOpen] = useState(false)

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild className="lg:hidden">
        <Button variant="ghost" size="icon" className="fixed bottom-4 left-4 z-50 bg-primary text-white shadow-lg hover:bg-primary/90">
          <Menu className="h-5 w-5" />
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="p-0 w-[250px]">
        <SidebarContent />
      </SheetContent>
    </Sheet>
  )
}

export default Sidebar
