"use client";

import { Suspense } from "react";
import { AnalysisForm } from "./AnalysisForm";

interface AnalysisFormWrapperProps {
  initialUrl?: string;
}

function AnalysisFormSuspended({ initialUrl }: AnalysisFormWrapperProps) {
  return <AnalysisForm initialUrl={initialUrl} />;
}

export function AnalysisFormWrapper({ initialUrl }: AnalysisFormWrapperProps) {
  return (
    <Suspense fallback={
      <div className="text-sm text-muted-foreground">
        Loading form...
      </div>
    }>
      <AnalysisFormSuspended initialUrl={initialUrl} />
    </Suspense>
  );
}

