import { HealthResponse } from '@/types/api';

interface StatusCardProps {
  healthData: HealthResponse;
  renderedAt: string;
}

export function StatusCard({ healthData, renderedAt }: StatusCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-8">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-gray-800">Server Status</h2>
        <span
          className={`px-3 py-1 rounded-full text-sm font-medium ${
            healthData.status === 'healthy'
              ? 'bg-green-100 text-green-800'
              : healthData.status === 'unhealthy'
              ? 'bg-red-100 text-red-800'
              : 'bg-yellow-100 text-yellow-800'
          }`}
        >
          {healthData.status === 'healthy'
            ? '✅ Healthy'
            : healthData.status === 'unhealthy'
            ? '❌ Unhealthy'
            : '❓ Unknown'}
        </span>
      </div>

      {/* Server Information Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Version</label>
            <p className="text-gray-900 font-mono">{healthData.version}</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Build Date</label>
            <p className="text-gray-900">
              {healthData.build_date !== 'unknown'
                ? new Date(healthData.build_date).toLocaleString()
                : 'Unknown'}
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Git Commit</label>
            <p className="text-gray-900 font-mono text-sm">{healthData.git_commit}</p>
          </div>
        </div>

        <div className="space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Server Time</label>
            <p className="text-gray-900">{new Date(healthData.timestamp).toLocaleString()}</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Uptime</label>
            <p className="text-gray-900">{healthData.uptime}</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-600 mb-1">Page Rendered At</label>
            <p className="text-gray-900">{new Date(renderedAt).toLocaleString()}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
