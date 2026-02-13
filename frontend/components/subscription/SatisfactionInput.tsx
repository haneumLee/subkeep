'use client';

import { useState } from 'react';

interface SatisfactionInputProps {
  value: number;
  onChange: (value: number) => void;
  disabled?: boolean;
}

export default function SatisfactionInput({ value, onChange, disabled }: SatisfactionInputProps) {
  const [hoverValue, setHoverValue] = useState<number | null>(null);

  const displayValue = hoverValue ?? value;

  return (
    <div className="flex items-center gap-2">
      <div className="flex gap-1">
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            type="button"
            disabled={disabled}
            onClick={() => onChange(star)}
            onMouseEnter={() => setHoverValue(star)}
            onMouseLeave={() => setHoverValue(null)}
            className="text-2xl transition-all disabled:cursor-not-allowed disabled:opacity-50"
          >
            <span className={displayValue >= star ? 'text-yellow-400' : 'text-gray-300'}>
              {displayValue >= star ? '★' : '☆'}
            </span>
          </button>
        ))}
      </div>
      <span className="text-sm text-gray-600">
        {displayValue > 0 ? `${displayValue}점` : '미평가'}
      </span>
    </div>
  );
}
