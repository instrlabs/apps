export default function AppsPage() {
  const tools = [
    {
      key: 'compress',
      title: 'Compress',
      desc: 'Reduce image file size while preserving quality.',
      href: '#',
      icon: '🗜️',
    },
    {
      key: 'resize',
      title: 'Resize',
      desc: 'Change width/height, keep aspect ratio or set exact.',
      href: '#',
      icon: '📐',
    },
    {
      key: 'crop',
      title: 'Crop',
      desc: 'Trim images to a selected region or aspect ratio.',
      href: '#',
      icon: '✂️',
    },
    {
      key: 'convert',
      title: 'Convert',
      desc: 'Convert between JPG, PNG, WEBP, and more.',
      href: '#',
      icon: '🔄',
    },
    {
      key: 'rotate',
      title: 'Rotate',
      desc: 'Rotate images by 90°, 180°, or custom angles.',
      href: '#',
      icon: '🧭',
    },
    {
      key: 'flip',
      title: 'Flip',
      desc: 'Flip images horizontally or vertically.',
      href: '#',
      icon: '🔁',
    },
    {
      key: 'watermark',
      title: 'Watermark',
      desc: 'Add text or image watermarks with controls.',
      href: '#',
      icon: '💧',
    },
    {
      key: 'optimize',
      title: 'Optimize',
      desc: 'Auto-optimize images for web performance.',
      href: '#',
      icon: '⚡',
    },
  ];

  return (
    <div className="container mx-auto p-4">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Image Tools</h1>
        <p className="text-sm text-muted">Quick utilities to process your images.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {tools.map((tool) => (
          <a
            key={tool.key}
            href={tool.href}
            className="group block rounded-xl bg-card border border-border p-4"
          >
            <div className="flex items-start gap-3">
              <div className="text-2xl leading-none select-none">
                <span aria-hidden>{tool.icon}</span>
              </div>
              <div>
                <h2 className="font-medium group-hover:text-primary">{tool.title}</h2>
                <p className="mt-1 text-sm text-muted">{tool.desc}</p>
              </div>
            </div>
          </a>
        ))}
      </div>
    </div>
  );
}
