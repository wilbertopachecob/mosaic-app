import { ChangeEvent, MouseEvent, useState } from "react";

type APIResponse = {
  mosaicImg: string;
  duration: number;
};

function App() {
  const [file, setFile] = useState<File>();
  const [tileSize, setTileSize] = useState<string>("10");

  const [mosaicImg, setMosaicImg] = useState<string>();
  const [duration, setDuration] = useState<number>();

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const target = event.target;
    if (target && target.files!.length > 0) {
      setFile(target.files![0]);
    }
  }

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
          <div className="col-sm">
            <form>
              <input
                type="file"
                name="imgUpload"
                id="imgUpload"
                onChange={handleChange}
                className="form-control-file"
              />
              <div className="form-group">
                <label htmlFor="tileSize">Select tile size</label>
                <select
                  name="tileSize"
                  id="tileSize"
                  onChange={(e) => setTileSize(e.target.value)}
                  value={tileSize}
                  className="form-control"
                >
                  <option value="10" selected>
                    10
                  </option>
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
                className="btn btn-primary"
              >
                Upload
              </button>
            </form>
          </div>
          <div className="col-sm">
            {mosaicImg && (
              <div className="d-flex flex-column">
                <img
                  src={`data:image/jpeg;base64,${mosaicImg}`}
                  alt="mosaic"
                  width="25%"
                />
                <span>{duration} seconds</span>
                <a
                  href={`data:image/jpeg;base64,${mosaicImg}`}
                  download={`mosaic-version-${file?.name}`}
                  className="btn btn-success"
                  style={{ width: "fit-content" }}
                >
                  Download
                </a>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
