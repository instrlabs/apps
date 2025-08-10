// WebSocket service for handling real-time notifications
class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectInterval: number = 5000; // 5 seconds
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 10;
  private listeners: Map<string, ((data: any) => void)[]> = new Map();
  private url: string;

  constructor() {
    // Get WebSocket URL from environment or use default
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = process.env.NEXT_PUBLIC_NOTIFICATION_SERVICE_HOST || window.location.hostname;
    const port = process.env.NEXT_PUBLIC_NOTIFICATION_SERVICE_PORT || '3030';
    const path = process.env.NEXT_PUBLIC_NOTIFICATION_SERVICE_PATH || '/ws';
    
    this.url = `${protocol}//${host}:${port}${path}`;
  }

  // Connect to the WebSocket server
  connect(): void {
    if (this.socket) {
      return;
    }

    try {
      this.socket = new WebSocket(this.url);

      this.socket.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
      };

      this.socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.handleMessage(data);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
        this.socket = null;
        this.attemptReconnect();
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.socket?.close();
      };
    } catch (error) {
      console.error('Error connecting to WebSocket:', error);
      this.attemptReconnect();
    }
  }

  // Attempt to reconnect to the WebSocket server
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

  // Disconnect from the WebSocket server
  disconnect(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
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
}

// Create a singleton instance
const websocketService = new WebSocketService();

export default websocketService;
