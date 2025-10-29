"use client";

import { useState, FormEvent, useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Search, Loader2 } from "lucide-react";

interface AnalysisFormProps {
  initialUrl?: string;
}

export function AnalysisForm({ initialUrl = "" }: AnalysisFormProps) {
  const [url, setUrl] = useState(initialUrl);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const router = useRouter();
  const searchParams = useSearchParams();
  const submitTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Update URL when searchParams change (e.g., browser back/forward)
  useEffect(() => {
    const urlParam = searchParams.get("url");
    if (urlParam) {
      setUrl(decodeURIComponent(urlParam));
    } else {
      // Reset form when URL param is removed (e.g., logo clicked)
      setUrl("");
    }
    // Reset submitting state when URL changes
    setIsSubmitting(false);
    // Clear any pending timeout
    if (submitTimeoutRef.current) {
      clearTimeout(submitTimeoutRef.current);
      submitTimeoutRef.current = null;
    }
  }, [searchParams]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (submitTimeoutRef.current) {
        clearTimeout(submitTimeoutRef.current);
      }
    };
  }, []);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!url.trim() || isSubmitting) return;

    const trimmedUrl = url.trim();
    const currentUrlParam = searchParams.get("url");

    setIsSubmitting(true);

    // If submitting the same URL, use refresh() to force server-side re-fetch
    // Otherwise, navigate to the new URL
    if (currentUrlParam === trimmedUrl) {
      router.refresh();
      // Reset after a short delay since refresh doesn't trigger searchParams change
      submitTimeoutRef.current = setTimeout(() => {
        setIsSubmitting(false);
        submitTimeoutRef.current = null;
      }, 2000);
    } else {
      router.push(`/?url=${encodeURIComponent(trimmedUrl)}`);
      // Safety timeout: reset submitting state after 5 seconds if navigation hasn't completed
      submitTimeoutRef.current = setTimeout(() => {
        setIsSubmitting(false);
        submitTimeoutRef.current = null;
      }, 5000);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="url">Website URL</Label>
        <div className="flex gap-2">
          <Input
            id="url"
            type="url"
            placeholder="https://example.com"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            disabled={isSubmitting}
            className="flex-1"
          />
          <Button type="submit" disabled={isSubmitting || !url.trim()}>
            {isSubmitting ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Analyzing...
              </>
            ) : (
              <>
                <Search className="h-4 w-4" />
                Analyze
              </>
            )}
          </Button>
        </div>
      </div>
    </form>
  );
}
