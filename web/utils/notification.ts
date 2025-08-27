import sseService from '../services/sse';

/**
 * Initialize the SSE connection for real-time notifications
 * This should be called when the application starts or when a user logs in
 */
export function initializeNotifications(): void {
  // Check if we're in a browser environment
  if (typeof window === 'undefined') {
    return;
  }

  // Connect to the SSE service
  sseService.connect();

  // Set up reconnection on visibility change
  // This helps reconnect when a user returns to the tab after it's been inactive
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible' && !sseService.isConnected()) {
      sseService.connect();
    }
  });
}

/**
 * Add a listener for job notifications
 * @param callback Function to call when a notification is received
 */
export function addJobNotificationListener(callback: (data: any) => void): void {
  sseService.addListener('job_notification', callback);
}

/**
 * Remove a job notification listener
 * @param callback The callback function to remove
 */
export function removeJobNotificationListener(callback: (data: any) => void): void {
  sseService.removeListener('job_notification', callback);
}

/**
 * Add a listener for a specific job status
 * @param status The job status to listen for (e.g., 'completed', 'failed')
 * @param callback Function to call when a notification with this status is received
 */
export function addJobStatusListener(status: string, callback: (data: any) => void): void {
  sseService.addListener(status, callback);
}

/**
 * Remove a job status listener
 * @param status The job status
 * @param callback The callback function to remove
 */
export function removeJobStatusListener(status: string, callback: (data: any) => void): void {
  sseService.removeListener(status, callback);
}
