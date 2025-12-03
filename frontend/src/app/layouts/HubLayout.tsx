import { AppShell, type AppShellNavLink } from '@0xhub/ui';
import { Suspense } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { LoadingState } from '@/components/LoadingState';

const navLinks: AppShellNavLink[] = [
  {
    label: 'Projects',
    href: '/',
  },
];

export const HubLayout = (): JSX.Element => {
  const location = useLocation();
  const navigate = useNavigate();

  const navigation = navLinks.map((item) => ({
    ...item,
    onClick: item.href ? () => navigate(item.href!) : undefined,
    isActive: item.href === '/' ? location.pathname === '/' : location.pathname.startsWith(item.href),
  }));

  return (
    <AppShell
      brand={{
        name: '0xHub',
        description: 'Homelab Knowledge Base',
      }}
      navigation={navigation}
      footer={<p>&copy; {new Date().getFullYear()} 0xHub. Crafted for self-hosted labs.</p>}
    >
      <Suspense fallback={<LoadingState message="Loading view…" />}>
        <Outlet />
      </Suspense>
    </AppShell>
  );
};


