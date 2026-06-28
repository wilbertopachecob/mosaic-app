import React, { useCallback } from "react";
import { useTranslation } from "react-i18next";
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

type MosaicImgContainerProps = {
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
  const { t } = useTranslation();

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
    : t("result.noSource");
  const blendPercent = `${Math.round(Number(blend) * 100)}%`;

  return (
    <div className="result-stack">
      <div className={`result-frame ${mosaicImg ? "has-result" : ""}`}>
        {isLoading ? (
          <div className="result-state">
            <Loader2 className="spin" size={40} aria-hidden="true" />
            <strong>{t("result.building")}</strong>
            <span>{t("result.buildingHint")}</span>
          </div>
        ) : mosaicImg ? (
          <img
            src={`data:image/jpeg;base64,${mosaicImg}`}
            alt={t("result.mosaicAlt")}
          />
        ) : (
          <div className="result-state">
            <Image size={44} aria-hidden="true" />
            <strong>
              {hasSourceImage ? t("result.ready") : t("result.preview")}
            </strong>
            <span>
              {hasSourceImage
                ? t("result.readyHint")
                : t("result.previewHint")}
            </span>
          </div>
        )}
      </div>

      <dl className="result-metadata" aria-label={t("result.metadataAria")}>
        <div>
          <dt>
            <Clock3 size={15} aria-hidden="true" />
            {t("result.time")}
          </dt>
          <dd>{mosaicImg ? formatDuration(duration) : "-"}</dd>
        </div>
        <div>
          <dt>
            <FileImage size={15} aria-hidden="true" />
            {t("result.file")}
          </dt>
          <dd>{displayFileName}</dd>
        </div>
        <div>
          <dt>
            <Grid3X3 size={15} aria-hidden="true" />
            {t("result.tile")}
          </dt>
          <dd>{tileSize}px</dd>
        </div>
        <div>
          <dt>
            <SlidersHorizontal size={15} aria-hidden="true" />
            {t("result.blend")}
          </dt>
          <dd>{blendPercent}</dd>
        </div>
      </dl>

      <div className="result-actions">
        <button
          onClick={handleDownload}
          className="download-action"
          disabled={!mosaicImg}
          title={t("result.downloadTitle")}
        >
          <Download size={18} aria-hidden="true" />
          {t("result.download")}
        </button>

        {onReset && (
          <button
            onClick={onReset}
            className="secondary-action"
            type="button"
            title={t("result.resetTitle")}
          >
            <RefreshCw size={16} aria-hidden="true" />
            {t("result.reset")}
          </button>
        )}
      </div>
    </div>
  );
};

export default MosaicImgContainer;
