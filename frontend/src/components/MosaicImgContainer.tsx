import imgPlaceholder from "../assets/img/img_placeholder.png";

const MosaicImgContainer = ({
  mosaicImg,
  duration,
  fileName,
}: {
  mosaicImg: string | undefined;
  duration: number;
  fileName: string | undefined;
}) => {
  return (
    <div className="col-sm d-flex flex-column">
      <img
        src={mosaicImg ? `data:image/jpeg;base64,${mosaicImg}` : imgPlaceholder}
        alt="mosaic"
      />
      <span>{duration} seconds</span>
      <a
        href={mosaicImg ? `data:image/jpeg;base64,${mosaicImg}` : "#"}
        download={`mosaic-version-${fileName}`}
        className={`btn btn-success ${mosaicImg ? "" : "disabled"}`}
        style={{ width: "fit-content" }}
      >
        Download
      </a>
    </div>
  );
};

export default MosaicImgContainer;
