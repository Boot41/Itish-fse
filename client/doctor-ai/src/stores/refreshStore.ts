import { create } from 'zustand';

export const useRefreshStore = create<{
  refreshPatients: boolean;
  refreshTranscripts: boolean;
  triggerPatientRefresh: () => void;
  triggerTranscriptRefresh: () => void;
}>((set) => ({
  refreshPatients: false,
  refreshTranscripts: false,
  triggerPatientRefresh: () => set((state) => ({ refreshPatients: !state.refreshPatients })),
  triggerTranscriptRefresh: () => set((state) => ({ refreshTranscripts: !state.refreshTranscripts })),
}));
