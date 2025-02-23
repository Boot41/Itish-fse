// import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';

const ImageComponent = ({ src, alt }: { src: string; alt: string }) => (
  <div>
    <img src={src} alt={alt} />
  </div>
);

describe('Public Assets', () => {
  test('renders stethoscope.svg', () => {
    const { getByAltText } = render(
      <ImageComponent src="/stethoscope.svg" alt="Stethoscope" />
    );
    expect(getByAltText('Stethoscope')).toBeInTheDocument();
  });

  test('renders vite.svg', () => {
    const { getByAltText } = render(
      <ImageComponent src="/vite.svg" alt="Vite" />
    );
    expect(getByAltText('Vite')).toBeInTheDocument();
  });
});
