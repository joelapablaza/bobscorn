'use client';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Loader2, ShoppingCart } from 'lucide-react';
import { useCornPurchase } from '@/app/hooks/use-corn-purchase';
import { useState, useEffect } from 'react';

export default function CornPurchase() {
  const { purchaseCorn, isLoading, response, purchaseCount } =
    useCornPurchase();
  const [showResponse, setShowResponse] = useState(false);

  useEffect(() => {
    if (response) {
      setShowResponse(true);

      // Only auto-hide success messages, not errors
      if (response.success) {
        const timer = setTimeout(() => {
          setShowResponse(false);
        }, 3000);
        return () => clearTimeout(timer);
      }
    }
  }, [response]);

  const handlePurchase = async () => {
    await purchaseCorn();
  };

  const getAlertStyles = () => {
    if (!response) return '';

    if (response.status === 429) {
      return 'bg-amber-100 border-amber-200';
    }
    // Success gets green
    else if (response.success) {
      return 'bg-green-100 border-green-200';
    }
    // Other errors get red
    else {
      return 'bg-red-100 border-red-200';
    }
  };

  // Determine the text color based on the response
  const getTextStyles = () => {
    if (!response) return '';

    if (response.status === 429) {
      return 'text-amber-800';
    } else if (response.success) {
      return 'text-green-800';
    } else {
      return 'text-red-800';
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-4">
        <div className="text-amber-800 font-medium">Fresh corn available!</div>
        {purchaseCount > 0 && (
          <div className="flex items-center text-green-700">
            <ShoppingCart className="h-4 w-4 mr-1" />
            <span>{purchaseCount} purchased</span>
          </div>
        )}
      </div>

      <Button
        onClick={handlePurchase}
        disabled={isLoading}
        className="w-full bg-yellow-500 hover:bg-yellow-600 text-white"
      >
        {isLoading ? (
          <>
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            Purchasing...
          </>
        ) : (
          "Buy Bob's Corn 🌽"
        )}
      </Button>

      {response && showResponse && (
        <Alert className={`transition-opacity ${getAlertStyles()}`}>
          <AlertDescription className={getTextStyles()}>
            {response.error || response.message}
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
