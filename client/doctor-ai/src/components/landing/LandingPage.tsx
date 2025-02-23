import { FC, useState } from 'react';
import { Link } from 'react-router-dom';
import { Mic, Users, BarChart, ArrowRight } from 'lucide-react';
import AuthModal from '../auth/AuthModal';

const LandingPage: FC = () => {
  const [isAuthModalOpen, setIsAuthModalOpen] = useState(false);

  const openAuthModal = () => setIsAuthModalOpen(true);
  const closeAuthModal = () => setIsAuthModalOpen(false);

  return (
    <div className="bg-[#1C1C1C] min-h-screen">
      {/* Auth Modal */}
      <AuthModal 
        isOpen={isAuthModalOpen} 
        onClose={closeAuthModal} 
      />

      {/* Hero Section */}
      <div className="container mx-auto px-4 pt-24 pb-20">
        <div className="text-center max-w-4xl mx-auto">
          <div 
            className="inline-flex items-center bg-[#2a2a2a] border border-[#3a3a3a] rounded-full px-6 py-2 text-sm mb-8 text-white transform hover:scale-105 transition-all duration-300 cursor-pointer"
          >
            <span className="bg-gradient-to-r from-[#3ECF8E] via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text font-semibold">
              Medical Documentation • Simplified
            </span>
          </div>
          
          <h1 className="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-105 transition-all duration-500">
            Build in a weekend
          </h1>
          <h2 className="text-4xl md:text-6xl font-bold mb-8 text-[#3ECF8E] opacity-90 transform hover:scale-105 transition-all duration-500">
            Scale to millions
          </h2>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed">
            Transform your medical practice with AI-powered transcription and intelligent patient management
          </p>
          <div className="flex gap-6 justify-center">
            <button 
              onClick={openAuthModal}
              className="group bg-gradient-to-r from-[#3ECF8E] to-[#4BB563] text-gray-950 px-8 py-4 rounded-md font-medium transition-all duration-300 hover:shadow-lg hover:shadow-[#3ECF8E]/20 transform hover:-translate-y-1"
            >
              Get Started
              <ArrowRight className="inline-block ml-2 w-5 h-5 transform group-hover:translate-x-1 transition-transform duration-300" />
            </button>
            <Link 
              to="/about" 
              className="group bg-[#2a2a2a] border border-[#3a3a3a] hover:border-[#3ECF8E] text-white px-8 py-4 rounded-md font-medium transition-all duration-300 hover:shadow-lg hover:shadow-[#3ECF8E]/10 transform hover:-translate-y-1"
            >
              Learn More
              <ArrowRight className="inline-block ml-2 w-5 h-5 transform group-hover:translate-x-1 transition-transform duration-300" />
            </Link>
          </div>
        </div>

        {/* Features Section */}
        <div className="mt-32 grid md:grid-cols-3 gap-8">
          <FeatureCard 
            title="AI Transcription"
            description="Convert medical conversations into accurate, structured reports instantly"
            icon={<Mic className="w-8 h-8" />}
          />
          <FeatureCard 
            title="Patient Management"
            description="Efficiently organize and track patient information and history"
            icon={<Users className="w-8 h-8" />}
          />
          <FeatureCard 
            title="Smart Analytics"
            description="Gain insights from your practice with detailed statistics and reports"
            icon={<BarChart className="w-8 h-8" />}
          />
        </div>

        {/* How It Works Section */}
        <div className="mt-32 text-center">
          <h2 className="text-3xl font-bold mb-16 text-white bg-gradient-to-r from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text">
            How It Works
          </h2>
          <div className="grid md:grid-cols-4 gap-8">
            <StepCard 
              step={1}
              title="Record"
              description="Record your patient consultations"
            />
            <StepCard 
              step={2}
              title="Transcribe"
              description="AI converts speech to text"
            />
            <StepCard 
              step={3}
              title="Analyze"
              description="Get AI-enhanced medical reports"
            />
            <StepCard 
              step={4}
              title="Manage"
              description="Track patient progress easily"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

const FeatureCard: FC<{ title: string; description: string; icon: React.ReactNode }> = ({ 
  title, 
  description, 
  icon 
}) => (
  <div className="group bg-[#2a2a2a] p-8 rounded-xl border border-[#3a3a3a] hover:border-[#3ECF8E] transition-all duration-300 hover:shadow-lg hover:shadow-[#3ECF8E]/10 transform hover:-translate-y-2">
    <div className="text-[#3ECF8E] mb-4 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
    <h3 className="text-xl font-bold text-white mb-3 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 leading-relaxed">{description}</p>
  </div>
);

const StepCard: FC<{ step: number; title: string; description: string }> = ({ 
  step, 
  title, 
  description 
}) => (
  <div className="group bg-[#2a2a2a] p-8 rounded-xl border border-[#3a3a3a] hover:border-[#3ECF8E] transition-all duration-300 hover:shadow-lg hover:shadow-[#3ECF8E]/10 transform hover:-translate-y-2">
    <div className="text-4xl font-bold text-[#3ECF8E] mb-4 transform group-hover:scale-110 transition-transform duration-300">{step}</div>
    <h3 className="text-xl font-bold text-white mb-3 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 leading-relaxed">{description}</p>
  </div>
);

export default LandingPage;
