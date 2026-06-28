import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Images, Wand2 } from "lucide-react";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";
import ErrorMessage from "./components/ErrorMessage";
import HeaderControls from "./components/HeaderControls";
import Logo from "./components/Logo";
import { parseJson } from "./utils/api";
import { formatFileSize } from "./utils/format";
import "./App.css";

/** Successful mosaic generation response from the upload API. */
type APIResponse = {
  mosaicImg: string;
  duration: number;
};

/** Error payload returned by the upload API. */
type APIError = {
  error: string;
  message: string;
  code: number;
};

/** Tracks the overall upload and generation lifecycle in the UI. */
type AppState = "idle" | "loading" | "success" | "error";

/**
 * Root application view: image upload, mosaic settings, generation, and result display.
 */
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

  /** Replaces the selected source image and clears any prior result. */
  const handleFileChange = useCallback((nextFile: File) => {
    setFile(nextFile);
    setMosaicImg(null);
    setDuration(0);
    setError(null);
    setAppState("idle");
  }, []);

  /** Maps a successful or failed upload response into UI state. */
  const handleResponse = useCallback(
    async (response: Response) => {
      if (!response.ok) {
        const errorData = await parseJson<APIError>(response);
        if (errorData?.message) {
          throw new Error(errorData.message);
        }
        if ([500, 502, 503, 504].includes(response.status)) {
          throw new Error(t("input.serverUnreachable"));
        }
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await parseJson<APIResponse>(response);
      if (!data) {
        throw new Error(t("input.unexpectedError"));
      }
      setMosaicImg(data.mosaicImg);
      setDuration(data.duration);
      setAppState("success");
      setError(null);
    },
    [t]
  );

  /** Surfaces API and network failures to the user. */
  const handleError = useCallback(
    (error: Error) => {
      console.error("API Error:", error);
      setError(error.message || t("input.unexpectedError"));
      setAppState("error");
    },
    [t]
  );

  /** Posts the selected image and mosaic settings to the upload API. */
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

  /** Clears the current session and returns the workspace to its initial state. */
  const handleReset = useCallback(() => {
    setFile(null);
    setMosaicImg(null);
    setDuration(0);
    setError(null);
    setAppState("idle");
  }, []);

  const formatSize = useCallback(
    (bytes: number): string =>
      formatFileSize(bytes, [
        t("units.bytes"),
        t("units.kb"),
        t("units.mb"),
        t("units.gb"),
      ]),
    [t]
  );

  const isLoading = appState === "loading";

  return (
    <main className="app-shell">
      <div className="workspace">
        <header className="app-header">
          <div>
            <Logo />
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
                  <dd>{formatSize(file.size)}</dd>
                </div>
              </dl>
            )}

            <UploadForm
              selectedTileSize={tileSize}
              selectedBlend={blend}
              isBtnDisabled={!file || isLoading}
              handleSubmit={handleSubmit}
              handleFileChange={handleFileChange}
              handleTileSizeChange={setTileSize}
              handleBlendChange={setBlend}
              isLoading={isLoading}
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
              isLoading={isLoading}
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
