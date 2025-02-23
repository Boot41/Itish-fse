import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

// Encryption utility
const encrypt = async (data: string) => {
  try {
    const salt = crypto.getRandomValues(new Uint8Array(16));
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const encoder = new TextEncoder();
    const dataBuffer = encoder.encode(data);

    const baseKey = await crypto.subtle.importKey(
      'raw', 
      salt, 
      { name: 'PBKDF2' }, 
      false, 
      ['deriveBits', 'deriveKey']
    );

    const key = await crypto.subtle.deriveKey(
      {
        name: 'PBKDF2',
        salt: salt,
        iterations: 100000,
        hash: 'SHA-256'
      },
      baseKey,
      { name: 'AES-GCM', length: 256 },
      true,
      ['encrypt', 'decrypt']
    );

    const encrypted = await crypto.subtle.encrypt(
      { name: 'AES-GCM', iv: iv },
      key,
      dataBuffer
    );

    const combinedArray = new Uint8Array(salt.length + iv.length + encrypted.byteLength);
    combinedArray.set(salt);
    combinedArray.set(iv, salt.length);
    combinedArray.set(new Uint8Array(encrypted), salt.length + iv.length);
    
    const result = btoa(String.fromCharCode.apply(null, Array.from(combinedArray)));
    console.log('Encryption successful:', result);
    return result;
  } catch (error) {
    console.error('Encryption error', error);
    return null;
  }
};

const decrypt = async (encryptedData: string | null) => {
  if (!encryptedData) return null;

  try {
    const combinedArray = new Uint8Array(
      atob(encryptedData)
        .split('')
        .map(char => char.charCodeAt(0))
    );
    
    const salt = combinedArray.slice(0, 16);
    const iv = combinedArray.slice(16, 28);
    const data = combinedArray.slice(28);

    const baseKey = await crypto.subtle.importKey(
      'raw', 
      salt, 
      { name: 'PBKDF2' }, 
      false, 
      ['deriveBits', 'deriveKey']
    );

    const key = await crypto.subtle.deriveKey(
      {
        name: 'PBKDF2',
        salt: salt,
        iterations: 100000,
        hash: 'SHA-256'
      },
      baseKey,
      { name: 'AES-GCM', length: 256 },
      true,
      ['encrypt', 'decrypt']
    );

    const decrypted = await crypto.subtle.decrypt(
      { name: 'AES-GCM', iv: iv },
      key,
      data
    );

    const decoder = new TextDecoder();
    const result = decoder.decode(decrypted);
    console.log('Decryption successful:', result);
    return result;
  } catch (error) {
    console.error('Decryption error', error);
    return null;
  }
};

interface AuthState {
  email: string | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (email: string, token: string) => Promise<void>;
  clearAuth: () => void;
  getDecryptedEmail: () => Promise<string | null>;
  getDecryptedToken: () => Promise<string | null>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      email: null,
      token: null,
      isAuthenticated: false,
      
      setAuth: async (email: string, token: string) => {
        const encryptedEmail = await encrypt(email);
        const encryptedToken = await encrypt(token);
        
        if (encryptedEmail && encryptedToken) {
          set({ 
            email: encryptedEmail, 
            token: encryptedToken,
            isAuthenticated: true
          });
          console.log('Auth state stored:', { email: encryptedEmail, token: encryptedToken });
        } else {
          console.error('Failed to encrypt auth data');
        }
      },
      
      clearAuth: () => {
        console.log('Clearing auth state');
        set({ 
          email: null, 
          token: null,
          isAuthenticated: false 
        });
      },
      
      getDecryptedEmail: async () => {
        const { email } = get();
        const decryptedEmail = await decrypt(email);
        console.log('Decrypted email:', decryptedEmail);
        return decryptedEmail;
      },
      
      getDecryptedToken: async () => {
        const { token } = get();
        const decryptedToken = await decrypt(token);
        console.log('Decrypted token:', decryptedToken);
        return decryptedToken;
      }
    }),
    {
      name: 'doctor-auth-storage',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({ 
        email: state.email, 
        token: state.token,
        isAuthenticated: state.isAuthenticated
      }),
      onRehydrateStorage: (state) => {
        console.log('Rehydrating auth state:', state);
      }
    }
  )
);
