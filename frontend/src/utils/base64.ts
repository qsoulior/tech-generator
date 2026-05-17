export function fromBase64(data: string): string {
  const bin = atob(data)
  const bytes = Uint8Array.from(bin, (m) => m.codePointAt(0) ?? 0)
  return new TextDecoder().decode(bytes)
}

export function toBase64(data: string): string {
  const encoded = new TextEncoder().encode(data)
  const bin = Array.from(encoded, (byte) => String.fromCodePoint(byte)).join("")
  return btoa(bin)
}