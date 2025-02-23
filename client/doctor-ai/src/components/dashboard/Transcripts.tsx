import React, { useState, useEffect } from 'react';
import { RefreshCw, AlertCircle, Download, Search, Pencil, Trash2 } from 'lucide-react';
import { Spinner } from '../ui/Spinner';
import { EmptyState } from '../ui/EmptyState';
import { useRefreshStore } from '../../stores/refreshStore';
import { useAuthStore } from '../../stores/authStore';
import { toast } from 'sonner';
import { useLocation } from 'react-router-dom';
import { AnimatePresence } from 'framer-motion';
import EditTranscriptionModal from './EditTranscriptionModal';
import ConfirmationModal from '../../pages/ConfirmationModal';
import MedicalSuggestions from './MedicalSuggestions';

export interface Transcript {
  id: string;
  patientId: string;
  patientName: string;
  timestamp: string;
  type: string;
  status: string;
  transcriptUrl?: string;
}

interface TranscriptCardProps {
  transcript: Transcript;
  image: string;
  onDownload: (transcript: Transcript) => void;
  onEdit: (transcript: Transcript) => void;
  onDelete: (transcript: Transcript) => void;
}

const TranscriptCard: React.FC<TranscriptCardProps> = ({ transcript, image, onDownload, onEdit, onDelete }) => {
  // Format date and time separately for better readability
  const date = new Date(transcript.timestamp);
  const formattedDate = date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
  const formattedTime = date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  });

  return (
    <div 
      className="relative w-[280px] h-[360px] bg-[#2a2a2a] rounded-lg border border-[#3a3a3a] group transition-all duration-300 ease-in-out transform hover:scale-[1.02] hover:border-emerald-500/50 cursor-pointer overflow-hidden"
    >
      {/* Top gradient bar */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-emerald-400 to-emerald-600" />
      
      {/* Content */}
      <div className="relative h-full flex flex-col p-4">
        {/* Avatar section */}
        <div className="flex justify-center mb-2">
          <div className="relative">
            <img alt="Patient Profile" className="w-16 h-16 rounded-full border-4 border-[#3ECF8E]/30 hover:border-[#3ECF8E] transition-all" src={image} />
            <div className="absolute -inset-1 bg-emerald-500/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300 blur-xl rounded-full" />
          </div>
        </div>

        {/* Transcript info */}
        <div className="text-center space-y-2 flex-1">
          <h3 className="text-lg font-bold text-white group-hover:text-emerald-400 transition-colors duration-300 truncate px-2">
            {transcript.patientName}
          </h3>

          {/* Status badge */}
          <div className="inline-flex items-center px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm">
            {transcript.type}
          </div>

          {/* Info grid */}
          <div className="grid grid-cols-2 gap-2 mt-2">
            <div className="bg-[#333] rounded-lg p-2 border border-[#3a3a3a] group-hover:border-emerald-500/20 transition-colors duration-300">
              <div className="text-xs text-gray-400 mb-0.5">Date</div>
              <div className="text-white font-medium text-sm">{formattedDate}</div>
            </div>
            <div className="bg-[#333] rounded-lg p-2 border border-[#3a3a3a] group-hover:border-emerald-500/20 transition-colors duration-300">
              <div className="text-xs text-gray-400 mb-0.5">Time</div>
              <div className="text-white font-medium text-sm">{formattedTime}</div>
            </div>
          </div>
        </div>

        {/* Action buttons - Fixed at bottom with glass effect */}
        <div className="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-[#2a2a2a] via-[#2a2a2a] to-transparent">
          <div className="grid grid-cols-3 gap-2 px-2">
            <button
              onClick={(e) => {
                e.stopPropagation();
                onDownload(transcript);
              }}
              className="flex items-center justify-center py-2 rounded-lg bg-[#333] hover:bg-emerald-500/10 border border-[#3a3a3a] hover:border-emerald-500/30 transition-all duration-300 group/download"
            >
              <Download className="w-4 h-4 text-emerald-500 transition-transform duration-300 group-hover/download:scale-110" />
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onEdit(transcript);
              }}
              className="flex items-center justify-center py-2 rounded-lg bg-[#333] hover:bg-emerald-500/10 border border-[#3a3a3a] hover:border-emerald-500/30 transition-all duration-300 group/edit"
            >
              <Pencil className="w-4 h-4 text-emerald-500 transition-transform duration-300 group-hover/edit:scale-110" />
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onDelete(transcript);
              }}
              className="flex items-center justify-center py-2 rounded-lg bg-[#333] hover:bg-red-500/10 border border-[#3a3a3a] hover:border-red-500/30 transition-all duration-300 group/delete"
            >
              <Trash2 className="w-4 h-4 text-red-400 transition-transform duration-300 group-hover/delete:scale-110" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

const Transcripts: React.FC = () => {
  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const searchQueryFromParams = query.get('search');

  const [transcripts, setTranscripts] = useState<Transcript[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const { refreshTranscripts, triggerTranscriptRefresh } = useRefreshStore();
  const { getDecryptedEmail } = useAuthStore();
  const [doctorEmail, setDoctorEmail] = useState<string | null>(null);

  // State for editing transcription
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedTranscriptionId, setSelectedTranscriptionId] = useState<string | null>(null);
  const [selectedTranscriptionData, setSelectedTranscriptionData] = useState<{ text: string; report: string } | null>(null);

  const [transcriptionToDelete, setTranscriptionToDelete] = useState<Transcript | null>(null);

  // 1. Set search query from URL params (if available)
  useEffect(() => {
    if (searchQueryFromParams) {
      setSearchQuery(searchQueryFromParams);
    }
  }, [searchQueryFromParams]);

  // 2. Decrypt the doctor's email on mount
  useEffect(() => {
    const decryptEmail = async () => {
      try {
        const decrypted = await getDecryptedEmail();
        setDoctorEmail(decrypted);
      } catch (error) {
        console.error('Failed to decrypt email:', error);
        toast.error('Failed to retrieve doctor email');
      }
    };
    decryptEmail();
  }, [getDecryptedEmail]);

  // 3. Fetch transcripts when doctorEmail or refreshTranscripts changes
  useEffect(() => {
    const fetchTranscripts = async () => {
      if (!doctorEmail) {
        setError('No doctor email found');
        setLoading(false);
        return;
      }
      setLoading(true);
      try {
        const response = await fetch('http://localhost:8080/transcription/', {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
            'Origin': import.meta.env.VITE_FRONTEND_URL,
            'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
          },
          body: JSON.stringify({ email: doctorEmail }),
        });
        if (!response.ok) {
          throw new Error(`Failed to fetch transcripts: ${await response.text()}`);
        }
        const data = await response.json();
        const transcriptsToSet = data.transcripts || [];
        const patientNames = data.patients || [];
        if (transcriptsToSet.length > 0) {
          const formattedTranscripts = transcriptsToSet.map((transcript: any, index: number) => ({
            id: transcript.id,
            patientId: transcript.patientId,
            patientName: patientNames[index] || 'Unknown Patient',
            timestamp: transcript.timestamp,
            type: 'Consultation',
            status: 'completed',
            transcriptUrl: transcript.transcriptUrl,
          }));
          setTranscripts(formattedTranscripts);
          setError(null);
        } else {
          setTranscripts([]);
        }
      } catch (err) {
        console.error('Error fetching transcripts:', err);
        setError(err instanceof Error ? err.message : 'An unknown error occurred');
        toast.error('Failed to fetch transcripts');
      } finally {
        setLoading(false);
      }
    };
    if (doctorEmail) {
      fetchTranscripts();
    }
  }, [doctorEmail, refreshTranscripts]);

  // 4. Filter transcripts client-side based on search query
  const filteredTranscripts = transcripts.filter((transcript) =>
    transcript.patientName.toLowerCase().startsWith(searchQuery.toLowerCase())
  );

  // Download handler
  const handleDownloadTranscript = async (transcript: Transcript) => {
    if (!transcript.transcriptUrl) {
      toast.error('Transcript URL not found');
      return;
    }
    try {
      const response = await fetch(`http://localhost:8080${transcript.transcriptUrl}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Origin': import.meta.env.VITE_FRONTEND_URL,
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: await getDecryptedEmail(),
          patient: { id: transcript.patientId },
        }),
      });
      if (!response.ok) {
        throw new Error('Failed to download transcript');
      }
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.style.display = 'none';
      a.href = url;
      a.download = `${transcript.patientName}_transcript.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      a.remove();
      toast.success(`Transcript for ${transcript.patientName} downloaded successfully`);
    } catch (error) {
      console.error('Error downloading transcript:', error);
      toast.error(`Failed to download transcript for ${transcript.patientName}`);
    }
  };

  // Inside Transcripts.tsx, within your handleEditTranscription:
const handleEditTranscription = async (transcript: Transcript) => {
  try {
    const response = await fetch(`http://localhost:8080/transcription/${transcript.id}/getTranscriptionByID`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
      },
      body: JSON.stringify({ transcription_id: transcript.id }),
    });
    if (!response.ok) {
      throw new Error('Failed to fetch transcription data');
    }
    const data = await response.json();
    const fetchedData = {
      text: data.text || '',
      report: data.report || '',
    };
    // Set data first, then show modal
    setSelectedTranscriptionId(transcript.id);
    setSelectedTranscriptionData(fetchedData);
    // Show modal after data is set
    setTimeout(() => setShowEditModal(true), 0);
  } catch (error) {
    console.error('Error fetching transcription data:', error);
    toast.error('Failed to fetch transcription data');
  }
};

  // Delete handler for transcription (using DELETE)
  const handleDeleteTranscript = async () => {
    if (!transcriptionToDelete) return;

    try {
      const response = await fetch(`http://localhost:8080/transcription/${transcriptionToDelete.id}/delete`, {
        method: 'POST',  // Changed from DELETE to POST to match backend route
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'Origin': import.meta.env.VITE_FRONTEND_URL,
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify({
          transcription_id: transcriptionToDelete.id
        })
      });
      if (!response.ok) {
        throw new Error('Failed to delete transcript');
      }
      toast.success(`Transcript for ${transcriptionToDelete.patientName} deleted successfully`);
      setTranscriptionToDelete(null);
      triggerTranscriptRefresh();
    } catch (error) {
      console.error('Error deleting transcript:', error);
      toast.error('Failed to delete transcript');
    }
  };

  const transcriptImages = [
    'https://i.pravatar.cc/150?img=1',
    'https://i.pravatar.cc/150?img=2',
    'https://i.pravatar.cc/150?img=3',
    'https://i.pravatar.cc/150?img=4',
    'https://i.pravatar.cc/150?img=5'
  ];

  return (<>
    <div className="min-h-screen bg-[#1C1C1C]">
      <div className="max-w-[120rem] mx-auto p-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8 max-w-[90rem] mx-auto">
          <div className="flex items-center space-x-4">
            <h2 className="text-3xl font-bold text-white flex items-center gap-3">
              Transcripts
              {!loading && transcripts.length > 0 && (
                <span className="text-lg text-gray-400">({transcripts.length})</span>
              )}
            </h2>
          </div>
          
          <div className="flex items-center space-x-4">
            <div className="relative w-96">
              <input
                type="text"
                placeholder="Search transcripts by patient name..."
                className="w-full border border-[#3a3a3a] rounded-lg p-2 pl-10 bg-[#2a2a2a] text-white focus:outline-none focus:ring-2 focus:ring-[#3ECF8E]/50"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            </div>
            
            <button 
              onClick={async () => {
                setLoading(true);
                if (doctorEmail) {
                  try {
                    const response = await fetch('http://localhost:8080/transcription/', {
                      method: 'POST',
                      credentials: 'include',
                      headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
                      },
                      body: JSON.stringify({ email: doctorEmail })
                    });
                    if (!response.ok) {
                      throw new Error('Failed to fetch transcripts');
                    }
                    const data = await response.json();
                    const transcriptsToSet = data.transcripts || [];
                    const patientNames = data.patients || [];
                    if (transcriptsToSet.length > 0) {
                      const formattedTranscripts = transcriptsToSet.map((transcript: any, index: number) => ({
                        id: transcript.id,
                        patientId: transcript.patientId,
                        patientName: patientNames[index] || 'Unknown Patient',
                        timestamp: transcript.timestamp,
                        type: 'Consultation',
                        status: 'completed',
                        transcriptUrl: transcript.transcriptUrl,
                      }));
                      setTranscripts(formattedTranscripts);
                      setError(null);
                    } else {
                      setTranscripts([]);
                    }
                  } catch (err) {
                    console.error('Error refreshing transcripts:', err);
                    setError(err instanceof Error ? err.message : 'An unknown error occurred');
                    toast.error('Failed to refresh transcripts');
                  } finally {
                    setLoading(false);
                  }
                }
              }}
              disabled={loading}
              className="bg-[#3ECF8E] text-black px-4 py-2 rounded-lg hover:bg-[#45E6A0] transition-colors duration-300 flex items-center space-x-2 disabled:opacity-50"
            >
              <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
              <span>Refresh</span>
            </button>
          </div>
        </div>


          {loading ? (
          <div className="flex justify-center items-center h-64">
            <Spinner />
          </div>
        ) : error ? (
          <EmptyState
            icon={<AlertCircle className="w-8 h-8 text-red-500" />}
            title="Error Loading Transcripts"
            description={error}
          />
        ) : filteredTranscripts.length > 0 ? (
          <div className="flex relative min-h-[calc(100vh-12rem)] pl-50">
            <div className="-mr-8">
              <MedicalSuggestions />
            </div>
            <div className="flex-1 -ml-26">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-x-10 gap-y-8 max-w-[900px] mx-auto">
                {filteredTranscripts.map((transcript, index) => (
                  <TranscriptCard 
                    key={transcript.id} 
                    transcript={transcript} 
                    image={transcriptImages[index % transcriptImages.length]}
                    onDownload={handleDownloadTranscript}
                    onEdit={handleEditTranscription}
                    onDelete={() => setTranscriptionToDelete(transcript)}
                  />
                ))}
              </div>
            </div>
          </div>
        ) : (
          <EmptyState
            title={searchQuery ? 'No Matching Transcripts' : 'No Transcripts Available'}
            description={searchQuery 
              ? `No transcripts found matching "${searchQuery}". Try a different search term.`
              : 'Start by recording a consultation to generate transcripts.'}
            action={
              searchQuery && (
                <button
                  onClick={() => setSearchQuery('')}
                  className="text-[#3ECF8E] hover:text-[#45E6A0] transition-colors duration-300"
                >
                  Clear search
                </button>
              )
            }
          />
        )}
        </div>
      </div>

      <AnimatePresence>
        {showEditModal && (
          <EditTranscriptionModal
            show={showEditModal}
            handleClose={() => setShowEditModal(false)}
            transcriptionId={selectedTranscriptionId || ''}
            onRefresh={triggerTranscriptRefresh}
            initialData={selectedTranscriptionData}
          />
        )}
      </AnimatePresence>

      <AnimatePresence>
        {transcriptionToDelete && (
          <ConfirmationModal
            isOpen={!!transcriptionToDelete}
            onClose={() => setTranscriptionToDelete(null)}
            onConfirm={handleDeleteTranscript}
            patientName={transcriptionToDelete.patientName}
          />
        )}
      </AnimatePresence>

    </>
  );
};

export default Transcripts;
