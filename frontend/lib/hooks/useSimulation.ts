import { useMutation, useQueryClient } from '@tanstack/react-query';

import { post } from '@/lib/api';
import type {
  AddSimulationRequest,
  ApplySimulationRequest,
  CancelSimulationRequest,
  CombinedSimulationRequest,
  SimulationResult,
} from '@/types';

/**
 * 해지 시뮬레이션 훅
 * 선택한 구독들을 해지했을 때의 금액 변화를 시뮬레이션
 */
export function useCancelSimulation() {
  return useMutation({
    mutationFn: (data: CancelSimulationRequest) =>
      post<SimulationResult>('/simulation/cancel', data),
  });
}

/**
 * 추가 시뮬레이션 훅
 * 새 구독을 추가했을 때의 금액 변화를 시뮬레이션
 */
export function useAddSimulation() {
  return useMutation({
    mutationFn: (data: AddSimulationRequest) => post<SimulationResult>('/simulation/add', data),
  });
}

/**
 * 통합 시뮬레이션 훅
 * 해지 + 추가를 동시에 시뮬레이션
 */
export function useCombinedSimulation() {
  return useMutation({
    mutationFn: (data: CombinedSimulationRequest) =>
      post<SimulationResult>('/simulation/combined', data),
  });
}

/**
 * 시뮬레이션 적용 훅
 * 시뮬레이션한 변경사항을 실제로 적용
 */
export function useApplySimulation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ApplySimulationRequest) => post<void>('/simulation/apply', data),
    onSuccess: () => {
      // 구독 목록과 대시보드 데이터 새로고침
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

/**
 * 시뮬레이션 되돌리기 훅
 * 최근 적용한 시뮬레이션을 되돌림
 */
export function useUndoSimulation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => post<void>('/simulation/undo', undefined),
    onSuccess: () => {
      // 구독 목록과 대시보드 데이터 새로고침
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}
