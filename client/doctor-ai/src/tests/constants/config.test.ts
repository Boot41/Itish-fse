import { FRONTEND_URL } from '../../constants/config';

describe('Config Constants', () => {
  test('should have expected properties', () => {
    expect(FRONTEND_URL).toHaveProperty('apiUrl'); // Adjust based on actual properties
    expect(FRONTEND_URL).toHaveProperty('timeout'); // Adjust based on actual properties
  });
});
