// import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import EditPatientModal from '../../components/dashboard/EditPatientModal';

describe('EditPatientModal Component', () => {
  test('renders EditPatientModal', () => {
    render(<EditPatientModal show={true} handleClose={() => {}} patientId="1" doctorEmail="doctor@example.com" onRefresh={() => {}} />);
    expect(screen.getByText(/edit patient/i)).toBeInTheDocument();
  });

  test('allows user to edit patient details', () => {
    render(<EditPatientModal show={true} handleClose={() => {}} patientId="1" doctorEmail="doctor@example.com" onRefresh={() => {}} />);
    const nameInput = screen.getByLabelText(/name/i) as HTMLInputElement;
    const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
    fireEvent.change(nameInput, { target: { value: 'Jane Doe' } });
    fireEvent.change(emailInput, { target: { value: 'jane@example.com' } });

    expect(nameInput.value).toBe('Jane Doe');
    expect(emailInput.value).toBe('jane@example.com');
  });
});
