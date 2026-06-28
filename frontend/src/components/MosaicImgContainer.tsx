import React, { useCallback } from "react";
import {
  Clock3,
  Download,
  FileImage,
  Grid3X3,
  Image,
  Loader2,
  RefreshCw,
  SlidersHorizontal,
} from "lucide-react";

interface MosaicImgContainerProps {
  mosaicImg: string | null;
  duration: number;
  fileName: string | undefined;
  tileSize: string;
  blend: string;
  isLoading: boolean;
  hasSourceImage: boolean;
  onReset?: () => void;
}

const MosaicImgContainer: React.FC<MosaicImgContainerProps> = ({
  mosaicImg,
  duration,
  fileName,
  tileSize,
  blend,
  isLoading,
  hasSourceImage,
  onReset,
}) => {
  const handleDownload = useCallback(() => {
    if (!mosaicImg) return;

    const link = document.createElement("a");
    link.href = `data:image/jpeg;base64,${mosaicImg}`;
    link.download = `mosaic-${fileName || "image"}.jpg`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }, [mosaicImg, fileName]);

  const formatDuration = useCallback((seconds: number) => {
    if (seconds < 1) {
      return `${Math.round(seconds * 1000)}ms`;
    }
    return `${seconds.toFixed(2)}s`;
  }, []);

  const displayFileName = fileName
    ? fileName.substring(0, 28) + (fileName.length > 28 ? "..." : "")
    : "No source selected";
  const blendPercent = `${Math.round(Number(blend) * 100)}%`;

  return (
    <div className="result-stack">
      <div className={`result-frame ${mosaicImg ? "has-result" : ""}`}>
        {isLoading ? (
          <div className="result-state">
            <Loader2 className="spin" size={40} aria-hidden="true" />
            <strong>Building your mosaic</strong>
            <span>Matching tiles and blending the final image.</span>
          </div>
        ) : mosaicImg ? (
          <img
            src={`data:image/jpeg;base64,${mosaicImg}`}
            alt="Generated mosaic result"
          />
        ) : (
          <div className="result-state">
            <Image size={44} aria-hidden="true" />
            <strong>{hasSourceImage ? "Ready to generate" : "Result preview"}</strong>
            <span>
              {hasSourceImage
                ? "Your mosaic will appear here after generation."
                : "Select an image and tune the controls to begin."}
            </span>
          </div>
        )}
      </div>

      <dl className="result-metadata" aria-label="Mosaic output details">
        <div>
          <dt>
            <Clock3 size={15} aria-hidden="true" />
            Time
          </dt>
          <dd>{mosaicImg ? formatDuration(duration) : "-"}</dd>
        </div>
        <div>
          <dt>
            <FileImage size={15} aria-hidden="true" />
            File
          </dt>
          <dd>{displayFileName}</dd>
        </div>
        <div>
          <dt>
            <Grid3X3 size={15} aria-hidden="true" />
            Tile
          </dt>
          <dd>{tileSize}px</dd>
        </div>
        <div>
          <dt>
            <SlidersHorizontal size={15} aria-hidden="true" />
            Blend
          </dt>
          <dd>{blendPercent}</dd>
        </div>
      </dl>

      <div className="result-actions">
        <button
          onClick={handleDownload}
          className="download-action"
          disabled={!mosaicImg}
          title="Download mosaic image"
        >
          <Download size={18} aria-hidden="true" />
          Download mosaic
        </button>

        {onReset && (
          <button
            onClick={onReset}
            className="secondary-action"
            type="button"
            title="Start over with a new image"
          >
            <RefreshCw size={16} aria-hidden="true" />
            Reset
          </button>
        )}
      </div>
    </div>
  );
};

export default MosaicImgContainer;
