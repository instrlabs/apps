'use client';

/**
 * Triggers a browser download from an ArrayBuffer by creating a temporary Blob URL.
 *
 * @param buffer The ArrayBuffer of the file content
 * @param fileName The name suggested for the downloaded file
 * @param mimeType Optional MIME type of the file. Defaults to application/octet-stream
 */
export function downloadFromArrayBuffer(buffer: ArrayBuffer, fileName: string, mimeType: string = 'application/octet-stream') {
  const blob = new Blob([buffer], { type: mimeType });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = fileName;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}
