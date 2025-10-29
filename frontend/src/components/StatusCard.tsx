import { HealthResponse } from "@/types/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  CheckCircle2,
  XCircle,
  AlertCircle,
  Server,
  Calendar,
  GitCommit,
  Clock,
  Activity,
} from "lucide-react";

interface StatusCardProps {
  healthData: HealthResponse;
  renderedAt: string;
}

function StatusBadge({ status }: { status: string }) {
  if (status === "healthy") {
    return (
      <Badge variant="default" className="bg-green-500 hover:bg-green-600">
        <CheckCircle2 className="h-3 w-3 mr-1" />
        Healthy
      </Badge>
    );
  }
  if (status === "unhealthy") {
    return (
      <Badge variant="destructive">
        <XCircle className="h-3 w-3 mr-1" />
        Unhealthy
      </Badge>
    );
  }
  return (
    <Badge variant="secondary">
      <AlertCircle className="h-3 w-3 mr-1" />
      Unknown
    </Badge>
  );
}

export function StatusCard({ healthData, renderedAt }: StatusCardProps) {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <Server className="h-5 w-5" />
            Server Status
          </CardTitle>
          <StatusBadge status={healthData.status} />
        </div>
        <CardDescription>Server health and system information</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-4">
            <div>
              <Label className="text-muted-foreground flex items-center gap-2 mb-2">
                <Activity className="h-4 w-4" />
                Version
              </Label>
              <p className="text-lg font-semibold font-mono">
                {healthData.version}
              </p>
            </div>

            <div>
              <Label className="text-muted-foreground flex items-center gap-2 mb-2">
                <Calendar className="h-4 w-4" />
                Build Date
              </Label>
              <p className="text-lg font-semibold">
                {healthData.build_date !== "unknown"
                  ? new Date(healthData.build_date).toLocaleString()
                  : "Unknown"}
              </p>
            </div>

            <div>
              <Label className="text-muted-foreground flex items-center gap-2 mb-2">
                <GitCommit className="h-4 w-4" />
                Git Commit
              </Label>
              <p className="text-sm font-mono text-muted-foreground break-all">
                {healthData.git_commit}
              </p>
            </div>
          </div>

          <div className="space-y-4">
            <div>
              <Label className="text-muted-foreground flex items-center gap-2 mb-2">
                <Clock className="h-4 w-4" />
                Server Time
              </Label>
              <p className="text-lg font-semibold">
                {new Date(healthData.timestamp).toLocaleString()}
              </p>
            </div>

            <div>
              <Label className="text-muted-foreground flex items-center gap-2 mb-2">
                <Activity className="h-4 w-4" />
                Uptime
              </Label>
              <p className="text-lg font-semibold">{healthData.uptime}</p>
            </div>

            <div>
              <Label className="text-muted-foreground flex items-center gap-2 mb-2">
                <Clock className="h-4 w-4" />
                Page Rendered At
              </Label>
              <p className="text-lg font-semibold">
                {new Date(renderedAt).toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
