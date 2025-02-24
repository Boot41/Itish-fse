import React, { useState } from 'react';
import { Mail, Phone, MapPin, Send } from 'lucide-react';

interface FormData {
  name: string;
  email: string;
  subject: string;
  message: string;
}

const ContactInfo: React.FC<{
  icon: React.ReactNode;
  title: string;
  content: string;
  link?: string;
}> = ({ icon, title, content, link }) => (
  <div className="flex items-start space-x-4">
    <div className="text-[#3ECF8E] transform group-hover:scale-110 transition-transform duration-300">
      {icon}
    </div>
    <div>
      <h3 className="text-white font-semibold mb-1">{title}</h3>
      {link ? (
        <a 
          href={link}
          className="text-gray-400 hover:text-[#3ECF8E] transition-colors duration-300"
        >
          {content}
        </a>
      ) : (
        <p className="text-gray-400">{content}</p>
      )}
    </div>
  </div>
);

const Contact: React.FC = () => {
  const [formData, setFormData] = useState<FormData>({
    name: '',
    email: '',
    subject: '',
    message: ''
  });

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState<{
    type: 'success' | 'error' | null;
    message: string;
  }>({ type: null, message: '' });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    try {
      // Here you would typically make an API call to your backend
      // For now, we'll simulate a successful submission
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setSubmitStatus({
        type: 'success',
        message: 'Thank you for your message. We\'ll get back to you soon!'
      });
      
      // Clear the form
      setFormData({
        name: '',
        email: '',
        subject: '',
        message: ''
      });
    } catch (error) {
      setSubmitStatus({
        type: 'error',
        message: 'Something went wrong. Please try again later.'
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#1C1C1C]">
      <div className="container mx-auto px-4 pt-24 pb-20">
        <div className="text-center max-w-4xl mx-auto mb-16">
          <div 
            className="inline-flex items-center bg-[#2a2a2a] border border-[#3a3a3a] rounded-full px-6 py-2 text-sm mb-8 text-white transform hover:scale-[1.01] transition-all duration-300 cursor-pointer will-change-transform transform-gpu"
          >
            <span className="bg-gradient-to-r from-[#3ECF8E] via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text font-semibold">
              Contact • Get in Touch
            </span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Contact Us
          </h1>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            Have questions? We'd love to hear from you. Send us a message and we'll respond as soon as possible
          </p>
        </div>

        <div className="max-w-6xl mx-auto grid md:grid-cols-3 gap-12">
          {/* Contact Form */}
          <div className="md:col-span-2">
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="name" className="block text-white mb-2">Name</label>
                  <input
                    type="text"
                    id="name"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    className="w-full px-4 py-3 rounded-lg bg-[#2a2a2a] border border-[#3a3a3a] text-white focus:outline-none focus:border-[#3ECF8E] transition-colors duration-300"
                    placeholder="Your name"
                  />
                </div>
                <div>
                  <label htmlFor="email" className="block text-white mb-2">Email</label>
                  <input
                    type="email"
                    id="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    className="w-full px-4 py-3 rounded-lg bg-[#2a2a2a] border border-[#3a3a3a] text-white focus:outline-none focus:border-[#3ECF8E] transition-colors duration-300"
                    placeholder="your@email.com"
                  />
                </div>
              </div>
              
              <div>
                <label htmlFor="subject" className="block text-white mb-2">Subject</label>
                <input
                  type="text"
                  id="subject"
                  name="subject"
                  value={formData.subject}
                  onChange={handleChange}
                  required
                  className="w-full px-4 py-3 rounded-lg bg-[#2a2a2a] border border-[#3a3a3a] text-white focus:outline-none focus:border-[#3ECF8E] transition-colors duration-300"
                  placeholder="What is this regarding?"
                />
              </div>

              <div>
                <label htmlFor="message" className="block text-white mb-2">Message</label>
                <textarea
                  id="message"
                  name="message"
                  value={formData.message}
                  onChange={handleChange}
                  required
                  rows={6}
                  className="w-full px-4 py-3 rounded-lg bg-[#2a2a2a] border border-[#3a3a3a] text-white focus:outline-none focus:border-[#3ECF8E] transition-colors duration-300 resize-none"
                  placeholder="Your message..."
                />
              </div>

              {submitStatus.type && (
                <div className={`p-4 rounded-lg ${
                  submitStatus.type === 'success' ? 'bg-green-900/50 text-green-300' : 'bg-red-900/50 text-red-300'
                }`}>
                  {submitStatus.message}
                </div>
              )}

              <button
                type="submit"
                disabled={isSubmitting}
                className={`
                  w-full md:w-auto px-8 py-4 rounded-lg bg-gradient-to-r from-[#3ECF8E] to-[#4BB563]
                  text-white font-semibold flex items-center justify-center space-x-2
                  transform hover:scale-[1.01] transition-all duration-300
                  disabled:opacity-50 disabled:cursor-not-allowed
                `}
              >
                <span>Send Message</span>
                <Send className="w-5 h-5" />
              </button>
            </form>
          </div>

          {/* Contact Information */}
          <div className="space-y-8">
            <div className="glass-card p-8">
              <h2 className="text-2xl font-bold text-white mb-6">Contact Information</h2>
              <div className="space-y-6">
                <ContactInfo
                  icon={<Mail className="w-6 h-6" />}
                  title="Email"
                  content="contact@doctorai.com"
                  link="mailto:contact@doctorai.com"
                />
                <ContactInfo
                  icon={<Phone className="w-6 h-6" />}
                  title="Phone"
                  content="+1 (888) 123-4567"
                  link="tel:+18881234567"
                />
                <ContactInfo
                  icon={<MapPin className="w-6 h-6" />}
                  title="Address"
                  content="123 Innovation Drive, Silicon Valley, CA 94025"
                />
              </div>
            </div>

            <div className="glass-card p-8">
              <h2 className="text-2xl font-bold text-white mb-6">Office Hours</h2>
              <div className="space-y-2 text-gray-400">
                <p>Monday - Friday: 9 AM - 6 PM PST</p>
                <p>Saturday: 10 AM - 4 PM PST</p>
                <p>Sunday: Closed</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Contact;
