export const useLogger = () => {
  return {
    debug: (...args: any[]) => console.debug(...args),
    warn: (...args: any[]) => console.warn(...args),
    info: (...args: any[]) => console.info(...args),
    error: (...args: any[]) => console.error(...args)
  };
}
