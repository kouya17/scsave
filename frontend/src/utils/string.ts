export function clip(text: string | undefined, max: number): string {
  if (!text) {
    return 'null'
  }
  if (text.length > max) {
    return text.substring(0, max) + '...'
  }
  return text
}