import React, { ChangeEvent, MouseEvent, useCallback } from "react";

interface UploadFormProps {
  selectedTileSize: string;
  isBtnDisabled: boolean;
  isLoading: boolean;
  handleSubmit: () => void;
  handleFileChange: (file: File) => void;
  handleTileSizeChange: (tile: string) => void;
}

const UploadForm: React.FC<UploadFormProps> = ({
  selectedTileSize,
  isBtnDisabled,
  isLoading,
  handleSubmit,
  handleFileChange,
  handleTileSizeChange,
}) => {
  // Handle form submission
  const onSubmit = useCallback((e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    handleSubmit();
  }, [handleSubmit]);

  // Handle file selection
  const onFileChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target;
    if (target && target.files && target.files.length > 0) {
      const file = target.files[0];
      
      // Validate file type
      if (!file.type.startsWith('image/')) {
        alert('Please select an image file');
        return;
      }
      
      // Validate file size (10MB limit)
      const maxSize = 10 * 1024 * 1024; // 10MB
      if (file.size > maxSize) {
        alert('File size must be less than 10MB');
        return;
      }
      
      handleFileChange(file);
    }
  }, [handleFileChange]);

  // Handle tile size change
  const onTileSizeChange = useCallback((event: ChangeEvent<HTMLSelectElement>) => {
    handleTileSizeChange(event.target.value);
  }, [handleTileSizeChange]);

  return (
    <form className="mt-3">
      {/* File Upload Section */}
      <div className="mb-3">
        <label htmlFor="imgUpload" className="form-label fw-bold">
          <i className="fas fa-image me-2"></i>
          Select Image
        </label>
        <input
          type="file"
          name="imgUpload"
          id="imgUpload"
          onChange={onFileChange}
          className="form-control"
          accept="image/*"
          disabled={isLoading}
        />
        <div className="form-text">
          Supported formats: JPG, PNG, GIF, BMP, TIFF, WebP (Max: 10MB)
        </div>
      </div>

      {/* Tile Size Selection */}
      <div className="mb-3">
        <label htmlFor="tileSize" className="form-label fw-bold">
          <i className="fas fa-th me-2"></i>
          Tile Size
        </label>
        <select
          name="tileSize"
          id="tileSize"
          onChange={onTileSizeChange}
          value={selectedTileSize}
          className="form-select"
          disabled={isLoading}
        >
          <option value="5">5px - Very Fine Detail</option>
          <option value="10">10px - Fine Detail</option>
          <option value="15">15px - Medium Detail</option>
          <option value="20">20px - Standard</option>
          <option value="25">25px - Coarse Detail</option>
          <option value="30">30px - Very Coarse</option>
          <option value="50">50px - Large Tiles</option>
          <option value="100">100px - Very Large Tiles</option>
        </select>
        <div className="form-text">
          Smaller tiles create more detailed mosaics but take longer to process
        </div>
      </div>

      {/* Submit Button */}
      <div className="d-grid">
        <button
          onClick={onSubmit}
          type="submit"
          className={`btn btn-primary ${isLoading ? 'disabled' : ''}`}
          disabled={isBtnDisabled || isLoading}
        >
          {isLoading ? (
            <>
              <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              Generating Mosaic...
            </>
          ) : (
            <>
              <i className="fas fa-magic me-2"></i>
              Generate Mosaic
            </>
          )}
        </button>
      </div>

      {/* Processing Info */}
      {isLoading && (
        <div className="mt-3 text-center">
          <small className="text-muted">
            <i className="fas fa-info-circle me-1"></i>
            Processing may take a few moments depending on image size and tile size
          </small>
        </div>
      )}
    </form>
  );
};

export default UploadForm;
