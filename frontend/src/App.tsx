import { ChangeEvent, MouseEvent, useEffect, useRef, useState } from "react";
import imgPlaceholder from "./assets/img/img_placeholder.png";

type APIResponse = {
  mosaicImg: string;
  duration: number;
};

function App() {
  const [file, setFile] = useState<File>();
  const [tileSize, setTileSize] = useState<string>("10");

  const [mosaicImg, setMosaicImg] = useState<string>();
  const [duration, setDuration] = useState<number>(0);
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

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const target = event.target;
    if (target && target.files!.length > 0) {
      setFile(target.files![0]);
    }
  };

  const handleResponse = async (response: Response) => {
    try {
      const data = (await response.json()) as APIResponse;
      setMosaicImg(data.mosaicImg);
      setDuration(data.duration);
    } catch (error) {
      console.log(error);
    }
  };

  const handleError = () => {};

  const handleSubmit = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    if (file) {
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

  return (
    <div className="d-flex justify-content-center mt-5">
      <div className="container">
        <div className="row">
          <h1>Mosaic App</h1>
        </div>
        <div className="row">
          <div className="col-sm d-flex flex-column">
            <img
              src={imgPlaceholder}
              id="preview"
              alt="preview"
              ref={previewImg}
            />
            <form>
              <input
                type="file"
                name="imgUpload"
                id="imgUpload"
                onChange={handleChange}
                className="form-control-file mt-1"
              />
              <div className="form-group mt-1">
                <label htmlFor="tileSize">Select tile size</label>
                <select
                  name="tileSize"
                  id="tileSize"
                  onChange={(e) => setTileSize(e.target.value)}
                  value={tileSize}
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
                onClick={handleSubmit}
                type="submit"
                className="btn btn-primary mt-1"
                disabled={!file}
              >
                Upload
              </button>
            </form>
          </div>
          <div className="col-sm d-flex flex-column">
            <img
              src={
                mosaicImg
                  ? `data:image/jpeg;base64,${mosaicImg}`
                  : imgPlaceholder
              }
              alt="mosaic"
            />
            <span>{duration} seconds</span>
            <a
              href={mosaicImg ? `data:image/jpeg;base64,${mosaicImg}` : "#"}
              download={`mosaic-version-${file?.name}`}
              className={`btn btn-success ${mosaicImg ? "" : "disabled"}`}
              style={{ width: "fit-content" }}
            >
              Download
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
