import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals }) => {
  if (!locals.userToken) {
    throw error(401, 'Tidak ada sesi aktif');
  }
  return json({ token: locals.userToken });
};

export const POST: RequestHandler = async ({ request, cookies }) => {
  const body = await request.json();
  const token = body.token as string;

  if (!token) {
    return json({ success: false, message: 'Token is required' }, { status: 400 });
  }

  cookies.set('__session_token', token, {
    path: '/',
    httpOnly: true,
    secure: false,
    sameSite: 'strict',
    maxAge: 60 * 60 * 24
  });

  return json({ success: true });
};

export const DELETE: RequestHandler = async ({ cookies }) => {
  cookies.delete('__session_token', { path: '/' });

  return json({ success: true });
};
