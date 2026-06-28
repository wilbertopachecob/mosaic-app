import React, { useCallback, useEffect, useState } from "react";
import { Images, Sparkles, Wand2 } from "lucide-react";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";
import ErrorMessage from "./components/ErrorMessage";
import "./App.css";

interface APIResponse {
  mosaicImg: string;
  duration: number;
}

interface APIError {
  error: string;
  message: string;
  code: number;
}

type AppState = "idle" | "loading" | "success" | "error";

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [tileSize, setTileSize] = useState<string>("20");
  const [blend, setBlend] = useState<string>("0.42");
  const [mosaicImg, setMosaicImg] = useState<string | null>(null);
  const [duration, setDuration] = useState<number>(0);
  const [appState, setAppState] = useState<AppState>("idle");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!file) {
      setPreviewUrl(null);
      return;
    }

    const objectUrl = URL.createObjectURL(file);
    setPreviewUrl(objectUrl);

    return () => URL.revokeObjectURL(objectUrl);
  }, [file]);

  const handleFileChange = useCallback((nextFile: File) => {
    setFile(nextFile);
    setMosaicImg(null);
    setDuration(0);
    setError(null);
    setAppState("idle");
  }, []);

  const handleResponse = useCallback(async (response: Response) => {
    if (!response.ok) {
      const errorData: APIError = await response.json();
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }

    const data: APIResponse = await response.json();
    setMosaicImg(data.mosaicImg);
    setDuration(data.duration);
    setAppState("success");
    setError(null);
  }, []);

  const handleError = useCallback((error: Error) => {
    console.error("API Error:", error);
    setError(error.message || "An unexpected error occurred");
    setAppState("error");
  }, []);

  const handleSubmit = useCallback(async () => {
    if (!file) {
      setError("Select an image before generating a mosaic.");
      setAppState("error");
      return;
    }

    setAppState("loading");
    setError(null);

    try {
      const formData = new FormData();
      formData.append("imgUpload", file);
      formData.append("fileName", file.name);
      formData.append("tileSize", tileSize);
      formData.append("blend", blend);

      const response = await fetch("/api/file/upload", {
        method: "POST",
        body: formData,
      });

      await handleResponse(response);
    } catch (error) {
      handleError(error as Error);
    }
  }, [file, tileSize, blend, handleResponse, handleError]);

  const handleReset = useCallback(() => {
    setFile(null);
    setMosaicImg(null);
    setDuration(0);
    setError(null);
    setAppState("idle");
  }, []);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  };

  return (
    <main className="app-shell">
      <div className="workspace">
        <header className="app-header">
          <div>
            <div className="eyebrow">
              <Sparkles size={16} aria-hidden="true" />
              Photo mosaic workspace
            </div>
            <h1 className="app-title">Mosaic Generator</h1>
            <p className="app-subtitle">
              Upload an image, tune the mosaic texture, and generate a polished
              downloadable result from your local tile library.
            </p>
          </div>
          <div className="header-metric" aria-label="Generator method">
            <Images size={18} aria-hidden="true" />
            <span>Photo tiles + source blend</span>
          </div>
        </header>

        <section className="workspace-grid" aria-label="Mosaic generator workspace">
          <section className="panel input-panel" aria-labelledby="input-heading">
            <div className="panel-heading">
              <div>
                <p className="section-kicker">Step 1</p>
                <h2 id="input-heading">Source image</h2>
              </div>
              <Wand2 size={20} aria-hidden="true" />
            </div>

            <div className={`preview-frame ${previewUrl ? "has-image" : ""}`}>
              {previewUrl ? (
                <img src={previewUrl} alt="Selected source preview" />
              ) : (
                <div className="preview-empty">
                  <Images size={42} aria-hidden="true" />
                  <strong>No image selected</strong>
                  <span>Your preview will appear here before generation.</span>
                </div>
              )}
            </div>

            {file && (
              <dl className="file-summary" aria-label="Selected file details">
                <div>
                  <dt>File</dt>
                  <dd>{file.name}</dd>
                </div>
                <div>
                  <dt>Size</dt>
                  <dd>{formatFileSize(file.size)}</dd>
                </div>
              </dl>
            )}

            <UploadForm
              selectedTileSize={tileSize}
              selectedBlend={blend}
              isBtnDisabled={!file || appState === "loading"}
              handleSubmit={handleSubmit}
              handleFileChange={handleFileChange}
              handleTileSizeChange={setTileSize}
              handleBlendChange={setBlend}
              isLoading={appState === "loading"}
            />

            {error && (
              <ErrorMessage
                message={error}
                onDismiss={() => {
                  setError(null);
                  setAppState("idle");
                }}
              />
            )}
          </section>

          <section className="panel output-panel" aria-labelledby="output-heading">
            <div className="panel-heading">
              <div>
                <p className="section-kicker">Step 2</p>
                <h2 id="output-heading">Generated mosaic</h2>
              </div>
            </div>

            <MosaicImgContainer
              duration={duration}
              mosaicImg={mosaicImg}
              fileName={file?.name}
              tileSize={tileSize}
              blend={blend}
              isLoading={appState === "loading"}
              hasSourceImage={Boolean(file)}
              onReset={handleReset}
            />
          </section>
        </section>
      </div>
    </main>
  );
}

export default App;
