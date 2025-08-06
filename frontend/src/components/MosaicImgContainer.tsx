import React, { useCallback } from "react";
import imgPlaceholder from "../assets/img/img_placeholder.png";

interface MosaicImgContainerProps {
  mosaicImg: string | null;
  duration: number;
  fileName: string | undefined;
  onReset?: () => void;
}

const MosaicImgContainer: React.FC<MosaicImgContainerProps> = ({
  mosaicImg,
  duration,
  fileName,
  onReset,
}) => {
  // Handle download
  const handleDownload = useCallback(() => {
    if (!mosaicImg) return;
    
    const link = document.createElement('a');
    link.href = `data:image/jpeg;base64,${mosaicImg}`;
    link.download = `mosaic-${fileName || 'image'}.jpg`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }, [mosaicImg, fileName]);

  // Format duration for display
  const formatDuration = useCallback((seconds: number) => {
    if (seconds < 1) {
      return `${Math.round(seconds * 1000)}ms`;
    }
    return `${seconds.toFixed(2)}s`;
  }, []);

  return (
    <div className="d-flex flex-column h-100">
      {/* Mosaic Image */}
      <div className="text-center mb-3">
        <img
          src={mosaicImg ? `data:image/jpeg;base64,${mosaicImg}` : imgPlaceholder}
          alt="Mosaic Result"
          className="img-fluid rounded shadow-sm"
          style={{ 
            maxHeight: '400px', 
            objectFit: 'contain',
            border: '2px solid #dee2e6'
          }}
        />
      </div>

      {/* Processing Stats */}
      <div className="card mb-3">
        <div className="card-body py-2">
          <div className="row text-center">
            <div className="col-6">
              <small className="text-muted d-block">Processing Time</small>
              <strong className="text-primary">{formatDuration(duration)}</strong>
            </div>
            <div className="col-6">
              <small className="text-muted d-block">Original File</small>
              <strong className="text-secondary">
                {fileName ? fileName.substring(0, 20) + (fileName.length > 20 ? '...' : '') : 'Unknown'}
              </strong>
            </div>
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="d-grid gap-2 mt-auto">
        <button
          onClick={handleDownload}
          className={`btn btn-success ${!mosaicImg ? 'disabled' : ''}`}
          disabled={!mosaicImg}
          title="Download mosaic image"
        >
          <i className="fas fa-download me-2"></i>
          Download Mosaic
        </button>
        
        {onReset && (
          <button
            onClick={onReset}
            className="btn btn-outline-secondary"
            title="Start over with a new image"
          >
            <i className="fas fa-redo me-2"></i>
            Create New Mosaic
          </button>
        )}
      </div>

      {/* Success Message */}
      {mosaicImg && (
        <div className="alert alert-success mt-3" role="alert">
          <i className="fas fa-check-circle me-2"></i>
          <strong>Success!</strong> Your mosaic has been generated successfully.
        </div>
      )}
    </div>
  );
};

export default MosaicImgContainer;
