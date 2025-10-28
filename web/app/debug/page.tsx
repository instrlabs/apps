export default function DebugPage() {
  return(
    <div className="container mx-auto p-10">
      <div className="flex flex-col gap-6">
        <div className="bg-white/5 border border-white/20 rounded-lg p-2">
          <span className="text-sm">Something went wrong</span>
        </div>
        <div className="bg-red-500/10 border border-red-500 rounded-lg p-2">
          <span className="text-sm text-red-500">Something went wrong</span>
        </div>
        <div className="bg-blue-500/10 border border-blue-500 rounded-lg p-2">
          <span className="text-sm text-blue-500">Something went wrong</span>
        </div>
        <div className="bg-yellow-500/10 border border-yellow-500 rounded-lg p-2">
          <span className="text-sm text-yellow-500">Something went wrong</span>
        </div>
        <div className="bg-green-500/10 border border-green-500 rounded-lg p-2">
          <span className="text-sm text-green-500">Something went wrong</span>
        </div>
      </div>
    </div>
  )
}
