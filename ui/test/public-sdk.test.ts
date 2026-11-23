import { describe, expect, it } from 'vitest';
import manifest from '../package.json';

const packageManifest: {
  scripts: Record<string, string | undefined>;
  devDependencies: Record<string, string | undefined>;
  dependencies: Record<string, string | undefined>;
} = manifest;

describe('public SDK package and build contract', () => {
  it('uses the registry package and canonical build without consumer workarounds', () => {
    expect(packageManifest.scripts.postinstall).toBeUndefined();
    expect(packageManifest.scripts.build).toBe('inari-ui-ext build ./src/index.tsx --name inari-ext-argocd');
    expect(packageManifest.devDependencies.webpack).toBeUndefined();
    expect(packageManifest.devDependencies['ts-loader']).toBeUndefined();
    expect(packageManifest.dependencies['@7k-inari/ui-plugin-sdk']).toBe('0.1.9');
  });
});
