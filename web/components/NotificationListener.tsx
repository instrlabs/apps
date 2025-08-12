import { useEffect, useState } from 'react';
import { 
  initializeNotifications, 
  addJobNotificationListener, 
  removeJobNotificationListener 
} from '../utils/notification';

interface Notification {
  id: string;
  status: string;
  message: string;
  timestamp: string;
}

interface NotificationListenerProps {
  onNotification?: (notification: Notification) => void;
}

/**
 * Component that listens for SSE notifications and can display them or pass them to a parent component
 */
const NotificationListener: React.FC<NotificationListenerProps> = ({ onNotification }) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [connectionStatus, setConnectionStatus] = useState<'connected' | 'disconnected'>('disconnected');

  useEffect(() => {
    // Initialize SSE connection
    initializeNotifications();

    // Set up event listeners for connection status
    const handleConnected = () => setConnectionStatus('connected');
    const handleDisconnected = () => setConnectionStatus('disconnected');

    window.addEventListener('sse:connected', handleConnected);
    window.addEventListener('sse:disconnected', handleDisconnected);

    // Handle notifications
    const handleNotification = (data: any) => {
      const notification: Notification = {
        id: data.id,
        status: data.status,
        message: data.message || `Job ${data.id} is ${data.status}`,
        timestamp: new Date().toISOString()
      };

      // Update local state
      setNotifications(prev => [notification, ...prev].slice(0, 10)); // Keep only the 10 most recent

      // Call the callback if provided
      if (onNotification) {
        onNotification(notification);
      }
    };

    // Add listener for job notifications
    addJobNotificationListener(handleNotification);

    // Clean up on unmount
    return () => {
      window.removeEventListener('sse:connected', handleConnected);
      window.removeEventListener('sse:disconnected', handleDisconnected);
      removeJobNotificationListener(handleNotification);
    };
  }, [onNotification]);

  // If this component is only used for listening and passing notifications up,
  // you can return null or a minimal UI showing connection status
  return (
    <div className="notification-listener">
      <div className={`connection-status ${connectionStatus}`}>
        SSE Status: {connectionStatus}
      </div>
      
      {notifications.length > 0 && (
        <div className="notifications-list">
          <h3>Recent Notifications</h3>
          <ul>
            {notifications.map((notification) => (
              <li key={`${notification.id}-${notification.timestamp}`} className={`notification ${notification.status}`}>
                <span className="notification-time">
                  {new Date(notification.timestamp).toLocaleTimeString()}
                </span>
                <span className="notification-message">
                  {notification.message}
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default NotificationListener;
