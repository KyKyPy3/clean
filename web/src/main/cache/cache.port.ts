export interface CacheContent {
  content: string;
}

export interface CacheReader {
  getItem: (key: string) => CacheContent;
}

export interface CacheWriter {
  save<T = unknown>(key: string, content: T): void
}