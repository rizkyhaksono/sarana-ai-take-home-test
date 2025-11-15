'use client'

import { useState } from 'react'
import { useGetLogs, useGetLogById } from '@/services/logs/logs.service'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ChevronLeft, ChevronRight, Search, Loader2 } from 'lucide-react'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

export default function LogsPage() {
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [searchInput, setSearchInput] = useState('')
  const [selectedLogId, setSelectedLogId] = useState<number | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)

  const { data, isLoading, error } = useGetLogs({
    page,
    limit: 5,
    search,
    sort_by: 'datetime',
    order: 'DESC',
  })

  const { data: logDetail, isLoading: isLoadingDetail } = useGetLogById(
    selectedLogId || 0,
    dialogOpen && selectedLogId !== null
  )

  const handleSearch = () => {
    setSearch(searchInput)
    setPage(1)
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch()
    }
  }

  const handleOpenDialog = (logId: number) => {
    setSelectedLogId(logId)
    setDialogOpen(true)
  }

  const handleCloseDialog = () => {
    setDialogOpen(false)
    setSelectedLogId(null)
  }

  const getStatusBadgeVariant = (statusCode: number) => {
    if (statusCode >= 200 && statusCode < 300) return 'default'
    if (statusCode >= 400 && statusCode < 500) return 'destructive'
    if (statusCode >= 500) return 'destructive'
    return 'default'
  }

  const getMethodBadgeColor = (method: string) => {
    const colors: Record<string, string> = {
      'GET': 'bg-blue-500',
      'POST': 'bg-green-500',
      'PUT': 'bg-yellow-500',
      'DELETE': 'bg-red-500',
      'PATCH': 'bg-purple-500',
    }
    return colors[method] || 'bg-gray-500'
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const formatJSON = (jsonString: string) => {
    try {
      return JSON.stringify(JSON.parse(jsonString), null, 2)
    } catch {
      return jsonString
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-heading mb-2">Application Logs</h1>
        <p className="text-muted-foreground">
          View and search through all application request logs
        </p>
      </div>

      {/* Search Bar */}
      <Card className="mb-6">
        <CardContent className="pt-6">
          <div className="flex gap-2">
            <Input
              placeholder="Search logs..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyPress={handleKeyPress}
              className="flex-1"
            />
            <Button onClick={handleSearch}>
              <Search className="h-4 w-4 mr-2" />
              Search
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Logs List */}
      {isLoading && (
        <div className="text-center py-12">
          <p className="text-muted-foreground">Loading logs...</p>
        </div>
      )}

      {error && (
        <Card>
          <CardContent className="pt-6">
            <p className="text-destructive">Error loading logs: {error.message}</p>
          </CardContent>
        </Card>
      )}

      {data && (
        <>
          <div className="space-y-4 mb-6">
            {data.data.logs.length === 0 ? (
              <Card>
                <CardContent className="pt-6 text-center">
                  <p className="text-muted-foreground">No logs found</p>
                </CardContent>
              </Card>
            ) : (
              data.data.logs.map((log) => (
                <Card
                  key={log.id}
                  className="cursor-pointer hover:shadow-md transition-shadow"
                  onClick={() => handleOpenDialog(log.id)}
                >
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <Badge className={`${getMethodBadgeColor(log.method)} text-white`}>
                          {log.method}
                        </Badge>
                        <span className="font-mono text-sm">{log.endpoint}</span>
                      </div>
                      <Badge variant={getStatusBadgeVariant(log.status_code)}>
                        {log.status_code}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex justify-between text-sm text-muted-foreground">
                      <span>Request at: {formatDate(log.datetime)}</span>
                      <span>ID: {log.id}</span>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>

          {/* Pagination */}
          {data.data.total_pages > 1 && (
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <p className="text-sm text-muted-foreground">
                    Page {page} of {data.data.total_pages} ({data.data.total} total logs)
                  </p>
                  <div className="flex gap-2">
                    <Button
                      variant="neutral"
                      size="sm"
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page === 1}
                    >
                      <ChevronLeft className="h-4 w-4" />
                      Previous
                    </Button>
                    <Button
                      variant="neutral"
                      size="sm"
                      onClick={() => setPage((p) => Math.min(data.data.total_pages, p + 1))}
                      disabled={page === data.data.total_pages}
                    >
                      Next
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </>
      )}

      {/* Log Detail Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {isLoadingDetail ? (
                "Loading Log Details..."
              ) : logDetail ? (
                <div className="flex items-center gap-3">
                  <Badge className={`${getMethodBadgeColor(logDetail.method)} text-white`}>
                    {logDetail.method}
                  </Badge>
                  <span className="font-mono">{logDetail.endpoint}</span>
                  <Badge variant={getStatusBadgeVariant(logDetail.status_code)}>
                    {logDetail.status_code}
                  </Badge>
                </div>
              ) : (
                "Log Details"
              )}
            </DialogTitle>
            <DialogDescription asChild>
              <div>
                {isLoadingDetail ? (
                  <div className="flex items-center justify-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                  </div>
                ) : logDetail ? (
                  <div className="space-y-4 mt-4">
                    <div>
                      <h4 className="font-semibold mb-2 text-foreground">Datetime</h4>
                      <p className="text-sm">{formatDate(logDetail.datetime)}</p>
                    </div>

                    <div>
                      <h4 className="font-semibold mb-2 text-foreground">Log ID</h4>
                      <p className="text-sm">{logDetail.id}</p>
                    </div>

                    {logDetail.headers && (
                      <div>
                        <h4 className="font-semibold mb-2 text-foreground">Headers</h4>
                        <pre className="text-xs bg-muted p-3 rounded-md overflow-x-auto">
                          {formatJSON(logDetail.headers)}
                        </pre>
                      </div>
                    )}

                    {logDetail.request_body && (
                      <div>
                        <h4 className="font-semibold mb-2 text-foreground">Request Body</h4>
                        <pre className="text-xs bg-muted p-3 rounded-md overflow-x-auto">
                          {formatJSON(logDetail.request_body)}
                        </pre>
                      </div>
                    )}

                    {logDetail.response_body && (
                      <div>
                        <h4 className="font-semibold mb-2 text-foreground">Response Body</h4>
                        <pre className="text-xs bg-muted p-3 rounded-md overflow-x-auto">
                          {formatJSON(logDetail.response_body)}
                        </pre>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <p className="text-muted-foreground">Failed to load log details</p>
                  </div>
                )}
              </div>
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="neutral" onClick={handleCloseDialog}>Close</Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}