import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Bot, MessageCircle, Search, User } from 'lucide-react';
import { chatApi } from '../../api/chat';
import { useAuthStore } from '../../stores/authStore';

interface Suggestion {
  text: string;
  type: 'medicine' | 'general' | 'user' | 'error';
  details?: string;
  patientName?: string;
}

interface ExpandedState {
  [key: number]: boolean;
}

const MedicalSuggestions: React.FC = () => {
  const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
  const [loading, setLoading] = useState(true);
  const [input, setInput] = useState('');
  const [expandedStates, setExpandedStates] = useState<ExpandedState>({});

  const fetchMedicineSuggestions = async (query: string) => {
    setLoading(true);
    try {
      const prompt = `Given the medical condition or symptoms "${query}", suggest 6 potential medications or treatments. For each suggestion:
1. Start with a dash (-)
2. First line: Name of medication/treatment
3. Second line: Brief description of how it works and common usage

Keep each suggestion concise but informative.`;

      const { response } = await chatApi.sendMessage(prompt);

      // Parse the response into individual suggestions
      const suggestions = response.split('\n');
      const newSuggestions: Suggestion[] = [];

      for (let i = 0; i < suggestions.length; i++) {
        const line = suggestions[i].trim();
        if (line.startsWith('-')) {
          const text = line.substring(1).trim();
          const details = suggestions[i + 1]?.trim();
          newSuggestions.push({
            text,
            details,
            type: 'medicine' as const,
          });
          i++; // Skip the details line
        }
      }

      setSuggestions(newSuggestions);
    } catch (error) {
      handleError(error, 'fetching medicine suggestions');
    } finally {
      setLoading(false);
    }
  };

  const extractUserQuery = async (query: string) => {
    setLoading(true);
    try {
      const prompt = `Analyze this query: "${query}"
Extract the following information:
1. Patient name (if mentioned)
2. Main symptoms or condition
3. Any specific concerns

Format the response as:
Patient: [name or "Not mentioned"]
Symptoms: [main symptoms/condition]
Concerns: [specific concerns or "None mentioned"]`;

      const { response } = await chatApi.sendMessage(prompt);
      const lines = response.split('\n');
      
      let patientName = "Not mentioned";
      let symptoms = "";
      let concerns = "";

      lines.forEach((line: string) => {
        if (line.startsWith('Patient:')) patientName = line.replace('Patient:', '').trim();
        if (line.startsWith('Symptoms:')) symptoms = line.replace('Symptoms:', '').trim();
        if (line.startsWith('Concerns:')) concerns = line.replace('Concerns:', '').trim();
      });

      // Add user query as a suggestion
      setSuggestions(prev => [{
        text: symptoms,
        type: 'user',
        details: concerns !== "None mentioned" ? concerns : undefined,
        patientName: patientName !== "Not mentioned" ? patientName : undefined
      }, ...prev]);

      // Fetch medicine suggestions based on symptoms
      await fetchMedicineSuggestions(symptoms);

      // Fetch patient history
      if (patientName !== "Not mentioned") {
        await fetchPatientHistory(patientName);
      }
    } catch (error) {
      handleError(error, 'processing user query');
    } finally {
      setLoading(false);
    }
  };

  const fetchMedicalTrivia = async () => {
    setLoading(true);
    try {
      const prompt = `Share 6 fascinating medical facts or trivia. For each fact:
1. Start with a dash (-)
2. First line: Brief, intriguing medical fact as a title
3. Second line: Detailed explanation with interesting context

Make them engaging and educational, focusing on surprising or lesser-known medical discoveries and phenomena.`;

      const { response } = await chatApi.sendMessage(prompt);

      const facts = response.split('\n');
      const triviaFacts: Suggestion[] = [];

      for (let i = 0; i < facts.length; i++) {
        const line = facts[i].trim();
        if (line.startsWith('-')) {
          const text = line.substring(1).trim();
          const details = facts[i + 1]?.trim();
          triviaFacts.push({
            text,
            details,
            type: 'general' as const,
          });
          i++; // Skip the details line
        }
      }

      setSuggestions(triviaFacts);
    } catch (error) {
      handleError(error, 'fetching medical trivia');
    } finally {
      setLoading(false);
    }
  };

  const fetchPatientHistory = async (patientName: string) => {
    setLoading(true);
    try {
      // Get doctor's email from auth store
      const email = await useAuthStore.getState().getDecryptedEmail();
      if (!email) {
        throw new Error('Not authenticated');
      }

      console.log('Fetching transcripts for patient:', patientName, 'doctor:', email);
      
      // Fetch transcripts from backend
      const response = await fetch(`http://localhost:8080/patients/transcript`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: patientName,
          email: email
        })
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch transcripts: ${response.statusText}`);
      }
      
      const transcripts = await response.json();
      console.log('Received transcripts:', transcripts);

      // Generate summary using Groq
      const prompt = `Analyze these medical transcripts and provide a comprehensive summary:
${JSON.stringify(transcripts, null, 2)}

Format the summary as:
Medical History: [key medical events and conditions]
Current Medications: [list of current medications]
Allergies: [known allergies]
Recent Visits: [summary of recent visits]
Key Observations: [important medical observations]
Recommendations: [any recommended follow-ups or actions]`;

      const { response: summary } = await chatApi.sendMessage(prompt);
      console.log('Generated summary:', summary);

      // Parse the summary
      const lines = summary.split('\n');
      let history = '';
      let medications = '';
      let allergies = '';
      let visits = '';
      let observations = '';
      let recommendations = '';

      lines.forEach((line: string) => {
        if (line.startsWith('Medical History:')) history = line.replace('Medical History:', '').trim();
        if (line.startsWith('Current Medications:')) medications = line.replace('Current Medications:', '').trim();
        if (line.startsWith('Allergies:')) allergies = line.replace('Allergies:', '').trim();
        if (line.startsWith('Recent Visits:')) visits = line.replace('Recent Visits:', '').trim();
        if (line.startsWith('Key Observations:')) observations = line.replace('Key Observations:', '').trim();
        if (line.startsWith('Recommendations:')) recommendations = line.replace('Recommendations:', '').trim();
      });

      // Add the summary as a suggestion
      setSuggestions(prev => [{
        text: `Patient Summary: ${patientName}`,
        type: 'user',
        details: [
          history && `📋 Medical History:\n\n${history}`,
          medications && `💊 Current Medications:\n\n${medications}`,
          allergies && `⚠️ Allergies:\n\n${allergies}`,
          visits && `🏥 Recent Visits:\n\n${visits}`,
          observations && `👁️ Key Observations:\n\n${observations}`,
          recommendations && `📝 Recommendations:\n\n${recommendations}`
        ].filter(Boolean).join('\n\n'),
        patientName
      }, ...prev]);

    } catch (error) {
      handleError(error, 'fetching patient history');
    } finally {
      setLoading(false);
    }
  };

  const handleError = (error: any, context: string) => {
    console.error(`Error in ${context}:`, error);
    setSuggestions(prev => [{
      text: `Error in ${context}`,
      type: 'error',
      details: error.message || 'An unexpected error occurred'
    }, ...prev]);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim()) return;

    console.log('=== Starting Query Processing ===');
    console.log('Input:', input);
    console.time('queryProcessing');

    try {
      // Check for "tell me about [name]" pattern first
      const tellMeAboutMatch = input.match(/tell me about\s+["']?([^"']+)["']?/i);
      if (tellMeAboutMatch) {
        const patientName = tellMeAboutMatch[1].trim();
        console.log('Detected "tell me about" pattern:', {
          originalQuery: input,
          patientName,
          timestamp: new Date().toISOString()
        });
        await fetchPatientHistory(patientName);
      } else {
        await extractUserQuery(input);
      }
      setInput('');
    } catch (error) {
      handleError(error, 'processing form submission');
    }

    console.timeEnd('queryProcessing');
    console.log('=== Query Processing Complete ===');
  };

  useEffect(() => {
    fetchMedicalTrivia();
  }, []);

  return (
    <div className="w-[500px] h-[calc(100vh-16rem)] sticky top-4">
      <motion.div
        className="h-full bg-[#2a2a2a] rounded-xl border border-[#3a3a3a] overflow-hidden flex flex-col shadow-2xl shadow-emerald-500/5"
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.3 }}
      >
        {/* Header */}
        <div className="px-6 py-5 border-b border-[#3a3a3a] flex items-center justify-between bg-[#2a2a2a] bg-gradient-to-r from-[#2a2a2a] to-[#2f2f2f]">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-emerald-500/10">
              <Bot className="w-5 h-5 text-emerald-500" />
            </div>
            <h3 className="font-semibold text-white text-lg">Medical Assistant</h3>
          </div>
        </div>

        {/* Search Input */}
        <div className="p-6 border-b border-[#3a3a3a] bg-gradient-to-r from-[#2a2a2a] to-[#2f2f2f]">
          <div className="relative group">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && input.trim()) {
                  handleSubmit(e);
                }
              }}
              placeholder="Describe symptoms or ask about a patient..."
              className="w-full bg-[#333] text-white rounded-xl pl-11 pr-4 py-3 text-sm border border-[#3a3a3a] group-hover:border-[#4a4a4a] focus:border-emerald-500/50 focus:outline-none transition-all duration-300 placeholder:text-gray-500"
              autoComplete="off"
            />
            <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500 group-hover:text-gray-400 transition-colors duration-300" />
          </div>
        </div>

        {/* Suggestions Area */}
        <div className="flex-1 overflow-y-auto px-6 py-4 space-y-4 bg-gradient-to-b from-[#2a2a2a] to-[#2f2f2f]">
          {loading ? (
            <div className="text-gray-400 flex items-center justify-center h-full">
              <div className="flex items-center gap-3">
                <div className="animate-spin rounded-full h-5 w-5 border-2 border-emerald-500 border-t-transparent"></div>
                <span>Loading suggestions...</span>
              </div>
            </div>
          ) : (
            suggestions.map((suggestion, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                className="flex items-start gap-3 group relative"
                onClick={() => {
                  if (suggestion.type === 'general') {
                    setInput(suggestion.text);
                    extractUserQuery(suggestion.text);
                  } else {
                    setExpandedStates((prev) => ({ ...prev, [index]: !prev[index] }));
                  }
                }}
              >
                <motion.div
                  className="mt-1.5"
                  animate={{
                    scale: expandedStates[index] ? 1.1 : 1,
                    rotate: expandedStates[index] ? 180 : 0,
                  }}
                  transition={{ duration: 0.3 }}
                >
                  <div
                    className={`p-1.5 rounded-lg transition-all duration-300 ${
                      suggestion.type === 'medicine'
                        ? 'bg-emerald-500/10 group-hover:bg-emerald-500/20'
                        : suggestion.type === 'user'
                        ? 'bg-blue-500/10 group-hover:bg-blue-500/20'
                        : suggestion.type === 'error'
                        ? 'bg-red-500/10 group-hover:bg-red-500/20'
                        : 'bg-[#333] group-hover:bg-[#3a3a3a]'
                    }`}
                  >
                    {suggestion.type === 'medicine' ? (
                      <MessageCircle
                        className="w-3.5 h-3.5 transition-colors duration-300 text-emerald-500"
                      />
                    ) : suggestion.type === 'user' ? (
                      <User
                        className="w-3.5 h-3.5 transition-colors duration-300 text-blue-500"
                      />
                    ) : suggestion.type === 'error' ? (
                      <div
                        className="w-3.5 h-3.5 transition-colors duration-300 text-red-500"
                      >
                        ⚠️
                      </div>
                    ) : (
                      <Search
                        className="w-3.5 h-3.5 transition-colors duration-300 text-emerald-400 group-hover:text-emerald-300"
                      />
                    )}
                  </div>
                </motion.div>
                <motion.div
                  className={`flex-1 px-4 py-3 rounded-xl text-white text-sm cursor-pointer backdrop-blur-sm ${
                    suggestion.type === 'general' ? 'hover:scale-[1.02]' : ''
                  } ${
                    suggestion.type === 'medicine'
                      ? 'bg-emerald-500/10 hover:bg-emerald-500/15'
                      : suggestion.type === 'user'
                      ? 'bg-blue-500/10 hover:bg-blue-500/15'
                      : suggestion.type === 'error'
                      ? 'bg-red-500/10 hover:bg-red-500/15'
                      : 'bg-[#333] hover:bg-[#3a3a3a]'
                  } ${expandedStates[index] ? 'ring-1 ring-emerald-500/20' : ''}`}
                  layout
                  transition={{
                    layout: { duration: 0.3 },
                    default: { ease: 'easeInOut' },
                  }}
                >
                  <div className="flex items-center justify-between">
                    <div className="font-medium">
                      {suggestion.text}
                      {suggestion.patientName && (
                        <span className="ml-2 text-xs text-blue-400">
                          Patient: {suggestion.patientName}
                        </span>
                      )}
                    </div>
                    <motion.div
                      className="text-xs text-emerald-500/70 ml-2"
                      animate={{ opacity: expandedStates[index] ? 1 : 0.5 }}
                      transition={{ duration: 0.2 }}
                    >
                      {expandedStates[index] ? 'Less' : 'More'}
                    </motion.div>
                  </div>

                  <AnimatePresence>
                    {expandedStates[index] && suggestion.details && (
                      <motion.div
                        initial={{ opacity: 0, height: 0, y: -10 }}
                        animate={{ opacity: 1, height: 'auto', y: 0 }}
                        exit={{ opacity: 0, height: 0, y: -10 }}
                        transition={{ duration: 0.3, ease: 'easeInOut' }}
                      >
                        <div className="mt-2 pt-2 border-t border-[#3a3a3a] text-gray-400 text-xs leading-relaxed">
                          {suggestion.details}
                        </div>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </motion.div>
              </motion.div>
            ))
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-[#3a3a3a] bg-gradient-to-r from-[#2a2a2a] to-[#2f2f2f]">
          <div className="flex items-center justify-center gap-2 text-xs text-gray-400">
            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
            Enter symptoms or patient details for AI assistance
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default MedicalSuggestions;
