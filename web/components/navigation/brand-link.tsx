import Link from "next/link";
import Icon from "@/components/icon";

export default function BrandLink() {
  return (
    <Link
      href="/"
      aria-label="Go to home"
      className="inline-flex items-center gap-2 hover:opacity-90 transition-opacity"
    >
      <Icon name="logo" size={40} />
    </Link>
  );
}
