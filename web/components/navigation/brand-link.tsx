import Link from "next/link";
import LogoIcon from "@/components/icons/logo-icon";

export default function BrandLink() {
  return (
    <Link
      href="/"
      aria-label="Go to home"
      className="inline-flex items-center gap-2 hover:opacity-90 transition-opacity"
    >
      <LogoIcon size={40} />
    </Link>
  );
}
