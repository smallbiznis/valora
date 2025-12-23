import { useEffect, useMemo, useState } from "react"
import { Link, useParams } from "react-router-dom"

import { api } from "@/api/client"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

type Price = {
  id: string
  product_id: string
  name?: string
  code?: string
  pricing_model: string
  billing_interval: string
  billing_interval_count: number
}

type PriceAmount = {
  id: string
  price_id: string
  meter_id?: string | null
  currency: string
  unit_amount_cents: number
  minimum_amount_cents?: number | null
  maximum_amount_cents?: number | null
}

type Meter = {
  id: string
  name?: string
  code?: string
}

const formatInterval = (interval: string, count: number) => {
  const base = interval.toLowerCase()
  if (!count || count === 1) {
    return `Every ${base}`
  }
  return `Every ${count} ${base}${count === 1 ? "" : "s"}`
}

const formatPricingModel = (model: string) =>
  model === "FLAT" ? "Flat" : "Usage-based"

const formatAmount = (amountCents: number, currency: string) => {
  const safeCurrency = currency?.toUpperCase() || "USD"
  try {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: safeCurrency,
    }).format(amountCents / 100)
  } catch {
    return `${amountCents} ${safeCurrency}`
  }
}

export default function OrgProductDetailPage() {
  const { orgId, productId } = useParams()
  const [product, setProduct] = useState<unknown | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [prices, setPrices] = useState<Price[]>([])
  const [priceAmounts, setPriceAmounts] = useState<PriceAmount[]>([])
  const [meters, setMeters] = useState<Meter[]>([])
  const [pricingLoading, setPricingLoading] = useState(true)
  const [pricingError, setPricingError] = useState<string | null>(null)

  useEffect(() => {
    if (!orgId || !productId) {
      setIsLoading(false)
      return
    }

    let isMounted = true
    setIsLoading(true)
    setError(null)

    api
      .get(`/products/${productId}`, { params: { organization_id: orgId } })
      .then((response) => {
        if (!isMounted) return
        setProduct(response.data?.data ?? null)
      })
      .catch((err) => {
        if (!isMounted) return
        setError(err?.message ?? "Unable to load product.")
      })
      .finally(() => {
        if (!isMounted) return
        setIsLoading(false)
      })

    return () => {
      isMounted = false
    }
  }, [orgId, productId])

  useEffect(() => {
    if (!orgId || !productId) {
      setPricingLoading(false)
      return
    }

    let isMounted = true
    setPricingLoading(true)
    setPricingError(null)

    Promise.all([
      api.get("/prices", { params: { organization_id: orgId } }),
      api.get("/price_amounts", { params: { organization_id: orgId } }),
      api.get("/meters", { params: { organization_id: orgId } }),
    ])
      .then(([priceRes, amountRes, meterRes]) => {
        if (!isMounted) return
        const allPrices = priceRes.data?.data ?? []
        const filteredPrices = allPrices.filter(
          (item: Price) => String(item.product_id) === String(productId)
        )
        setPrices(filteredPrices)
        setPriceAmounts(amountRes.data?.data ?? [])
        setMeters(meterRes.data?.data ?? [])
      })
      .catch((err) => {
        if (!isMounted) return
        setPricingError(err?.message ?? "Unable to load pricing.")
      })
      .finally(() => {
        if (!isMounted) return
        setPricingLoading(false)
      })

    return () => {
      isMounted = false
    }
  }, [orgId, productId])

  const metersById = useMemo(() => {
    const map = new Map<string, Meter>()
    meters.forEach((meter) => {
      map.set(String(meter.id), meter)
    })
    return map
  }, [meters])

  const productRecord =
    product && typeof product === "object" ? (product as Record<string, unknown>) : null
  const productName =
    productRecord && typeof productRecord.name === "string"
      ? productRecord.name
      : "Product detail"
  const productCode =
    productRecord && typeof productRecord.code === "string"
      ? productRecord.code
      : "-"
  const productStatus =
    productRecord && typeof productRecord.active === "boolean"
      ? productRecord.active
        ? "Active"
        : "Archived"
      : "Unknown"

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">{productName}</h1>
          <p className="text-muted-foreground text-sm">
            Configure pricing, usage, and advanced settings for this product.
          </p>
        </div>
        {orgId && (
          <Button asChild variant="outline">
            <Link to={`/orgs/${orgId}/products`}>Back to products</Link>
          </Button>
        )}
      </div>
      {isLoading && (
        <div className="text-muted-foreground text-sm">Loading product...</div>
      )}
      {error && <div className="text-destructive text-sm">{error}</div>}
      {!isLoading && !error && (
        <Tabs defaultValue="overview" className="space-y-4">
          <TabsList className="flex w-full flex-wrap justify-start">
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="pricing">Pricing</TabsTrigger>
            <TabsTrigger value="usage">Usage & Meters</TabsTrigger>
            <TabsTrigger value="advanced">Advanced</TabsTrigger>
          </TabsList>

          {/* Pricing stays in product detail so the sidebar remains intent-driven, not table-driven. */}
          <TabsContent value="overview" className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle>Summary</CardTitle>
                  <CardDescription>Core details for this product.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-2 text-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-muted-foreground">Code</span>
                    <span>{productCode}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-muted-foreground">Status</span>
                    <span>{productStatus}</span>
                  </div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Customer view</CardTitle>
                  <CardDescription>What customers see when selecting this plan.</CardDescription>
                </CardHeader>
                <CardContent className="text-sm text-muted-foreground">
                  Add a name, description, and highlights to make the plan easy to understand.
                </CardContent>
              </Card>
            </div>
            <Card>
              <CardHeader>
                <CardTitle>Raw product data</CardTitle>
                <CardDescription>Reference payload for support and debugging.</CardDescription>
              </CardHeader>
              <CardContent>
                <pre className="bg-muted overflow-auto rounded-md p-4 text-xs">
                  {JSON.stringify(product, null, 2)}
                </pre>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="pricing" className="space-y-4">
            {pricingLoading && (
              <div className="text-muted-foreground text-sm">Loading pricing...</div>
            )}
            {pricingError && <div className="text-destructive text-sm">{pricingError}</div>}
            {!pricingLoading && !pricingError && prices.length === 0 && (
              <div className="rounded-lg border border-dashed p-6 text-center text-muted-foreground text-sm">
                No prices configured yet.
              </div>
            )}
            {!pricingLoading && !pricingError && prices.length > 0 && (
              <div className="space-y-4">
                {prices.map((price) => {
                  const amounts = priceAmounts.filter(
                    (amount) => String(amount.price_id) === String(price.id)
                  )
                  const currencyList = Array.from(
                    new Set(amounts.map((amount) => amount.currency))
                  )
                  const currencySummary =
                    currencyList.length > 0 ? currencyList.join(", ") : ""
                  const intervalLabel = formatInterval(
                    price.billing_interval,
                    price.billing_interval_count
                  )
                  const modelLabel = formatPricingModel(price.pricing_model)

                  return (
                    <Card key={price.id}>
                      <CardHeader>
                        <CardTitle>{price.name || price.code || "Price"}</CardTitle>
                        <CardDescription>
                          {modelLabel} | {intervalLabel}
                          {currencySummary ? ` | ${currencySummary}` : ""}
                        </CardDescription>
                      </CardHeader>
                      <CardContent className="space-y-3">
                        <div className="text-sm text-muted-foreground">Rates</div>
                        {amounts.length === 0 && (
                          <div className="text-muted-foreground text-sm">
                            No rates yet.
                          </div>
                        )}
                        {amounts.length > 0 && (
                          <div className="space-y-2">
                            {amounts.map((amount) => {
                              const meter = amount.meter_id
                                ? metersById.get(String(amount.meter_id))
                                : null
                              const meterLabel = meter
                                ? meter.name || meter.code || meter.id
                                : null
                              const minLabel =
                                amount.minimum_amount_cents != null
                                  ? ` | Min ${formatAmount(
                                      amount.minimum_amount_cents,
                                      amount.currency
                                    )}`
                                  : ""
                              const maxLabel =
                                amount.maximum_amount_cents != null
                                  ? ` | Max ${formatAmount(
                                      amount.maximum_amount_cents,
                                      amount.currency
                                    )}`
                                  : ""

                              return (
                                <div
                                  key={amount.id}
                                  className="rounded-md border px-4 py-3 text-sm"
                                >
                                  <div className="font-medium">
                                    {formatAmount(amount.unit_amount_cents, amount.currency)}
                                  </div>
                                  <div className="text-muted-foreground text-xs">
                                    {meterLabel
                                      ? `Meter: ${meterLabel}`
                                      : "Flat rate"}
                                    {minLabel}
                                    {maxLabel}
                                  </div>
                                </div>
                              )
                            })}
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  )
                })}
              </div>
            )}
          </TabsContent>

          <TabsContent value="usage" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Usage meters</CardTitle>
                <CardDescription>Attach meters that drive usage-based billing.</CardDescription>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                {orgId ? (
                  <Link className="text-primary hover:underline" to={`/orgs/${orgId}/meter`}>
                    Manage meters
                  </Link>
                ) : (
                  "Select meters to track usage for this product."
                )}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Usage settings</CardTitle>
                <CardDescription>Align reporting windows and aggregation.</CardDescription>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Configure how usage is summarized and displayed on invoices.
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="advanced" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Price tiers</CardTitle>
                <CardDescription>Define tiered pricing behavior.</CardDescription>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Set volume tiers, overage behavior, and pricing curves.
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Limits</CardTitle>
                <CardDescription>Guardrails for usage and billing.</CardDescription>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Add soft or hard limits to protect customers and systems.
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Internal flags</CardTitle>
                <CardDescription>Operational settings for admins.</CardDescription>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                Toggle internal states without exposing them to customers.
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      )}
    </div>
  )
}
