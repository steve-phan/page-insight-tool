import { fetchHealthData } from "@/lib/data/health";
import { StatusCard } from "@/components/StatusCard";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AlertCircle, Lightbulb } from "lucide-react";

export default async function StatusPage() {
  const { data: healthData, error, timestamp } = await fetchHealthData();

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
      <main className="container mx-auto px-4 py-12 max-w-4xl">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold mb-4">
            🚀 Page Insight Tool Status
          </h1>
          <p className="text-lg text-muted-foreground">Server health and system information</p>
        </div>

        {/* Status Card Component */}
        {healthData && <StatusCard healthData={healthData} renderedAt={timestamp} />}

        {/* Error Notice (if any) */}
        {error && (
          <Card className="mt-6 border-yellow-200 bg-yellow-50 dark:bg-yellow-950 dark:border-yellow-800">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-yellow-800 dark:text-yellow-200">
                <AlertCircle className="h-5 w-5" />
                Data Fetch Warning
              </CardTitle>
              <CardDescription className="text-yellow-700 dark:text-yellow-300">
                Using fallback data due to an error
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-yellow-800 dark:text-yellow-200">{error}</p>
            </CardContent>
          </Card>
        )}

        {/* SSR Demonstration Note */}
        <Card className="mt-6 border-blue-200 bg-blue-50 dark:bg-blue-950 dark:border-blue-800">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-blue-800 dark:text-blue-200">
              <Lightbulb className="h-5 w-5" />
              Production-Grade SSR Architecture
            </CardTitle>
            <CardDescription className="text-blue-700 dark:text-blue-300">
              Demonstrating best practices for server-side rendering
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ul className="list-disc list-inside space-y-2 text-sm text-blue-800 dark:text-blue-200">
              <li>Data fetching separated from UI components</li>
              <li>Centralized API client with singleton pattern</li>
              <li>Proper error handling and fallback data</li>
              <li>Reusable components for maintainability</li>
              <li>Type-safe data validation</li>
            </ul>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
