import React from 'react';
import { screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import UploadForm from '../UploadForm';
import { renderWithProviders } from '../../testUtils';

// Mock functions
const mockHandleSubmit = jest.fn();
const mockHandleFileChange = jest.fn();
const mockHandleTileSizeChange = jest.fn();
const mockHandleBlendChange = jest.fn();

const defaultProps = {
  selectedTileSize: '20',
  selectedBlend: '0.42',
  isBtnDisabled: false,
  isLoading: false,
  handleSubmit: mockHandleSubmit,
  handleFileChange: mockHandleFileChange,
  handleTileSizeChange: mockHandleTileSizeChange,
  handleBlendChange: mockHandleBlendChange,
};

describe('UploadForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders all form elements', () => {
    renderWithProviders(<UploadForm {...defaultProps} />);
    
    expect(screen.getByLabelText(/select image/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/tile size/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/source blend/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /generate mosaic/i })).toBeInTheDocument();
  });

  it('displays correct tile size options', () => {
    renderWithProviders(<UploadForm {...defaultProps} />);
    
    const select = screen.getByRole('combobox');
    expect(select).toHaveValue('20');
    
    // Check for some key options
    expect(screen.getByText('5px - Very fine')).toBeInTheDocument();
    expect(screen.getByText('20px - Standard')).toBeInTheDocument();
    expect(screen.getByText('100px - Poster blocks')).toBeInTheDocument();
  });

  it('calls handleSubmit when form is submitted', () => {
    renderWithProviders(<UploadForm {...defaultProps} />);
    
    const submitButton = screen.getByRole('button', { name: /generate mosaic/i });
    fireEvent.click(submitButton);
    
    expect(mockHandleSubmit).toHaveBeenCalledTimes(1);
  });

  it('calls handleTileSizeChange when tile size is changed', () => {
    renderWithProviders(<UploadForm {...defaultProps} />);
    
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '50' } });
    
    expect(mockHandleTileSizeChange).toHaveBeenCalledWith('50');
  });

  it('calls handleBlendChange when blend is changed', () => {
    renderWithProviders(<UploadForm {...defaultProps} />);

    const blend = screen.getByLabelText(/source blend/i);
    fireEvent.change(blend, { target: { value: '0.55' } });

    expect(mockHandleBlendChange).toHaveBeenCalledWith('0.55');
  });

  it('explains source blend in help text', () => {
    renderWithProviders(<UploadForm {...defaultProps} />);

    expect(screen.getByText(/purer tiles/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/blend help/i)).toHaveAttribute('title', expect.stringContaining('original image'));
  });

  it('calls handleFileChange when file is selected', async () => {
    renderWithProviders(<UploadForm {...defaultProps} />);
    
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    const input = screen.getByLabelText(/select image/i);
    
    fireEvent.change(input, { target: { files: [file] } });
    
    await waitFor(() => {
      expect(mockHandleFileChange).toHaveBeenCalledWith(file);
    });
  });

  it('shows loading state when isLoading is true', () => {
    renderWithProviders(<UploadForm {...defaultProps} isLoading={true} />);
    
    expect(screen.getByText(/generating mosaic/i)).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('disables submit button when isBtnDisabled is true', () => {
    renderWithProviders(<UploadForm {...defaultProps} isBtnDisabled={true} />);
    
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('disables form elements when loading', () => {
    renderWithProviders(<UploadForm {...defaultProps} isLoading={true} />);
    
    expect(screen.getByLabelText(/select image/i)).toBeDisabled();
    expect(screen.getByRole('combobox')).toBeDisabled();
    expect(screen.getByLabelText(/source blend/i)).toBeDisabled();
  });

  it('shows processing info when loading', () => {
    renderWithProviders(<UploadForm {...defaultProps} isLoading={true} />);
    
    expect(screen.getByText(/processing may take a few moments/i)).toBeInTheDocument();
  });
}); 
