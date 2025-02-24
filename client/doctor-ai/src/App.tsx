import React, { Suspense, lazy, Component } from 'react';
import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from './stores/authStore';

// Lazy load components to improve performance
const MainLayout = lazy(() => import('./components/layout/MainLayout'));
const Dashboard = lazy(() => import('./components/dashboard/Dashboard'));
const LandingPage = lazy(() => import('./components/landing/LandingPage'));
const About = lazy(() => import('./pages/About'));
const Features = lazy(() => import('./pages/Features'));
const Documentation = lazy(() => import('./pages/Documentation'));
const Privacy = lazy(() => import('./pages/Privacy'));
const Terms = lazy(() => import('./pages/Terms'));
const Support = lazy(() => import('./pages/Support'));
const Contact = lazy(() => import('./pages/Contact'));
const Patients = lazy(() => import('./components/dashboard/Patients'));
const Logs = lazy(() => import('./components/dashboard/Transcripts'));

// Loading Fallback Component
const LoadingFallback = () => (
  <div className="flex justify-center items-center min-h-screen bg-dark-background">
    <div className="animate-pulse text-primary">Loading...</div>
  </div>
);

// Error Boundary Component
class ErrorBoundary extends Component<
  { children: React.ReactNode },
  { hasError: boolean }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(_: Error) {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error:', error);
    console.error('Error Info:', errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="flex justify-center items-center min-h-screen bg-dark-background text-red-500">
          Something went wrong. Please try again later.
        </div>
      );
    }

    return this.props.children;
  }
}

// Protected Route Component
const ProtectedRoute = () => {
  const { isAuthenticated } = useAuthStore();
  return isAuthenticated ? <Outlet /> : <Navigate to="/" replace />;
};

const App: React.FC = () => {
  const { isAuthenticated } = useAuthStore();

  return (
    <ErrorBoundary>
      <Suspense fallback={<LoadingFallback />}>
        <BrowserRouter>
          <Routes>
            <Route element={<MainLayout><Outlet /></MainLayout>}>
              {/* Landing Page */}
              <Route 
                path="/" 
                element={
                  isAuthenticated ? (
                    <Navigate to="/dashboard" replace />
                  ) : (
                    <LandingPage />
                  )
                } 
              />
              
              {/* Public Routes */}
              <Route path="/about" element={<About />} />
              <Route path="/features" element={<Features />} />
              <Route path="/documentation" element={<Documentation />} />
              <Route path="/privacy" element={<Privacy />} />
              <Route path="/terms" element={<Terms />} />
              <Route path="/support" element={<Support />} />
              <Route path="/contact" element={<Contact />} />
              
              {/* Protected Routes */}
              <Route element={<ProtectedRoute />}>
                <Route path="/dashboard" element={<Dashboard />} />
                <Route path="/patients" element={<Patients />} />
                <Route path="/logs" element={<Logs />} />
              </Route>
              
              {/* Catch-all route */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </Suspense>
    </ErrorBoundary>
  );
};

export default App;
