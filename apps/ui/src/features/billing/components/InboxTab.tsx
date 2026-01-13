import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useInbox, useClaimAssignment } from "../hooks/useIA"
import { formatCurrency } from "../utils/formatting"
import { cn } from "@/lib/utils"

export function InboxTab() {
  const { data, isLoading, error } = useInbox()
  const claimMutation = useClaimAssignment()

  const handleClaim = (entityType: string, entityId: string) => {
    claimMutation.mutate({
      entity_type: entityType as "invoice" | "customer",
      entity_id: entityId,
      assignment_ttl_minutes: 120, // 2 hours
    })
  }

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex h-64 items-center justify-center text-destructive">
        Failed to load inbox
      </div>
    )
  }

  const items = data?.items || []

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          {items.length} {items.length === 1 ? "item" : "items"} need attention
        </p>
      </div>

      <div className="rounded-xl border border-border/60 bg-card shadow-sm overflow-hidden">
        <Table>
          <TableHeader className="bg-muted/30">
            <TableRow className="hover:bg-transparent">
              <TableHead className="font-medium">Issue Type</TableHead>
              <TableHead className="font-medium">Entity</TableHead>
              <TableHead className="font-medium">Customer</TableHead>
              <TableHead className="font-medium">Amount / Exposure</TableHead>
              <TableHead className="font-medium">Days Overdue</TableHead>
              <TableHead className="text-right font-medium">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center text-muted-foreground">
                  No items in inbox. All systems nominal! 🎉
                </TableCell>
              </TableRow>
            ) : (
              items.map((item) => {
                const isLoading = claimMutation.isPending
                const isFailedPayment = item.category === "failed_payment"

                return (
                  <TableRow key={`${item.entity_type}-${item.entity_id}`} className="group">
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className={cn(
                          "border",
                          isFailedPayment
                            ? "bg-red-500/10 text-red-600 border-red-500/20"
                            : "bg-amber-500/10 text-amber-600 border-amber-500/20"
                        )}
                      >
                        {item.category.replace(/_/g, " ")}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-xs">
                      {item.entity_name}
                    </TableCell>
                    <TableCell className="font-medium">
                      {item.customer_name || "N/A"}
                    </TableCell>
                    <TableCell className="font-mono text-xs tabular-nums">
                      {item.amount_due !== undefined
                        ? formatCurrency(item.amount_due, item.currency || data?.currency)
                        : "N/A"}
                    </TableCell>
                    <TableCell>
                      {item.days_overdue !== undefined ? (
                        <span
                          className={cn(
                            "font-medium",
                            item.days_overdue > 30
                              ? "text-red-600"
                              : item.days_overdue > 14
                                ? "text-amber-600"
                                : "text-muted-foreground"
                          )}
                        >
                          {item.days_overdue} days
                        </span>
                      ) : (
                        "N/A"
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="sm"
                        variant="default"
                        className="h-7 text-xs opacity-0 transition-opacity group-hover:opacity-100"
                        disabled={isLoading}
                        onClick={() => handleClaim(item.entity_type, item.entity_id)}
                      >
                        {isLoading ? (
                          <Loader2 className="h-3 w-3 animate-spin" />
                        ) : (
                          "Claim"
                        )}
                      </Button>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
