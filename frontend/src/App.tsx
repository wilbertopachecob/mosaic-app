import React, { useEffect, useRef, useState, useCallback } from "react";
import imgPlaceholder from "./assets/img/img_placeholder.png";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";
import LoadingSpinner from "./components/LoadingSpinner";
import ErrorMessage from "./components/ErrorMessage";

// Type definitions for better type safety
interface APIResponse {
  mosaicImg: string;
  duration: number;
}

interface APIError {
  error: string;
  message: string;
  code: number;
}

// App states for better UX
type AppState = 'idle' | 'loading' | 'success' | 'error';

function App() {
  // State management
  const [file, setFile] = useState<File | null>(null);
  const [tileSize, setTileSize] = useState<string>("20");
  const [mosaicImg, setMosaicImg] = useState<string | null>(null);
  const [duration, setDuration] = useState<number>(0);
  const [appState, setAppState] = useState<AppState>('idle');
  const [error, setError] = useState<string | null>(null);
  
  // Refs
  const previewImg = useRef<HTMLImageElement>(null);

  // Update preview image when file changes
  useEffect(() => {
    if (file && previewImg.current) {
      const reader = new FileReader();
      
      reader.onload = (e) => {
        if (previewImg.current && e.target?.result) {
          previewImg.current.src = e.target.result as string;
        }
      };
      
      reader.onerror = () => {
        setError("Failed to read the selected file");
        setAppState('error');
      };
      
      reader.readAsDataURL(file);
    }
  }, [file]);

  // Handle API response
  const handleResponse = useCallback(async (response: Response) => {
    if (!response.ok) {
      const errorData: APIError = await response.json();
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }
    
    const data: APIResponse = await response.json();
    setMosaicImg(data.mosaicImg);
    setDuration(data.duration);
    setAppState('success');
    setError(null);
  }, []);

  // Handle API errors
  const handleError = useCallback((error: Error) => {
    console.error("API Error:", error);
    setError(error.message || "An unexpected error occurred");
    setAppState('error');
  }, []);

  // Handle form submission
  const handleSubmit = useCallback(async () => {
    if (!file) {
      setError("Please select a file first");
      setAppState('error');
      return;
    }

    setAppState('loading');
    setError(null);

    try {
      const formData = new FormData();
      formData.append("imgUpload", file);
      formData.append("fileName", file.name);
      formData.append("tileSize", tileSize);

      const response = await fetch("/api/file/upload", {
        method: "POST",
        body: formData,
      });

      await handleResponse(response);
    } catch (error) {
      handleError(error as Error);
    }
  }, [file, tileSize, handleResponse, handleError]);

  // Reset the application state
  const handleReset = useCallback(() => {
    setFile(null);
    setMosaicImg(null);
    setDuration(0);
    setAppState('idle');
    setError(null);
    if (previewImg.current) {
      previewImg.current.src = imgPlaceholder;
    }
  }, []);

  return (
    <div className="min-vh-100 bg-light">
      <div className="container py-5">
        {/* Header */}
        <div className="row mb-4">
          <div className="col-12 text-center">
            <h1 className="display-4 fw-bold text-primary mb-2">
              🎨 Mosaic Generator
            </h1>
            <p className="lead text-muted">
              Transform your images into beautiful mosaics using our AI-powered tile matching
            </p>
          </div>
        </div>

        {/* Main Content */}
        <div className="row g-4">
          {/* Left Column - Upload and Preview */}
          <div className="col-lg-6">
            <div className="card shadow-sm h-100">
              <div className="card-header bg-primary text-white">
                <h5 className="mb-0">
                  <i className="fas fa-upload me-2"></i>
                  Upload Image
                </h5>
              </div>
              <div className="card-body d-flex flex-column">
                {/* Image Preview */}
                <div className="text-center mb-3">
                  <img
                    src={imgPlaceholder}
                    alt="Preview"
                    ref={previewImg}
                    className="img-fluid rounded shadow-sm"
                    style={{ maxHeight: '300px', objectFit: 'contain' }}
                  />
                </div>

                {/* Upload Form */}
                <UploadForm
                  selectedTileSize={tileSize}
                  isBtnDisabled={!file || appState === 'loading'}
                  handleSubmit={handleSubmit}
                  handleFileChange={setFile}
                  handleTileSizeChange={setTileSize}
                  isLoading={appState === 'loading'}
                />

                {/* Error Display */}
                {error && (
                  <ErrorMessage 
                    message={error} 
                    onDismiss={() => setError(null)} 
                  />
                )}
              </div>
            </div>
          </div>

          {/* Right Column - Mosaic Result */}
          <div className="col-lg-6">
            <div className="card shadow-sm h-100">
              <div className="card-header bg-success text-white">
                <h5 className="mb-0">
                  <i className="fas fa-image me-2"></i>
                  Mosaic Result
                </h5>
              </div>
              <div className="card-body d-flex flex-column">
                {appState === 'loading' && (
                  <div className="text-center py-5">
                    <LoadingSpinner />
                    <p className="mt-3 text-muted">Generating your mosaic...</p>
                  </div>
                )}
                
                {appState === 'success' && mosaicImg && (
                  <MosaicImgContainer
                    duration={duration}
                    mosaicImg={mosaicImg}
                    fileName={file?.name}
                    onReset={handleReset}
                  />
                )}
                
                {appState === 'idle' && (
                  <div className="text-center py-5 text-muted">
                    <i className="fas fa-image fa-3x mb-3"></i>
                    <p>Upload an image to generate your mosaic</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="row mt-5">
          <div className="col-12 text-center">
            <p className="text-muted small">
              Built with React, TypeScript, and Go • 
              <a href="https://github.com/your-repo" className="text-decoration-none ms-1">
                View on GitHub
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
