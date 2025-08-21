'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { verifyToken } from '@/services/auth';
import { useNotification } from '@/components/notification';
import { useSSECache } from '@/hooks/useSSECache';
import { useOverlay } from '@/hooks/useOverlay';

type AppEvent = {
  type?: string;
  [key: string]: unknown;
};

export default function HomeContent() {
  const { showNotification } = useNotification();

  const [token, setToken] = useState<string | null>(null);
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const t = window.localStorage.getItem('auth_token');
      if (t) setToken(t);
    }
  }, []);

  // Hook up SSE with a small cache if token exists
  const { events, lastEvent, isConnected, error, clear } = useSSECache<AppEvent>(token || undefined, {
    maxCacheSize: 50,
  });

  useEffect(() => {
    const checkAuth = async () => {
      const { error } = await verifyToken();

      if (error) showNotification(error, "error", 5000);
    };

    checkAuth().then();
  }, [showNotification]);

  const connectionBadge = useMemo(() => {
    if (!token) return <span className="text-muted">No token</span>;
    if (error) return <span className="text-primary">Error</span>;
    return (
      <span className={isConnected ? 'text-primary' : 'text-muted'}>
        {isConnected ? 'Connected' : 'Connecting...'}
      </span>
    );
  }, [token, isConnected, error]);

  const { isLeftOpen, isRightOpen, toggleLeft, toggleRight, leftWidth, rightWidth, setLeftWidth, setRightWidth } = useOverlay();

  return (
    <div className="container mx-auto p-4 space-y-6">
      {/* Overlay controls (demo) */}
      <div className="border border-border rounded p-4 bg-card/60">
        <div className="flex flex-wrap items-center gap-3">
          <button type="button" className="px-3 py-1 rounded bg-foreground/5 hover:bg-foreground/10" onClick={() => { toggleLeft('page:demo-left'); }}>
            {isLeftOpen ? 'Hide' : 'Show'} Left
          </button>
          <label className="text-sm">Left width
            <input
              type="number"
              className="ml-2 w-24 px-2 py-1 border border-border rounded"
              value={leftWidth}
              min={0}
              max={2000}
              onChange={(e) => setLeftWidth(parseInt(e.target.value || '0', 10))}
            />
          </label>

          <button type="button" className="px-3 py-1 rounded bg-foreground/5 hover:bg-foreground/10" onClick={() => toggleRight('page:demo')}>
            {isRightOpen ? 'Hide' : 'Show'} Right
          </button>
          <label className="text-sm">Right width
            <input
              type="number"
              className="ml-2 w-24 px-2 py-1 border border-border rounded"
              value={rightWidth}
              min={0}
              max={2000}
              onChange={(e) => setRightWidth(parseInt(e.target.value || '0', 10))}
            />
          </label>
        </div>
      </div>

      {/* Simple SSE demo panel */}
      <div className="border border-border rounded p-4 bg-card/50">
        <div className="flex items-center justify-between mb-2">
          <h2 className="font-semibold">Live jobs (SSE)</h2>
          <div className="text-sm">Status: {connectionBadge}</div>
        </div>
        {!token ? (
          <p className="text-sm text-muted">No auth token found. Log in to start receiving live updates.</p>
        ) : (
          <>
            {error && <p className="text-sm text-primary">{error}</p>}
            <div className="flex items-center gap-2 mb-2">
              <button
                type="button"
                className="text-xs px-2 py-1 rounded bg-foreground/5 hover:bg-foreground/10"
                onClick={clear}
              >
                Clear
              </button>
            </div>
            <div className="space-y-2 max-h-64 overflow-auto">
              {events.length === 0 ? (
                <p className="text-sm text-muted">No events yet.</p>
              ) : (
                events.slice(-10).map((evt, idx) => (
                  <pre key={idx} className="text-xs bg-card p-2 rounded border border-border">
                    {JSON.stringify(evt, null, 2)}
                  </pre>
                ))
              )}
            </div>
            {lastEvent && (
              <div className="mt-2 text-xs text-muted">
                Last event type: {typeof lastEvent === 'object' && lastEvent ? (lastEvent as AppEvent).type || 'unknown' : 'n/a'}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
