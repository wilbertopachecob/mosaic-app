import { useEffect, useRef, useState } from "react";
import imgPlaceholder from "./assets/img/img_placeholder.png";
import MosaicImgContainer from "./components/MosaicImgContainer";
import UploadForm from "./components/UploadForm";

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

  const handleSubmit = () => {
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
            <UploadForm
              selectedTileSize={tileSize}
              isBtnDisabled={!file}
              handleSubmit={handleSubmit}
              handleFileChange={setFile}
              handleTileSizeChange={setTileSize}
            />
          </div>
          <MosaicImgContainer
            duration={duration}
            mosaicImg={mosaicImg}
            fileName={file?.name}
          />
        </div>
      </div>
    </div>
  );
}

export default App;
