import { ChangeEvent, MouseEvent, useState, useRef } from "react";

type UploadFormProps = {
  selectedTileSize: string;
  isBtnDisabled: boolean;
  handleSubmit: () => void;
  handleFileChange: (file: File) => void;
  handleTileSizeChange: (tile: string) => void;
  isLoading: boolean;
};

const UploadForm = ({
  selectedTileSize,
  isBtnDisabled,
  handleSubmit,
  handleFileChange,
  handleTileSizeChange,
  isLoading,
}: UploadFormProps) => {
  const [isDragOver, setIsDragOver] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const onSubmit = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    handleSubmit();
  };

  const onFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target;
    if (target && target.files && target.files.length > 0) {
      handleFileChange(target.files[0]);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      const file = files[0];
      if (file.type.startsWith('image/')) {
        handleFileChange(file);
      }
    }
  };

  const handleFileSelect = () => {
    fileInputRef.current?.click();
  };

  const tileSizeOptions = [
    { value: "5", label: "5px - High Detail" },
    { value: "10", label: "10px - Fine Detail" },
    { value: "15", label: "15px - Medium Detail" },
    { value: "20", label: "20px - Balanced" },
    { value: "25", label: "25px - Coarse Detail" },
    { value: "50", label: "50px - Large Tiles" },
    { value: "100", label: "100px - Very Large Tiles" },
  ];

  return (
    <form className="mt-3">
      {/* File Upload Area */}
      <div
        className={`file-upload-area ${isDragOver ? 'dragover' : ''}`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={handleFileSelect}
      >
        <div className="file-upload-icon">
          <i className="fas fa-cloud-upload-alt"></i>
        </div>
        <div className="file-upload-text">
          <strong>Click to upload</strong> or drag and drop
        </div>
        <div className="file-upload-hint">
          PNG, JPG, GIF up to 10MB
        </div>
      </div>

      {/* Hidden file input */}
      <input
        ref={fileInputRef}
        type="file"
        name="imgUpload"
        id="imgUpload"
        onChange={onFileChange}
        className="form-control"
        style={{ display: 'none' }}
        accept="image/*"
      />

      {/* Tile Size Selection */}
      <div className="form-group mt-3">
        <label htmlFor="tileSize" className="form-label">
          <i className="fas fa-th"></i>
          Tile Size
        </label>
        <select
          name="tileSize"
          id="tileSize"
          onChange={(e) => handleTileSizeChange(e.target.value)}
          value={selectedTileSize}
          className="form-select"
          disabled={isLoading}
        >
          {tileSizeOptions.map(option => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        <small style={{ color: 'var(--text-secondary)', fontSize: '0.85rem', marginTop: '0.5rem', display: 'block' }}>
          Smaller tiles create more detailed mosaics but take longer to process
        </small>
      </div>

      {/* Generate Button */}
      <button
        onClick={onSubmit}
        type="submit"
        className="btn btn-primary mt-3"
        disabled={isBtnDisabled}
        style={{ width: '100%' }}
      >
        {isLoading ? (
          <>
            <div className="loading-spinner"></div>
            Generating Mosaic...
          </>
        ) : (
          <>
            <i className="fas fa-magic"></i>
            Generate Mosaic
          </>
        )}
      </button>

      {/* Processing Info */}
      {isLoading && (
        <div className="mt-3" style={{
          textAlign: 'center',
          color: 'var(--text-secondary)',
          fontSize: '0.9rem'
        }}>
          <i className="fas fa-cog fa-spin"></i>
          Processing your image... This may take a few moments.
        </div>
      )}
    </form>
  );
};

export default UploadForm;
