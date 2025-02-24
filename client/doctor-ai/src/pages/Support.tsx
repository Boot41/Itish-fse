import React from 'react';
import { Mail, Phone, MessageCircle, BookOpen, Clock, Users, HelpCircle, Zap } from 'lucide-react';

const ContactCard: React.FC<{
  icon: React.ReactNode;
  title: string;
  description: string;
  action: {
    text: string;
    link: string;
  };
}> = ({ icon, title, description, action }) => (
  <div className="group glass-card p-8 hover:border-[#3ECF8E] transition-all duration-300 ease-in-out transform hover:scale-[1.01] hover:shadow-lg hover:shadow-[#3ECF8E]/20 overflow-hidden will-change-transform transform-gpu">
    <div className="text-[#3ECF8E] mb-6 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
    <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 mb-6 group-hover:text-gray-300 transition-colors duration-300">{description}</p>
    <a 
      href={action.link}
      className="inline-flex items-center text-[#3ECF8E] hover:text-[#4BB563] transition-colors duration-300"
      target="_blank"
      rel="noopener noreferrer"
    >
      {action.text}
      <Zap className="w-5 h-5 ml-2" />
    </a>
  </div>
);

const SupportSection: React.FC<{ 
  icon: React.ReactNode; 
  title: string; 
  description: string; 
  items: string[] 
}> = ({ icon, title, description, items }) => (
  <div className="group glass-card p-8 hover:border-[#3ECF8E] transition-all duration-300 ease-in-out transform hover:scale-[1.01] hover:shadow-lg hover:shadow-[#3ECF8E]/20 overflow-hidden will-change-transform transform-gpu">
    <div className="text-[#3ECF8E] mb-6 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
    <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 mb-6 group-hover:text-gray-300 transition-colors duration-300">{description}</p>
    <ul className="space-y-3">
      {items.map((item, index) => (
        <li key={index} className="flex items-center space-x-2 group-hover:translate-x-1 transition-transform duration-300">
          <Zap className="w-5 h-5 text-[#3ECF8E]" />
          <span className="text-gray-300">{item}</span>
        </li>
      ))}
    </ul>
  </div>
);

const Support: React.FC = () => {
  const contactMethods = [
    {
      icon: <Mail className="w-12 h-12" />,
      title: 'Email Support',
      description: 'Get in touch with our dedicated support team via email.',
      action: {
        text: 'support@doctorai.com',
        link: 'mailto:support@doctorai.com'
      }
    },
    {
      icon: <Phone className="w-12 h-12" />,
      title: 'Phone Support',
      description: 'Speak directly with our support representatives.',
      action: {
        text: '+1 (888) 123-4567',
        link: 'tel:+18881234567'
      }
    },
    {
      icon: <MessageCircle className="w-12 h-12" />,
      title: 'Live Chat',
      description: 'Chat with our support team in real-time.',
      action: {
        text: 'Start Chat',
        link: '#chat'
      }
    }
  ];

  const supportResources = [
    {
      icon: <BookOpen className="w-12 h-12" />,
      title: 'Knowledge Base',
      description: 'Access our comprehensive help articles and guides.',
      items: [
        'Getting started guides',
        'Feature tutorials',
        'Troubleshooting tips',
        'Best practices'
      ]
    },
    {
      icon: <Clock className="w-12 h-12" />,
      title: 'Support Hours',
      description: 'When you can reach our support team.',
      items: [
        'Monday - Friday: 9 AM - 8 PM EST',
        'Saturday: 10 AM - 6 PM EST',
        'Sunday: 12 PM - 5 PM EST',
        'Emergency support available 24/7'
      ]
    },
    {
      icon: <Users className="w-12 h-12" />,
      title: 'Community',
      description: 'Join our community of healthcare professionals.',
      items: [
        'Discussion forums',
        'User groups',
        'Feature requests',
        'Success stories'
      ]
    },
    {
      icon: <HelpCircle className="w-12 h-12" />,
      title: 'FAQs',
      description: 'Quick answers to common questions.',
      items: [
        'Account management',
        'Billing inquiries',
        'Technical issues',
        'Feature usage'
      ]
    }
  ];

  return (
    <div className="min-h-screen bg-[#1C1C1C]">
      <div className="container mx-auto px-4 pt-24 pb-20">
        <div className="text-center max-w-4xl mx-auto mb-16">
          <div 
            className="inline-flex items-center bg-[#2a2a2a] border border-[#3a3a3a] rounded-full px-6 py-2 text-sm mb-8 text-white transform hover:scale-[1.01] transition-all duration-300 cursor-pointer will-change-transform transform-gpu"
          >
            <span className="bg-gradient-to-r from-[#3ECF8E] via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text font-semibold">
              Support • Help Center
            </span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Customer Support
          </h1>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            We're here to help! Get in touch with our support team or browse our help resources
          </p>
        </div>

        {/* Contact Methods */}
        <div className="grid md:grid-cols-3 gap-8 mb-16 animate-fade-in-up">
          {contactMethods.map((method, index) => (
            <div key={index} className="transform transition-all duration-500" style={{ animationDelay: `${index * 150}ms` }}>
              <ContactCard 
                icon={method.icon} 
                title={method.title} 
                description={method.description} 
                action={method.action}
              />
            </div>
          ))}
        </div>

        {/* Support Resources */}
        <div className="grid md:grid-cols-2 gap-8 animate-fade-in-up">
          {supportResources.map((resource, index) => (
            <div key={index} className="transform transition-all duration-500" style={{ animationDelay: `${(index + 3) * 150}ms` }}>
              <SupportSection 
                icon={resource.icon} 
                title={resource.title} 
                description={resource.description} 
                items={resource.items} 
              />
            </div>
          ))}
        </div>

        {/* Emergency Support Notice */}
        <div className="mt-16 text-center">
          <p className="text-gray-400 max-w-3xl mx-auto leading-relaxed">
            For urgent matters requiring immediate attention, please call our emergency support line:
            <br />
            <a 
              href="tel:+18881234567" 
              className="text-[#3ECF8E] hover:text-[#4BB563] transition-colors duration-300 font-semibold"
            >
              +1 (888) 123-4567
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Support;
