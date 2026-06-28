import { screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom";
import MosaicImgContainer from "../MosaicImgContainer";
import { renderWithProviders } from "../../testUtils";

const defaultProps = {
  mosaicImg: null as string | null,
  duration: 0,
  fileName: undefined as string | undefined,
  tileSize: "20",
  blend: "0.42",
  isLoading: false,
  hasSourceImage: false,
};

describe("MosaicImgContainer", () => {
  it("shows the preview placeholder when no mosaic exists", () => {
    renderWithProviders(<MosaicImgContainer {...defaultProps} />);

    expect(screen.getByText(/result preview/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /download mosaic/i })).toBeDisabled();
  });

  it("shows the ready state when a source image is selected", () => {
    renderWithProviders(
      <MosaicImgContainer {...defaultProps} hasSourceImage />
    );

    expect(screen.getByText(/ready to generate/i)).toBeInTheDocument();
  });

  it("shows loading state while generating", () => {
    renderWithProviders(
      <MosaicImgContainer {...defaultProps} isLoading hasSourceImage />
    );

    expect(screen.getByText(/building your mosaic/i)).toBeInTheDocument();
  });

  it("renders mosaic metadata when a result is available", () => {
    renderWithProviders(
      <MosaicImgContainer
        {...defaultProps}
        mosaicImg="abc123"
        duration={1.25}
        fileName="vacation-photo.jpg"
        hasSourceImage
      />
    );

    expect(screen.getByRole("img", { name: /generated mosaic/i })).toBeInTheDocument();
    expect(screen.getByText("1.25s")).toBeInTheDocument();
    expect(screen.getByText("vacation-photo.jpg")).toBeInTheDocument();
    expect(screen.getByText("20px")).toBeInTheDocument();
    expect(screen.getByText("42%")).toBeInTheDocument();
  });

  it("truncates long filenames in metadata", () => {
    const longName = "a".repeat(35) + ".jpg";
    renderWithProviders(
      <MosaicImgContainer
        {...defaultProps}
        mosaicImg="abc123"
        duration={0.5}
        fileName={longName}
        hasSourceImage
      />
    );

    expect(screen.getByText("a".repeat(28) + "...")).toBeInTheDocument();
    expect(screen.getByText("500ms")).toBeInTheDocument();
  });

  it("calls onReset when the reset button is clicked", () => {
    const onReset = jest.fn();
    renderWithProviders(
      <MosaicImgContainer
        {...defaultProps}
        mosaicImg="abc123"
        hasSourceImage
        onReset={onReset}
      />
    );

    fireEvent.click(screen.getByRole("button", { name: /^reset$/i }));
    expect(onReset).toHaveBeenCalledTimes(1);
  });
});
