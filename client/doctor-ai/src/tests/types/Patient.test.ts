import { Patient } from '../../types/Patient';

describe('Patient Type', () => {
  test('should have expected properties', () => {
    const patient: Patient = {
      id: '1',
      name: 'John Doe',
      age: 30,
      gender: 'Male',
      createdAt: '2025-01-01'
    };
    expect(patient).toHaveProperty('id');
    expect(patient).toHaveProperty('name');
    expect(patient).toHaveProperty('age');
    expect(patient).toHaveProperty('gender');

    expect(patient).toHaveProperty('created_at');
  });
});
