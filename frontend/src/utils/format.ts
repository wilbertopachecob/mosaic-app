/**
 * Formats a duration in seconds for display in the result metadata.
 * Values under one second are shown in milliseconds; otherwise in seconds.
 */
export function formatDuration(seconds: number): string {
  if (seconds < 1) {
    return `${Math.round(seconds * 1000)}ms`;
  }
  return `${seconds.toFixed(2)}s`;
}

/**
 * Truncates a filename for display, appending an ellipsis when it exceeds `maxLength`.
 */
export function truncateFileName(
  name: string,
  maxLength = 28
): string {
  if (name.length <= maxLength) {
    return name;
  }
  return name.substring(0, maxLength) + "...";
}

/**
 * Formats a byte count using the given unit labels (bytes, KB, MB, GB).
 */
export function formatFileSize(
  bytes: number,
  units: [string, string, string, string]
): string {
  if (bytes === 0) {
    return `0 ${units[0]}`;
  }
  const bytesPerUnit = 1024;
  const unitIndex = Math.floor(Math.log(bytes) / Math.log(bytesPerUnit));
  return `${parseFloat((bytes / Math.pow(bytesPerUnit, unitIndex)).toFixed(2))} ${units[unitIndex]}`;
}
