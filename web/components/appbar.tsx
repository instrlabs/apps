import clsx from "clsx";
import { ChatBubbleOvalLeftIcon } from "@heroicons/react/24/outline";
import OutlinedButton from "./outlined-button";

const AppBar = () => {
  return (
    <div className={clsx("absolute h-[56px] w-full", "border-b border-gray-300", "pl-[300px]")}>
      <div className="h-[56px] px-6 flex items-center">
        <h1 className="text-lg font-bold">Home</h1>
        <div className="ml-auto flex items-center">
          <OutlinedButton icon={<ChatBubbleOvalLeftIcon className="h-4 w-4 text-black" />}>
            Feedback
          </OutlinedButton>
        </div>
      </div>
    </div>
  );
};

export default AppBar;
