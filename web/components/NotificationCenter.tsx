import React, { useEffect } from 'react';
import { useNotifications } from '../hooks/useNotifications';

interface NotificationCenterProps {
  maxNotifications?: number;
}

/**
 * NotificationCenter component displays real-time notifications from the notification service
 * 
 * Example usage:
 * ```tsx
 * <NotificationCenter maxNotifications={5} />
 * ```
 */
const NotificationCenter: React.FC<NotificationCenterProps> = ({ 
  maxNotifications = 5 
}) => {
  const { notifications, connected, clearNotifications } = useNotifications();
  
  // Display only the most recent notifications up to maxNotifications
  const recentNotifications = notifications.slice(0, maxNotifications);
  
  // Get status color based on notification status
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'text-green-500';
      case 'failed':
        return 'text-red-500';
      case 'processing':
        return 'text-blue-500';
      default:
        return 'text-gray-500';
    }
  };
  
  // Get connection status indicator color
  const connectionColor = connected ? 'bg-green-500' : 'bg-red-500';
  
  return (
    <div className="bg-white shadow-md rounded-lg p-4 max-w-md">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold">Notifications</h2>
        <div className="flex items-center">
          <span className="text-sm mr-2">
            {connected ? 'Connected' : 'Disconnected'}
          </span>
          <div className={`w-3 h-3 rounded-full ${connectionColor}`}></div>
        </div>
      </div>
      
      {recentNotifications.length > 0 ? (
        <div>
          <ul className="space-y-2">
            {recentNotifications.map((notification, index) => (
              <li 
                key={`${notification.id}-${index}`} 
                className="border-b border-gray-100 pb-2"
              >
                <div className="flex justify-between">
                  <span className="font-medium">Job {notification.id}</span>
                  <span className={`${getStatusColor(notification.status)}`}>
                    {notification.status}
                  </span>
                </div>
              </li>
            ))}
          </ul>
          
          <div className="mt-4 text-right">
            <button
              onClick={clearNotifications}
              className="text-sm text-blue-500 hover:text-blue-700"
            >
              Clear all
            </button>
          </div>
        </div>
      ) : (
        <div className="text-center py-4 text-gray-500">
          No notifications yet
        </div>
      )}
    </div>
  );
};

export default NotificationCenter;
