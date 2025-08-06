import React, { useEffect, useRef, useState, useCallback } from "react";
import imgPlaceholder from "./assets/img/img_placeholder.png";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";
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

  // Format file size for display
  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="App">
      <div className="container">
        {/* Header */}
        <header className="app-header fade-in">
          <h1 className="app-title">
            <i className="fas fa-palette"></i>
            Mosaic Generator
          </h1>
          <p className="app-subtitle">
            Transform your images into beautiful mosaics using AI-powered tile matching
          </p>
        </header>

        {/* Main Content */}
        <div className="main-content">
          {/* Upload Section */}
          <div className="card slide-up">
            <h2 className="card-title">
              <i className="fas fa-upload"></i>
              Upload Image
            </h2>
            
            <div className="image-container">
              <img
                src={imgPlaceholder}
                id="preview"
                alt="preview"
                ref={previewImg}
              />
            </div>

            {file && (
              <div className="stats-container slide-up">
                <div className="stat-item">
                  <i className="fas fa-file-image"></i>
                  <span>File: {file.name}</span>
                </div>
                <div className="stat-item">
                  <i className="fas fa-weight-hanging"></i>
                  <span>Size: {formatFileSize(file.size)}</span>
                </div>
              </div>
            )}

            <UploadForm
              selectedTileSize={tileSize}
              isBtnDisabled={!file || appState === 'loading'}
              handleSubmit={handleSubmit}
              handleFileChange={setFile}
              handleTileSizeChange={setTileSize}
              isLoading={appState === 'loading'}
            />

            {error && (
              <ErrorMessage 
                message={error} 
                onDismiss={() => {
                  setError(null);
                  setAppState('idle');
                }}
              />
            )}
          </div>

          {/* Result Section */}
          <div className="card slide-up">
            <h2 className="card-title">
              <i className="fas fa-image"></i>
              Mosaic Result
            </h2>
            
            <MosaicImgContainer
              duration={duration}
              mosaicImg={mosaicImg}
              fileName={file?.name}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
