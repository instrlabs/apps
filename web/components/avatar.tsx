import React from "react";

interface AvatarProps {
  character: string;
  size?: "sm" | "md" | "lg";
}

const Avatar: React.FC<AvatarProps> = ({ character, size = "md" }) => {
  // Get the first character and convert to uppercase
  const displayChar = character.charAt(0).toUpperCase();

  // Generate a deterministic color based on the character
  const charCode = displayChar.charCodeAt(0);

  // Define color pairs (background and text colors)
  const colorPairs = [
    { bg: "bg-red-200", text: "text-red-700" },
    { bg: "bg-blue-200", text: "text-blue-700" },
    { bg: "bg-green-200", text: "text-green-700" },
    { bg: "bg-yellow-200", text: "text-yellow-700" },
    { bg: "bg-purple-200", text: "text-purple-700" },
    { bg: "bg-pink-200", text: "text-pink-700" },
    { bg: "bg-indigo-200", text: "text-indigo-700" },
    { bg: "bg-gray-200", text: "text-gray-700" },
  ];

  // Select color pair based on character code
  const colorIndex = charCode % colorPairs.length;
  const { bg, text } = colorPairs[colorIndex];

  // Determine size classes
  const sizeClasses = {
    sm: "h-8 w-8 text-xs",
    md: "h-10 w-10 text-sm",
    lg: "h-14 w-14 text-base",
  };

  return (
    <div
      className={`flex items-center justify-center rounded-full ${bg} ${text} font-bold border border-white shadow-sm ${sizeClasses[size]}`}
    >
      {displayChar}
    </div>
  );
};

export default Avatar;
