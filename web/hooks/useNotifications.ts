import { useEffect, useState } from 'react';
import websocketService from '../services/websocket';

// Type for job notification
export interface JobNotification {
  id: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
}

// Hook for using notifications in components
export function useNotifications() {
  const [notifications, setNotifications] = useState<JobNotification[]>([]);
  const [connected, setConnected] = useState<boolean>(false);

  useEffect(() => {
    // Connect to WebSocket when component mounts
    websocketService.connect();
    
    // Add listener for connection status
    const handleOpen = () => {
      setConnected(true);
    };
    
    const handleClose = () => {
      setConnected(false);
    };
    
    // Add listener for job notifications
    const handleNotification = (data: JobNotification) => {
      setNotifications(prev => [data, ...prev].slice(0, 50)); // Keep last 50 notifications
    };
    
    // Register event listeners
    websocketService.addListener('job_notification', handleNotification);
    
    // Add event listeners for connection status
    if (typeof window !== 'undefined') {
      window.addEventListener('websocket:connected', handleOpen);
      window.addEventListener('websocket:disconnected', handleClose);
    }
    
    // Clean up on unmount
    return () => {
      websocketService.removeListener('job_notification', handleNotification);
      
      if (typeof window !== 'undefined') {
        window.removeEventListener('websocket:connected', handleOpen);
        window.removeEventListener('websocket:disconnected', handleClose);
      }
    };
  }, []);
  
  // Function to clear notifications
  const clearNotifications = () => {
    setNotifications([]);
  };
  
  return {
    notifications,
    connected,
    clearNotifications
  };
}
