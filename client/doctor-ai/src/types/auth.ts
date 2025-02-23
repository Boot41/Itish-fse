export interface Doctor {
  ID: string;
  Name: string;
  Specialization: string;
  Email: string;
  Phone: string;
  Password: string;
  CreatedAt: string;
  Patients: any[] | null;
  Transcriptions: any[] | null;
  Statistics: any | null;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  success: boolean;
  message: string;
  data?: {
    token?: string;
    user?: Doctor;
  };
}
