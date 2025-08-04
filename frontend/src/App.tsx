import { useEffect, useRef, useState } from "react";
import imgPlaceholder from "./assets/img/img_placeholder.png";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";

type APIResponse = {
  mosaicImg: string;
  duration: number;
};

function App() {
  const [file, setFile] = useState<File>();
  const [tileSize, setTileSize] = useState<string>("10");
  const [mosaicImg, setMosaicImg] = useState<string>();
  const [duration, setDuration] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const previewImg = useRef<HTMLImageElement>(null);

  useEffect(() => {
    if (file && previewImg && previewImg.current !== null) {
      const reader = new FileReader();

      reader.onload = function (e) {
        previewImg.current!.setAttribute("src", e.target?.result as string);
      };

      reader.readAsDataURL(file as Blob);
    }
  }, [file]);

  const handleResponse = async (response: Response) => {
    try {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = (await response.json()) as APIResponse;
      setMosaicImg(data.mosaicImg);
      setDuration(data.duration);
      setError("");
    } catch (error) {
      console.error("Error processing response:", error);
      setError("Failed to generate mosaic. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleError = (error: any) => {
    console.error("Network error:", error);
    setError("Network error. Please check your connection and try again.");
    setIsLoading(false);
  };

  const handleSubmit = () => {
    if (file) {
      setIsLoading(true);
      setError("");
      const url = "/api/file/upload";
      const formData = new FormData();
      formData.append("imgUpload", file);
      formData.append("fileName", file.name);
      formData.append("tileSize", tileSize!);

      fetch(url, { method: "POST", body: formData })
        .then(handleResponse)
        .catch(handleError);
    }
  };

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
              isBtnDisabled={!file || isLoading}
              handleSubmit={handleSubmit}
              handleFileChange={setFile}
              handleTileSizeChange={setTileSize}
              isLoading={isLoading}
            />

            {error && (
              <div className="error-message slide-up" style={{
                color: 'var(--error-color)',
                padding: '1rem',
                marginTop: '1rem',
                backgroundColor: 'rgba(255, 107, 107, 0.1)',
                borderRadius: '8px',
                border: '1px solid var(--error-color)'
              }}>
                <i className="fas fa-exclamation-triangle"></i>
                {error}
              </div>
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
              isLoading={isLoading}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
