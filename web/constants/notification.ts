export const NOTIFICATION_BASE_URL =
  process.env.NEXT_PUBLIC_NOTIFICATION_URL || 'http://notification.localhost';

export const NOTIFICATION_JOBS_URL = `${NOTIFICATION_BASE_URL}/jobs`;
