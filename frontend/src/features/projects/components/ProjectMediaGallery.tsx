import type { MediaAsset } from '@0xhub/api';

export interface ProjectMediaGalleryProps {
  media: MediaAsset[];
}

const isImage = (asset: MediaAsset): boolean =>
  Boolean(asset.mimeType?.startsWith('image/') || asset.storagePath.match(/\.(png|jpe?g|gif|webp|avif)$/i));

export const ProjectMediaGallery = ({ media }: ProjectMediaGalleryProps): JSX.Element | null => {
  if (!media.length) {
    return null;
  }

  const imageAssets = media.filter(isImage);
  const otherAssets = media.filter((asset) => !isImage(asset));

  return (
    <div className="space-y-6">
      {imageAssets.length ? (
        <section className="space-y-4">
          <header className="flex items-center justify-between">
            <h3 className="text-xl font-semibold text-text">Gallery</h3>
            <span className="text-xs uppercase tracking-[0.28em] text-text-muted">
              {imageAssets.length} {imageAssets.length === 1 ? 'asset' : 'assets'}
            </span>
          </header>
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            {imageAssets.map((asset) => (
              <figure
                key={asset.id}
                className="group overflow-hidden rounded-3xl border border-border/70 bg-surface-raised shadow-surface-sm transition hover:shadow-surface-md"
              >
                <img
                  src={asset.storagePath}
                  alt={asset.description ?? asset.originalFilename}
                  className="h-56 w-full object-cover transition duration-500 group-hover:scale-105"
                  loading="lazy"
                  decoding="async"
                />
                {(asset.description || asset.originalFilename) && (
                  <figcaption className="px-4 py-3 text-xs text-text-muted">
                    {asset.description ?? asset.originalFilename}
                  </figcaption>
                )}
              </figure>
            ))}
          </div>
        </section>
      ) : null}

      {otherAssets.length ? (
        <section className="space-y-4">
          <header className="flex items-center justify-between">
            <h3 className="text-xl font-semibold text-text">Downloads</h3>
            <span className="text-xs uppercase tracking-[0.28em] text-text-muted">
              {otherAssets.length} file{otherAssets.length === 1 ? '' : 's'}
            </span>
          </header>
          <ul className="grid gap-3 sm:grid-cols-2">
            {otherAssets.map((asset) => (
              <li key={asset.id}>
                <a
                  href={asset.storagePath}
                  target="_blank"
                  rel="noreferrer"
                  className="flex items-center justify-between gap-3 rounded-2xl border border-border/70 bg-surface px-4 py-3 text-sm text-text transition hover:border-brand/70 hover:text-brand"
                >
                  <span className="truncate">{asset.originalFilename}</span>
                  <span className="text-xs uppercase tracking-[0.28em] text-text-muted">
                    {asset.mimeType ?? 'Download'}
                  </span>
                </a>
              </li>
            ))}
          </ul>
        </section>
      ) : null}
    </div>
  );
};


