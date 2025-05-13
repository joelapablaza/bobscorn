'use client';

import { useState } from 'react';
import { purchaseCornFromBob } from '@/app/services/corn-service';

type PurchaseResponse = {
  success?: boolean;
  message?: string;
  error?: string;
  status?: number;
};

export function useCornPurchase() {
  const [isLoading, setIsLoading] = useState(false);
  const [response, setResponse] = useState<PurchaseResponse | null>(null);
  const [purchaseCount, setPurchaseCount] = useState(0);

  const purchaseCorn = async () => {
    try {
      setIsLoading(true);

      const result = await purchaseCornFromBob();
      setResponse(result);

      if (result.success && result.message) {
        setPurchaseCount((prev) => prev + 1);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return {
    purchaseCorn,
    isLoading,
    response,
    purchaseCount,
  };
}
