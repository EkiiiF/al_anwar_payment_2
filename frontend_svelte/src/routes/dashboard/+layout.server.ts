import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
  const token = locals.userToken;

  if (!token) {
    throw redirect(302, '/login');
  }

  return {
    hasSession: !!locals.userToken
  };
};
