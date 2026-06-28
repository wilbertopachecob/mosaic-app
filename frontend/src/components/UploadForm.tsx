import React, {
  ChangeEvent,
  DragEvent,
  MouseEvent,
  useCallback,
  useState,
} from "react";
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

const tileOptions = [
  { value: "5", label: "5px - Very fine" },
  { value: "10", label: "10px - Fine" },
  { value: "15", label: "15px - Medium" },
  { value: "20", label: "20px - Standard" },
  { value: "25", label: "25px - Coarse" },
  { value: "30", label: "30px - Very coarse" },
  { value: "50", label: "50px - Large" },
  { value: "100", label: "100px - Poster blocks" },
];

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
  const [isDragging, setIsDragging] = useState(false);

  const validateAndSetFile = useCallback(
    (file: File) => {
      if (!file.type.startsWith("image/")) {
        alert("Please select an image file");
        return;
      }

      const maxSize = 10 * 1024 * 1024;
      if (file.size > maxSize) {
        alert("File size must be less than 10MB");
        return;
      }

      handleFileChange(file);
    },
    [handleFileChange]
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
            <strong>Select image</strong>
            <small>Drag a JPG, PNG, or WebP here. Max 10MB.</small>
          </span>
          <span className="upload-action">
            <ImagePlus size={16} aria-hidden="true" />
            Browse
          </span>
        </label>
      </div>

      <div className="settings-grid">
        <div className="control-group">
          <label htmlFor="tileSize" className="field-label">
            <SlidersHorizontal size={16} aria-hidden="true" />
            Tile size
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
          <p className="field-help">Smaller tiles preserve more detail.</p>
        </div>

        <div className="control-group">
          <div className="field-label-row">
            <label htmlFor="blend" className="field-label">
              <Wand2 size={16} aria-hidden="true" />
              Source blend
            </label>
            <span
              className="tooltip-icon"
              role="img"
              aria-label="Blend help"
              title="Controls how much of the original image is blended over the mosaic. Lower values show purer photo tiles; higher values look more like the source image."
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
            Lower is purer tiles; higher preserves faces and edges.
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
            Generating mosaic
          </>
        ) : (
          <>
            <Wand2 size={18} aria-hidden="true" />
            Generate mosaic
          </>
        )}
      </button>

      {isLoading && (
        <p className="processing-note">
          Processing may take a few moments depending on image size and tile size.
        </p>
      )}
    </form>
  );
};

export default UploadForm;
