'use client';

import { useFolders } from '@/lib/hooks/useFolders';

interface FolderPanelProps {
  selectedFolderId: string | null;
  onSelectFolder: (folderId: string | null) => void;
}

export default function FolderPanel({ selectedFolderId, onSelectFolder }: FolderPanelProps) {
  const { data: folders, isLoading } = useFolders();

  return (
    <div className="w-48 shrink-0 rounded-lg border-2 border-gray-200 bg-white p-3">
      <h3 className="mb-3 text-sm font-semibold text-gray-700">폴더</h3>

      {isLoading ? (
        <div className="animate-pulse space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-8 rounded bg-gray-200" />
          ))}
        </div>
      ) : (
        <ul className="space-y-1">
          {/* All */}
          <li>
            <button
              onClick={() => onSelectFolder(null)}
              className={`w-full rounded-md px-3 py-2 text-left text-sm font-medium transition-colors ${
                selectedFolderId === null
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              <div className="flex items-center gap-2">
                <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                </svg>
                전체
              </div>
            </button>
          </li>

          {/* Folders */}
          {folders?.map((folder) => (
            <li key={folder.id}>
              <button
                onClick={() => onSelectFolder(folder.id)}
                className={`w-full rounded-md px-3 py-2 text-left text-sm font-medium transition-colors ${
                  selectedFolderId === folder.id
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-50'
                }`}
              >
                <div className="flex items-center gap-2">
                  <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                  </svg>
                  {folder.name}
                </div>
              </button>
            </li>
          ))}

          {(!folders || folders.length === 0) && (
            <li className="px-3 py-2 text-xs text-gray-400">
              폴더가 없습니다
            </li>
          )}
        </ul>
      )}
    </div>
  );
}
