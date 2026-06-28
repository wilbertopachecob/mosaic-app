import { formatDuration, formatFileSize, truncateFileName } from "../format";

describe("formatDuration", () => {
  it("shows milliseconds for sub-second durations", () => {
    expect(formatDuration(0.5)).toBe("500ms");
    expect(formatDuration(0.001)).toBe("1ms");
  });

  it("shows seconds with two decimal places for longer durations", () => {
    expect(formatDuration(1)).toBe("1.00s");
    expect(formatDuration(2.456)).toBe("2.46s");
  });
});

describe("truncateFileName", () => {
  it("returns the name unchanged when within the limit", () => {
    expect(truncateFileName("photo.jpg")).toBe("photo.jpg");
    expect(truncateFileName("a".repeat(28))).toBe("a".repeat(28));
  });

  it("truncates and appends ellipsis when over the limit", () => {
    const longName = "a".repeat(35) + ".jpg";
    expect(truncateFileName(longName)).toBe("a".repeat(28) + "...");
  });
});

describe("formatFileSize", () => {
  const units: [string, string, string, string] = ["B", "KB", "MB", "GB"];

  it("formats zero bytes", () => {
    expect(formatFileSize(0, units)).toBe("0 B");
  });

  it("formats bytes, kilobytes, and megabytes", () => {
    expect(formatFileSize(512, units)).toBe("512 B");
    expect(formatFileSize(1024, units)).toBe("1 KB");
    expect(formatFileSize(1536, units)).toBe("1.5 KB");
    expect(formatFileSize(1048576, units)).toBe("1 MB");
  });
});
