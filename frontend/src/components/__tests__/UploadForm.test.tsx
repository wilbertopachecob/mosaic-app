import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import UploadForm from '../UploadForm';

// Mock functions
const mockHandleSubmit = jest.fn();
const mockHandleFileChange = jest.fn();
const mockHandleTileSizeChange = jest.fn();

const defaultProps = {
  selectedTileSize: '20',
  isBtnDisabled: false,
  isLoading: false,
  handleSubmit: mockHandleSubmit,
  handleFileChange: mockHandleFileChange,
  handleTileSizeChange: mockHandleTileSizeChange,
};

describe('UploadForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders all form elements', () => {
    render(<UploadForm {...defaultProps} />);
    
    expect(screen.getByLabelText(/select image/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/tile size/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /generate mosaic/i })).toBeInTheDocument();
  });

  it('displays correct tile size options', () => {
    render(<UploadForm {...defaultProps} />);
    
    const select = screen.getByRole('combobox');
    expect(select).toHaveValue('20');
    
    // Check for some key options
    expect(screen.getByText('5px - Very Fine Detail')).toBeInTheDocument();
    expect(screen.getByText('20px - Standard')).toBeInTheDocument();
    expect(screen.getByText('100px - Very Large Tiles')).toBeInTheDocument();
  });

  it('calls handleSubmit when form is submitted', () => {
    render(<UploadForm {...defaultProps} />);
    
    const submitButton = screen.getByRole('button', { name: /generate mosaic/i });
    fireEvent.click(submitButton);
    
    expect(mockHandleSubmit).toHaveBeenCalledTimes(1);
  });

  it('calls handleTileSizeChange when tile size is changed', () => {
    render(<UploadForm {...defaultProps} />);
    
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '50' } });
    
    expect(mockHandleTileSizeChange).toHaveBeenCalledWith('50');
  });

  it('calls handleFileChange when file is selected', async () => {
    render(<UploadForm {...defaultProps} />);
    
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    const input = screen.getByLabelText(/select image/i);
    
    fireEvent.change(input, { target: { files: [file] } });
    
    await waitFor(() => {
      expect(mockHandleFileChange).toHaveBeenCalledWith(file);
    });
  });

  it('shows loading state when isLoading is true', () => {
    render(<UploadForm {...defaultProps} isLoading={true} />);
    
    expect(screen.getByText(/generating mosaic/i)).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('disables submit button when isBtnDisabled is true', () => {
    render(<UploadForm {...defaultProps} isBtnDisabled={true} />);
    
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('disables form elements when loading', () => {
    render(<UploadForm {...defaultProps} isLoading={true} />);
    
    expect(screen.getByLabelText(/select image/i)).toBeDisabled();
    expect(screen.getByRole('combobox')).toBeDisabled();
  });

  it('shows processing info when loading', () => {
    render(<UploadForm {...defaultProps} isLoading={true} />);
    
    expect(screen.getByText(/processing may take a few moments/i)).toBeInTheDocument();
  });
}); 