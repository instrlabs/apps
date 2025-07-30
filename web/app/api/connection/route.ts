import { NextResponse } from 'next/server';

export async function GET() {
  try {
    // Get the gateway URL from environment variables or use a default
    const gatewayUrl = 'http://gateway-service:3000';
    
    // Make a request to the gateway's health endpoint
    const response = await fetch(`${gatewayUrl}/health`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Set a timeout to avoid hanging requests
      signal: AbortSignal.timeout(5000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { 
          status: 'error', 
          message: `Gateway service returned status: ${response.status}` 
        },
        { status: 503 }
      );
    }

    const data = await response.json();
    
    return NextResponse.json({ 
      status: 'ok',
      message: 'Successfully connected to gateway service',
      gateway: data
    });
  } catch (error) {
    console.error('Error connecting to gateway service:', error);
    
    return NextResponse.json(
      { 
        status: 'error', 
        message: 'Failed to connect to gateway service',
        error: error instanceof Error ? error.message : String(error)
      },
      { status: 503 }
    );
  }
}