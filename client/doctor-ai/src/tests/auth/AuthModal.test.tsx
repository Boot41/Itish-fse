// import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import AuthModal from '../../components/auth/AuthModal';

describe('AuthModal Component', () => {
  test('renders AuthModal', () => {
    render(<AuthModal isOpen={true} onClose={() => {}} />);
    expect(screen.getByText(/login/i)).toBeInTheDocument();
  });

  test('allows user to enter email and password', () => {
    render(<AuthModal isOpen={true} onClose={() => {}} />);
    const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
    const passwordInput = screen.getByLabelText(/password/i) as HTMLInputElement;
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });

    expect(emailInput.value).toBe('test@example.com');
    expect(passwordInput.value).toBe('password123');
  });

  test('submits the form', () => {
    const handleSubmit = jest.fn();
    render(<AuthModal isOpen={true} onClose={() => {}} />);

    fireEvent.click(screen.getByText(/submit/i));
    expect(handleSubmit).toHaveBeenCalled();
  });
});
