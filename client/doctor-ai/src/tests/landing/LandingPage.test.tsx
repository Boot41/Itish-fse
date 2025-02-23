import { describe, test, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import LandingPage from '../../components/landing/LandingPage';

describe('LandingPage Component', () => {
  test('renders LandingPage', () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );
    // Update the text matcher to something that actually exists on the page:
    expect(screen.getByText(/build in a weekend/i)).toBeInTheDocument();
  });
});
