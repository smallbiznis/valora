import { useEffect, useMemo, useState } from "react"
import { NavLink, Outlet, useNavigate, useParams } from "react-router-dom"

import {api} from "@/api/client"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Separator } from "@/components/ui/separator"
import { Skeleton } from "@/components/ui/skeleton"
import { useAuthStore } from "@/stores/authStore"
import { useOrgStore } from "@/stores/orgStore"
import { cn } from "@/lib/utils"

const navItems = [
  { label: "Dashboard", path: "dashboard" },
  { label: "Products", path: "products" },
  { label: "Meters", path: "meter" },
  { label: "Customers", path: "customers" },
  { label: "Subscriptions", path: "subscriptions" },
  { label: "Invoices", path: "invoices" },
  { label: "Settings", path: "settings" },
]

export function AppShell() {
  const { orgId } = useParams<{ orgId: string }>()
  const navigate = useNavigate()
  const { currentOrg, setCurrentOrg, setOrgs } = useOrgStore()
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    let isMounted = true
    if (!orgId) return
    setIsLoading(true)
    api
      .get(`/orgs/${orgId}`)
      .then((res) => {
        if (!isMounted) return
        const org = res.data?.org ?? { id: orgId, name: `Org ${orgId}` }
        setCurrentOrg(org)
        setIsLoading(false)
      })
      .catch(() => {
        if (!isMounted) return
        setIsLoading(false)
        navigate("/orgs", { replace: true })
      })

    api
      .get("/me/orgs")
      .then((res) => {
        if (!isMounted) return
        setOrgs(res.data?.orgs ?? [])
      })
      .catch(() => {
        /* ignore */
      })

    return () => {
      isMounted = false
    }
  }, [navigate, orgId, setCurrentOrg, setOrgs])

  const sidebar = useMemo(() => {
    if (!orgId) return null
    return navItems.map((item) => ({
      ...item,
      to: `/orgs/${orgId}/${item.path}`,
    }))
  }, [orgId])

  if (isLoading || !currentOrg) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-muted/40">
        <div className="space-y-3">
          <Skeleton className="h-10 w-72" />
          <Skeleton className="h-4 w-52" />
        </div>
      </div>
    )
  }

  return (
    <div className="grid min-h-screen grid-cols-[260px_1fr] bg-muted/20">
      <aside className="border-r bg-white">
        <div className="px-4 py-5">
          <div className="text-lg font-semibold">Valora</div>
          <div className="text-sm text-muted-foreground">{currentOrg.name}</div>
        </div>
        <Separator />
        <nav className="flex flex-col gap-1 p-3">
          {sidebar?.map((item) => (
            <NavLink
              key={item.path}
              to={item.to}
              className={({ isActive }) =>
                cn(
                  "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-muted",
                  isActive && "bg-muted text-primary"
                )
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <div className="flex flex-col">
        <Topbar />
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

function Topbar() {
  const { orgs, currentOrg, setCurrentOrg } = useOrgStore()
  const navigate = useNavigate()
  const logout = useAuthStore((s) => s.logout)

  const initial = currentOrg?.name?.[0]?.toUpperCase() ?? "O"

  return (
    <header className="flex items-center justify-between border-b bg-white px-6 py-4">
      <div className="flex items-center gap-3">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" className="gap-2">
              <span className="font-medium">{currentOrg?.name}</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            <DropdownMenuLabel>Workspaces</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {orgs.map((org) => (
              <DropdownMenuItem
                key={org.id}
                onSelect={() => {
                  setCurrentOrg(org)
                  navigate(`/orgs/${org.id}/dashboard`)
                }}
              >
                {org.name}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem onSelect={() => navigate("/onboarding")}>
              + New workspace
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Avatar className="cursor-pointer">
            <AvatarFallback>{initial}</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>Account</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem onSelect={() => navigate("/orgs")}>
            Account settings
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onSelect={async () => {
              await logout()
              navigate("/login", { replace: true })
            }}
          >
            Logout
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
  )
}
