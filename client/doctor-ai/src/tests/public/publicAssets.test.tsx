import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import React from 'react';

interface ImageComponentProps {
  src: string;
  alt: string;
  testId: string;
}

const ImageComponent: React.FC<ImageComponentProps> = ({ src, alt, testId }) => (
  <div>
    <img src={src} alt={alt} data-testid={testId} />
  </div>
);

describe('Public Assets', () => {
  test('renders stethoscope.svg', () => {
    const { getByRole } = render(
      <ImageComponent 
        src="/stethoscope.svg" 
        alt="Stethoscope" 
        testId="stethoscope" 
      />
    );
    expect(getByRole('img')).toBeInTheDocument();
    expect(getByRole('img')).toHaveAttribute('src', '/stethoscope.svg');
  });

  test('renders vite.svg', () => {
    const { getByRole } = render(
      <ImageComponent 
        src="/vite.svg" 
        alt="Vite" 
        testId="vite" 
      />
    );
    expect(getByRole('img')).toBeInTheDocument();
    expect(getByRole('img')).toHaveAttribute('src', '/vite.svg');
  });
});
