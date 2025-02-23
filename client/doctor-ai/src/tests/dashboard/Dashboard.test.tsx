// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Dashboard from '../../components/dashboard/Dashboard';

describe('Dashboard Component', () => {
  test('renders Dashboard', () => {
    render(<Dashboard />);
    expect(screen.getByText(/dashboard/i)).toBeInTheDocument();
  });

  test('displays patient list', () => {
    render(<Dashboard />);
    expect(screen.getByText(/patients/i)).toBeInTheDocument();
  });

  test('displays transcripts list', () => {
    render(<Dashboard />);
    expect(screen.getByText(/transcripts/i)).toBeInTheDocument();
  });
});
