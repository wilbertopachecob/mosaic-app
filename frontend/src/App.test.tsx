import { screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from './App';
import { renderWithProviders } from './testUtils';

const mockFetch = jest.fn();
global.fetch = mockFetch;

beforeEach(() => {
  mockFetch.mockReset();
  URL.createObjectURL = jest.fn(() => 'blob:preview');
  URL.revokeObjectURL = jest.fn();
});

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

test('disables generate until an image is selected', () => {
  renderWithProviders(<App />);

  expect(screen.getByRole('button', { name: /generate mosaic/i })).toBeDisabled();
  expect(mockFetch).not.toHaveBeenCalled();
});

test('displays file size after selecting an image', async () => {
  renderWithProviders(<App />);

  const file = new File(['x'.repeat(2048)], 'photo.jpg', { type: 'image/jpeg' });
  fireEvent.change(screen.getByLabelText(/select image/i), {
    target: { files: [file] },
  });

  await waitFor(() => {
    expect(screen.getAllByText('photo.jpg').length).toBeGreaterThan(0);
    expect(screen.getByText(/2 kb/i)).toBeInTheDocument();
  });
});

test('submits the selected image and shows the mosaic result', async () => {
  mockFetch.mockResolvedValue({
    ok: true,
    text: async () =>
      JSON.stringify({ mosaicImg: 'abc123', duration: 1.5 }),
  });

  renderWithProviders(<App />);

  const file = new File(['image'], 'photo.jpg', { type: 'image/jpeg' });
  fireEvent.change(screen.getByLabelText(/select image/i), {
    target: { files: [file] },
  });

  fireEvent.click(screen.getByRole('button', { name: /generate mosaic/i }));

  expect(await screen.findByRole('img', { name: /generated mosaic/i })).toHaveAttribute(
    'src',
    'data:image/jpeg;base64,abc123'
  );
  expect(screen.getByText('1.50s')).toBeInTheDocument();
  expect(mockFetch).toHaveBeenCalledWith('/api/file/upload', expect.objectContaining({
    method: 'POST',
  }));
});

test('shows API error messages from failed uploads', async () => {
  mockFetch.mockResolvedValue({
    ok: false,
    status: 400,
    text: async () =>
      JSON.stringify({ error: 'invalid_image', message: 'Failed to decode image' }),
  });

  renderWithProviders(<App />);

  const file = new File(['image'], 'photo.jpg', { type: 'image/jpeg' });
  fireEvent.change(screen.getByLabelText(/select image/i), {
    target: { files: [file] },
  });

  fireEvent.click(screen.getByRole('button', { name: /generate mosaic/i }));

  expect(await screen.findByRole('alert')).toHaveTextContent(
    'Failed to decode image'
  );
});
