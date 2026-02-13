'use client';

import { ProfileSection } from '@/components/settings/ProfileSection';
import { CategoryManager } from '@/components/settings/CategoryManager';

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900">설정</h2>
      <ProfileSection />
      <CategoryManager />
    </div>
  );
}
