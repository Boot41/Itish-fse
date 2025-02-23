import { describe, test, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import ConfirmationModal from '../../pages/ConfirmationModal';

// Mock framer-motion to avoid animation-related issues in tests
jest.mock('framer-motion', () => ({
  motion: {
    div: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  },
}));

// Mock lucide-react icons
jest.mock('lucide-react', () => ({
  X: () => <div data-testid="close-icon">X</div>,
}));

describe('ConfirmationModal Component', () => {
  test('renders ConfirmationModal when open', () => {
    render(
      <ConfirmationModal 
        isOpen={true} 
        onClose={() => {}} 
        onConfirm={() => {}} 
        patientName="John Doe" 
      />
    );
    expect(screen.getByTestId('close-icon')).toBeInTheDocument();
    expect(screen.getByText(/John Doe/)).toBeInTheDocument();
  });

  test('does not render when closed', () => {
    render(
      <ConfirmationModal 
        isOpen={false} 
        onClose={() => {}} 
        onConfirm={() => {}} 
        patientName="John Doe" 
      />
    );
    expect(screen.queryByTestId('close-icon')).not.toBeInTheDocument();
  });
});
