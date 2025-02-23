// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import {Spinner} from '../../components/ui/Spinner';

describe('Spinner Component', () => {
  test('renders Spinner', () => {
    render(<Spinner />);
    expect(screen.getByRole('progressbar')).toBeInTheDocument();
  });
});
