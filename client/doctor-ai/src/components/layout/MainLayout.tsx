import { FC, ReactNode, useState, useEffect } from 'react';
import { Link, NavLink, useLocation, useNavigate } from 'react-router-dom';
import { Stethoscope, Menu, X, LogOut, Users, ClipboardList, Github, Twitter, Linkedin } from 'lucide-react';
import { useAuthStore } from '../../stores/authStore';
import AuthModal from '../auth/AuthModal';

interface MainLayoutProps {
  children: ReactNode;
}

const MainLayout: FC<MainLayoutProps> = ({ children }) => {
  const { isAuthenticated, clearAuth } = useAuthStore();
  const [isOpen, setIsOpen] = useState(false);
  const [isScrolled, setIsScrolled] = useState(false);
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();

  const handleSignOut = () => {
    clearAuth();
    navigate('/', { replace: true });
  };

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    setIsOpen(false);
    window.scrollTo(0, 0);
  }, [location]);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  return (
    <div className="min-h-screen">
      <header className={`header fixed top-0 left-0 right-0 bg-[#1C1C1C] z-40 ${
        isScrolled ? 'scrolled' : ''
      }`}>
        <nav className="container mx-auto px-4 py-4">
          <div className="flex justify-between items-center">
            <Link to="/" className="logo flex items-center space-x-2 hover:opacity-80 transition-opacity duration-200">
              <Stethoscope className="logo-icon w-7 h-7 text-emerald-500" />
              <span className="logo-text text-xl font-semibold bg-gradient-to-r from-emerald-400 to-emerald-600 text-transparent bg-clip-text">Pocket Doc</span>
            </Link>

            {/* Desktop Navigation */}
            <div className="hidden md:flex items-center space-x-8">
              {isAuthenticated && (
                <div className="flex items-center space-x-8">
                  <NavLink 
                    to="/patients" 
                    className={({ isActive }) => `
                      flex items-center 
                      ${isActive ? 'text-[#3ECF8E]' : ''}
                    `}
                  >
                    <Users className="w-5 h-5 mr-2" />
                    Patients
                  </NavLink>
                  <NavLink 
                    to="/logs" 
                    className={({ isActive }) => `
                      flex items-center 
                      ${isActive ? 'text-[#3ECF8E]' : ''}
                    `}
                  >
                    <ClipboardList className="w-5 h-5 mr-2" />
                    Transcripts
                  </NavLink>
                  <button
                    onClick={handleSignOut}
                    className="btn btn-primary flex items-center"
                  >
                    <LogOut className="w-5 h-5 mr-2" />
                    Sign Out
                  </button>
                </div>
              )}
              {isAuthenticated ? (
                <></>
              ) : (
                <>
                  <NavLink to="/about" className="nav-link">About</NavLink>
                  <NavLink to="/features" className="nav-link">Features</NavLink>
                  <button className="btn btn-primary" onClick={() => setIsAuthModalOpen(true)}>
                    <span>Sign In</span>
                  </button>
                </>
              )}
            </div>

            {/* Mobile menu button */}
            <button 
              onClick={() => setIsOpen(!isOpen)}
              className="md:hidden p-2 rounded-lg hover:bg-white/5 transition-colors"
            >
              {isOpen ? (
                <X size={24} className="text-primary" />
              ) : (
                <Menu size={24} className="text-white" />
              )}
            </button>
          </div>

          {/* Mobile Navigation */}
          <div className={`mobile-menu ${isOpen ? 'block' : 'hidden'}`}>
            <div className="flex flex-col space-y-2 pb-4">
              {isAuthenticated && (
                <div className="space-y-4 px-4">
                  <div className="flex flex-col space-y-2">
                    <NavLink 
                      to="/patients"
                      className={({ isActive }) => `
                        mobile-menu-link 
                        w-full 
                        text-left 
                        flex items-center 
                        ${isActive ? 'text-[#3ECF8E]' : ''}
                      `}
                      onClick={() => setIsOpen(false)}
                    >
                      <Users className="w-5 h-5 inline-block mr-1" />
                      Patients
                    </NavLink>
                    <NavLink 
                      to="/logs"
                      className={({ isActive }) => `
                        mobile-menu-link 
                        w-full 
                        text-left 
                        flex items-center 
                        ${isActive ? 'text-[#3ECF8E]' : ''}
                      `}
                      onClick={() => setIsOpen(false)}
                    >
                      <ClipboardList className="w-5 h-5 inline-block mr-1 nav-link" />
                      Transcripts
                    </NavLink>
                    <button 
                      className="mobile-menu-link w-full text-left flex items-center justify-between" 
                      onClick={() => {
                        handleSignOut();
                        setIsOpen(false);
                      }}
                    >
                      <span className="flex items-center">
                        <LogOut className="w-5 h-5 inline-block mr-1" />
                        Sign Out
                      </span>
                    </button>
                  </div>
                </div>
              )}
              {isAuthenticated ? (
                <></>
              ) : (
                <>
                  <NavLink 
                    to="/about"
                    className={({ isActive }) => `
                      mobile-menu-link 
                      w-full 
                      text-left 
                      px-4
                      ${isActive ? 'text-[#3ECF8E]' : ''}
                    `}
                    onClick={() => setIsOpen(false)}
                  >
                    About
                  </NavLink>
                  <NavLink 
                    to="/features"
                    className={({ isActive }) => `
                      mobile-menu-link 
                      w-full 
                      text-left 
                      px-4
                      ${isActive ? 'text-[#3ECF8E]' : ''}
                    `}
                    onClick={() => setIsOpen(false)}
                  >
                    Features
                  </NavLink>
                  <button 
                    className="mobile-menu-link w-full text-left flex items-center justify-between px-4" 
                    onClick={() => {
                      setIsAuthModalOpen(true);
                      setIsOpen(false);
                    }}
                  >
                    <span>Sign In</span>
                  </button>
                </>
              )}
            </div>
          </div>
        </nav>
      </header>

      <main className="pt-20 min-h-[calc(100vh-64px-320px)]">
        {children}
      </main>

      <footer className="footer">
        <div className="container mx-auto px-4">
          <div className="footer-grid">
            <div className="footer-section">
              <Link to="/" className="logo flex items-center space-x-2 mb-4 hover:opacity-80 transition-opacity duration-200">
                <Stethoscope className="logo-icon w-7 h-7 text-emerald-500" />
                <span className="logo-text text-xl font-semibold bg-gradient-to-r from-emerald-400 to-emerald-600 text-transparent bg-clip-text">Pocket Doc</span>
              </Link>
              <p className="footer-text mb-6">
                Revolutionizing healthcare through intelligent AI solutions. 
                Empowering medical professionals with cutting-edge technology 
                to enhance patient care and streamline documentation.
              </p>
            </div>
            <div className="footer-section">
              <h4 className="footer-heading">Quick Links</h4>
              <ul className="footer-links">
                <li><Link to="/" className="footer-link">Home</Link></li>
                <li><Link to="/features" className="footer-link">Features</Link></li>
                <li><Link to="/about" className="footer-link">About Us</Link></li>
                <li><Link to="/contact" className="footer-link">Contact</Link></li>
              </ul>
            </div>
            <div className="footer-section">
              <h4 className="footer-heading">Resources</h4>
              <ul className="footer-links">
                <li><NavLink to="/documentation" className="footer-link">Documentation</NavLink></li>
                <li><NavLink to="/privacy" className="footer-link">Privacy Policy</NavLink></li>
                <li><NavLink to="/terms" className="footer-link">Terms of Service</NavLink></li>
                <li><NavLink to="/support" className="footer-link">Support</NavLink></li>
              </ul>
            </div>
            <div className="footer-section col-span-full border-t border-dark-lighter pt-8 mt-8 text-center">
              <div className="footer-social-icons mb-4">
                <a href="https://twitter.com/doctorAI" className="footer-social-link" target="_blank" rel="noopener noreferrer">
                  <Twitter className="w-6 h-6" />
                </a>
                <a href="https://linkedin.com/company/doctorAI" className="footer-social-link" target="_blank" rel="noopener noreferrer">
                  <Linkedin className="w-6 h-6" />
                </a>
                <a href="https://github.com/doctorAI" className="footer-social-link" target="_blank" rel="noopener noreferrer">
                  <Github className="w-6 h-6" />
                </a>
              </div>
              <p className="footer-text">
                &copy; {new Date().getFullYear()} DoctorAI. All rights reserved. 
                Powered by advanced AI and machine learning technologies.
              </p>
            </div>
          </div>
        </div>
      </footer>

      {isAuthModalOpen && (
        <AuthModal 
          isOpen={isAuthModalOpen} 
          onClose={() => setIsAuthModalOpen(false)} 
        />
      )}
    </div>
  );
};

export default MainLayout;
