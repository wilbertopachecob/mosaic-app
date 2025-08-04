const MosaicImgContainer = ({
  mosaicImg,
  duration,
  fileName,
  isLoading,
}: {
  mosaicImg: string | undefined;
  duration: number;
  fileName: string | undefined;
  isLoading: boolean;
}) => {
  const formatDuration = (seconds: number): string => {
    if (seconds < 1) return `${(seconds * 1000).toFixed(0)}ms`;
    return `${seconds.toFixed(2)}s`;
  };

  const getDownloadFileName = (): string => {
    if (!fileName) return 'mosaic.jpg';
    const nameWithoutExt = fileName.split('.').slice(0, -1).join('.');
    return `mosaic-${nameWithoutExt}.jpg`;
  };

  return (
    <div className="d-flex flex-column">
      {/* Image Container */}
      <div className="image-container">
        {isLoading ? (
          <div className="image-placeholder">
            <div className="loading-spinner" style={{ width: '40px', height: '40px', borderWidth: '4px' }}></div>
            <div>Generating your mosaic...</div>
          </div>
        ) : mosaicImg ? (
          <img
            src={`data:image/jpeg;base64,${mosaicImg}`}
            alt="Generated mosaic"
            style={{ maxWidth: '100%', maxHeight: '400px', objectFit: 'contain' }}
          />
        ) : (
          <div className="image-placeholder">
            <i className="fas fa-image"></i>
            <div>Your mosaic will appear here</div>
            <small>Upload an image and click "Generate Mosaic" to get started</small>
          </div>
        )}
      </div>

      {/* Stats and Actions */}
      {mosaicImg && (
        <div className="stats-container slide-up">
          <div className="stat-item">
            <i className="fas fa-clock"></i>
            <span>Processing time: <span className="stat-value">{formatDuration(duration)}</span></span>
          </div>
          <div className="stat-item">
            <i className="fas fa-check-circle"></i>
            <span>Status: <span className="stat-value">Complete</span></span>
          </div>
        </div>
      )}

      {/* Download Button */}
      {mosaicImg && (
        <div className="mt-3 slide-up">
          <a
            href={`data:image/jpeg;base64,${mosaicImg}`}
            download={getDownloadFileName()}
            className="btn btn-success"
            style={{ width: '100%' }}
          >
            <i className="fas fa-download"></i>
            Download Mosaic
          </a>
        </div>
      )}

      {/* Tips */}
      {!mosaicImg && !isLoading && (
        <div className="mt-3" style={{
          padding: '1rem',
          backgroundColor: 'rgba(102, 126, 234, 0.1)',
          borderRadius: '8px',
          border: '1px solid rgba(102, 126, 234, 0.2)',
          color: 'var(--text-secondary)',
          fontSize: '0.9rem'
        }}>
          <div style={{ fontWeight: '600', marginBottom: '0.5rem', color: 'var(--primary-color)' }}>
            <i className="fas fa-lightbulb"></i>
            Tips for best results:
          </div>
          <ul style={{ margin: '0', paddingLeft: '1.5rem' }}>
            <li>Use high-quality images for better results</li>
            <li>Smaller tile sizes create more detailed mosaics</li>
            <li>Images with good contrast work best</li>
            <li>Processing time depends on image size and tile size</li>
          </ul>
        </div>
      )}
    </div>
  );
};

export default MosaicImgContainer;
