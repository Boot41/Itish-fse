import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { X, Eye, EyeOff } from 'lucide-react';
import { authApi } from '../../api/auth';
import type { Doctor } from '../../types/auth';
import { useAuthStore } from '../../stores/authStore';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onClose }) => {
  const navigate = useNavigate();
  const [isLogin, setIsLogin] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isClosing, setIsClosing] = useState(false);
  const [formData, setFormData] = useState<Doctor>({
    ID: '',
    Name: '',
    Specialization: '',
    Email: '',
    Phone: '',
    Password: '',
    CreatedAt: '',
    Patients: null,
    Transcriptions: null,
    Statistics: null,
  });
  const [showPassword, setShowPassword] = useState(false);

  const handleClose = () => {
    setIsClosing(true);
    setTimeout(onClose, 300); // Match animation duration
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const response = isLogin
        ? await authApi.login({ email: formData.Email, password: formData.Password })
        : await authApi.signup(formData);

      console.log('Auth response:', response); // Debug log

      if (response.success) {
        // Store email in auth store
        await useAuthStore.getState().setAuth(formData.Email, response.data?.token || '');
        setSuccess(response.message || (isLogin ? 'Login successful' : 'Account created successfully'));
        handleClose();
        // Navigate to dashboard after successful auth
        navigate('/dashboard', { replace: true });
      } else {
        setError(response.message || 'Authentication failed');
      } 
    } catch (error) {
      setError('An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const field = e.target.name.charAt(0).toUpperCase() + e.target.name.slice(1);
    setFormData({ ...formData, [field]: e.target.value });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/80 backdrop-blur-md transition-opacity duration-300" 
        onClick={handleClose} 
      />
      <div 
        className={`relative w-full max-w-md rounded-2xl auth-modal-bg p-8 shadow-2xl border ${isClosing ? 'modal-exit' : 'modal-enter'}`}
        style={{ boxShadow: '0 15px 50px rgba(0, 255, 0, 0.1)' }}
      >
        {/* Close button */}
        <button
          onClick={handleClose}
          className="absolute right-4 top-4 text-gray-400 hover:text-white hover:rotate-90 transition-all duration-300"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="mb-8 text-center animate-fadeIn">
          <h2 className="text-2xl font-bold text-white mb-2">
            {isLogin ? 'Welcome Back' : 'Create Account'}
          </h2>
          <p className="text-gray-400">
            {isLogin
              ? 'Sign in to access your account'
              : 'Join us to experience AI-powered healthcare'}
          </p>
        </div>

        <form 
          key={isLogin ? 'login' : 'signup'}
          onSubmit={handleSubmit} 
          className="space-y-6 animate-fadeIn delay-150 form-switch"
        >
          {!isLogin && (
            <>
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-1">
                  Full Name
                </label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={formData.Name}
                  onChange={handleChange}
                  className="auth-input"
                  placeholder="Dr. John Doe"
                  required
                />
              </div>
              <div>
                <label htmlFor="specialization" className="block text-sm font-medium text-gray-300 mb-1">
                  Specialization
                </label>
                <input
                  type="text"
                  id="specialization"
                  name="specialization"
                  value={formData.Specialization}
                  onChange={handleChange}
                  className="auth-input"
                  placeholder="Cardiology"
                  required
                />
              </div>
              <div>
                <label htmlFor="phone" className="block text-sm font-medium text-gray-300 mb-1">
                  Phone Number
                </label>
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  value={formData.Phone}
                  onChange={handleChange}
                  className="auth-input"
                  placeholder="1234567890"
                  required
                />
              </div>
            </>
          )}
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-1">
              Email
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.Email}
              onChange={handleChange}
              className="auth-input"
              placeholder="doctor@example.com"
              required
            />
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-300 mb-1">
              Password
            </label>
            <div className="relative">
              <input
                type={showPassword ? 'text' : 'password'}
                id="password"
                name="password"
                value={formData.Password}
                onChange={handleChange}
                className="auth-input"
                placeholder="••••••••"
                required
              />
              <button
                type="button"
                className={`password-toggle ${isLogin ? 'login-form' : ''}`}
                onClick={() => setShowPassword(!showPassword)}
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>

          {(error || success) && (
            <p className={`text-sm ${error ? 'text-red-400' : 'text-green-400'}`}>
              {error || success}
            </p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full btn btn-primary relative overflow-hidden group hover:scale-[1.02] active:scale-[0.98] transition-all duration-300"
          >
            <span className={`${loading ? 'invisible' : ''}`}>
              {isLogin ? 'Sign In' : 'Create Account'}
            </span>
            {loading && (
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-5 h-5 border-2 border-dark border-t-transparent rounded-full animate-spin" />
              </div>
            )}
          </button>
        </form>

        <div className="mt-8 text-center animate-fadeIn delay-300">
          <p className="text-gray-400">
            {isLogin ? "Don't have an account?" : 'Already have an account?'}
            <button
              onClick={() => setIsLogin(!isLogin)}
              className="ml-2 text-primary hover:text-primary-light transition-colors"
            >
              {isLogin ? 'Sign Up' : 'Sign In'}
            </button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default AuthModal;
