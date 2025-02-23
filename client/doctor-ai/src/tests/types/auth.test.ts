import type { Doctor } from '../../types/auth';

describe('Auth Type', () => {
  test('should have expected properties', () => {
    const doctor: Doctor = {
      ID: '1',
      Name: 'John Doe',
      Email: 'john@example.com',
      Password: 'hashedpassword',
      Phone: '1234567890',
      Specialization: 'General',
      CreatedAt: '2025-01-01',
      Patients: [],
      Transcriptions: [],
      Statistics: null
    };
    expect(doctor).toHaveProperty('ID');
    expect(doctor).toHaveProperty('Name');
    expect(doctor).toHaveProperty('Email');
    expect(doctor).toHaveProperty('Password');
    expect(doctor).toHaveProperty('Phone');
    expect(doctor).toHaveProperty('Specialization');
    expect(doctor).toHaveProperty('CreatedAt');
    expect(doctor).toHaveProperty('Patients');
    expect(doctor).toHaveProperty('Transcriptions');
    expect(doctor).toHaveProperty('Statistics');
  });
});
