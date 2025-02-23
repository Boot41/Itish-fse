import React, { useState, useRef, useEffect } from 'react';
import { Send, Bot, Maximize2, Minimize2, X } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { chatApi } from '../../api/chat';

interface Message {
  role: 'user' | 'assistant';
  content: string;
}

interface ChatWindowProps {
  className?: string;
  onClose?: () => void;
}

const ChatWindow: React.FC<ChatWindowProps> = ({ className, onClose }) => {
  const [isFullScreen, setIsFullScreen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([
    {
      role: 'assistant',
      content: 'Hello! I\'m an AI assistant designed to help with medical transcription and analysis. \nHow can I assist you today?'
    }
  ]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || isLoading) return;

    const userMessage = input.trim();
    setInput('');
    setMessages(prev => [...prev, { role: 'user', content: userMessage }]);
    setIsLoading(true);

    try {
      const data = await chatApi.sendMessage(userMessage);
      setMessages(prev => [...prev, { role: 'assistant', content: data.response }]);
    } catch (error) {
      console.error('Failed to get response:', error);
      setMessages(prev => [...prev, {
        role: 'assistant',
        content: 'Sorry, I encountered an error. Please try again.'
      }]);
    } finally {
      setIsLoading(false);
    }
  };

  console.log({className})

  return (
    <AnimatePresence mode="sync">
      {/* Chat Window Container */}
      <motion.div 
        key="chat-window"
        className={`${!isFullScreen ? 'relative w-full h-[300px]' : 'fixed inset-0 z-50'}`}
      >
        {/* Overlay (only in fullscreen) */}
        {isFullScreen && (
          <motion.div
            key="overlay"
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2, ease: 'easeInOut' }}
            onClick={() => setIsFullScreen(false)}
          />
        )}
        
        {/* Chat Window */}
        <motion.div 
          key="chat-content"
          className={`${isFullScreen ? 'absolute top-[55%] left-1/2 -translate-x-1/2 -translate-y-1/2 w-[80vw] h-[80vh]' : 'w-full h-full'} bg-[#2a2a2a] rounded-2xl border border-[#3a3a3a] overflow-hidden flex flex-col`}
          layout
          initial={false}
          animate={{ 
            scale: 1,
            opacity: 1,
            y: 0,
            transition: {
              type: 'spring',
              stiffness: 400,
              damping: 35,
              mass: 0.5
            }
          }}
          exit={{ 
            scale: 0.98,
            opacity: 0,
            transition: {
              duration: 0.2,
              ease: 'easeInOut'
            }
          }}
        >
          {/* Header */}
          <div className="p-4 border-b border-[#3a3a3a] flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Bot className="w-5 h-5 text-emerald-500" />
              <h3 className="font-semibold text-white">AI Assistant</h3>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setIsFullScreen(!isFullScreen)}
                className="p-2 hover:bg-[#3a3a3a] rounded-lg transition-colors duration-200"
              >
                {isFullScreen ? (
                  <Minimize2 className="w-4 h-4 text-gray-400" />
                ) : (
                  <Maximize2 className="w-4 h-4 text-gray-400" />
                )}
              </button>
              {isFullScreen && (
                <button
                  onClick={onClose}
                  className="p-2 hover:bg-[#3a3a3a] rounded-lg transition-colors duration-200"
                >
                  <X className="w-4 h-4 text-gray-400" />
                </button>
              )}
            </div>
          </div>

          {/* Messages */}
          <motion.div 
            className="flex-1 overflow-y-auto p-4 space-y-4"
            layout
          >
            {messages.map((message, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 10, scale: 0.95 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                transition={{ 
                  type: 'spring',
                  stiffness: 500,
                  damping: 25,
                  mass: 0.5
                }}
                className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <motion.div
                  layout
                  className={`max-w-[80%] rounded-lg px-4 py-2 ${
                    message.role === 'user'
                      ? 'bg-emerald-500 text-black'
                      : 'bg-[#333] text-white'
                  }`}
                >
                  {message.content}
                </motion.div>
              </motion.div>
            ))}
            <div ref={messagesEndRef} />
          </motion.div>

          {/* Input */}
          <form onSubmit={handleSubmit} className="p-4 border-t border-[#3a3a3a]">
            <div className="flex gap-2">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Type your message..."
                className="flex-1 bg-[#333] text-white rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-emerald-500"
                disabled={isLoading}
              />
              <button
                type="submit"
                disabled={isLoading}
                className="bg-emerald-500 text-black rounded-lg px-4 py-2 hover:bg-emerald-600 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Send className="w-5 h-5" />
              </button>
              
            </div>
          </form>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
};

export default ChatWindow;
