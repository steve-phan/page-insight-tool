import { AnalysisForm } from "@/components/AnalysisForm";
import { AnalysisResults } from "@/components/AnalysisResults";
import { Home24Logo } from "@/components/Home24Logo";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { fetchAnalysisData } from "@/lib/data/analysis";

interface HomeProps {
  searchParams: Promise<{ url?: string }> | { url?: string };
}

export default async function Home({ searchParams }: HomeProps) {
  // Handle both Promise and non-Promise searchParams for Next.js compatibility
  const params =
    searchParams instanceof Promise ? await searchParams : searchParams;
  // Next.js automatically decodes searchParams, so url is already decoded
  const url = params.url;

  // Fetch analysis data directly if URL is provided (SSR)
  let analysisData = null;
  let analysisError = null;
  if (url) {
    const result = await fetchAnalysisData(url);
    analysisData = result.data;
    analysisError = result.error;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
      <main className="container mx-auto px-4 py-6 max-w-6xl">
        <div className="text-center mb-6">
          <h1 className="text-3xl font-bold mb-2 flex items-center justify-center gap-3">
            <Home24Logo />
            Page Insight Tool
          </h1>
          <p className="text-sm text-muted-foreground">
            Analyze any webpage to extract HTML structure, headings, links, and
            more
          </p>
        </div>

        <Card className="mb-6">
          <CardHeader className="pb-4">
            <CardTitle className="text-lg">Analyze a URL</CardTitle>
            <CardDescription className="text-sm">
              Enter a URL to analyze its HTML structure and content
            </CardDescription>
          </CardHeader>
          <CardContent>
            <AnalysisForm initialUrl={url || ""} />
          </CardContent>
        </Card>

        {analysisError && (
          <Card className="mt-6 border-destructive bg-destructive/10">
            <CardHeader className="pb-4">
              <CardTitle className="text-destructive text-lg">Error</CardTitle>
              <CardDescription className="text-sm">
                Failed to analyze the URL
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-destructive">{analysisError}</p>
            </CardContent>
          </Card>
        )}

        {analysisData && <AnalysisResults data={analysisData} />}
      </main>
    </div>
  );
}
