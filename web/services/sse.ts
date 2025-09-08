// SSE (Server-Sent Events) service for handling real-time notifications
class SSEService {
  private eventSource: EventSource | null = null;
  private reconnectInterval: number = 5000; // 5 seconds
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 10;
  private listeners: Map<string, ((data: any) => void)[]> = new Map();
  private baseUrl: string;
  private isConnecting: boolean = false;

  constructor() {
    this.baseUrl = `${process.env.NEXT_PUBLIC_API_URL}/notification/jobs`;
  }

  private getAuthToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  }

  connect(): void {
    if (this.eventSource || this.isConnecting) {
      return;
    }

    this.isConnecting = true;

    try {
      // Get the authentication token and add it to the URL as a query parameter
      const token = this.getAuthToken();
      if (!token) {
        console.error('No authentication token available');
        this.isConnecting = false;
        return;
      }

      const url = `${this.baseUrl}?token=${token}`;

      this.eventSource = new EventSource(url);

      this.eventSource.onopen = () => {
        console.log('SSE connected');
        this.reconnectAttempts = 0;
        this.isConnecting = false;

        // Dispatch a custom event for connection status
        if (typeof window !== 'undefined') {
          window.dispatchEvent(new Event('sse:connected'));
        }
      };

      this.eventSource.addEventListener('message', (event) => {
        try {
          const data = JSON.parse(event.data);
          this.handleMessage(data);
        } catch (error) {
          console.error('Error parsing SSE message:', error);
        }
      });

      this.eventSource.addEventListener('connected', (event) => {
        console.log('SSE connection established:', event.data);
      });

      this.eventSource.onerror = (error) => {
        console.error('SSE error:', error);
        this.disconnect();
        this.isConnecting = false;
        this.attemptReconnect();
      };
    } catch (error) {
      console.error('Error connecting to SSE:', error);
      this.isConnecting = false;
      this.attemptReconnect();
    }
  }

  // Attempt to reconnect to the SSE server
  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);

    setTimeout(() => {
      this.connect();
    }, this.reconnectInterval);
  }

  // Disconnect from the SSE server
  disconnect(): void {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;

      // Dispatch a custom event for disconnection status
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new Event('sse:disconnected'));
      }
    }
  }

  // Add a listener for a specific event type
  addListener(eventType: string, callback: (data: any) => void): void {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, []);
    }

    this.listeners.get(eventType)?.push(callback);
  }

  // Remove a listener for a specific event type
  removeListener(eventType: string, callback: (data: any) => void): void {
    if (!this.listeners.has(eventType)) {
      return;
    }

    const listeners = this.listeners.get(eventType) || [];
    const index = listeners.indexOf(callback);

    if (index !== -1) {
      listeners.splice(index, 1);
    }
  }

  // Handle incoming messages and dispatch to appropriate listeners
  private handleMessage(data: any): void {
    // For job notifications, we expect data to have id and status fields
    if (data.id && data.status) {
      // Dispatch to all listeners for the specific status
      const statusListeners = this.listeners.get(data.status) || [];
      statusListeners.forEach(callback => callback(data));

      // Also dispatch to general job notification listeners
      const generalListeners = this.listeners.get('job_notification') || [];
      generalListeners.forEach(callback => callback(data));
    }
  }

  // Check if currently connected
  isConnected(): boolean {
    return this.eventSource !== null && this.eventSource.readyState === EventSource.OPEN;
  }
}

// Create a singleton instance
const sseService = new SSEService();

export default sseService;
