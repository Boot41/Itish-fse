import { useAuthStore } from '../../stores/authStore';

describe('Auth Store', () => {
  test('should have expected properties', () => {
    const authStore = useAuthStore();
    expect(authStore).toHaveProperty('isAuthenticated');
    expect(authStore).toHaveProperty('clearAuth');
  });
});
