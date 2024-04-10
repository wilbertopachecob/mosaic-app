import { ChangeEvent, MouseEvent } from "react";

type UploadFormProps = {
  selectedTileSize: string;
  isBtnDisabled: boolean;
  handleSubmit: () => void;
  handleFileChange: (file: File) => void;
  handleTileSizeChange: (tile: string) => void;
};

const UploadForm = ({
  selectedTileSize,
  isBtnDisabled,
  handleSubmit,
  handleFileChange,
  handleTileSizeChange,
}: UploadFormProps) => {
  const onSubmit = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    handleSubmit();
  };

  const onFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target;
    if (target && target.files!.length > 0) {
      handleFileChange(target.files![0]);
    }
  };

  return (
    <form className="mt-3">
      <input
        type="file"
        name="imgUpload"
        id="imgUpload"
        onChange={onFileChange}
        className="form-control"
      />
      <div className="form-group mt-1">
        <label htmlFor="tileSize">Select tile size</label>
        <select
          name="tileSize"
          id="tileSize"
          onChange={(e) => handleTileSizeChange(e.target.value)}
          value={selectedTileSize}
          className="form-select"
        >
          <option value="10">10</option>
          <option value="15">15</option>
          <option value="20">20</option>
          <option value="25">25</option>
          <option value="50">50</option>
          <option value="100">100</option>
        </select>
      </div>
      <button
        onClick={onSubmit}
        type="submit"
        className="btn btn-primary mt-1"
        disabled={isBtnDisabled}
      >
        Upload
      </button>
    </form>
  );
};

export default UploadForm;
