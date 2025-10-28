import { fetchHealthData } from "@/lib/data/health";
import { StatusCard } from "@/components/StatusCard";

export default async function StatusPage() {
  const { data: healthData, error, timestamp } = await fetchHealthData();

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            🚀 Page Insight Tool Status
          </h1>
          <p className="text-gray-600">Server health and system information</p>
        </div>

        {/* Status Card Component */}
        <StatusCard healthData={healthData!} renderedAt={timestamp} />

        {/* Error Notice (if any) */}
        {error && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-8">
            <div className="flex">
              <div className="flex-shrink-0">
                <span className="text-yellow-400">⚠️</span>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-yellow-800">
                  Data Fetch Warning
                </h3>
                <div className="mt-2 text-sm text-yellow-700">
                  <p>Using fallback data due to: {error}</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* SSR Demonstration Note */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <span className="text-blue-400">💡</span>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">
                Production-Grade SSR Architecture
              </h3>
              <div className="mt-2 text-sm text-blue-700">
                <ul className="list-disc list-inside space-y-1">
                  <li>Data fetching separated from UI components</li>
                  <li>Centralized API client with singleton pattern</li>
                  <li>Proper error handling and fallback data</li>
                  <li>Reusable components for maintainability</li>
                  <li>Type-safe data validation</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
