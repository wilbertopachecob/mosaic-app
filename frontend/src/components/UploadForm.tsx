import React, { ChangeEvent, MouseEvent, useCallback } from "react";

interface UploadFormProps {
  selectedTileSize: string;
  selectedBlend: string;
  isBtnDisabled: boolean;
  isLoading: boolean;
  handleSubmit: () => void;
  handleFileChange: (file: File) => void;
  handleTileSizeChange: (tile: string) => void;
  handleBlendChange: (blend: string) => void;
}

const UploadForm: React.FC<UploadFormProps> = ({
  selectedTileSize,
  selectedBlend,
  isBtnDisabled,
  isLoading,
  handleSubmit,
  handleFileChange,
  handleTileSizeChange,
  handleBlendChange,
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

  // Handle source blend change
  const onBlendChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    handleBlendChange(event.target.value);
  }, [handleBlendChange]);

  const blendPercent = Math.round(Number(selectedBlend) * 100);

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

      {/* Source Blend Selection */}
      <div className="mb-3">
        <label htmlFor="blend" className="form-label fw-bold">
          <i className="fas fa-adjust me-2"></i>
          Source Blend
          <i
            className="fas fa-circle-info ms-2 text-muted"
            role="img"
            aria-label="Blend help"
            title="Controls how much of the original image is blended over the mosaic. Lower values show purer photo tiles; higher values look more like the source image."
          ></i>
        </label>
        <div className="d-flex align-items-center gap-3">
          <input
            type="range"
            name="blend"
            id="blend"
            min="0"
            max="0.75"
            step="0.01"
            value={selectedBlend}
            onChange={onBlendChange}
            className="form-range"
            disabled={isLoading}
            aria-describedby="blendHelp"
          />
          <span className="badge bg-secondary">{blendPercent}%</span>
        </div>
        <div id="blendHelp" className="form-text">
          0% is a pure tile mosaic. Higher values preserve faces, edges, and shadows like online mosaic services.
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
