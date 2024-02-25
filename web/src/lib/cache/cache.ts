import type { CacheContent, CacheReader, CacheWriter } from "@/src/main/cache/cache.port"

export class LocalStorageWriter implements CacheWriter {
  save(key: string, content: unknown): void {
    localStorage.setItem(key, this.adapt(content));
  }

  private adapt(content: unknown): string {
    return JSON.stringify(content);
  }
}

export class LocalStorageReader implements CacheReader {
  getItem(key: string): CacheContent {
    const cacheResult = localStorage.getItem(key);

    return this.adapt(cacheResult);
  }

  private adapt(cacheResult?: string | null): CacheContent {
    return {
      content: cacheResult || JSON.stringify([]),
    };
  }
}