import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from './App';

test('renders mosaic generator title', () => {
  render(<App />);
  const titleElement = screen.getByText(/mosaic generator/i);
  expect(titleElement).toBeInTheDocument();
});

test('renders upload form', () => {
  render(<App />);
  const uploadLabel = screen.getByText(/select image/i);
  expect(uploadLabel).toBeInTheDocument();
});

test('renders tile size selector', () => {
  render(<App />);
  const tileSizeLabel = screen.getByText(/tile size/i);
  expect(tileSizeLabel).toBeInTheDocument();
});

test('renders generate button', () => {
  render(<App />);
  const generateButton = screen.getByText(/generate mosaic/i);
  expect(generateButton).toBeInTheDocument();
});
