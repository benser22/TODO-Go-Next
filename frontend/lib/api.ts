// Cliente del API: una sola capa tipada que el resto de la app consume.

const BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly code: string,
  ) {
    super(message);
  }
}

export interface Todo {
  id: string;
  title: string;
  completed: boolean;
  createdAt: string;
  updatedAt: string;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  let res: Response;
  try {
    res = await fetch(BASE + path, {
      ...init,
      headers: { "Content-Type": "application/json", ...init?.headers },
    });
  } catch {
    // Servidor caído o CORS: mensaje propio, nunca el texto crudo.
    throw new ApiError("No se pudo conectar con el servidor", 0, "NETWORK");
  }

  if (res.status === 204) return undefined as T;

  const body = res.headers.get("content-type")?.includes("json")
    ? await res.json()
    : null;

  if (!res.ok) {
    throw new ApiError(
      body?.error?.message ?? "Respuesta inesperada del servidor",
      res.status,
      body?.error?.code ?? "UNKNOWN",
    );
  }
  return body as T;
}

export const api = {
  list: () => request<{ todos: Todo[] }>("/api/v1/todos"),

  create: (title: string) =>
    request<Todo>("/api/v1/todos", {
      method: "POST",
      body: JSON.stringify({ title }),
    }),

  update: (id: string, patch: { title?: string; completed?: boolean }) =>
    request<Todo>(`/api/v1/todos/${id}`, {
      method: "PATCH",
      body: JSON.stringify(patch),
    }),

  remove: (id: string) => request<void>(`/api/v1/todos/${id}`, { method: "DELETE" }),
};
