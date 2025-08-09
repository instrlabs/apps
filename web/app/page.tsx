'use client';

import { useEffect, useState } from 'react';
import { verifyToken } from '@/services/auth';
import { useNotification } from '@/components/notification';
import { NOTIFICATION_DURATION } from '@/components/ui-styles';

export default function Home() {
  const [user, setUser] = useState<null | { [key: string]: unknown }>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { showNotification } = useNotification();

  useEffect(() => {
    const checkAuth = async () => {
      setIsLoading(true);

      const { data, error } = await verifyToken();

      if (error) showNotification(error, "error", NOTIFICATION_DURATION);
      else if (data) setUser(data.data.user);

      setIsLoading(false);
    };

    checkAuth().then();
  }, []);

  return (
      <div className="container mx-auto p-4">
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
      </div>
  );
}
