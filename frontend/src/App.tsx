import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Images, Sparkles, Wand2 } from "lucide-react";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";
import ErrorMessage from "./components/ErrorMessage";
import HeaderControls from "./components/HeaderControls";
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
  const { t } = useTranslation();
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

  const handleError = useCallback(
    (error: Error) => {
      console.error("API Error:", error);
      setError(error.message || t("input.unexpectedError"));
      setAppState("error");
    },
    [t]
  );

  const handleSubmit = useCallback(async () => {
    if (!file) {
      setError(t("input.selectBeforeGenerate"));
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
  }, [file, tileSize, blend, handleResponse, handleError, t]);

  const handleReset = useCallback(() => {
    setFile(null);
    setMosaicImg(null);
    setDuration(0);
    setError(null);
    setAppState("idle");
  }, []);

  const formatFileSize = useCallback(
    (bytes: number): string => {
      if (bytes === 0) return `0 ${t("units.bytes")}`;
      const k = 1024;
      const sizes = [
        t("units.bytes"),
        t("units.kb"),
        t("units.mb"),
        t("units.gb"),
      ];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
    },
    [t]
  );

  return (
    <main className="app-shell">
      <div className="workspace">
        <header className="app-header">
          <div>
            <div className="eyebrow">
              <Sparkles size={16} aria-hidden="true" />
              {t("header.eyebrow")}
            </div>
            <h1 className="app-title">{t("header.title")}</h1>
            <p className="app-subtitle">{t("header.subtitle")}</p>
          </div>
          <HeaderControls />
        </header>

        <section
          className="workspace-grid"
          aria-label={t("workspace.ariaLabel")}
        >
          <section className="panel input-panel" aria-labelledby="input-heading">
            <div className="panel-heading">
              <div>
                <p className="section-kicker">{t("input.step")}</p>
                <h2 id="input-heading">{t("input.title")}</h2>
              </div>
              <Wand2 size={20} aria-hidden="true" />
            </div>

            <div className={`preview-frame ${previewUrl ? "has-image" : ""}`}>
              {previewUrl ? (
                <img src={previewUrl} alt={t("input.previewAlt")} />
              ) : (
                <div className="preview-empty">
                  <Images size={42} aria-hidden="true" />
                  <strong>{t("input.noImage")}</strong>
                  <span>{t("input.previewHint")}</span>
                </div>
              )}
            </div>

            {file && (
              <dl className="file-summary" aria-label={t("input.title")}>
                <div>
                  <dt>{t("input.file")}</dt>
                  <dd>{file.name}</dd>
                </div>
                <div>
                  <dt>{t("input.size")}</dt>
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
                <p className="section-kicker">{t("output.step")}</p>
                <h2 id="output-heading">{t("output.title")}</h2>
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
