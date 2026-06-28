import React from "react";
import App from "@/App";
import ErrorMessage from "@/components/ErrorMessage";
import HeaderControls from "@/components/HeaderControls";
import MosaicImgContainer from "@/components/MosaicImgContainer";
import UploadForm from "@/components/UploadForm";
import { renderWithProviders } from "@/testUtils";
import { expectNoA11yViolations } from "@/testUtils/a11y";

const uploadFormProps = {
  selectedTileSize: "20",
  selectedBlend: "0.42",
  isBtnDisabled: false,
  isLoading: false,
  handleSubmit: vi.fn(),
  handleFileChange: vi.fn(),
  handleTileSizeChange: vi.fn(),
  handleBlendChange: vi.fn(),
};

describe("accessibility", () => {
  it("App has no axe violations on initial render", async () => {
    await expectNoA11yViolations(renderWithProviders(<App />));
  });

  it("UploadForm has no axe violations", async () => {
    await expectNoA11yViolations(
      renderWithProviders(<UploadForm {...uploadFormProps} />)
    );
  });

  it("UploadForm has no axe violations while loading", async () => {
    await expectNoA11yViolations(
      renderWithProviders(
        <UploadForm {...uploadFormProps} isLoading isBtnDisabled />
      )
    );
  });

  it("HeaderControls has no axe violations", async () => {
    await expectNoA11yViolations(
      renderWithProviders(<HeaderControls />)
    );
  });

  it("ErrorMessage has no axe violations", async () => {
    await expectNoA11yViolations(
      renderWithProviders(
        <ErrorMessage message="Something went wrong" onDismiss={vi.fn()} />
      )
    );
  });

  it("MosaicImgContainer has no axe violations in preview state", async () => {
    await expectNoA11yViolations(
      renderWithProviders(
        <MosaicImgContainer
          mosaicImg={null}
          duration={0}
          tileSize="20"
          blend="0.42"
          isLoading={false}
          hasSourceImage={false}
          onReset={vi.fn()}
        />
      )
    );
  });

  it("MosaicImgContainer has no axe violations with a generated result", async () => {
    await expectNoA11yViolations(
      renderWithProviders(
        <MosaicImgContainer
          mosaicImg="abc123"
          duration={1.2}
          fileName="photo.jpg"
          tileSize="20"
          blend="0.42"
          isLoading={false}
          hasSourceImage
          onReset={vi.fn()}
        />
      )
    );
  });

  it("MosaicImgContainer has no axe violations while loading", async () => {
    await expectNoA11yViolations(
      renderWithProviders(
        <MosaicImgContainer
          mosaicImg={null}
          duration={0}
          tileSize="20"
          blend="0.42"
          isLoading
          hasSourceImage
        />
      )
    );
  });
});
