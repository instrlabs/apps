import { useEffect, useState } from 'react';
import sseService from '../services/sse';

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
    // Connect to SSE when component mounts
    sseService.connect();
    
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
    sseService.addListener('job_notification', handleNotification);
    
    // Add event listeners for connection status
    if (typeof window !== 'undefined') {
      window.addEventListener('sse:connected', handleOpen);
      window.addEventListener('sse:disconnected', handleClose);
    }
    
    // Clean up on unmount
    return () => {
      sseService.removeListener('job_notification', handleNotification);
      
      if (typeof window !== 'undefined') {
        window.removeEventListener('sse:connected', handleOpen);
        window.removeEventListener('sse:disconnected', handleClose);
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
