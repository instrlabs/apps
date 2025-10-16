import { getImageInstructions } from "@/services/images";
import { redirect } from "next/navigation";
import { format } from "date-fns";

export default async function History() {
  const { success, data } = await getImageInstructions();
  // if (!success) redirect("/");

  return (
    <div className="flex w-[300px] flex-col gap-4">
      <div className="flex flex-col gap-2">
        <h4 className="text-sm">Histories</h4>
      </div>
      <div className="bg-primary-black shadow-primary rounded-lg">
        <div className="flex flex-col gap-2 p-4">
          {data!.instructions.map((instruction) => (
            <div
              key={instruction.id}
              className="group flex items-center justify-between rounded-lg cursor-pointer"
            >
              <span className="text-sm text-white/70 group-hover:text-white/40">
                {format(new Date(instruction.created_at), "HH:mm dd/MM")}
              </span>
              <button
                className="text-sm text-gray-400 group-hover:text-white/40"
                aria-label="Delete"
              >
                ✕
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
