"use client";

import { useEffect, useMemo, useRef, useState } from "react";
// import { NOTIFICATION_JOBS_URL } from "@/constants/notification";

export type SSEEvent = any;

export interface UseSSECacheOptions {
  eventTypes?: string[];
  maxCacheSize?: number;
}

export interface UseSSECacheResult<T = SSEEvent> {
  events: T[];
  lastEvent: T | null;
  isConnected: boolean;
  error: string | null;
  clear: () => void;
}

interface TokenStore<T = SSEEvent> {
  eventSource: EventSource | null;
  isConnected: boolean;
  error: string | null;
  events: T[];
  listeners: Set<() => void>;
  maxCacheSize: number;
}

const registry: Map<string, TokenStore> = new Map();

function ensureStore(token: string, maxCacheSize: number): TokenStore {
  let store = registry.get(token);

  if (!store) {
    store = {
      eventSource: null,
      isConnected: false,
      error: null,
      events: [],
      listeners: new Set(),
      maxCacheSize,
    };

    registry.set(token, store);
  } else {
    store.maxCacheSize = Math.max(store.maxCacheSize ?? 0, maxCacheSize);
  }

  return store;
}

function connectIfNeeded(token: string, store: TokenStore) {
  if (typeof window === "undefined") return;
  if (store.eventSource) return;
  if (!token) return;

  const url = `/api/sse/jobs?token=${encodeURIComponent(token)}`;
  const eventSource = new EventSource(url, { withCredentials: true });
  store.eventSource = eventSource;

  eventSource.onopen = () => {
    store.isConnected = true;
    store.error = null;
    notify(store);
  };

  eventSource.onmessage = (event: MessageEvent) => {
    try {
      const data = JSON.parse(event.data);
      pushEvent(store, data);
    } catch (e) {
      pushEvent(store, event.data as any);
    }
  };

  eventSource.onerror = () => {
    store.isConnected = false;
    store.error = "SSE connection error";
    notify(store);
  };
}

function pushEvent(store: TokenStore, data: SSEEvent) {
  store.events.push(data);

  if (store.events.length > (store.maxCacheSize || 100)) {
    store.events.splice(0, store.events.length - (store.maxCacheSize || 100));
  }

  notify(store);
}

function notify(store: TokenStore) {
  store.listeners.forEach((cb) => cb());
}

function teardownIfNoListeners(token: string, store: TokenStore) {
  if (store.listeners.size === 0) {
    if (store.eventSource) {
      store.eventSource.close();
      store.eventSource = null;
    }

    registry.delete(token);
  }
}

export function useSSECache<T = SSEEvent>(
  token: string | null | undefined,
  options: UseSSECacheOptions = {}
): UseSSECacheResult<T> {
  const { eventTypes, maxCacheSize = 100 } = options;

  const [tick, setTick] = useState(0);
  const tokenKey = token ?? "";
  const storeRef = useRef<TokenStore>(null);

  useEffect(() => {
    if (!tokenKey) return;
    const store = ensureStore(tokenKey, maxCacheSize);
    storeRef.current = store;
    connectIfNeeded(tokenKey, store);

    const listener = () => setTick((x) => x + 1);
    store.listeners.add(listener);

    listener();

    return () => {
      store.listeners.delete(listener);
      teardownIfNoListeners(tokenKey, store);
    };
  }, [tokenKey, maxCacheSize]);

  const derived = useMemo(() => {
    const store = storeRef.current;
    if (!store) {
      return {
        events: [] as T[],
        lastEvent: null as T | null,
        isConnected: false,
        error: null as string | null,
      };
    }

    let events = store.events as T[];
    if (eventTypes && eventTypes.length > 0) {
      events = events.filter((e: any) => eventTypes.includes(e?.type));
    }
    const lastEvent = events.length > 0 ? (events[events.length - 1] as T) : null;
    return {
      events,
      lastEvent,
      isConnected: store.isConnected,
      error: store.error,
    };
  }, [tick, eventTypes]);

  const clear = useMemo(() => {
    return () => {
      const store = storeRef.current;
      if (!store) return;
      store.events = [];
      notify(store);
    };
  }, []);

  return {
    ...derived,
    clear,
  } as UseSSECacheResult<T>;
}
