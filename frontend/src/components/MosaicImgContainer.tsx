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
import { formatDuration, truncateFileName } from "@/utils/format";

/** Props for the mosaic result panel, including metadata and actions. */
type MosaicImgContainerProps = {
  mosaicImg: string | null;
  duration: number;
  fileName: string | undefined;
  tileSize: string;
  blend: string;
  isLoading: boolean;
  hasSourceImage: boolean;
  onReset?: () => void;
};

/**
 * Displays the generated mosaic, generation metadata, and download/reset actions.
 */
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

  /** Triggers a browser download of the base64-encoded mosaic JPEG. */
  const handleDownload = useCallback(() => {
    if (!mosaicImg) return;

    const link = document.createElement("a");
    link.href = `data:image/jpeg;base64,${mosaicImg}`;
    link.download = `mosaic-${fileName || "image"}.jpg`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }, [mosaicImg, fileName]);

  const displayFileName = fileName
    ? truncateFileName(fileName)
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
