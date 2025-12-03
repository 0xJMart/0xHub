import type {
  MediaAsset,
  Project,
  ProjectListResponse,
  ProjectSummary,
  ProjectTag,
  Tag,
} from '@0xhub/api';

const makeTag = (tag: Partial<ProjectTag> & Pick<ProjectTag, 'id' | 'name' | 'displayName'>): Tag => ({
  color: '#38bdf8',
  ...tag,
});

export const mockTags: Tag[] = [
  makeTag({
    id: 'tag-automation',
    name: 'automation',
    displayName: 'Automation',
    color: '#34d399',
  }),
  makeTag({
    id: 'tag-k3s',
    name: 'k3s',
    displayName: 'K3s',
    color: '#60a5fa',
  }),
  makeTag({
    id: 'tag-media',
    name: 'media',
    displayName: 'Media',
    color: '#fbbf24',
  }),
];

const makeMedia = (media: Partial<MediaAsset> & Pick<MediaAsset, 'id' | 'storagePath' | 'originalFilename'>): MediaAsset => ({
  status: 'available',
  ...media,
});

export const mockProjects: Project[] = [
  {
    id: 'project-pihole',
    title: 'Pi-hole Gateway',
    slug: 'pi-hole-gateway',
    summary: 'Network-wide DNS filtering and telemetry reduction with Pi-hole and Unbound.',
    status: 'active',
    category: { id: 'cat-network', name: 'Networking', slug: 'networking', description: 'Network services' },
    tags: [mockTags[0], mockTags[1]],
    heroMedia: makeMedia({
      id: 'media-pihole-hero',
      storagePath: 'https://images.0xhub.dev/pi-hole-dashboard.png',
      originalFilename: 'pi-hole-dashboard.png',
      description: 'Pi-hole query dashboard with live metrics.',
    }),
    description: `## Overview

This Pi-hole deployment fronts the homelab network with DNS filtering and analytics. It runs within a lightweight Kubernetes cluster using a dedicated service account and persistent storage.

## Highlights

- Kubernetes deployment with health probes and rolling updates.
- Unbound upstream resolver for encrypted DNS.
- Automatic gravity updates using a CronJob.

## Hardware

- Intel NUC i5 with 16GB RAM.
- Dual NIC for VLAN segregation.
`,
    links: [
      { id: 'link-pihole-docs', linkType: 'docs', url: 'https://pi-hole.net/', label: 'Pi-hole Docs' },
      { id: 'link-pihole-chart', linkType: 'repo', url: 'https://github.com/MoJo2600/pihole-kubernetes', label: 'Helm Chart' },
    ],
    media: [
      makeMedia({
        id: 'media-pihole-1',
        storagePath: 'https://images.0xhub.dev/pi-hole-overview.png',
        originalFilename: 'pi-hole-overview.png',
        description: 'Service topology diagram.',
      }),
      makeMedia({
        id: 'media-pihole-logs',
        storagePath: 'https://files.0xhub.dev/pi-hole-logs.txt',
        originalFilename: 'pi-hole-logs.txt',
        description: 'Sample query logs for troubleshooting.',
        mimeType: 'text/plain',
      }),
    ],
  },
  {
    id: 'project-argocd',
    title: 'GitOps Control Plane',
    slug: 'gitops-control-plane',
    summary: 'Argo CD powered GitOps workflows managing self-hosted services.',
    status: 'active',
    category: { id: 'cat-platform', name: 'Platform', slug: 'platform', description: 'Platform tooling' },
    tags: [mockTags[1]],
    heroMedia: makeMedia({
      id: 'media-argocd-hero',
      storagePath: 'https://images.0xhub.dev/argocd-dashboard.png',
      originalFilename: 'argocd-dashboard.png',
      description: 'Argo CD applications overview.',
    }),
    description: `## Workflow

1. Applications defined as Helm charts committed to the gitops repository.
2. Argo CD watches the repository on the main branch.
3. Sync waves ensure databases come online before workloads.

## Secrets

- External Secrets Operator pulls values from Vault.
- Sensitive values never live inside the repository.
`,
    links: [
      { id: 'link-argocd', linkType: 'docs', url: 'https://argo-cd.readthedocs.io/', label: 'Argo CD Docs' },
    ],
    media: [
      makeMedia({
        id: 'media-argocd-rollout',
        storagePath: 'https://images.0xhub.dev/argocd-rollout.png',
        originalFilename: 'argocd-rollout.png',
        description: 'Application deployment pipeline overview.',
      }),
    ],
  },
  {
    id: 'project-plex',
    title: 'Plex Media Stack',
    slug: 'plex-media-stack',
    summary: 'Transcoding-optimized Plex stack with automated library management.',
    status: 'draft',
    category: { id: 'cat-media', name: 'Media', slug: 'media' },
    tags: [mockTags[2]],
    heroMedia: makeMedia({
      id: 'media-plex-hero',
      storagePath: 'https://images.0xhub.dev/plex-library.png',
      originalFilename: 'plex-library.png',
      description: 'Custom Plex dashboard with On Deck view.',
    }),
    description: `Planned rebuild of the Plex stack to support hardware transcoding and improved metadata management. Drafting the architecture before rollout.`,
    links: [
      { id: 'link-plex', linkType: 'docs', url: 'https://www.plex.tv/media-server-downloads/', label: 'Plex Media Server' },
    ],
    media: [],
  },
];

export const toProjectSummary = (project: Project): ProjectSummary => ({
  id: project.id,
  title: project.title,
  slug: project.slug,
  summary: project.summary,
  status: project.status,
  category: project.category,
  tags: project.tags,
  heroMedia: project.heroMedia,
});

const projectSummaries: ProjectSummary[] = mockProjects.map(toProjectSummary);

export const listProjects = (filters: Partial<Record<'tag' | 'category' | 'search', string>> = {}): ProjectListResponse => {
  let filtered = [...mockProjects];

  if (filters.tag) {
    filtered = filtered.filter((project) =>
      project.tags.some((tag) => tag.name === filters.tag),
    );
  }

  if (filters.category) {
    filtered = filtered.filter((project) => project.category?.slug === filters.category);
  }

  if (filters.search) {
    const term = filters.search.toLowerCase();
    filtered = filtered.filter(
      (project) =>
        project.title.toLowerCase().includes(term) ||
        (project.summary?.toLowerCase().includes(term) ?? false) ||
        project.tags.some((tag) => tag.displayName.toLowerCase().includes(term)),
    );
  }

  return {
    items: filtered.map(toProjectSummary),
    total: filtered.length,
  };
};

export const findProjectBySlug = (slug: string): Project | undefined =>
  mockProjects.find((project) => project.slug === slug);

export const allProjectSummaries = (): ProjectSummary[] => projectSummaries;


