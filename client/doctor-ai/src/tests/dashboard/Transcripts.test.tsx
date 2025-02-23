// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Transcripts from '../../components/dashboard/Transcripts';

describe('Transcripts Component', () => {
  test('renders Transcripts', () => {
    render(<Transcripts />);
    expect(screen.getByText(/transcripts/i)).toBeInTheDocument();
  });

  test('displays transcripts list', () => {
    render(<Transcripts />);
    // Assuming there is a mock transcript data or a way to test the list
    expect(screen.getByText(/sample transcript text/i)).toBeInTheDocument(); // Adjust based on actual transcript content
  });
});
