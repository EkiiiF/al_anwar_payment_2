import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('$env/static/public', () => ({ PUBLIC_API_URL: 'http://localhost:8080' }));
vi.mock('$app/navigation', () => ({ goto: vi.fn() }));

function mockFetch(status: number, body: unknown) {
  return vi.fn().mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body)
  } as Response);
}

describe('apiClient', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    localStorage.clear();
  });

  it('mengirim Authorization header jika ada token', async () => {
    const { apiClient, tokenStore } = await import('$lib/api');
    tokenStore.set('test-token');

    const fetchMock = mockFetch(200, { status: 'ok', message: '', data: {} });
    vi.stubGlobal('fetch', fetchMock);

    await apiClient.get('/api/test');

    const calledWith = fetchMock.mock.calls[0][1] as RequestInit;
    expect((calledWith.headers as Record<string, string>)['Authorization']).toBe('Bearer test-token');
  });

  it('melempar error jika response tidak ok', async () => {
    const { apiClient } = await import('$lib/api');
    vi.stubGlobal('fetch', mockFetch(400, { status: 'error', message: 'Bad request', data: null }));

    await expect(apiClient.get('/api/test')).rejects.toThrow('Bad request');
  });

  it('redirect ke /login jika response 401', async () => {
    const { goto } = await import('$app/navigation');
    const { apiClient, tokenStore } = await import('$lib/api');
    tokenStore.set('expired-token');

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: false, status: 401,
      json: () => Promise.resolve({ message: 'Unauthorized', data: null })
    }));

    await expect(apiClient.get('/api/secure')).rejects.toThrow();
    expect(goto).toHaveBeenCalledWith('/login');
    expect(tokenStore.get()).toBeNull();
  });

  it('POST mengirimkan body JSON dengan benar', async () => {
    const { apiClient } = await import('$lib/api');
    const fetchMock = mockFetch(201, { status: 'ok', message: 'Created', data: { id: '123' } });
    vi.stubGlobal('fetch', fetchMock);

    await apiClient.post('/api/test', { name: 'Al Anwar', amount: 1000 });

    const calledWith = fetchMock.mock.calls[0][1] as RequestInit;
    expect(calledWith.method).toBe('POST');
    expect(JSON.parse(calledWith.body as string)).toEqual({ name: 'Al Anwar', amount: 1000 });
  });
});
