// import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import EditTranscriptionModal from '../../components/dashboard/EditTranscriptionModal';

describe('EditTranscriptionModal Component', () => {
  test('renders EditTranscriptionModal', () => {
    render(<EditTranscriptionModal show={true} handleClose={() => {}} transcriptionId="1" onRefresh={() => {}} />);
    expect(screen.getByText(/edit transcription/i)).toBeInTheDocument();
  });

  test('allows user to edit transcription details', () => {
    render(<EditTranscriptionModal show={true} handleClose={() => {}} transcriptionId="1" onRefresh={() => {}} />);
    const transcriptionInput = screen.getByLabelText(/transcription text/i) as HTMLTextAreaElement;
    fireEvent.change(transcriptionInput, { target: { value: 'Updated transcription text' } });

    expect(transcriptionInput.value).toBe('Updated transcription text');
  });
});
