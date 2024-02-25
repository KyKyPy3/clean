/**
 * Checks if value is empty. Accepts strings, numbers, booleans, objects and arrays.
 */
export function isEmpty(value: unknown): boolean {
  if (typeof value === 'number' || typeof value === 'boolean') {
    return false;
  }
  if (typeof value === 'undefined' || value === null) {
    return true;
  }
  if (value instanceof Date) {
    return false;
  }
  if (value instanceof Object && !Object.keys(value).length) {
    return true;
  }
  if (Array.isArray(value)) {
    if (value.length === 0) {
      return true;
    }
    if (value.every((item) => isEmpty(item))) {
      return true;
    }
  }
  return value === '';
}