"use client";

import { useCallback, useEffect, useState } from "react";

import { api, ApiError, type Todo } from "@/lib/api";

export function TodoApp() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [title, setTitle] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const refresh = useCallback(async () => {
    setError(null);
    try {
      setTodos((await api.list()).todos);
    } catch (err) {
      setError(errMessage(err));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    const clean = title.trim();
    if (!clean) return;

    setSubmitting(true);
    try {
      const created = await api.create(clean);
      setTodos((prev) => [created, ...prev]);
      setTitle("");
    } catch (err) {
      setError(errMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function toggle(todo: Todo) {
    setTodos((prev) =>
      prev.map((t) => (t.id === todo.id ? { ...t, completed: !t.completed } : t)),
    );
    try {
      const updated = await api.update(todo.id, { completed: !todo.completed });
      setTodos((prev) => prev.map((t) => (t.id === todo.id ? updated : t)));
    } catch (err) {
      setTodos((prev) => prev.map((t) => (t.id === todo.id ? todo : t))); // rollback
      setError(errMessage(err));
    }
  }

  async function remove(todo: Todo) {
    setTodos((prev) => prev.filter((t) => t.id !== todo.id));
    try {
      await api.remove(todo.id);
    } catch (err) {
      setError(errMessage(err));
      void refresh(); // la fuente de verdad es el servidor
    }
  }

  const pending = todos.filter((t) => !t.completed).length;

  return (
    <section>
      <p className="eyebrow">Cuaderno de tareas</p>
      <h1 className="title">
        Tasks <span className="count">· {pending} pendientes</span>
      </h1>

      <form className="new-form" onSubmit={onSubmit}>
        <input
          className="new-input"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="¿Qué hay que hacer?"
          maxLength={200}
          aria-label="Nueva tarea"
        />
        <button className="btn" type="submit" disabled={!title.trim()} aria-busy={submitting}>
          {submitting ? "Guardando…" : "Añadir"}
        </button>
      </form>

      {error && (
        <p className="error" role="alert">
          {error}
        </p>
      )}

      {loading ? (
        <div className="loading" aria-label="Cargando tareas" />
      ) : todos.length === 0 ? (
        <p className="empty">Nada por aquí. Escribe la primera tarea.</p>
      ) : (
        <ul className="ledger">
          {todos.map((todo) => (
            <TodoRow
              key={todo.id}
              todo={todo}
              onToggle={() => void toggle(todo)}
              onRemove={() => void remove(todo)}
            />
          ))}
        </ul>
      )}
    </section>
  );
}

function TodoRow({
  todo,
  onToggle,
  onRemove,
}: {
  todo: Todo;
  onToggle: () => void;
  onRemove: () => void;
}) {
  return (
    <li className={`row ${todo.completed ? "done" : ""}`}>
      <input
        className="check"
        type="checkbox"
        checked={todo.completed}
        onChange={onToggle}
        aria-label={`Marcar: ${todo.title}`}
      />
      <span className="text">{todo.title}</span>
      <button
        className="delete"
        onClick={onRemove}
        aria-label={`Borrar: ${todo.title}`}
        type="button"
      >
        ×
      </button>
    </li>
  );
}

function errMessage(err: unknown): string {
  if (err instanceof ApiError) {
    if (err.code === "INVALID_JSON") return "El servidor no entendió la petición.";
    if (err.code === "NETWORK") return err.message;
    return err.message;
  }
  return "Algo falló. Inténtalo otra vez.";
}
