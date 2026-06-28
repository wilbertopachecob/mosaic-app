import React, {
  ChangeEvent,
  DragEvent,
  MouseEvent,
  useCallback,
  useMemo,
  useState,
} from "react";
import { useTranslation } from "react-i18next";
import {
  HelpCircle,
  ImagePlus,
  Loader2,
  SlidersHorizontal,
  UploadCloud,
  Wand2,
} from "lucide-react";

interface UploadFormProps {
  selectedTileSize: string;
  selectedBlend: string;
  isBtnDisabled: boolean;
  isLoading: boolean;
  handleSubmit: () => void;
  handleFileChange: (file: File) => void;
  handleTileSizeChange: (tile: string) => void;
  handleBlendChange: (blend: string) => void;
}

const tileValues = ["5", "10", "15", "20", "25", "30", "50", "100"];

const UploadForm: React.FC<UploadFormProps> = ({
  selectedTileSize,
  selectedBlend,
  isBtnDisabled,
  isLoading,
  handleSubmit,
  handleFileChange,
  handleTileSizeChange,
  handleBlendChange,
}) => {
  const { t } = useTranslation();
  const [isDragging, setIsDragging] = useState(false);

  const tileOptions = useMemo(
    () =>
      tileValues.map((value) => ({
        value,
        label: t(`tileSize.${value}`),
      })),
    [t]
  );

  const validateAndSetFile = useCallback(
    (file: File) => {
      if (!file.type.startsWith("image/")) {
        alert(t("upload.invalidFile"));
        return;
      }

      const maxSize = 10 * 1024 * 1024;
      if (file.size > maxSize) {
        alert(t("upload.fileTooLarge"));
        return;
      }

      handleFileChange(file);
    },
    [handleFileChange, t]
  );

  const onSubmit = useCallback(
    (e: MouseEvent<HTMLButtonElement>) => {
      e.preventDefault();
      handleSubmit();
    },
    [handleSubmit]
  );

  const onFileChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (file) {
        validateAndSetFile(file);
      }
    },
    [validateAndSetFile]
  );

  const onDragOver = useCallback((event: DragEvent<HTMLLabelElement>) => {
    event.preventDefault();
    setIsDragging(true);
  }, []);

  const onDragLeave = useCallback(() => {
    setIsDragging(false);
  }, []);

  const onDrop = useCallback(
    (event: DragEvent<HTMLLabelElement>) => {
      event.preventDefault();
      setIsDragging(false);
      const file = event.dataTransfer.files?.[0];
      if (file) {
        validateAndSetFile(file);
      }
    },
    [validateAndSetFile]
  );

  const onTileSizeChange = useCallback(
    (event: ChangeEvent<HTMLSelectElement>) => {
      handleTileSizeChange(event.target.value);
    },
    [handleTileSizeChange]
  );

  const onBlendChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      handleBlendChange(event.target.value);
    },
    [handleBlendChange]
  );

  const blendPercent = Math.round(Number(selectedBlend) * 100);

  return (
    <form className="control-stack">
      <div className="control-group">
        <label
          htmlFor="imgUpload"
          className={`upload-zone ${isDragging ? "is-dragging" : ""}`}
          onDragOver={onDragOver}
          onDragLeave={onDragLeave}
          onDrop={onDrop}
        >
          <input
            type="file"
            name="imgUpload"
            id="imgUpload"
            onChange={onFileChange}
            accept="image/*"
            disabled={isLoading}
          />
          <span className="upload-icon">
            <UploadCloud size={24} aria-hidden="true" />
          </span>
          <span className="upload-copy">
            <strong>{t("upload.selectImage")}</strong>
            <small>{t("upload.dragHint")}</small>
          </span>
          <span className="upload-action">
            <ImagePlus size={16} aria-hidden="true" />
            {t("upload.browse")}
          </span>
        </label>
      </div>

      <div className="settings-grid">
        <div className="control-group">
          <label htmlFor="tileSize" className="field-label">
            <SlidersHorizontal size={16} aria-hidden="true" />
            {t("upload.tileSize")}
          </label>
          <select
            name="tileSize"
            id="tileSize"
            onChange={onTileSizeChange}
            value={selectedTileSize}
            className="field-control"
            disabled={isLoading}
          >
            {tileOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <p className="field-help">{t("upload.tileSizeHelp")}</p>
        </div>

        <div className="control-group">
          <div className="field-label-row">
            <label htmlFor="blend" className="field-label">
              <Wand2 size={16} aria-hidden="true" />
              {t("upload.sourceBlend")}
            </label>
            <span
              className="tooltip-icon"
              role="img"
              aria-label={t("upload.blendHelpAria")}
              title={t("upload.blendHelpTitle")}
            >
              <HelpCircle size={15} aria-hidden="true" />
            </span>
          </div>
          <div className="range-row">
            <input
              type="range"
              name="blend"
              id="blend"
              min="0"
              max="0.75"
              step="0.01"
              value={selectedBlend}
              onChange={onBlendChange}
              className="range-control"
              disabled={isLoading}
              aria-describedby="blendHelp"
            />
            <span className="range-value">{blendPercent}%</span>
          </div>
          <p id="blendHelp" className="field-help">
            {t("upload.blendHelp")}
          </p>
        </div>
      </div>

      <button
        onClick={onSubmit}
        type="submit"
        className="primary-action"
        disabled={isBtnDisabled || isLoading}
      >
        {isLoading ? (
          <>
            <Loader2 className="spin" size={18} aria-hidden="true" />
            {t("upload.generating")}
          </>
        ) : (
          <>
            <Wand2 size={18} aria-hidden="true" />
            {t("upload.generate")}
          </>
        )}
      </button>

      {isLoading && (
        <p className="processing-note">{t("upload.processing")}</p>
      )}
    </form>
  );
};

export default UploadForm;
