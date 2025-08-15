'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { verifyToken } from '@/services/auth';
import { NotificationProvider, Notification, useNotification } from '@/components/notification';
import { useSSECache } from '@/hooks/useSSECache';
import { OverlayProvider, useOverlay } from '@/hooks/useOverlay';
import OverlayTop from '@/components/overlay-top';
import OverlayLeft from '@/components/overlay-left';
import OverlayRight from '@/components/overlay-right';
import OverlayContent from '@/components/overlay-content';
import OverlayModal from '@/components/overlay-modal';

// Define a minimal event type used by SSE in this page
type AppEvent = {
  type?: string;
  [key: string]: unknown;
};

function HomeContent() {
  const [user, setUser] = useState<null | { [key: string]: unknown }>(null);
  const [isLoading, setIsLoading] = useState(false);
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
      setIsLoading(true);

      const { data, error } = await verifyToken();

      if (error) showNotification(error, "error", 5000);
      else if (data) setUser(data.data.user);

      setIsLoading(false);
    };

    checkAuth().then();
  }, [showNotification]);

  const connectionBadge = useMemo(() => {
    if (!token) return <span className="text-gray-500">No token</span>;
    if (error) return <span className="text-red-600">Error</span>;
    return (
      <span className={isConnected ? 'text-green-600' : 'text-amber-600'}>
        {isConnected ? 'Connected' : 'Connecting...'}
      </span>
    );
  }, [token, isConnected, error]);

  const { isLeftOpen, isRightOpen, toggleLeft, toggleRight, leftWidth, rightWidth, setLeftWidth, setRightWidth } = useOverlay();

  return (
    <div className="container mx-auto p-4 space-y-6">
      {/* Overlay controls (demo) */}
      <div className="border rounded p-4 bg-white/60">
        <div className="flex flex-wrap items-center gap-3">
          <button type="button" className="px-3 py-1 rounded bg-gray-100 hover:bg-gray-200" onClick={() => { toggleLeft('page:demo-left'); }}>
            {isLeftOpen ? 'Hide' : 'Show'} Left
          </button>
          <label className="text-sm">Left width
            <input
              type="number"
              className="ml-2 w-24 px-2 py-1 border rounded"
              value={leftWidth}
              min={0}
              max={2000}
              onChange={(e) => setLeftWidth(parseInt(e.target.value || '0', 10))}
            />
          </label>

          <button type="button" className="px-3 py-1 rounded bg-gray-100 hover:bg-gray-200" onClick={() => toggleRight('page:demo')}>
            {isRightOpen ? 'Hide' : 'Show'} Right
          </button>
          <label className="text-sm">Right width
            <input
              type="number"
              className="ml-2 w-24 px-2 py-1 border rounded"
              value={rightWidth}
              min={0}
              max={2000}
              onChange={(e) => setRightWidth(parseInt(e.target.value || '0', 10))}
            />
          </label>
        </div>
      </div>

      {isLoading ? (
        <p>Loading...</p>
      ) : user ? (
        <div>
          <h1 className="text-2xl font-bold mb-4">Welcome back!</h1>
          <pre className="bg-gray-100 p-4 rounded">
            {JSON.stringify(user, null, 2)}
          </pre>
        </div>
      ) : (
        <div>
          <h1 className="text-2xl font-bold mb-4">Hello World!</h1>
          <p>Please log in to see your profile information.</p>
        </div>
      )}

      {/* Simple SSE demo panel */}
      <div className="border rounded p-4 bg-white/50">
        <div className="flex items-center justify-between mb-2">
          <h2 className="font-semibold">Live jobs (SSE)</h2>
          <div className="text-sm">Status: {connectionBadge}</div>
        </div>
        {!token ? (
          <p className="text-sm text-gray-600">No auth token found. Log in to start receiving live updates.</p>
        ) : (
          <>
            {error && <p className="text-sm text-red-600">{error}</p>}
            <div className="flex items-center gap-2 mb-2">
              <button
                type="button"
                className="text-xs px-2 py-1 rounded bg-gray-100 hover:bg-gray-200"
                onClick={clear}
              >
                Clear
              </button>
            </div>
            <div className="space-y-2 max-h-64 overflow-auto">
              {events.length === 0 ? (
                <p className="text-sm text-gray-600">No events yet.</p>
              ) : (
                events.slice(-10).map((evt, idx) => (
                  <pre key={idx} className="text-xs bg-gray-50 p-2 rounded border">
                    {JSON.stringify(evt, null, 2)}
                  </pre>
                ))
              )}
            </div>
            {lastEvent && (
              <div className="mt-2 text-xs text-gray-500">
                Last event type: {typeof lastEvent === 'object' && lastEvent ? (lastEvent as AppEvent).type || 'unknown' : 'n/a'}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

export default function Home() {
  return (
    <NotificationProvider>
      <OverlayProvider>
        <OverlayContent>
          <HomeContent />
        </OverlayContent>
        <OverlayLeft />
        <OverlayRight />
        <OverlayTop />
        <OverlayModal />
        <Notification />
      </OverlayProvider>
    </NotificationProvider>
  );
}
