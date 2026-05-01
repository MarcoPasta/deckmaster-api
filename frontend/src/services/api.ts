const BASE_URL = import.meta.env.VITE_API_BASE_URL;
// || "http://localhost:8080/api";

type RequestOptions = {
  headers?: Record<string, string>;
  params?: Record<string, string>;
  body?: unknown;
  signal?: AbortSignal;
};

async function request<T>(
  method: string,
  path: string,
  token?: string | null,
  options: RequestOptions = {},
): Promise<T> {
  const url = new URL(path, BASE_URL);

  if (options.params) {
    for (const [key, value] of Object.entries(options.params)) {
      url.searchParams.set(key, value);
    }
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(url.toString(), {
    method,
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
    signal: options.signal,
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new ApiError(res.status, text || res.statusText);
  }

  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export class ApiError extends Error {
  constructor(
    public status: number,
    public body: string,
  ) {
    super(`API error ${status}: ${body}`);
    this.name = "ApiError";
  }
}

export const api = {
  get<T>(path: string, token?: string | null, options?: RequestOptions) {
    return request<T>("GET", path, token, options);
  },
  post<T>(path: string, token?: string | null, options?: RequestOptions) {
    return request<T>("POST", path, token, options);
  },
  put<T>(path: string, token?: string | null, options?: RequestOptions) {
    return request<T>("PUT", path, token, options);
  },
  patch<T>(path: string, token?: string | null, options?: RequestOptions) {
    return request<T>("PATCH", path, token, options);
  },
  delete<T>(path: string, token?: string | null, options?: RequestOptions) {
    return request<T>("DELETE", path, token, options);
  },
};
