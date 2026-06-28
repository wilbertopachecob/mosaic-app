import { screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from './App';
import { renderWithProviders } from './testUtils';

test('renders mosaic generator logo', () => {
  renderWithProviders(<App />);
  const logoElement = screen.getByRole('img', { name: /mosaic generator/i });
  expect(logoElement).toBeInTheDocument();
});

test('renders upload form', () => {
  renderWithProviders(<App />);
  const uploadLabel = screen.getByText(/select image/i);
  expect(uploadLabel).toBeInTheDocument();
});

test('renders tile size selector', () => {
  renderWithProviders(<App />);
  const tileSizeLabel = screen.getByText(/tile size/i);
  expect(tileSizeLabel).toBeInTheDocument();
});

test('renders source blend control', () => {
  renderWithProviders(<App />);
  const blendControl = screen.getByRole('slider', { name: /source blend/i });
  expect(blendControl).toBeInTheDocument();
});

test('renders generate button', () => {
  renderWithProviders(<App />);
  const generateButton = screen.getByText(/generate mosaic/i);
  expect(generateButton).toBeInTheDocument();
});
