import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_URL = 'http://backend:8080';

export const fallback: RequestHandler = async ({ request, params, url }) => {
  const path = url.pathname.replace('/api/v1/', '');
  const targetUrl = `${BACKEND_URL}/api/v1/${path}${url.search}`;

  try {
    const headers = new Headers(request.headers);
    
    headers.set('host', 'backend:8080');

    const res = await fetch(targetUrl, {
      method: request.method,
      headers: headers,
      body: request.method !== 'GET' && request.method !== 'HEAD' ? await request.arrayBuffer() : undefined,
      duplex: 'half'
    } as any);

    const data = await res.arrayBuffer();

    return new Response(data, {
      status: res.status,
      headers: {
        'Content-Type': res.headers.get('Content-Type') || 'application/json',
        'Cache-Control': 'no-cache'
      }
    });
  } catch (err: any) {
    console.error(`[Proxy Error] ${request.method} ${targetUrl}:`, err);
    throw error(502, `Gagal menghubungi backend: ${err.message}`);
  }
};
