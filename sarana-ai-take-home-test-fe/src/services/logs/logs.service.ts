import { Log, LogsResponse, PaginationParams } from "@/types";
import { useQuery } from "@tanstack/react-query";
import { BASE_API, baseHeaders } from "../api.config";

export const logsKeys = {
  all: ['logs'] as const,
  lists: () => [...logsKeys.all, 'list'] as const,
  list: (params?: PaginationParams) => [...logsKeys.lists(), params] as const,
  details: () => [...logsKeys.all, 'detail'] as const,
  detail: (id: number) => [...logsKeys.details(), id] as const,
};

export const useGetLogs = (params?: PaginationParams) => {
  return useQuery({
    queryKey: logsKeys.list(params),
    queryFn: async () => {
      const headers = await baseHeaders();
      const queryParams = new URLSearchParams();

      if (params?.search) queryParams.append('search', params.search);
      if (params?.sort_by) queryParams.append('sort_by', params.sort_by);
      if (params?.order) queryParams.append('order', params.order);
      if (params?.page) queryParams.append('page', params.page.toString());
      if (params?.limit) queryParams.append('limit', params.limit.toString());

      const query = queryParams.toString();
      const endpoint = query ? `/logs?${query}` : '/logs';

      const res = await fetch(`${BASE_API}${endpoint}`, {
        method: "GET",
        headers,
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to fetch logs");
      }

      return await res.json() as LogsResponse;
    },
    staleTime: 30000, // 30 seconds
  });
};

// Get single log
export const useGetLog = (id: number, enabled = true) => {
  return useQuery({
    queryKey: logsKeys.detail(id),
    queryFn: async () => {
      const headers = await baseHeaders();
      const res = await fetch(`${BASE_API}/logs/${id}`, {
        method: "GET",
        headers,
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to fetch log");
      }

      const data = await res.json();
      return data.data as Log;
    },
    enabled: !!id && enabled,
    staleTime: 30000,
  });
};
