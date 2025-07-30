import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const gatewayUrl = 'http://gateway-service:3000';
    
    const response = await fetch(`${gatewayUrl}/health`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      signal: AbortSignal.timeout(5000),
    });

    if (!response.ok) {
      return NextResponse.json({
        status: 'error',
        message: `Gateway service returned status: ${response.status}`
      });
    }

    return NextResponse.json({
      status: 'ok',
      message: 'Successfully connected to gateway service',
      gateway: 'ok'
    });
  } catch (error) {
    console.error('Error connecting to gateway service:', error);
    
    return NextResponse.json({
      status: 'error',
      message: 'Failed to connect to gateway service',
      error: error instanceof Error ? error.message : String(error)
    })
  }
}