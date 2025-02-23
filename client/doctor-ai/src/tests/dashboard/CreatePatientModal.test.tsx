// import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import CreatePatientModal from '../../components/dashboard/CreatePatientModal';

describe('CreatePatientModal Component', () => {
  test('renders CreatePatientModal', () => {
    render(<CreatePatientModal isOpen={true} onClose={() => {}} onPatientCreated={() => {}} />);
    expect(screen.getByText(/create patient/i)).toBeInTheDocument();
  });

  test('allows user to enter patient details', () => {
    render(<CreatePatientModal isOpen={true} onClose={() => {}} onPatientCreated={() => {}} />);
    const nameInput = screen.getByLabelText(/name/i) as HTMLInputElement;
    const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
    fireEvent.change(nameInput, { target: { value: 'John Doe' } });
    fireEvent.change(emailInput, { target: { value: 'john@example.com' } });

    expect(nameInput.value).toBe('John Doe');
    expect(emailInput.value).toBe('john@example.com');
  });
});
