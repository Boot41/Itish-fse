import React, { useState, useEffect } from 'react';
import { Search, Plus, Pencil, Trash2, AlertCircle } from 'lucide-react';
import CreatePatientModal from './CreatePatientModal';
import ConfirmationModal from '../../pages/ConfirmationModal';
import EditPatientModal from './EditPatientModal';
import { Patient } from '../../types/Patient';
import { useRefreshStore } from '../../stores/refreshStore';
import { useAuthStore } from '../../stores/authStore';
import { toast } from 'sonner';
import { Spinner } from '../ui/Spinner';
import { EmptyState } from '../ui/EmptyState';
import { useNavigate } from 'react-router-dom';
import { AnimatePresence } from 'framer-motion';
import MedicalSuggestions from './MedicalSuggestions'; // Import MedicalSuggestions component

interface PatientCardProps {
  patient: Patient;
  image: string;
  onEdit: (patient: Patient) => void;
  onRequestDelete: (patient: Patient) => void;
}

const PatientCard: React.FC<PatientCardProps> = ({ patient, image, onEdit, onRequestDelete }) => {
  const navigate = useNavigate();

  return (
    <div 
      className="relative w-[280px] h-[360px] bg-[#2a2a2a] rounded-lg border border-[#3a3a3a] group transition-all duration-300 ease-in-out transform hover:scale-[1.02] hover:border-emerald-500/50 cursor-pointer overflow-hidden"
      onClick={() => navigate(`/logs?search=${encodeURIComponent(patient.name)}`)}
    >
      {/* Top gradient bar */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-emerald-400 to-emerald-600" />
      
      {/* Content */}
      <div className="relative h-full flex flex-col p-6">
        {/* Avatar section */}
        <div className="flex justify-center mb-4">
          <div className="relative">
            <img alt="Patient Profile" className="w-20 h-20 rounded-full border-4 border-[#3ECF8E]/30 hover:border-[#3ECF8E] transition-all" src={image} />
            <div className="absolute -inset-1 bg-emerald-500/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300 blur-xl rounded-full" />
          </div>
        </div>

        {/* Patient info */}
        <div className="text-center space-y-3">
          <h3 className="text-xl font-bold text-white group-hover:text-emerald-400 transition-colors duration-300 truncate px-4">
            {patient.name}
          </h3>

          {/* Status badge */}
          <div className="inline-flex items-center px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm">
            Active Patient
          </div>

          {/* Info grid */}
          <div className="grid grid-cols-2 gap-3 mt-4 mb-16">
            <div className="bg-[#333] rounded-lg p-3 border border-[#3a3a3a] group-hover:border-emerald-500/20 transition-colors duration-300">
              <div className="text-xs text-gray-400 mb-1">Age</div>
              <div className="text-white font-medium">{patient.age}</div>
            </div>
            <div className="bg-[#333] rounded-lg p-3 border border-[#3a3a3a] group-hover:border-emerald-500/20 transition-colors duration-300">
              <div className="text-xs text-gray-400 mb-1">Gender</div>
              <div className="text-white font-medium">{patient.gender}</div>
            </div>
          </div>
        </div>

        {/* Action buttons - Fixed at bottom with glass effect */}
        <div className="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-[#2a2a2a] via-[#2a2a2a] to-transparent">
          <div className="grid grid-cols-2 gap-3 px-2">
            <button
              onClick={(e) => {
                e.stopPropagation();
                onEdit(patient);
              }}
              className="flex items-center justify-center gap-2 py-2 rounded-lg bg-[#333] hover:bg-emerald-500/10 border border-[#3a3a3a] hover:border-emerald-500/30 transition-all duration-300 group/edit"
            >
              <Pencil className="w-4 h-4 text-emerald-500 transition-transform duration-300 group-hover/edit:scale-110" />
              <span className="text-sm text-emerald-500">Edit</span>
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                onRequestDelete(patient);
              }}
              className="flex items-center justify-center gap-2 py-2 rounded-lg bg-[#333] hover:bg-red-500/10 border border-[#3a3a3a] hover:border-red-500/30 transition-all duration-300 group/delete"
            >
              <Trash2 className="w-4 h-4 text-red-400 transition-transform duration-300 group-hover/delete:scale-110" />
              <span className="text-sm text-red-400">Delete</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

const Patients: React.FC = () => {
  const [isCreatePatientModalOpen, setIsCreatePatientModalOpen] = useState(false);
  const [patients, setPatients] = useState<Patient[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [patientToDelete, setPatientToDelete] = useState<Patient | null>(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedPatientId, setSelectedPatientId] = useState<string | null>(null);
  const [selectedPatientData, setSelectedPatientData] = useState<{ name: string; age: string; gender: string } | null>(null);

  const { refreshPatients, triggerPatientRefresh } = useRefreshStore();
  const { getDecryptedEmail } = useAuthStore();
  const [doctorEmail, setDoctorEmail] = useState<string | null>(null);

  // const navigate = useNavigate();

  // When edit icon is clicked, fetch the patient data (both patient & transcription are returned)
  const handleEditPatient = async (patient: Patient) => {
    try {
      const response = await fetch(`http://localhost:8080/patients/${patient.id}/getPatient`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify({ email: doctorEmail, patient_id: patient.id }),
      });
      if (!response.ok) {
        throw new Error('Failed to fetch patient data');
      }
      const data = await response.json();
      // Assuming backend returns: { patient: { ... }, transcription: { ... } }
      const fetchedData = {
        name: data.patient.name || '',
        age: String(data.patient.age || ''),
        gender: data.patient.gender || '',
      };
      // Set data first
      setSelectedPatientId(patient.id);
      setSelectedPatientData(fetchedData);
      // Show modal in next tick to ensure data is set
      setTimeout(() => setShowEditModal(true), 0);
    } catch (error) {
      console.error('Error fetching patient data:', error);
      toast.error('Failed to fetch patient data');
    }
  };

  const handleRequestDelete = (patient: Patient) => {
    setPatientToDelete(patient);
  };

  // For deletion, use POST as per your backend
  const handleDeletePatient = async () => {
    if (!patientToDelete) return;
    try {
      const response = await fetch(`http://localhost:8080/patients/${patientToDelete.id}`, {
        method: 'POST', // Using POST for deletion
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify({ email: doctorEmail, patient_id: patientToDelete.id }),
      });
      if (!response.ok) {
        throw new Error('Failed to delete patient');
      }
      toast.success(`${patientToDelete.name} deleted successfully`);
      triggerPatientRefresh();
      setPatientToDelete(null);
    } catch (err) {
      console.error('Error deleting patient:', err);
      toast.error('Failed to delete patient');
    }
  };

  // Filter patients client-side
  const filteredPatients = patients.filter((patient) =>
    patient.name.toLowerCase().startsWith(searchQuery.toLowerCase())
  );

  // Decrypt the doctor's email on mount
  useEffect(() => {
    const decryptEmail = async () => {
      try {
        const decrypted = await getDecryptedEmail();
        setDoctorEmail(decrypted);
      } catch (err) {
        console.error('Failed to decrypt email:', err);
        toast.error('Failed to retrieve doctor email');
      }
    };
    decryptEmail();
  }, [getDecryptedEmail]);

  // Fetch patients when doctorEmail or refreshPatients changes
  useEffect(() => {
    const fetchPatients = async () => {
      if (!doctorEmail) {
        setError('No doctor email available');
        setLoading(false);
        return;
      }
      setLoading(true);
      try {
        const response = await fetch('http://localhost:8080/patients/patientsList', {
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
          throw new Error('Failed to fetch patients');
        }
        const data = await response.json();
        const patientsData = data.patients || data;
        if (Array.isArray(patientsData)) {
          const formattedPatients = patientsData.map((p: any) => ({
            id: p.id || p.ID,
            name: p.name || p.Name,
            age: p.age || p.Age,
            gender: p.gender || p.Gender,
            url: p.url || `/patients/${p.id || p.ID}`,
            email: p.email || p.Email,
            doctorId: p.doctorId || p.DoctorID,
            createdAt: p.createdAt || p.CreatedAt,
          }));
          setPatients(formattedPatients);
          setError(null);
        } else {
          setPatients([]);
        }
      } catch (err) {
        console.error('Error fetching patients:', err);
        setError(err instanceof Error ? err.message : 'An unknown error occurred');
        toast.error('Failed to fetch patients');
      } finally {
        setLoading(false);
      }
    };
    if (doctorEmail) {
      fetchPatients();
    }
  }, [doctorEmail, refreshPatients]);

  const handlePatientCreated = (newPatient: Patient) => {
    setPatients((prev) => [...prev, newPatient]);
    setIsCreatePatientModalOpen(false);
    toast.success(`Patient ${newPatient.name} created successfully`);
  };

  const patientImages = [
    'https://i.pravatar.cc/150?img=1',
    'https://i.pravatar.cc/150?img=2',
    'https://i.pravatar.cc/150?img=3',
    'https://i.pravatar.cc/150?img=4',
    'https://i.pravatar.cc/150?img=5'
  ];

  return (
    <div className="min-h-screen bg-[#1C1C1C] p-6">
      <div className="max-w-[90rem] mx-auto">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
          <div className="flex items-center space-x-4">
            <h2 className="text-3xl font-bold text-white flex items-center gap-3">
              Patients
              {!loading && patients.length > 0 && (
                <span className="text-lg text-gray-400">({patients.length})</span>
              )}
            </h2>
          </div>
          <div className="flex items-center space-x-4">
            <div className="relative w-96">
              <input
                type="text"
                placeholder="Search Patients..."
                className="w-full border border-[#3a3a3a] rounded-lg p-2 pl-10 bg-[#2a2a2a] text-white focus:outline-none focus:ring-2 focus:ring-[#3ECF8E]/50"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            </div>
            <button 
              onClick={() => setIsCreatePatientModalOpen(true)}
              className="bg-[#3ECF8E] text-black px-4 py-2 rounded-lg hover:bg-[#45E6A0] transition-colors duration-300 flex items-center space-x-2"
            >
              <Plus className="w-5 h-5" />
              <span>Create Patient</span>
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
            title="Error Loading Patients"
            description={error}
          />
        ) : filteredPatients.length > 0 ? (
          <div className="flex gap-24 px-4 pb-4 relative min-h-[calc(100vh-12rem)]">
            <MedicalSuggestions />
            <div className="flex-1">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-x-24 gap-y-10">
                {filteredPatients.map((patient, index) => (
                  <PatientCard 
                    key={patient.id}
                    patient={patient}
                    image={patientImages[index % patientImages.length]}
                    onEdit={handleEditPatient}
                    onRequestDelete={handleRequestDelete}
                  />
                ))}
              </div>
            </div>
          </div>
        ) : (
          <EmptyState
            title={searchQuery ? 'No Matching Patients' : 'No Patients Available'}
            description={
              searchQuery
                ? `No patients found matching "${searchQuery}". Try a different search term.`
                : 'Start by creating your first patient.'
            }
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

      <CreatePatientModal
        isOpen={isCreatePatientModalOpen}
        onClose={() => setIsCreatePatientModalOpen(false)}
        onPatientCreated={handlePatientCreated}
      />

      <EditPatientModal
        show={showEditModal}
        handleClose={() => setShowEditModal(false)}
        patientId={selectedPatientId || ''}
        doctorEmail={doctorEmail || ''}
        onRefresh={triggerPatientRefresh}
        initialData={selectedPatientData}
      />

      <AnimatePresence>
        {patientToDelete && (
          <ConfirmationModal 
            isOpen={!!patientToDelete}
            onClose={() => setPatientToDelete(null)}
            onConfirm={handleDeletePatient}
            patientName={patientToDelete.name}
          />
        )}
      </AnimatePresence>
    </div>
  );
};

export default Patients;
