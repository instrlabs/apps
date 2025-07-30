// Define the menu item types
export type IconName =
  | "HomeIcon"
  | "ClockIcon"
  | "PhotoIcon"
  | "ArrowsPointingInIcon"
  | "ArrowsPointingOutIcon"
  | "ScissorsIcon"
  | "ArrowPathIcon"
  | "ArrowPathRoundedSquareIcon"
  | "PencilSquareIcon";

export type BaseMenuItem = {
  type: string;
  icon?: IconName;
};

export type TextMenuItem = BaseMenuItem & {
  type: "text";
  value: string;
  href: string;
};

export type ParentMenuItem = BaseMenuItem & {
  type: "parent";
  value: string;
};

export type DividerMenuItem = BaseMenuItem & {
  type: "divider";
};

export type MenuItem = TextMenuItem | ParentMenuItem | DividerMenuItem;

export const sidebar_menus: MenuItem[] = [
  {
    type: "text",
    value: "Home",
    href: "/",
    icon: "HomeIcon",
  },
  {
    type: "text",
    value: "Histories",
    href: "/histories",
    icon: "ClockIcon",
  },
  {
    type: "divider",
  },
  {
    type: "parent",
    value: "Image",
  },
  {
    type: "text",
    value: "Compress",
    href: "/compress-image",
    icon: "ArrowsPointingInIcon",
  },
  {
    type: "text",
    value: "Resize",
    href: "/resize-image",
    icon: "ArrowsPointingOutIcon",
  },
  {
    type: "text",
    value: "Crop",
    href: "/crop-image",
    icon: "ScissorsIcon",
  },
  {
    type: "text",
    value: "Convert",
    href: "/convert-image",
    icon: "ArrowPathIcon",
  },
  {
    type: "text",
    value: "Rotate",
    href: "/rotate-image",
    icon: "ArrowPathRoundedSquareIcon",
  },
];
