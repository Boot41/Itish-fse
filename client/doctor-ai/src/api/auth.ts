import { Doctor, LoginCredentials, AuthResponse } from '../types/auth';

const API_URL = 'http://localhost:8080/auth';

export const authApi = {
  login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
    try {
      const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });
      
      const data = await response.json();
      console.log('Login Server response:', data); // Debug log

      if (!response.ok) {
        return {
          success: false,
          message: data.message || 'Login failed',
        };
      }

      return {
        success: true,
        message: 'Login successful',
        data: {
          token: data.token,
        },
      };
    } catch (error: any) {
      console.error('Login error:', error);
      return {
        success: false,
        message: error.message || 'Login failed',
      };
    }
  },

  signup: async (doctor: Partial<Doctor>): Promise<AuthResponse> => {
    try {
      // Ensure all required fields are present
      const signupData = {
        Name: doctor.Name || '',
        Specialization: doctor.Specialization || '',
        Email: doctor.Email || '',
        Phone: doctor.Phone || '',
        Password: doctor.Password || '',
        // Optional fields can be omitted or set to default values
        CreatedAt: doctor.CreatedAt || new Date().toISOString(),
        Patients: doctor.Patients || null,
        Transcriptions: doctor.Transcriptions || null,
        Statistics: doctor.Statistics || null,
      };

      const response = await fetch(`${API_URL}/signup`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(signupData),
      });

      const data = await response.json();
      console.log('Signup Server response:', data); // Debug log

      if (!response.ok) {
        return {
          success: false,
          message: data.message || 'Signup failed',
        };
      }

      return {
        success: true,
        message: 'Account created successfully',
        data: {
          token: data.token,
        },
      };
    } catch (error: any) {
      console.error('Signup error:', error);
      return {
        success: false,
        message: error.message || 'Signup failed',
      };
    }
  },
};
