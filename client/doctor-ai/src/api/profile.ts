import { Doctor } from '../types/auth';

const API_URL = 'http://localhost:8080/auth';

export const profileApi = {
  getProfile: async (email: string): Promise<{ success: boolean; data?: Doctor; message: string }> => {
    try {
      const response = await fetch(`${API_URL}/me`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'Origin': import.meta.env.VITE_FRONTEND_URL,
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();
      console.log('Profile response:', data);

      if (!response.ok) {
        return {
          success: false,
          message: data.message || 'Failed to fetch profile',
        };
      }

      return {
        success: true,
        message: 'Profile fetched successfully',
        data: data as Doctor,
      };
    } catch (error: any) {
      return {
        success: false,
        message: error.message || 'Failed to fetch profile',
      };
    }
  },

  updateProfile: async (email: string, profile: Partial<Doctor>): Promise<{ success: boolean; message: string }> => {
    try {
      const response = await fetch(`${API_URL}/me`, {
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'Origin': import.meta.env.VITE_FRONTEND_URL,
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify({ email, ...profile }),
      });

      const data = await response.json();
      console.log('Update profile response:', data);

      if (!response.ok) {
        return {
          success: false,
          message: data.message || 'Failed to update profile',
        };
      }

      return {
        success: true,
        message: 'Profile updated successfully',
      };
    } catch (error: any) {
      return {
        success: false,
        message: error.message || 'Failed to update profile',
      };
    }
  },
};
