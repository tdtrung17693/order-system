export interface PaginationQuery {
  pageIndex: number
  itemsPerPage: number
  filters?: Record<string, any>
}

export interface PaginationResponse<T> {
  items: T[]
  pageIndex: number
  itemsPerPage: number
  total: number
}

export function buildPaginationRequest(
  pageIndex: number,
  itemsPerPage: number,
  filters: Record<string, any> = {}
): PaginationQuery {
  return { pageIndex, itemsPerPage, filters }
}
