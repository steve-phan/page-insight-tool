import { Label } from "@/components/ui/label";
import { Search } from "lucide-react";

interface AnalysisFormProps {
  initialUrl?: string;
}

export function AnalysisForm({ initialUrl = "" }: AnalysisFormProps) {
  return (
    <form method="GET" action="/" className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="url">Website URL</Label>
        <div className="flex gap-2">
          <input
            id="url"
            name="url"
            type="url"
            placeholder="https://example.com"
            defaultValue={initialUrl}
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 flex-1"
            required
          />
          <button
            type="submit"
            className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
          >
            <Search className="h-4 w-4" />
            Analyze
          </button>
        </div>
      </div>
    </form>
  );
}
