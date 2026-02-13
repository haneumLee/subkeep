import type { ReactNode } from 'react';

import { AppLayout } from '@/components/layout/AppLayout';

export default function SettingsLayout({ children }: { children: ReactNode }) {
  return <AppLayout>{children}</AppLayout>;
}
