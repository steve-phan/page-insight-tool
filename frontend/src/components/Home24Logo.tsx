"use client";

import { useRouter } from "next/navigation";
import Image from "next/image";

interface LogoProps {
  className?: string;
}

export function Home24Logo({ className }: LogoProps) {
  const router = useRouter();

  const handleClick = () => {
    // Reset form by navigating to home without URL parameter
    router.push("/");
  };

  return (
    <button
      onClick={handleClick}
      className={`transition-opacity hover:opacity-80 focus:outline-none cursor-pointer ${className}`}
      aria-label="Reset form and go to home"
    >
      <Image
        src="/images/home24-logo.svg"
        alt="Home24 Logo"
        width={139}
        height={50}
        className="h-8 w-auto"
        priority
      />
    </button>
  );
}
