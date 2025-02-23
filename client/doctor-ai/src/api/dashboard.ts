import { dummyData } from '../data/dummyDashboard';
// import { useAuthStore } from '../stores/authStore';

const BASE_URL = 'http://localhost:8080/dashboard';

export const dashboardApi = {
  getDailyStatistics: async (days: number) => {
    return dummyData.dailyStatistics(days);
  },

  getMonthlyStatistics: async (months: number) => {
    return dummyData.monthlyStatistics(months);
  },

  getDashboardPatients: async (days: number) => {
    return dummyData.dashboardPatients(days);
  },

  getDashboardTranscripts: async (days: number) => {
    return dummyData.dashboardTranscripts(days);
  },

  getBusiestDays: async (days: number) => {
    return dummyData.busiestDays(days);
  },

  // Keep the profile API connected to the backend
  getProfile: async (email: string) => {
    const response = await fetch(`${BASE_URL}/profile`, {
      method: 'POST',
      body: JSON.stringify({ email }),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
      }
    });

    if (!response.ok) {
      throw new Error('Failed to fetch profile');
    }

    return response.json();
  },

  updateProfile: async (email: string, profile: any) => {
    const response = await fetch(`${BASE_URL}/profile`, {
      method: 'PUT',
      body: JSON.stringify({ email, ...profile }),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
      }
    });

    if (!response.ok) {
      throw new Error('Failed to update profile');
    }

    return response.json();
  },
};
