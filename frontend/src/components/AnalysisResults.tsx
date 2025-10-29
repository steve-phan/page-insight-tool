import { AnalysisResponse } from '@/types/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { FileText, Link2, Lock, Clock } from 'lucide-react';

interface AnalysisResultsProps {
  data: AnalysisResponse;
}

export function AnalysisResults({ data }: AnalysisResultsProps) {
  return (
    <div className="space-y-4 mt-6">
      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="text-lg flex items-center gap-2">
            <FileText className="h-5 w-5" />
            Page Information
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Label className="text-xs text-muted-foreground">HTML Version</Label>
              <p className="text-base font-semibold mt-1">{data.html_version}</p>
            </div>
            <div>
              <Label className="text-xs text-muted-foreground">Page Title</Label>
              <p className="text-base font-semibold mt-1 line-clamp-2">{data.page_title || 'N/A'}</p>
            </div>
            <div>
              <Label className="text-xs text-muted-foreground">Analysis Time</Label>
              <p className="text-base font-semibold flex items-center gap-2 mt-1">
                <Clock className="h-4 w-4" />
                {data.analysis_time_ms} ms
              </p>
            </div>
            <div>
              <Label className="text-xs text-muted-foreground">Login Form</Label>
              <p className="text-base font-semibold flex items-center gap-2 mt-1">
                {data.has_login_form ? (
                  <>
                    <Lock className="h-4 w-4 text-green-600 dark:text-green-400" />
                    <span className="text-green-600 dark:text-green-400">Detected</span>
                  </>
                ) : (
                  <span className="text-muted-foreground">Not detected</span>
                )}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="text-lg">Headings Structure</CardTitle>
          <CardDescription className="text-sm">Heading distribution on the page</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 md:grid-cols-6 gap-3">
            {Object.entries(data.headings).map(([level, count]) => (
              <div key={level} className="text-center">
                <div className="text-xl font-bold">{count}</div>
                <div className="text-xs text-muted-foreground uppercase">{level}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="text-lg flex items-center gap-2">
            <Link2 className="h-5 w-5" />
            Links Analysis
          </CardTitle>
          <CardDescription className="text-sm">Link breakdown by type</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div className="p-3 rounded-lg bg-blue-50 dark:bg-blue-950">
              <div className="text-xl font-bold text-blue-600 dark:text-blue-400">
                {data.links.internal}
              </div>
              <div className="text-xs text-muted-foreground">Internal Links</div>
            </div>
            <div className="p-3 rounded-lg bg-green-50 dark:bg-green-950">
              <div className="text-xl font-bold text-green-600 dark:text-green-400">
                {data.links.external}
              </div>
              <div className="text-xs text-muted-foreground">External Links</div>
            </div>
            <div className="p-3 rounded-lg bg-orange-50 dark:bg-orange-950">
              <div className="text-xl font-bold text-orange-600 dark:text-orange-400">
                {data.links.inaccessible}
              </div>
              <div className="text-xs text-muted-foreground">Inaccessible</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

