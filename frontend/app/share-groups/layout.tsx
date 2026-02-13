import type { ReactNode } from 'react';

import { AppLayout } from '@/components/layout/AppLayout';

export default function ShareGroupsLayout({ children }: { children: ReactNode }) {
  return <AppLayout>{children}</AppLayout>;
}
