import { Code, Database, Globe } from "lucide-react";
import { Link } from "react-router";
import okeyFrieren from "~/assets/okey-frieren.png";
import { Button } from "~/components/ui/button";
import { CompanyCarousel } from "~/components/CompanyCarousel";
import staticData from "~/data/problems.json";

export default function Home() {
  const companies = staticData.companies.length;
  const problems = staticData.problems.length;

  return (
    <main className="flex h-[calc(100vh-var(--header-height,4rem))] flex-col items-center justify-between">
      <div className="flex w-full flex-col items-center gap-4 px-6 pt-16">
        <div className="font-geist text-primary relative flex flex-col items-start text-3xl font-semibold sm:text-4xl md:text-5xl dark:drop-shadow-[0_4px_8px_rgba(0,0,0,0.5)]">
          <span>Company-Wise</span>
          <span>
            <span>
              <span className="bg-linear-to-b from-[#ffa116] via-[#ff9345] to-[#ff3d00] bg-clip-text text-transparent dark:drop-shadow-[0_4px_25px_rgba(255,100,0,0.6)]">
                Leetcode{" "}
              </span>
              <img
                src={okeyFrieren}
                alt="Okey Frieren"
                className="absolute right-0 bottom-8 -z-10 w-12 rotate-12 rounded-md sm:right-0 sm:bottom-9 sm:w-14 md:right-5 md:bottom-10 md:w-16 dark:brightness-[1] dark:drop-shadow-[0_4px_40px_rgba(255,255,255,0.3)]"
              />
            </span>
            <span className="dark:drop-shadow-[0_4px_12px_rgba(0,0,0,1)]">
              Problems
            </span>
          </span>
        </div>
      </div>

      <div className="my-auto">
        <div className="font-geist flex w-fit flex-col items-center px-6 py-8">
          <div className="mx-auto flex w-fit items-center gap-4 font-mono">
            <Link
              to="https://github.com/snehasishroy/leetcode-companywise-interview-questions"
              target="_blank"
            >
              <Button
                size="lg"
                className="bg-muted/50 mt-4 cursor-pointer text-lg font-normal shadow-none"
                variant="outline"
              >
                <Database />
                Database
              </Button>
            </Link>
            <Link
              to="https://github.com/shadil-rayyan/visor-leetcode"
              target="_blank"
            >
              <Button
                size="lg"
                className="bg-muted/50 mt-4 cursor-pointer text-lg font-normal shadow-none"
                variant="outline"
              >
                <Globe />
                Website
              </Button>
            </Link>
          </div>
          <Link to="/all-problems" className="w-full">
            <Button
              size="lg"
              className="font-space-grotesk mt-4 w-full cursor-pointer text-lg shadow-none"
              variant="default"
            >
              Start Solving <Code strokeWidth={3} />
            </Button>
          </Link>
        </div>
        <div className="flex w-full flex-1 flex-col">
          <div className="font-geist flex w-full flex-col items-center gap-4">
            <div className="mx-auto flex max-w-4xl justify-center gap-16 text-center">
              <div className="flex flex-row items-center gap-2">
                <div className="text-xl font-medium tracking-tight tabular-nums">
                  {companies.toLocaleString()}
                </div>
                <div className="text-muted-foreground text-sm">Companies</div>
              </div>
              <div className="flex flex-row items-center gap-2">
                <div className="text-xl font-semibold tabular-nums">
                  {problems.toLocaleString()}
                </div>
                <div className="text-muted-foreground text-sm">Problems</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="mt-auto w-full">
        <CompanyCarousel />
      </div>
    </main>
  );
}
