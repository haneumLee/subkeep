'use client';

import { useState } from 'react';

import { Modal } from '@/components/ui/Modal';
import { useFolders, useCreateFolder, useUpdateFolder, useDeleteFolder } from '@/lib/hooks/useFolders';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';

interface FolderManagerPopupProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function FolderManagerPopup({ isOpen, onClose }: FolderManagerPopupProps) {
  const { data: folders, isLoading } = useFolders();
  const createFolder = useCreateFolder();
  const updateFolder = useUpdateFolder();
  const deleteFolder = useDeleteFolder();

  const [newFolderName, setNewFolderName] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState('');
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const handleCreate = () => {
    if (!newFolderName.trim()) return;
    createFolder.mutate(
      { name: newFolderName.trim() },
      {
        onSuccess: () => setNewFolderName(''),
      }
    );
  };

  const handleUpdate = () => {
    if (!editingId || !editingName.trim()) return;
    updateFolder.mutate(
      { id: editingId, data: { name: editingName.trim() } },
      {
        onSuccess: () => {
          setEditingId(null);
          setEditingName('');
        },
      }
    );
  };

  const handleDelete = () => {
    if (!deletingId) return;
    deleteFolder.mutate(deletingId, {
      onSuccess: () => setDeletingId(null),
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
      <div className="w-full max-w-md rounded-lg bg-white shadow-xl">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 px-6 py-4">
          <h2 className="text-lg font-semibold text-gray-900">폴더 관리</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
            aria-label="닫기"
          >
            <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Add folder */}
        <div className="border-b border-gray-200 px-6 py-4">
          <div className="flex gap-2">
            <input
              type="text"
              value={newFolderName}
              onChange={(e) => setNewFolderName(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
              placeholder="새 폴더 이름"
              className="flex-1 rounded-lg border-2 border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none"
            />
            <button
              onClick={handleCreate}
              disabled={!newFolderName.trim() || createFolder.isPending}
              className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              추가
            </button>
          </div>
        </div>

        {/* Folder list */}
        <div className="max-h-80 overflow-y-auto px-6 py-4">
          {isLoading ? (
            <div className="flex justify-center py-8">
              <LoadingSpinner size="md" />
            </div>
          ) : folders && folders.length > 0 ? (
            <ul className="space-y-2">
              {folders.map((folder) => (
                <li
                  key={folder.id}
                  className="flex items-center justify-between rounded-lg border border-gray-200 p-3"
                >
                  {editingId === folder.id ? (
                    <div className="flex flex-1 gap-2">
                      <input
                        type="text"
                        value={editingName}
                        onChange={(e) => setEditingName(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && handleUpdate()}
                        className="flex-1 rounded border border-gray-300 px-2 py-1 text-sm focus:border-blue-500 focus:outline-none"
                        // eslint-disable-next-line jsx-a11y/no-autofocus
                        autoFocus
                      />
                      <button
                        onClick={handleUpdate}
                        className="text-sm font-medium text-blue-600 hover:text-blue-700"
                      >
                        저장
                      </button>
                      <button
                        onClick={() => setEditingId(null)}
                        className="text-sm font-medium text-gray-500 hover:text-gray-700"
                      >
                        취소
                      </button>
                    </div>
                  ) : (
                    <>
                      <div className="flex items-center gap-2">
                        <svg className="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                        </svg>
                        <span className="text-sm font-medium text-gray-900">{folder.name}</span>
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => {
                            setEditingId(folder.id);
                            setEditingName(folder.name);
                          }}
                          className="text-sm text-gray-500 hover:text-blue-600"
                        >
                          수정
                        </button>
                        <button
                          onClick={() => setDeletingId(folder.id)}
                          className="text-sm text-gray-500 hover:text-red-600"
                        >
                          삭제
                        </button>
                      </div>
                    </>
                  )}
                </li>
              ))}
            </ul>
          ) : (
            <p className="py-8 text-center text-sm text-gray-500">
              등록된 폴더가 없습니다
            </p>
          )}
        </div>

        {/* Footer */}
        <div className="border-t border-gray-200 px-6 py-4">
          <button
            onClick={onClose}
            className="w-full rounded-lg border-2 border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            닫기
          </button>
        </div>

        {/* Delete Confirm */}
        <Modal
          isOpen={deletingId !== null}
          onClose={() => setDeletingId(null)}
          title="폴더 삭제"
          confirmText="삭제"
          cancelText="취소"
          onConfirm={handleDelete}
          confirmVariant="danger"
          isConfirmLoading={deleteFolder.isPending}
        >
          <p>이 폴더를 삭제하시겠습니까? 폴더 내 구독은 미분류로 변경됩니다.</p>
        </Modal>
      </div>
    </div>
  );
}
