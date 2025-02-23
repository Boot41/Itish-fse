import { chatApi } from '../../api/chat';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';
import { vi } from 'vitest';

// Define your API base URL (used by your chatApi)
const API_BASE_URL = 'http://mock-api-url';
vi.stubEnv('VITE_API_BASE_URL', API_BASE_URL);
vi.stubEnv('VITE_GROQAPIKEY', 'mock-api-key');

// Set up the MSW server using the new http API
const server = setupServer(
  http.post(`${API_BASE_URL}/chat/completions`, async ({ request }) => {
    // Parse the JSON request body
    const body = (await request.json()) as {
      messages: Array<{ role: string; content: string }>;
    };

    // Check for valid authorization header
    const authHeader = request.headers.get('Authorization');
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    // Validate the messages array: ensure it exists, is an array, and contains at least two items
    if (!body.messages || !Array.isArray(body.messages) || body.messages.length < 2) {
      return HttpResponse.json({ error: 'Invalid message format' }, { status: 400 });
    }

    // Return a successful JSON response
    return HttpResponse.json(
      { response: 'This is a mock response from the AI medical assistant.' },
      { status: 200 }
    );
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('Chat API', () => {
  it('should successfully send a message and receive a response', async () => {
    const response = await chatApi.sendMessage('What are the symptoms of flu?');

    expect(response).toBeDefined();
    expect(response.response).toBe('This is a mock response from the AI medical assistant.');
  });

  it('should handle API errors gracefully', async () => {
    // Temporarily override the handler to simulate an internal server error
    server.use(
      http.post(`${API_BASE_URL}/chat/completions`, () => {
        return HttpResponse.json({ error: 'Internal server error' }, { status: 500 });
      })
    );

    await expect(chatApi.sendMessage('test message')).rejects.toThrow();
  });
});
