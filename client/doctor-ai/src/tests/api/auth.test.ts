import { authApi } from '../../api/auth';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

const server = setupServer(
  // Mock login endpoint
  http.post('http://localhost:8080/auth/login', async ({ request }) => {
    const body = (await request.json()) as { email: string; password: string };

    if (body.email === 'test@example.com' && body.password === 'password123') {
      return HttpResponse.json(
        {
          success: true,
          message: 'Login successful',
          data: { token: 'fake-jwt-token' }
        },
        { status: 200 }
      );
    }

    return HttpResponse.json(
      {
        success: false,
        message: 'Invalid credentials'
      },
      { status: 401 }
    );
  }),

  // Mock signup endpoint
  http.post('http://localhost:8080/auth/signup', async ({ request }) => {
    const body = (await request.json()) as { Email: string; Password: string; Name: string };

    if (body.Email && body.Password && body.Name) {
      return HttpResponse.json(
        {
          success: true,
          message: 'Signup successful',
          data: { token: 'fake-jwt-token' }
        },
        { status: 201 }
      );
    }

    return HttpResponse.json(
      {
        success: false,
        message: 'Invalid data'
      },
      { status: 400 }
    );
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('Auth API', () => {
  describe('login', () => {
    it('should successfully login with valid credentials', async () => {
      const response = await authApi.login({
        email: 'test@example.com',
        password: 'password123'
      });

      expect(response.success).toBe(true);
      expect(response.message).toBe('Login successful');
      expect(response.data?.token).toBe('fake-jwt-token');
    });

    it('should fail login with invalid credentials', async () => {
      const response = await authApi.login({
        email: 'wrong@example.com',
        password: 'wrongpass'
      });

      expect(response.success).toBe(false);
      expect(response.message).toBe('Invalid credentials');
    });
  });

  describe('signup', () => {
    it('should successfully create a new doctor account', async () => {
      const response = await authApi.signup({
        Email: 'newdoctor@example.com',
        Password: 'password123',
        Name: 'Dr. Smith'
      });

      expect(response.success).toBe(true);
      expect(response.message).toBe('Signup successful');
      expect(response.data?.token).toBe('fake-jwt-token');
    });

    it('should fail signup with invalid data', async () => {
      const response = await authApi.signup({
        Email: '',
        Password: '',
        Name: ''
      });

      expect(response.success).toBe(false);
      expect(response.message).toBe('Invalid data');
    });
  });
});
