import clsx from "clsx";
import Image from "next/image";
import Link from "next/link";
import { sidebar_menus, IconName, TextMenuItem, ParentMenuItem } from "@/constants/routes";
import Avatar from "./avatar";
import {
  HomeIcon,
  ClockIcon,
  PhotoIcon,
  ArrowsPointingInIcon,
  ArrowsPointingOutIcon,
  ScissorsIcon,
  ArrowPathIcon,
  ArrowPathRoundedSquareIcon,
  PencilSquareIcon,
} from "@heroicons/react/24/outline";

const getIconComponent = (iconName: IconName) => {
  switch (iconName) {
    case "HomeIcon":
      return HomeIcon;
    case "ClockIcon":
      return ClockIcon;
    case "PhotoIcon":
      return PhotoIcon;
    case "ArrowsPointingInIcon":
      return ArrowsPointingInIcon;
    case "ArrowsPointingOutIcon":
      return ArrowsPointingOutIcon;
    case "ScissorsIcon":
      return ScissorsIcon;
    case "ArrowPathIcon":
      return ArrowPathIcon;
    case "ArrowPathRoundedSquareIcon":
      return ArrowPathRoundedSquareIcon;
    case "PencilSquareIcon":
      return PencilSquareIcon;
    default:
      return HomeIcon; // Default icon
  }
};

const MenuItem = ({ item }: { item: TextMenuItem }) => {
  const Icon = item.icon ? getIconComponent(item.icon) : null;

  return (
    <Link
      href={item.href}
      className="flex items-center py-2.5 px-6 hover:bg-gray-200 transition-colors font-medium text-sm"
    >
      {Icon && <Icon className="h-5 w-5 mr-2" />}
      {item.value}
    </Link>
  );
};

const MenuParent = ({ item }: { item: ParentMenuItem }) => {
  const Icon = item.icon ? getIconComponent(item.icon) : null;

  return (
    <div className="flex items-center text-xs font-bold uppercase tracking-wider py-2 px-6 mt-2">
      {Icon && <Icon className="h-4 w-4 mr-2" />}
      {item.value}
    </div>
  );
};

// Menu item component for type "divider"
const MenuDivider = () => {
  return <div className="h-2 w-full"></div>;
};

const SidebarLeft = () => {
  return (
    <div
      className={clsx(
        "fixed inset-x-0 z-10",
        "h-screen w-[300px] border-r border-gray-300",
        "flex flex-col"
      )}
    >
      <div className="px-6 h-[56px] flex items-center">
        <Image alt="Logo" src="/logo.svg" width={84} height={24} />
      </div>
      <div>
        <nav className="space-y-1 py-2">
          {sidebar_menus.map((item, index) => {
            if (item.type === "text") {
              return <MenuItem key={index} item={item as TextMenuItem} />;
            } else if (item.type === "divider") {
              return <MenuDivider key={index} />;
            } else if (item.type === "parent") {
              return <MenuParent key={index} item={item as ParentMenuItem} />;
            }
            return null;
          })}
        </nav>
      </div>
      <div className="mt-auto py-2">
        <div className="hover:bg-gray-200 transition-colors">
          <div className="flex items-center space-x-2 px-6 py-2.5">
            <Avatar character="Artha Suryawan" size="md" />
            <div className="flex flex-col justify-center">
              <span className="font-medium text-sm">Artha Suryawan</span>
              <span className="text-gray-500 text-xs">My Workspace</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SidebarLeft;
